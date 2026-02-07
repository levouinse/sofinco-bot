package game

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type TicTacToePlugin struct {
	plugins.BasePlugin
	games map[int64]*TicTacToeGame
}

type TicTacToeGame struct {
	Board     [9]string
	PlayerX   int64
	PlayerO   int64
	Turn      string
	CreatedAt time.Time
}

var TicTacToeInstance *TicTacToePlugin

func init() {
	TicTacToeInstance = &TicTacToePlugin{
		games: make(map[int64]*TicTacToeGame),
	}
	plugins.Register(TicTacToeInstance)
}

func (p *TicTacToePlugin) Commands() []string { return []string{"tictactoe", "ttt"} }
func (p *TicTacToePlugin) Tags() []string     { return []string{"game"} }
func (p *TicTacToePlugin) Help() string       { return "Play Tic Tac Toe game" }
func (p *TicTacToePlugin) RequireLimit() bool { return false }

func (p *TicTacToePlugin) Execute(ctx *plugins.Context) error {
	chatID := ctx.Message.Chat.ID
	userID := ctx.Message.From.ID

	// Check if user already in a game
	for _, game := range p.games {
		if game.PlayerX == userID || game.PlayerO == userID {
			ctx.API.Send(tgbotapi.NewMessage(chatID, "❌ Kamu masih dalam permainan!"))
			return nil
		}
	}

	// Check if there's a waiting game
	var waitingGame *TicTacToeGame
	for id, game := range p.games {
		if game.PlayerO == 0 {
			waitingGame = game
			chatID = id
			break
		}
	}

	if waitingGame != nil {
		// Join existing game
		waitingGame.PlayerO = userID
		waitingGame.Turn = "X"

		board := p.renderBoard(waitingGame)
		msg := fmt.Sprintf("🎮 *Tic Tac Toe*\n\n%s\n\nGiliran: <a href=\"tg://user?id=%d\">Player X</a>\nKetik angka 1-9 untuk bermain\nKetik 'nyerah' untuk menyerah",
			board, waitingGame.PlayerX)

		reply := tgbotapi.NewMessage(chatID, msg)
		reply.ParseMode = "HTML"
		ctx.API.Send(reply)
	} else {
		// Create new game
		game := &TicTacToeGame{
			Board:     [9]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
			PlayerX:   userID,
			PlayerO:   0,
			Turn:      "X",
			CreatedAt: time.Now(),
		}
		p.games[chatID] = game

		board := p.renderBoard(game)
		msg := fmt.Sprintf("🎮 *Tic Tac Toe*\n\n%s\n\n⏳ Menunggu pemain lain...\nKetik /tictactoe untuk bergabung",
			board)

		reply := tgbotapi.NewMessage(chatID, msg)
		reply.ParseMode = "Markdown"
		ctx.API.Send(reply)

		// Auto-delete game after 5 minutes if no one joins
		go func() {
			time.Sleep(5 * time.Minute)
			if g, exists := p.games[chatID]; exists && g.PlayerO == 0 {
				delete(p.games, chatID)
				ctx.API.Send(tgbotapi.NewMessage(chatID, "⏰ Game dibatalkan karena tidak ada pemain yang bergabung"))
			}
		}()
	}

	return nil
}

func (p *TicTacToePlugin) renderBoard(game *TicTacToeGame) string {
	board := game.Board
	return fmt.Sprintf(
		"%s %s %s\n%s %s %s\n%s %s %s",
		p.cellEmoji(board[0]), p.cellEmoji(board[1]), p.cellEmoji(board[2]),
		p.cellEmoji(board[3]), p.cellEmoji(board[4]), p.cellEmoji(board[5]),
		p.cellEmoji(board[6]), p.cellEmoji(board[7]), p.cellEmoji(board[8]),
	)
}

func (p *TicTacToePlugin) cellEmoji(v string) string {
	emoji := map[string]string{
		"X": "❌", "O": "⭕",
		"1": "1️⃣", "2": "2️⃣", "3": "3️⃣",
		"4": "4️⃣", "5": "5️⃣", "6": "6️⃣",
		"7": "7️⃣", "8": "8️⃣", "9": "9️⃣",
	}
	if e, ok := emoji[v]; ok {
		return e
	}
	return v
}

func (p *TicTacToePlugin) HandleMove(chatID int64, userID int64, move string) (string, bool) {
	game, exists := p.games[chatID]
	if !exists {
		return "", false
	}

	// Check if it's player's turn
	currentPlayer := game.PlayerX
	if game.Turn == "O" {
		currentPlayer = game.PlayerO
	}

	if userID != currentPlayer {
		return "❌ Bukan giliran kamu!", true
	}

	// Handle surrender
	if strings.ToLower(move) == "nyerah" {
		winner := game.PlayerO
		if userID == game.PlayerO {
			winner = game.PlayerX
		}
		delete(p.games, chatID)
		return fmt.Sprintf("🏳️ <a href=\"tg://user?id=%d\">Player</a> menyerah!\n🏆 <a href=\"tg://user?id=%d\">Player</a> menang!", userID, winner), true
	}

	// Validate move
	pos := -1
	fmt.Sscanf(move, "%d", &pos)
	if pos < 1 || pos > 9 {
		return "❌ Pilih angka 1-9!", true
	}

	pos-- // Convert to 0-indexed
	if game.Board[pos] == "X" || game.Board[pos] == "O" {
		return "❌ Posisi sudah terisi!", true
	}

	// Make move
	game.Board[pos] = game.Turn

	// Check winner
	if winner := p.checkWinner(game); winner != "" {
		board := p.renderBoard(game)
		delete(p.games, chatID)
		winnerID := game.PlayerX
		if winner == "O" {
			winnerID = game.PlayerO
		}
		return fmt.Sprintf("🎮 *Tic Tac Toe*\n\n%s\n\n🏆 <a href=\"tg://user?id=%d\">Player %s</a> menang!", board, winnerID, winner), true
	}

	// Check draw
	draw := true
	for _, cell := range game.Board {
		if cell != "X" && cell != "O" {
			draw = false
			break
		}
	}
	if draw {
		board := p.renderBoard(game)
		delete(p.games, chatID)
		return fmt.Sprintf("🎮 *Tic Tac Toe*\n\n%s\n\n🤝 Seri!", board), true
	}

	// Switch turn
	if game.Turn == "X" {
		game.Turn = "O"
	} else {
		game.Turn = "X"
	}

	board := p.renderBoard(game)
	nextPlayer := game.PlayerX
	if game.Turn == "O" {
		nextPlayer = game.PlayerO
	}

	return fmt.Sprintf("🎮 *Tic Tac Toe*\n\n%s\n\nGiliran: <a href=\"tg://user?id=%d\">Player %s</a>", board, nextPlayer, game.Turn), true
}

func (p *TicTacToePlugin) checkWinner(game *TicTacToeGame) string {
	b := game.Board
	lines := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // cols
		{0, 4, 8}, {2, 4, 6}, // diagonals
	}

	for _, line := range lines {
		if b[line[0]] == b[line[1]] && b[line[1]] == b[line[2]] {
			if b[line[0]] == "X" || b[line[0]] == "O" {
				return b[line[0]]
			}
		}
	}
	return ""
}
