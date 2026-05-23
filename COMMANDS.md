# Ping-Bot Commands Guide

Welcome to the official command guide for Ping-Bot! 

By default, the bot listens for the `!` prefix (unless you changed `BOT_PREFIX` in your `.env` file). Here is a complete list of everything the bot can do.

## General & Utility

- **`!ping`**  
  Tests the bot's connection. The bot will instantly reply with a "pong" and the exact websocket latency in milliseconds.
- **`!help`**  
  Displays a brief in-game help menu with the main categories.
- **`!remind <duration> <message>`**  
  Sets a reminder. The bot will ping you after the specified time.  
  *Example:* `!remind 10m Check the oven`  
  *Example:* `!remind 2h Finish homework`
- **`!reminders`**  
  Lists all of your currently active reminders.
- **`!delreminder <id>`**  
  Deletes a specific reminder by its ID (which you can find using `!reminders`).

## Economy & Leveling

- **`!daily`**  
  Claims your daily allowance of free coins. Can be used once every 24 hours.
- **`!balance`** (or `!bal`)  
  Shows your current coin balance.
- **`!give <@user> <amount>`**  
  Transfers coins from your balance to another user.  
  *Example:* `!give @Kishan 50`
- **`!leaderboard`**  
  Displays the top 10 richest users in the server.
- **`!rank`**  
  Displays your current level, total XP, and progress to the next level. (You earn XP automatically by sending messages).

## Games & Gambling

- **`!coinflip <amount> <heads/tails>`**  
  Bet your coins on a 50/50 coin flip. Guess correctly to double your bet!  
  *Example:* `!coinflip 100 heads`
- **`!slots <amount>`**  
  Bet your coins on the slot machine. Match 3 emojis for a massive 10x payout, or 2 emojis for a 2x payout.  
  *Example:* `!slots 50`
- **`!blackjack <amount>`**  
  Play a fully interactive game of Blackjack against the dealer using Discord buttons (Hit/Stand). Win to double your bet!  
  *Example:* `!blackjack 200`
- **`!wordle start`**  
  Starts a new game of Wordle.
- **`!guess <5-letter-word>`**  
  Make a guess in your active Wordle game. The bot will respond with 🟩, 🟨, and ⬛ emojis. Win in 6 guesses for a massive coin prize!  
  *Example:* `!guess apple`
- **`!trivia`**  
  Starts a multiple-choice trivia question. Click the correct button to earn coins!

## Moderation

*(Note: These commands require the user to have the "Administrator" permission in the server).*

- **`!warn <@user> <reason>`**  
  Issues a formal warning to a user and logs it in the database. If a user receives 3 warnings, the bot will automatically timeout/mute them.  
  *Example:* `!warn @Spammer Stop posting links`
- **`!kick <@user> <reason>`**  
  Kicks the user from the server.
- **`!ban <@user> <reason>`**  
  Permanently bans the user from the server.
- **`!timeout <@user> <duration> <reason>`**  
  Mutes the user so they cannot send messages or join voice channels for the specified duration.  
  *Example:* `!timeout @Spammer 10m Chill out`

## Passive Features

You don't need to type commands for these to work:
- **Welcome Images**: Whenever someone joins the server, the bot automatically generates and uploads a beautiful custom PNG image welcoming them.
- **Auto-Mod**: The bot silently watches the chat. If anyone types a banned word or spams, the bot instantly deletes the message and issues a warning.
