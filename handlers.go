package main

import (
	"fmt"
	commandsbot "go-task-bot/commands_bot"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

func HandleStart(token string, chatId int64) {
	err := SendMessage(
		token,
		chatId,
		"Привет! Я Go Task Manager 🤖",
	)
	if err != nil {
		fmt.Println("Ошибка отправки: ", err)
	}
}

func HandleAdd(conn *pgx.Conn, token string, chatID int64, messageText string) {
	taskText := strings.TrimPrefix(messageText, "/add")
	taskText = strings.TrimSpace(taskText)
	if taskText == "" {
		SendMessage(token, chatID, "Введите название задачи!")
	} else {
		err := commandsbot.AddTask(conn, chatID, taskText)
		if err != nil {
			log.Println("Ошибка сохранения задачи:", err)
			SendMessage(token, chatID, "❌ Не удалось сохранить задачу")
			return
		} else {
			SendMessage(token,
				chatID,
				"✅ Задача добавлена: "+taskText)
		}
	}
}

func HandleList(conn *pgx.Conn, token string, chatID int64) {
	tasks, err := commandsbot.ListTasks(conn, chatID)
	if err != nil {
		log.Println("Ошибка запроса", err)
		SendMessage(token, chatID, "❌ Не удалось получить задачи")
		return
	}
	text := "📋 Твои задачи:\n"
	for _, task := range tasks {
		status := "❌"
		if task.Done {
			status = "✅"
		}
		line := fmt.Sprintf("%d. %s %s\n", task.ID, status, task.Text)
		text += line
	}

	if len(tasks) == 0 {
		SendMessage(token, chatID, "📭 Список задач пуст")
	} else {
		SendMessage(token, chatID, text)
	}
}

func HandleDone(conn *pgx.Conn, token string, chatID int64, messageText string) {
	taskIDText := strings.TrimSpace(strings.TrimPrefix(messageText, "/done"))
	taskID, err := strconv.Atoi(taskIDText)
	if err != nil || taskID <= 0 {
		SendMessage(token, chatID, "❌ Укажи номер задачи!")
		return
	}

	done, err := commandsbot.DoneTask(conn, taskID, chatID)
	if err != nil {
		log.Println("Ошибка обновления задач: ", err)
		SendMessage(token, chatID, "❌ Не удалось обновить задачу")
		return
	}

	if !done {
		SendMessage(token, chatID, "❌ Задача не найдена")
		return
	}
	SendMessage(token, chatID, "✅ Задача выполнена!")
}

func HandleDelete(conn *pgx.Conn, token string, chatID int64, messageText string) {
	taskIDText := strings.TrimSpace(strings.TrimPrefix(messageText, "/delete"))
	taskID, err := strconv.Atoi(taskIDText)
	if err != nil || taskID <= 0 {
		SendMessage(token, chatID, "❌ Укажи номер задачи!")
		return
	}

	deleted, err := commandsbot.DeleteTask(conn, taskID, chatID)
	if err != nil {
		log.Println("Ошибка удаления задач: ", err)
		SendMessage(token, chatID, "❌ Не удалось удалить задачу")
		return
	}

	if !deleted {
		SendMessage(token, chatID, "❌ Задач нет")
		return
	}
	SendMessage(token, chatID, "✅ Задача удалена!")
}

func HandleUpdate(conn *pgx.Conn, token string, update Update) {

	if update.Message.Text == "/start" {
		HandleStart(token, update.Message.Chat.ID)
		return
	}

	if strings.HasPrefix(update.Message.Text, "/add") {
		HandleAdd(conn, token, update.Message.Chat.ID, update.Message.Text)
		return
	}

	if update.Message.Text == "/list" {
		HandleList(conn, token, update.Message.Chat.ID)
		return
	}

	if strings.HasPrefix(update.Message.Text, "/done") {
		HandleDone(conn, token, update.Message.Chat.ID, update.Message.Text)
		return
	}

	if strings.HasPrefix(update.Message.Text, "/delete") {
		HandleDelete(conn, token, update.Message.Chat.ID, update.Message.Text)
		return
	}
}
