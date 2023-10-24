package user_db

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/user"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) user.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertUser(ctx context.Context, user *user.User) error {
	query := `
	INSERT INTO
		public."User" (email, full_name, phone, image, role)
	VALUES
		($1, $2, $3, $4, $5)
	RETURNING
		id
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, &user.Email, &user.FullName, &user.Phone, &user.Image, &user.Role)
	err := rows.Scan(&user.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Добавлен", fmt.Sprintf("%s c id=:%d", "пользователь", &user.Id))

	return nil
}

// // GetUser implements Repository.
func (r *repository) SelectUser(ctx context.Context, id int) (us user.User, err error) {
	query := `
		SELECT
			id, email, full_name, phone, image, role
		FROM
			public."User"
		WHERE
			id = $1`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return user.User{}, err
	}

	err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role)
	if err != nil {
		return user.User{}, err
	}
	return us, nil
}

// GetUsers implements Repository.
func (r *repository) SelectUsers(ctx context.Context) (users []user.User, err error) {
	query := `
	SELECT
		id, email, full_name, phone, image, role
	FROM
		public."User"`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var us user.User
	for rows.Next() {
		err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, us)
	}
	return users, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *user.User) error {
	query := `
	UPDATE
		public."User"
	SET
		email = $1, full_name = $2, phone = $3, image = $4, role = $5
	WHERE
		id = $6
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, &user.Email, &user.FullName, &user.Phone, &user.Image, &user.Role, &user.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Обновлен", fmt.Sprintf("%s c id=:%d", "пользователь", &user.Id))

	return nil
}

func (r *repository) DeleteUser(ctx context.Context, id int) error {
	query := `
	DELETE INTO
		public."User"
	WHERE
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Удален", fmt.Sprintf("%s c id=:%d", "пользователь", id))

	return nil
}
