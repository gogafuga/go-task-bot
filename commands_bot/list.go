package commandsbot

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Task struct {
	ID   int
	Text string
	Done bool
}

func ListTasks(conn *pgx.Conn, chatID int64) ([]Task, error) {
	var tasks []Task

	rows, err := conn.Query(
		context.Background(),
		"SELECT id, text, done FROM tasks WHERE user_id = $1 ORDER BY id",
		chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Text, &task.Done)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
