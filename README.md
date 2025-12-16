# Zee-Bot 🤖

A powerful Telegram Userbot built with Pyrogram, running in Docker.

## Features

- ⚡ **Fast & Async** - Built with Pyrogram and uvloop for maximum performance
- 🐳 **Dockerized** - Easy deployment with Docker Compose
- 🔌 **Plugin System** - Modular plugin architecture for easy extension
- 🔐 **String Session** - No interactive login needed, perfect for cloud deployment

## Available Commands

| Command | Description |
|---------|-------------|
| `.alive` | Check if the bot is running and show ping |
| `.purge` | Delete messages from replied message to current |

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Telegram API credentials from [my.telegram.org](https://my.telegram.org)
- Python 3.10+ (for generating session string)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/ifauzeee/Zee-Bot.git
   cd Zee-Bot
   ```

2. **Generate Session String**
   
   First, install Pyrogram locally:
   ```bash
   pip install pyrogram tgcrypto
   ```
   
   Then run the session generator:
   ```bash
   python gen_session.py
   ```
   - Enter your API ID & API Hash
   - Enter your phone number & OTP code
   - Copy the generated session string

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` and fill in your credentials:
   - `API_ID` - Your Telegram API ID
   - `API_HASH` - Your Telegram API Hash
   - `SESSION_STRING` - Paste the session string from step 2

4. **Build and run**
   ```bash
   docker-compose build
   docker-compose up -d
   ```
   The bot will start immediately without asking for login!

### View Logs

```bash
docker-compose logs -f
```

### Stop the Bot

```bash
docker-compose down
```

## Project Structure

```
Zee-Bot/
├── main.py              # Main entry point
├── config.py            # Configuration loader
├── gen_session.py       # Session string generator
├── plugins/             # Plugin modules
│   ├── alive.py         # Alive check command
│   └── purge.py         # Message purge command
├── Dockerfile           # Docker image definition
├── docker-compose.yml   # Docker Compose config
├── requirements.txt     # Python dependencies
└── .env.example         # Environment template
```

## Adding New Plugins

Create a new file in the `plugins/` directory:

```python
from pyrogram import Client, filters

@Client.on_message(filters.command("mycommand", prefixes=".") & filters.me)
async def my_command(client, message):
    await message.edit("Hello from my plugin!")
```

## License

MIT License - Feel free to use and modify!

## Author

Made with ❤️ by [ifauzeee](https://github.com/ifauzeee)
