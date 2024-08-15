package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/alecthomas/kong"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"gopkg.in/yaml.v3"
)

const (
	// Version vars. Will be set during build
	Version = "1.0.0"
	Timestamp = "2021-01-01T00:00:00Z"
	GitCommit = "0000000"
	Repo = "azhinu/telegram-question-cards-bot"
)
var (
	// Sessions map to store user sessions
	Sessions map[int64]Session
	Decks map[string][]string
	Lock sync.RWMutex
)
type Session struct {
	Deck string
	PlayingQuestinons []int
	DestroyAfter time.Time
}

// CLI
var cli struct {
	// flags
	Version 	  bool `name:"version" help:"Print version and quit"`
	Debug       bool `short:"d" help:"Enable debug log" env:"QC_BOT_DEBUG"`
	Token       string `short:"t" help:"Telegram bot token" env:"QC_BOT_TOKEN" placeholder:"201204456:AAFFJJ"`
	URL		 			string `short:"u" help:"Webhook URL." env:"QC_BOT_URL" placeholder:"https://example.com/bot-secret-url"`
	Port				int `short:"p" help:"Webhook port" env:"QC_BOT_PORT" default:"1443"`

	// args
	Decks      string `arg:"" optional:"" type:"existingfile" help:"File with decks to load"`
}

func loadDecks(filename string) (map[string][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var decks map[string][]string
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&decks); err != nil {
		return nil, err
	}

	return decks, nil
}

func cleanup() {
	for {
		for chatID, session := range Sessions {
			if time.Now().After(session.DestroyAfter) {
				Lock.Lock()
				delete(Sessions, chatID)
				Lock.Unlock()
			}
		}
		time.Sleep(10 * time.Minute)
	}
}

func main() {
		// parse cli
		ctx := kong.Parse(&cli,
		kong.Name("tg_question_cards_bot"),
		kong.Description("Run telegram bot to play question cards game"),
		kong.UsageOnError(),
		)
	if cli.Version {
		fmt.Println("Version:", Version, "GitCommit:", GitCommit, "Timestamp:", Timestamp)
		os.Exit(0)
	}

	if cli.Decks == "" {
		err := ctx.PrintUsage(true)
		if err != nil {
			fmt.Println("Failed printing usage:", err)
		}
		os.Exit(0)
	}

	err := error(nil)
	Decks, err = loadDecks(cli.Decks)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Init global vars
	Sessions = make(map[int64]Session)
	Lock = sync.RWMutex{}

	bot, err := telego.NewBot(cli.Token, telego.WithDefaultLogger(cli.Debug, true))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start bot with polling or with webhook
	var updates <-chan telego.Update
	if cli.URL != "" {
		// Set up a webhook on Telegram side
		_ = bot.SetWebhook(&telego.SetWebhookParams{
			URL: cli.URL,
		})

		// Receive information about webhook
		info, _ := bot.GetWebhookInfo()
		fmt.Printf("Webhook Info: %+v\n", info)

		// Get an update channel from webhook.
		updates, _ = bot.UpdatesViaWebhook("/bot" + bot.Token())

		// Start server for receiving requests from the Telegram
		go func() {
			_ = bot.StartWebhook(fmt.Sprint("localhost:", cli.Port))
		}()

		// Stop reviving updates from update channel and shutdown webhook server
		defer func() {
			_ = bot.StopWebhook()
		}()
	} else {
	updates, _ = bot.UpdatesViaLongPolling(nil)
	}
	botHandler, _ := th.NewBotHandler(bot, updates)

	botHandler.HandleMessage(Start, th.CommandEqual("start"))

	botHandler.HandleCallbackQuery(NextQuestion, th.CallbackDataEqual("next"))

	botHandler.HandleCallbackQuery(SelectDeck, th.AnyCallbackQuery())

	// Stop handling updates on exit
	defer botHandler.Stop()
	defer bot.StopLongPolling()

	// Start handling updates
	botHandler.Start()
	go cleanup()
}
