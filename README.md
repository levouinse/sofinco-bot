# Sofinco Bot

> Modern Telegram bot built with Go featuring modular plugin system and extensive features

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build](https://img.shields.io/badge/Build-Passing-success)](https://github.com/levouinse/sofinco-bot)

## ✨ Features

### 📥 Downloader
- **YouTube** - Download audio/video from YouTube
- **TikTok** - Download videos without watermark
- **Instagram** - Download posts and reels

### 🎮 Games
- **Math Quiz** - Multiple difficulty levels (noob to god)
- **Tic Tac Toe** - Classic game
- **Rock Paper Scissors** - Play with bot
- **Family 100** - Guess the answers

### 🛠 Tools
- **Remini** - AI image enhancer
- **Anime Converter** - Convert photos to anime style
- **Sticker Maker** - Create custom stickers

### 🔍 Stalker
- **Instagram Stalker** - View profile info
- **GitHub Stalker** - View GitHub profiles
- **TikTok Stalker** - View TikTok profiles

### 🤖 AI
- **ChatGPT** - Chat with AI
- **Ask AI** - Get answers from AI

### 👑 Owner Commands
- **Broadcast** - Send messages to all users
- **Add Premium** - Grant premium access
- **Statistics** - View bot stats

## 🚀 Quick Start

```bash
git clone https://github.com/levouinse/sofinco-bot.git
cd sofinco-bot
cp .env.example .env
# Edit .env with your credentials
make install && make build
./sofinco-bot
```

## Commands

### General
| Command | Description |
|---------|-------------|
| `/start` | Show main menu with buttons |
| `/menu` | Display menu categories |
| `/ping` | Check bot status and uptime |
| `/getid` | Get your user and chat ID |
| `/limit` | Check your daily limit |
| `/profile` | View your profile and stats |

### Downloader
| Command | Description |
|---------|-------------|
| `/play <query>` | Download audio from YouTube |
| `/ytv <url>` | Download video from YouTube |
| `/tiktok <url>` | Download TikTok video |

### Games
| Command | Description |
|---------|-------------|
| `/math <mode>` | Math quiz (noob/easy/medium/hard/master/grandmaster/legendary/mythic/god) |

### Tools
| Command | Description |
|---------|-------------|
| `/remini` | Enhance image quality (reply to photo) |
| `/jadianime` | Convert photo to anime style (reply to photo) |

### Stalker
| Command | Description |
|---------|-------------|
| `/igstalk <username>` | View Instagram profile |
| `/ghstalk <username>` | View GitHub profile |

### AI
| Command | Description |
|---------|-------------|
| `/ai <question>` | Chat with AI |
| `/ask <question>` | Ask AI anything |

### Owner Only
| Command | Description |
|---------|-------------|
| `/broadcast` | Send message to all users (reply to message) |
| `/addprem <user_id>` | Grant premium access |

## Requirements

- Go 1.21 or higher
- Telegram Bot Token from [@BotFather](https://t.me/BotFather)
- BetaBotz API Key from [api.betabotz.eu.org](https://api.betabotz.eu.org)

## Installation

See [INSTALL.md](INSTALL.md) for detailed installation instructions.

### Docker

```bash
docker-compose up -d
```

### Manual

```bash
make install
make build
./sofinco-bot
```

## Configuration

Create `.env` file:

```env
BOT_TOKEN=your_telegram_bot_token
OWNER_ID=your_telegram_user_id
OWNER_USERNAME=your_username
API_KEY=your_betabotz_api_key
AKSES_KEY=your_akses_key
```

## Project Structure

```
sofinco-bot/
├── cmd/bot/main.go              # Entry point
├── internal/
│   ├── bot/                     # Bot core
│   ├── config/                  # Configuration
│   ├── database/                # Database operations
│   ├── downloader/              # YouTube downloader
│   └── handlers/                # Message handlers
├── Dockerfile                   # Docker build
├── docker-compose.yml           # Docker Compose
└── Makefile                     # Build automation
```

## Development

```bash
make dev        # Run in development mode
make test       # Run tests
make clean      # Clean build artifacts
```

## Deployment

### Production Build

```bash
make deploy
```

### Systemd Service

```bash
sudo systemctl enable sofinco-bot
sudo systemctl start sofinco-bot
```

### Docker

```bash
docker-compose up -d
```

## Documentation

- [Installation Guide](INSTALL.md) - Comprehensive installation instructions
- [Project Structure](STRUCTURE.md) - Architecture documentation
- [Quick Reference](QUICKREF.md) - Command reference
- [Summary](SUMMARY.md) - Complete rewrite summary

## Author

**levouinse**
- GitHub: [@levouinse](https://github.com/levouinse)
- Repository: [sofinco-bot](https://github.com/levouinse/sofinco-bot)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Support

For issues and questions, please open an issue on [GitHub](https://github.com/levouinse/sofinco-bot/issues).

---

Made with ❤️ by [levouinse](https://github.com/levouinse)
