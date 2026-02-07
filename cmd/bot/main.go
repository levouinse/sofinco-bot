package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/levouinse/sofinco-bot/internal/bot"
	"github.com/levouinse/sofinco-bot/internal/config"
	"github.com/levouinse/sofinco-bot/internal/database"
)

func printBanner(botUsername string) {
	banner := `
╔═══════════════════════════════════════════════╗
║                                               ║
║   ███████╗ ██████╗ ███████╗██╗███╗   ██╗ ██████╗  ██████╗     ║
║   ██╔════╝██╔═══██╗██╔════╝██║████╗  ██║██╔════╝ ██╔═══██╗    ║
║   ███████╗██║   ██║█████╗  ██║██╔██╗ ██║██║      ██║   ██║    ║
║   ╚════██║██║   ██║██╔══╝  ██║██║╚██╗██║██║      ██║   ██║    ║
║   ███████║╚██████╔╝██║     ██║██║ ╚████║╚██████╗ ╚██████╔╝    ║
║   ╚══════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═══╝ ╚═════╝  ╚═════╝     ║
║                                               ║
║              Telegram Bot v1.0                ║
║         Modern Button-Based Interface         ║
║                                               ║
╚═══════════════════════════════════════════════╝
`
	fmt.Println(banner)
	fmt.Printf("🤖 Bot Username: @%s\n", botUsername)
	fmt.Printf("✅ Status: Running\n")
	fmt.Printf("⏰ Started: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	db, err := database.New("data/sofinco.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	printBanner(botAPI.Self.UserName)

	b := bot.New(botAPI, db, cfg)
	b.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down bot...")
}
