# Future Implementation Plans

This document outlines planned future enhancements and architectural extensions for the `discord-ping` bot. 

---

## 1. Prometheus & Grafana Telemetry (Priority)

Expose a `/metrics` HTTP endpoint to allow real-time monitoring of the bot's performance, health, and latency statistics using Prometheus and Grafana.

### Metrics to Track
* **`discord_ping_websocket_latency_milliseconds`** (Gauge)
  * Tracks the current WebSocket heartbeat latency between the bot process and Discord Gateway.
* **`discord_ping_commands_total`** (Counter)
  * Labeled by `command` (e.g., `ping`, `about`, `help`, `setprefix`) and `status` (`success`, `error`).
  * Tracks overall usage patterns.
* **`discord_ping_ratelimit_hits_total`** (Counter)
  * Labeled by `channel_id` or `guild_id`.
  * Monitors user spam rate and identifies if specific servers are hitting limits excessively.
* **`discord_ping_api_roundtrip_milliseconds`** (Histogram)
  * Tracks the distribution of REST HTTP request round-trip speeds (e.g., time to send/edit messages).

### Proposed Architecture
1. Import `github.com/prometheus/client_golang/prometheus/promhttp`.
2. Start a background HTTP listener inside `main.go` on a dedicated telemetry port (e.g., `:8080`):
   ```go
   go func() {
       http.Handle("/metrics", promhttp.Handler())
       if err := http.ListenAndServe(":8080", nil); err != nil {
           slog.Error("Metrics server failed", "error", err)
       }
   }()
   ```
3. Update `isRateLimited` and `cmdPing` to increment counters/gauges on event execution.

