package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Start(bot *telego.Bot, msg telego.Message) {
	// Init user
	Lock.Lock()
	Sessions[msg.Chat.ID] = Session{
		Deck: "",
		DestroyAfter: time.Now().Add(time.Hour),
	}
	Lock.Unlock()


	// Generate keyboard with decks
	var rows [][]telego.InlineKeyboardButton
	for deckName := range Decks {
		button := tu.InlineKeyboardRow(tu.InlineKeyboardButton(deckName).WithCallbackData(deckName))
		rows = append(rows, button)
	}
	keyboard := tu.InlineKeyboard(rows...)

	_, _ = bot.SendMessage(tu.Message(
		tu.ID(msg.Chat.ID),
		fmt.Sprintf("Hello %s!\n Choose the deck.\n\n Note: Your session will be reset after 1 hour of inactivity", msg.From.FirstName),
	).WithReplyMarkup(keyboard),
	)
}

func SelectDeck(bot *telego.Bot, query telego.CallbackQuery) {
	if _, exists := Decks[query.Data]; exists {
		// Shuffle deck
		questionID := make([]int, len(Decks[query.Data]))
		for i := 0; i < len(questionID); i++ {
			questionID[i] = i
		}
		rand.Shuffle(len(questionID), func(i, j int) {
			questionID[i], questionID[j] = questionID[j], questionID[i]
		})

		// Update session
		Lock.Lock()
		Sessions[query.Message.GetChat().ID] = Session{
			Deck: query.Data,
			PlayingQuestinons: questionID,
			DestroyAfter: time.Now().Add(time.Hour),
		}
		Lock.Unlock()

		// Send message
		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("Show question").WithCallbackData("next")),
			)
		_, _ = bot.SendMessage(
			tu.Message(tu.ID(query.Message.GetChat().ID), fmt.Sprint("Playing " + query.Data + " deck.")).WithReplyMarkup(keyboard))

		// Answer callback query
		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Done"))
		return
	}

	_, _ = bot.SendMessage(tu.Message(tu.ID(query.Message.GetChat().ID), "Ошибочка вышла\n Use /start to reset"))
	_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Done"))
}

func NextQuestion(bot *telego.Bot, query telego.CallbackQuery) {
	session := Sessions[query.Message.GetChat().ID]

	// Check if deck is empty
	if len(session.PlayingQuestinons) == 0 {
		_, _ = bot.SendMessage(tu.Message(tu.ID(query.Message.GetChat().ID), "Deck is empty\n Use /start to reset"))
		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Done"))
		return
	}

	questionID := session.PlayingQuestinons[0]
	question := Decks[session.Deck][questionID]

	// Remove asked question from the deck
	Lock.Lock()
	Sessions[query.Message.GetChat().ID] = Session{
		Deck: session.Deck,
		PlayingQuestinons: session.PlayingQuestinons[1:],
		DestroyAfter: time.Now().Add(time.Hour),
	}
	Lock.Unlock()

	// Send message
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Next question").WithCallbackData("next"),
			),
		)
	_, _ = bot.SendMessage(tu.Message(tu.ID(query.Message.GetChat().ID), question).WithReplyMarkup(keyboard))
	_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Done"))
}
