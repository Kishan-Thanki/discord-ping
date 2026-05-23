# go-discord-ping

[![Build Status](https://github.com/Kishan-Thanki/discord-ping/actions/workflows/ci.yml/badge.svg)](https://github.com/Kishan-Thanki/discord-ping/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kishan-Thanki/discord-ping)](https://goreportcard.com/report/github.com/Kishan-Thanki/discord-ping)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A hyper-optimized, zero-allocation Discord utility bot built in Go. Features native image rendering, auto-moderation, leveling, a virtual economy, and lightning-fast SQLite prepared statements.

## Project Philosophy: Production-Ready Learning

> [!NOTE]
> This repository is fundamentally an **Exploration and Learning Project** engineered to strict production-ready standards. It was built to demonstrate high-performance Go patterns, zero-allocation memory management, and advanced concurrency.
>
> If you find flaws, want to fix bugs, or simply want to learn how these complex systems interact, you are highly encouraged to explore the codebase, open issues, and submit pull requests. Please consider this a living, educational repository for building production-grade Discord applications!

## Features

- **Zero-Allocation Core**: Extensively optimized using `strconv` and direct string concatenation to bypass `fmt` reflection overhead, ensuring near-zero garbage collection pauses.
- **SQLite Database**: Uses a local, high-performance WAL-mode SQLite database with prepared statements to track users, economies, and warnings per server.
- **Virtual Economy & Leveling**: Users earn XP and Coins by participating in chat. Features include `!daily`, `!balance`, `!leaderboard`, `!rank`, `!coinflip`, and `!blackjack`.
- **Mini-Games**: Play fully interactive games like `!wordle` and trivia right in Discord. Safe for high-concurrency environments using `sync.Mutex` locks.
- **Native Image Rendering**: Generates beautiful custom welcome images on-the-fly using `fogleman/gg` whenever a new user joins a server.
- **Auto-Moderation**: Includes spam/bad-word filters and a strict 3-strike warning system (`!warn`) that automatically timeouts repeat offenders.

## External Dependencies

This project relies on the following external libraries:

1. **[DiscordGo](https://github.com/bwmarrin/discordgo)**: Handles the complex WebSocket connections to Discord's gateway and provides bindings for the Discord REST API.
2. **[godotenv](https://github.com/joho/godotenv)**: Loads environment variables from the `.env` file.
3. **[modernc.org/sqlite](https://gitlab.com/cznic/sqlite)**: A CGo-free SQLite driver for Go, ensuring lightning-fast local data storage without requiring a C compiler.
4. **[fogleman/gg](https://github.com/fogleman/gg)**: A 2D rendering library in Go used to draw the custom welcome images.

## Setup & Running

1. Ensure you have Go installed (v1.23.4 or higher).
2. Create a `.env` file in the root directory and add your bot token and default prefix:

   ```env
   TOKEN=your_discord_bot_token_here
   BOT_PREFIX=!
   ```

3. Obtain your `TOKEN` from the [Discord Developer Portal](https://discord.com/developers/applications).
4. Run the bot:

   ```bash
   go run main.go
   ```

5. Or build the optimized binary:

   ```bash
   go build -o ping-bot .
   ./ping-bot
   ```

## Usage

For a complete list of all available commands, games, and moderation tools, please refer to the [Commands Guide](COMMANDS.md).

## Contributing

We welcome community contributions! Please review our [Contributing Guidelines](CONTRIBUTING.md) to understand the strict performance and architectural standards required for this repository before submitting a Pull Request.

## Policies

Please review our [Privacy Policy](PRIVACY_POLICY.md) and [Terms of Service](TERMS_OF_SERVICE.md) regarding data collection and moderation actions.

## License

This project is licensed under the [MIT License](LICENSE).
