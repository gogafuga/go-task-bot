package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	urlpkg "net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Task struct {
	ID   int
	Text string
	Done bool
}

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
		updates, err := getUpdate(token, offset)
		if err != nil {
			fmt.Println("Ошибка получения сообщений: ", err)
			time.Sleep(3 * time.Second)
			continue
		}
		for _, update := range updates {
			offset = update.UpdateID + 1
			fmt.Println("Сообщение: ", update.Message.Text)

			if update.Message.Text == "/start" {
				err := sendMessage(
					token,
					update.Message.Chat.ID,
					"Привет! Я Go Task Manager 🤖",
				)
				if err != nil {
					fmt.Println("Ошибка отправки: ", err)
				}
			}

			if strings.HasPrefix(update.Message.Text, "/add") {
				taskText := strings.TrimPrefix(update.Message.Text, "/add")
				taskText = strings.TrimSpace(taskText)
				if taskText == "" {
					sendMessage(token, update.Message.Chat.ID, "Введите название задачи!")
				} else {
					_, err := conn.Exec(
						context.Background(),
						"INSERT INTO tasks (user_id, text) VALUES ($1, $2)",
						update.Message.Chat.ID,
						taskText,
					)
					if err != nil {
						log.Println("Ошибка сохранения задачи:", err)
						sendMessage(token, update.Message.Chat.ID, "❌ Не удалось сохранить задачу")
						continue
					} else {
						sendMessage(token,
							update.Message.Chat.ID,
							"✅ Задача добавлена: "+taskText)
					}
				}
			}

			if update.Message.Text == "/list" {
				rows, err := conn.Query(
					context.Background(),
					"SELECT id, text, done FROM tasks WHERE user_id = $1 ORDER BY id",
					update.Message.Chat.ID,
				)
				if err != nil {
					log.Println("Ошибка запроса", err)
					sendMessage(token, update.Message.Chat.ID, "❌ Не удалось получить задачи")
					continue
				}
				text := "📋 Твои задачи:\n"
				count := 0
				for rows.Next() {
					var id int
					var taskText string
					var done bool
					err := rows.Scan(&id, &taskText, &done)
					if err != nil {
						log.Println("Ошибка чтения", err)
						continue
					}
					status := "❌"
					if done {
						status = "✅"
					}
					line := fmt.Sprintf("%d. %s %s\n", id, status, taskText)
					text += line
					count++
				}
				readErr := rows.Err()
				rows.Close()

				if readErr != nil {
					log.Println("Ошибка чтения строк:", readErr)
					sendMessage(token, update.Message.Chat.ID, "❌ Не удалось получить список задач")
					continue
				}

				rows.Close()
				if count == 0 {
					sendMessage(token, update.Message.Chat.ID, "📭 Список задач пуст")
				} else {
					sendMessage(token, update.Message.Chat.ID, text)
				}

			}

			if strings.HasPrefix(update.Message.Text, "/done") {
				taskIDText := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/done"))
				taskID, err := strconv.Atoi(taskIDText)
				if err != nil || taskID <= 0 {
					sendMessage(token, update.Message.Chat.ID, "❌ Укажи номер задачи!")
					continue
				}

				result, err := conn.Exec(
					context.Background(),
					"UPDATE tasks SET done = TRUE WHERE id = $1 AND user_id = $2",
					taskID,
					update.Message.Chat.ID,
				)
				if err != nil {
					log.Println("Ошибка обновления задач: ", err)
					sendMessage(token, update.Message.Chat.ID, "❌ Не удалось обновить задачу")
					continue
				}

				if result.RowsAffected() == 0 {
					sendMessage(token, update.Message.Chat.ID, "❌ Задача не найдена")
					continue
				}
				sendMessage(token, update.Message.Chat.ID, "✅ Задача выполнена!")
			}

			if strings.HasPrefix(update.Message.Text, "/delete") {
				taskIDText := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/delete"))
				taskID, err := strconv.Atoi(taskIDText)
				if err != nil || taskID <= 0 {
					sendMessage(token, update.Message.Chat.ID, "❌ Укажи номер задачи!")
					continue
				}

				result, err := conn.Exec(
					context.Background(),
					"DELETE FROM tasks WHERE id = $1 AND user_id = $2",
					taskID,
					update.Message.Chat.ID,
				)
				if err != nil {
					log.Println("Ошибка удаления задач: ", err)
					sendMessage(token, update.Message.Chat.ID, "❌ Не удалось удалить задачу")
					continue
				}

				if result.RowsAffected() == 0 {
					sendMessage(token, update.Message.Chat.ID, "❌ Задач нет")
					continue
				}
				sendMessage(token, update.Message.Chat.ID, "✅ Задача удалена!")
			}
		}
	}
}

func getUpdate(token string, offset int) ([]Update, error) {
	baseURL := "https://api.telegram.org/bot" + token + "/getUpdates"
	params := urlpkg.Values{}
	params.Set("offset", strconv.Itoa(offset))
	params.Set("timeout", "30")
	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response TelegramResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Result, nil
}
func sendMessage(token string, chatID int64, message string) error {
	baseURL := "https://api.telegram.org/bot" + token + "/sendMessage"
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(chatID, 10))
	params.Set("text", message)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
