package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Chat struct {
	ID int64 `json:"ID"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type BotInfo struct {
	FirstName string `json:"first_name"`
	Username  string
}

type TelegramResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL не задан")
	}

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения!", err)
	}
	defer conn.Close(context.Background())

	offset := 0
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		fmt.Println("Токен не найден!")
		return
	}
	fmt.Println("Go Task Bot запущен!")

	for {
		updates, err := GetUpdate(token, offset)
		if err != nil {
			fmt.Println("Ошибка почения сообщений: ", err)
			time.Sleep(3 * time.Second)
			continue
		}
		for _, update := range updates {
			offset = update.UpdateID + 1
			fmt.Println("Сообщение: ", update.Message.Text)

			HandleUpdate(conn, token, update)
		}
	}
}
