# discord-ping

[![Build Status](https://github.com/Kishan-Thanki/discord-ping/actions/workflows/ci.yml/badge.svg)](https://github.com/Kishan-Thanki/discord-ping/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance, lightweight Discord diagnostic and utility bot built in Go. Engineered to strictly measure WebSocket heartbeat and API round-trip latencies with near-zero overhead.

## Project Philosophy: Production-Ready Diagnostics

NOTE: This repository is fundamentally an *Exploration and Learning Project* engineered to strict production-ready standards. It was built to demonstrate high-performance Go patterns and zero-allocation memory management.

## Features

* *Zero-Allocation Core*: Extensively optimized using `strconv` and `strings.Builder` to bypass `fmt` reflection overhead, ensuring near-zero garbage collection pauses.
* *Precise Diagnostics*: Instantly calculates Discord API round-trip latency, WebSocket heartbeat, and message transit times.
* *Lightweight Design*: Stripped down to the bare essentials, ensuring instant startup times and minimal memory footprint.

## External Dependencies

This project relies on the following external libraries:

. https://github.com/bwmarrin/discordgo[*DiscordGo*]: Handles the complex WebSocket connections to Discord's gateway and provides bindings for the Discord REST API.
. https://gitlab.com/cznic/sqlite[*modernc.org/sqlite*]: A CGo-free SQLite driver used for storing server-specific configurations (like custom prefixes).

## Usage

For a complete list of all available commands, please refer to the [Commands Guide](COMMANDS.md).

## Policies

Please review our [Privacy Policy](PRIVACY_POLICY.md) and [Terms of Service](TERMS_OF_SERVICE.md) regarding data collection and moderation actions.

## License

This project is licensed under the https://opensource.org/licenses/MIT[MIT License].
