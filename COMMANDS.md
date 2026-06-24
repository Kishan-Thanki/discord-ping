# Ping-Bot Commands Guide

Welcome to the official command guide for Ping-Bot!

By default, the bot listens for the `!` prefix (unless you changed `BOT_PREFIX` in your `.env` file). Here is a complete list of everything the bot can do.

## General & Utility

* *`!ping`* +
  Tests the bot's connection. The bot will instantly reply with a "pong" and the exact websocket latency in milliseconds.
* *`!help`* +
  Displays a brief in-game help menu with the main categories.
* *`!version`* (or `!about`) +
  Shows the bot's current version and description.
* *`!uptime`* +
  Shows how long the bot has been running since its last restart.
* *`!serverinfo`* +
  Displays detailed information about the current server (name, owner, member count, channels, roles, and creation date).
* *`!avatar [@user]`* +
  Displays the full-size avatar of a user. If no user is mentioned, it shows your own avatar.
* *`!remind <duration> <message>`* +
  Sets a reminder. The bot will DM you after the specified time. +
  _Example:_ `!remind 10m Check the oven` +
  _Example:_ `!remind 2h Finish homework`

## Economy & Leveling

* *`!daily`* +
  Claims your daily allowance of 100 free coins. Can be used once every 24 hours.
* *`!balance [@user]`* +
  Shows the current coin balance of yourself or a mentioned user.
* *`!give <@user> <amount>`* +
  Transfers coins from your balance to another user. +
  _Example:_ `!give @Kishan 50`
* *`!leaderboard`* +
  Displays the top 10 users in the server ranked by XP.
* *`!rank`* +
  Displays your current level, total XP, and progress. (You earn XP automatically by sending messages).

## Games & Gambling

* *`!coinflip`* +
  Flips a coin — heads or tails. Pure 50/50 chance.
* *`!roll [max]`* +
  Rolls a dice. Defaults to 1-6 unless you specify a max value. +
  _Example:_ `!roll 20` (rolls 1-20)
* *`!slots <amount>`* +
  Bet your coins on the slot machine. Match 3 emojis for a massive 10x payout, or 2 emojis for a 2x payout. +
  _Example:_ `!slots 50`
* *`!blackjack <amount>`* +
  Play a fully interactive game of Blackjack against the dealer using Discord buttons (Hit/Stand). Win to double your bet! +
  _Example:_ `!blackjack 200`
* *`!wordle`* +
  Starts a new game of Wordle.
* *`!guess <5-letter-word>`* +
  Make a guess in your active Wordle game. The bot will respond with 🟩, 🟨, and ⬛ emojis. Win in fewer guesses for a bigger coin prize! +
  _Example:_ `!guess apple`
* *`!trivia`* +
  Starts a timed multiple-choice trivia question. Type the correct number to earn XP and coins!
* *`!8ball <question>`* +
  Ask the Magic 8-Ball a question and receive a mystical answer. +
  _Example:_ `!8ball Will I pass the exam?`
* *`!poll "Question" "Option 1" "Option 2" ...`* +
  Creates an interactive poll with emoji reactions. Max 10 options. +
  _Example:_ `!poll "Best language?" "Go" "Rust" "Python"`

## Moderation

_(Note: These commands require the user to have the relevant moderation permissions in the server)._

* *`!warn <@user> [reason]`* +
  Issues a formal warning to a user and logs it in the database. If a user receives 3 warnings, the bot will automatically timeout them for 10 minutes. +
  _Example:_ `!warn @Spammer Stop posting links`
* *`!warnings [@user]`* +
  Shows the total number of warnings for yourself or a mentioned user.
* *`!kick <@user> [reason]`* +
  Kicks the user from the server. Requires **Kick Members** permission.
* *`!ban <@user> [reason]`* +
  Permanently bans the user from the server. Requires **Ban Members** permission.
* *`!mute <@user>`* +
  Mutes the user for 10 minutes (they cannot send messages or join voice). Requires **Moderate Members** permission.
* *`!setprefix <new_prefix>`* +
  Changes the bot's command prefix for this server. Requires **Administrator** permission. +
  _Example:_ `!setprefix ?` (commands now use `?ping`, `?help`, etc.)

## Passive Features

You don't need to type commands for these to work:

* *Welcome Images*: Whenever someone joins the server, the bot automatically generates and uploads a beautiful custom PNG image welcoming them.
* *Auto-Mod*: The bot silently watches the chat. If anyone types a banned word or spams, the bot instantly deletes the message and issues a warning.
* *Passive XP*: Every message you send earns 15-25 XP. Level up automatically and get a congratulations embed!
