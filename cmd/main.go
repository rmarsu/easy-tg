package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rmarsu/easy-tg/src/bot"
	"github.com/rmarsu/easy-tg/src/types"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	bot, err := bot.New(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	initBotRoutes(bot)
	
	bot.Start()
}

func initBotRoutes(b *bot.Bot) {
	b.Add("/start", startHandler)
	b.Add("/hello", helloHandler)
	b.Add("/kawaii", kawaiiHandler)
	b.Add(types.PhotoType, wtfIsThis)
	b.Add(types.DocumentType, documentHandler)
} 

func helloHandler(bot *bot.Bot, upd *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.FromChat().ID, "Привет!! Как тебя зовут?")
	if err := bot.Send(msg); err != nil {
		fmt.Println(err)
	}
	ans := bot.WaitForMessage(upd)
	if ans != nil {
		msg := tgbotapi.NewMessage(ans.Message.From.ID, fmt.Sprintf("Отлично, я вас слышал! Вас зовут: %s", ans.Message.Text))
		if err := bot.Send(msg); err != nil {
			fmt.Println(err)
		}
	}
}

func startHandler(bot *bot.Bot, upd *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.FromChat().ID, "Привествую в бота")
	if err := bot.Send(msg); err != nil {
		fmt.Println(err)
	}
}

func kawaiiHandler(bot *bot.Bot, upd *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.FromChat().ID, "Докажи что ты няшка.")
	if err := bot.Send(msg); err != nil {
		fmt.Println(err)
	}
	ans := bot.WaitForMessage(upd)
	if ans != nil {
		if ans.Message.Text == "я няшка" {
			msg := tgbotapi.NewMessage(ans.Message.From.ID, "Ты действительно няшка? А если я проверю?")
			if err := bot.Send(msg); err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Second)
			if strings.Contains(ans.Message.From.FirstName, "няшка") {
				msg := tgbotapi.NewMessage(ans.Message.From.ID, "Да, ты няшка!")
				if err := bot.Send(msg); err != nil {
					fmt.Println(err)
				}
			} else {
				msg := tgbotapi.NewMessage(ans.Message.From.ID, "Извините, но я не могу быть таким уверенным. Вы не няшка.")
				if err := bot.Send(msg); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func wtfIsThis(bot *bot.Bot, upd *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Что это? Это фотография? втф.")
	if err := bot.Send(msg); err != nil {
		fmt.Println(err)
	}
}

func documentHandler(bot *bot.Bot, upd *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "люблю разные файлики. а ты?")
	if err := bot.Send(msg); err != nil {
		fmt.Println(err)
	}
	ans := bot.WaitForMessage(upd)
	if ans != nil {
		if strings.Contains(strings.ToLower(ans.Message.Text), "да") {
			msg := tgbotapi.NewMessage(ans.Message.From.ID, "хихихихи, братан!")
			if err := bot.Send(msg); err != nil {
				fmt.Println(err)
			}
		} else {
			msg := tgbotapi.NewMessage(ans.Message.From.ID, "плохо.")
			if err := bot.Send(msg); err != nil {
				fmt.Println(err)
			}
		}
	}
}
