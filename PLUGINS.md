# Sofinco Bot - Plugin System

## 📦 Integrated Features

All features from the JavaScript version have been rewritten in Go with a modular plugin system.

### 📥 Downloader Plugins
- `/play` - Download audio from YouTube
- `/ytv` - Download video from YouTube  
- `/tiktok` - Download TikTok videos (no watermark)
- `/ytmp3` - YouTube to MP3 converter

### 🎮 Game Plugins
- `/math [mode]` - Math quiz with multiple difficulty levels
  - Modes: noob, easy, medium, hard, master, grandmaster, legendary, mythic, god
- `/tictactoe` - Play Tic Tac Toe
- `/suit` - Rock Paper Scissors
- `/family100` - Family 100 game

### 🛠 Tools Plugins
- `/remini` - AI image enhancer (reply to photo)
- `/getid` - Get your Telegram user ID
- `/sticker` - Create sticker from image
- `/toimg` - Convert sticker to image

### 🤖 AI Plugins
- `/ai [question]` - Chat with AI
- `/ask [question]` - Ask AI anything
- `/chatgpt [question]` - ChatGPT integration

### ℹ️ Info Commands
- `/start` - Show main menu
- `/menu` - Display command categories
- `/ping` - Check bot status and latency
- `/limit` - Check your daily limit
- `/profile` - View your profile and stats
- `/stats` - Bot statistics

## 🔧 Plugin System Architecture

### Structure
```
internal/plugins/
├── plugins.go          # Base plugin interface
├── downloader/         # Download plugins
│   ├── play.go
│   └── tiktok.go
├── game/               # Game plugins
│   └── math.go
├── tools/              # Utility plugins
│   └── remini.go
└── ai/                 # AI plugins
    └── openai.go
```

### Creating New Plugin

```go
package myplugin

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type MyPlugin struct {
	plugins.BasePlugin
}

func init() {
	plugins.Register(&MyPlugin{})
}

func (p *MyPlugin) Commands() []string { 
	return []string{"mycommand", "alias"} 
}

func (p *MyPlugin) Tags() []string { 
	return []string{"category"} 
}

func (p *MyPlugin) Help() string { 
	return "Plugin description" 
}

func (p *MyPlugin) RequireLimit() bool { 
	return true 
}

func (p *MyPlugin) Execute(ctx *plugins.Context) error {
	// Your plugin logic here
	msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "Hello!")
	ctx.API.Send(msg)
	return nil
}
```

## 🎯 Features

### User System
- Daily limit system (30 requests/day)
- Premium users (unlimited access)
- XP and leveling system
- User profiles and statistics

### Database
- BoltDB embedded database
- Automatic user registration
- Persistent data storage

### API Integration
- BetaBotz API for downloaders
- AI chat integration
- Image processing tools

## 📊 User Stats

Each user has:
- **Limit**: Daily request limit (resets at 00:00 WIB)
- **XP**: Experience points from games
- **Level**: User level based on XP
- **Premium**: Premium status for unlimited access

## 🚀 Usage

1. Start the bot: `./sofinco-bot`
2. Send `/start` to the bot
3. Browse categories using inline buttons
4. Use commands with `/command [args]`

## 🔑 Configuration

Required environment variables in `.env`:
```env
BOT_TOKEN=your_telegram_bot_token
OWNER_ID=your_telegram_user_id
OWNER_USERNAME=your_username
API_KEY=your_betabotz_api_key
AKSES_KEY=your_akses_key
```

## 📝 Notes

- All JavaScript plugins have been converted to Go
- Plugin system is modular and extensible
- Automatic plugin registration via `init()`
- Built-in limit and premium checks
- Error handling and user feedback

## 🎨 UI Features

- Inline keyboard navigation
- Category-based command organization
- User stats display
- Markdown formatting
- Emoji icons for better UX

---

Made with ❤️ by [levouinse](https://github.com/levouinse)
