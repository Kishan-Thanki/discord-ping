# Future Roadmap & Architectural Audit

This document outlines the long-term roadmap and architectural critiques for the `discord-ping` bot based on Discord platform shifts and Go cloud-native development standards.

---

## Deeper Critiques & Future Progression

### 1. The "Message Intent" Critique (Crucial for Verification)

Currently, in `internal/bot/bot.go`, the bot requests `discordgo.IntentsMessageContent`. Furthermore, the database is primarily being used to store custom text-command prefixes (`SetPrefix`, `GetPrefix`).

- **The Problem:** Discord is aggressively phasing out text-based commands to protect user privacy.
- **The Fix:** **Deprecate text commands entirely.** Will move 100% to Slash Commands (`/ping`, `/about`). This allows to drop the Message Content intent, guarantees instant approval from Discord, and will completely eliminate the need for the `prefix` database table—simplifying architecture.

### 2. The Observability Critique (The "Pro" Feature)

A true diagnostic tool shouldn't just print metrics to a chat window;
it should integrate with standard infrastructure monitoring.

- **The Problem:** The latency is only visible point-in-time when a user explicitly runs the command.
- **The Fix:** Adhere strictly to Go backend standards by adding a **Prometheus Metrics Endpoint**. Expose an HTTP server on port `:8080` with a `/metrics` route that continuously exports the Discord WebSocket heartbeat and API latency. This would allow server owners to scrape the bot and graph Discord's real-time latency on their own Grafana dashboards.

### 3. The UX Critique (Ephemeral Diagnostics)

When a server is lagging, administrators tend to spam the `/ping` command repeatedly to see if latency is stabilizing.

- **The Problem:** The current `cmdPing` implementation sends a public embed. Spamming this command will completely clog up a community chat channel.
- **The Fix:** Will make the Slash Command response **Ephemeral** by default (only visible to the user who executed it) using `discordgo.MessageFlagsEphemeral`. Could add an optional parameter to the slash command (e.g., `/ping public:True`) for when users actually want to show off the bot's speed to the channel.

---

## Roadmap Phases

1. **Phase 1 (Current):** Accurate, 3-way latency splitting (WS, API, Transit) with rate-limiting.
2. **Phase 2:** Complete migration to Slash Commands, removal of Message Content Intent, and elimination of the prefix database.
3. **Phase 3:** Integration of Prometheus `/metrics` for external Grafana observability.
4. **Phase 4:** Expanding diagnostics to external IPs and HTTP endpoints.
