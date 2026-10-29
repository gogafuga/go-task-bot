package commandsbot

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func StatsTask(conn *pgx.Conn, chatID int64) (int, int, int, error) {
	var total int
	countTasks := conn.QueryRow(
		context.Background(),
		"SELECT COUNT(id) FROM tasks WHERE user_id = $1",
		chatID,
	)
	if err := countTasks.Scan(&total); err != nil {
		return 0, 0, 0, err
	}

	var done int
	doneTask := conn.QueryRow(
		context.Background(),
		"SELECT COUNT(id) FROM tasks WHERE user_id = $1 AND done = true",
		chatID,
	)
	if err := doneTask.Scan(&done); err != nil {
		return 0, 0, 0, err
	}

	var unDone int
	unDone = total - done

	return total, done, unDone, nil
}
