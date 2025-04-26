package service

import (
	"context"
	m "effective-mobile-task/internal/model"

	"github.com/jackc/pgx/v5"
)

func SavePerson(conn *pgx.Conn, u m.User) error {
	query := `
		INSERT INTO Users (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	var id int
	err := conn.QueryRow(context.Background(), query,
		u.Name,
		u.Surname,
		u.Patronymic,
		u.Age,
		u.Gender,
		u.Nationality,
	).Scan(&id)

	if err != nil {
		return err
	}
	return nil
}
