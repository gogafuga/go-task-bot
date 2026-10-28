package commandsbot

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func AddTask(conn *pgx.Conn, ID int64, text string) error {
	_, err := conn.Exec(
		context.Background(),
		"INSERT INTO tasks (user_id, text) VALUES ($1, $2)",
		ID,
		text)
	return err
}
