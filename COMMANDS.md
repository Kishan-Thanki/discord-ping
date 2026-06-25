# Ping-Bot Commands Guide

Welcome to the official command guide for Ping-Bot!

By default, the bot listens for the `!` prefix (unless you changed `BOT_PREFIX` in your `.env` file or used the `!setprefix` command). Here is a complete list of everything the bot can do.

## General & Utility

* *`!ping`* +
  Tests the bot's connection. The bot will instantly reply with the exact WebSocket heartbeat latency, API round-trip time, and message transit time in milliseconds.
* *`!help`* +
  Displays a brief in-game help menu.
* *`!version`* (or `!about`) +
  Shows the bot's current version and description.
* *`!uptime`* +
  Shows how long the bot has been running since its last restart.

## Configuration

_(Note: These commands require the user to have the relevant moderation permissions in the server)._

* *`!setprefix <new_prefix>`* +
  Changes the bot's command prefix for this server. Requires **Administrator** permission. +
  _Example:_ `!setprefix ?` (commands now use `?ping`, `?help`, etc.)
