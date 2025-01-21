package main

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rmarsu/easy-tg/src/bot"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	bot, err := bot.New(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Add("/start", startHandler)
	bot.Add("/hello", helloHandler)

	bot.Start()
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
