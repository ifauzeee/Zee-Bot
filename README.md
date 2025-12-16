<h1 align="center">
  <br>
  <a href="https://github.com/ifauzeee/Zee-Boat"><img src="https://img.icons8.com/clouds/200/000000/rocket.png" alt="Zee-Bot"></a>
  <br>
  ⚡️ Zee-Bot Ultimate <br>
</h1>

<p align="center">
  <a href="https://github.com/ifauzeee/Zee-Bot/stargazers">
    <img src="https://img.shields.io/github/stars/ifauzeee/Zee-Bot?style=social">
  </a>
  <a href="https://github.com/ifauzeee/Zee-Bot/network/members">
    <img src="https://img.shields.io/github/forks/ifauzeee/Zee-Bot?style=social">
  </a>
  <br>
  <a href="https://python.org">
    <img src="https://img.shields.io/badge/Python-3.10%2B-blue?style=for-the-badge&logo=python">
  </a>
  <a href="https://docs.pyrogram.org">
    <img src="https://img.shields.io/badge/Pyrogram-2.0-orange?style=for-the-badge&logo=telegram">
  </a>
  <a href="https://www.docker.com">
    <img src="https://img.shields.io/badge/Docker-Enabled-blue?style=for-the-badge&logo=docker">
  </a>
</p>

<p align="center">
  A powerful, aesthetic, and completely modular <b>Telegram Userbot</b> dedicated to enhancing your Telegram experience.<br>
  Built with performance and aesthetics in mind.
</p>

<hr>

## 🌟 Features

| 🔰 **Core** | 🛡 **Security** | 📹 **Media** | 🛠 **Utilities** |
| :--- | :--- | :--- | :--- |
| `Sync/Async` Architecture | **PM Permit** System | **Sticker** Stealer/Maker | **Dictionary** Lookup |
| **Plugin** Modular System | **Spam** Protection | **Carbon** Code Image | **Translator** (All Lang) |
| **Docker** Ready | **Admin** Tools (Ban/Kick) | **Screenshot** Website | **Currency** Converter |
| **String Session** Login | **AFK** Auto-Reply | **TTS** (Text to Speech) | **Notes** Saver |

## 🚀 Quick Start

### 1. Prerequisites
- [Python 3.10+](https://www.python.org/downloads/)
- [Docker](https://www.docker.com/) (Optional but Recommended)
- Telegram Account

### 2. Get Credentials
1. Go to [my.telegram.org](https://my.telegram.org)
2. Obtain `API_ID` and `API_HASH`

### 3. Deployment

#### 🐳 via Docker (Recommended)
This is the easiest way to run the bot 24/7.

```bash
# 1. Clone Repo
git clone https://github.com/ifauzeee/Zee-Bot.git
cd Zee-Bot

# 2. Generate Session
# Run this once on your local machine to login
pip install pyrogram tgcrypto
python gen_session.py

# 3. Configure
cp .env.example .env
# Edit .env and paste your API_ID, API_HASH, SESSION_STRING

# 4. Run
docker compose up -d --build
```

#### 🐍 via Local Python
```bash
# Install dependencies
pip install -r requirements.txt

# Configure .env as above

# Run
python main.py
```

## 🎮 Command List
Type `.help` in any chat to see the interactive menu.

### 🛡 Admin
- `.ban`, `.unban`, `.kick`, `.mute`
- `.pmpermit on/off` - Toggle security
- `.purge` - Delete messages

### 📹 Media
- `.kang` - Steal stickers
- `.carbon [code]` - Beautiful code snippets
- `.ss [url]` - Screenshot websites
- `.tts` / `.tr` - Audio & Translation

## 📂 Project Structure
```
Zee-Bot/
├── 📂 helpers/         # UI & Utility helper functions
├── 📂 plugins/         # 🧩 Modular plugins directory
│   ├── 📂 admin/       # Admin tools
│   ├── 📂 core/        # Core system (Alive, Help, PM Permit)
│   ├── 📂 media/       # Sticker, Carbon, Media tools
│   └── 📂 utils/       # Notes, AFK, Misc tools
├── 📄 config.py        # Safe configuration loader
├── 📄 main.py          # Bot entry point
└── 🐳 Dockerfile       # Container configuration
```

## 🤝 Contribution
Feel free to contribute!
1. Fork the repo
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License
Distributed under the MIT License. See `LICENSE` for more information.

<p align="center">
  Made with ❤️ by <a href="https://github.com/ifauzeee"><b>ifauzeee</b></a>
</p>
