package commandsbot

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func DeleteTask(conn *pgx.Conn, taskID int, chatID int64) (bool, error) {
	result, err := conn.Exec(
		context.Background(),
		"DELETE FROM tasks WHERE id = $1 AND user_id = $2",
		taskID,
		chatID,
	)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}
