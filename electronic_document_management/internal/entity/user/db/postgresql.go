package user_db

import (
	"context"
	"fmt"
	"time"

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
		public."User" (email, full_name, phone, image, role, password)
	VALUES
		($1, $2, $3, $4, $5, $6)
	RETURNING
		id
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, &user.Email, &user.FullName, &user.Phone, &user.Image, &user.Role, &user.Password)
	err := rows.Scan(&user.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Добавлен", fmt.Sprintf("%s c id=%d / %s", "пользователь", user.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) InsertUserRole(ctx context.Context, name string) error {
	query := `
	INSERT INTO public."Role" (name)
	VALUES ($1)
	`
	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, name)

	r.logger.LogEvents("Добавлена", fmt.Sprintf("%s - %s / %s", "роль", name, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

// // GetUser implements Repository.
func (r *repository) SelectUser(ctx context.Context, id int) (us user.User, err error) {
	query := `
		SELECT
			id, email, full_name, phone, image, role, password
		FROM
			public."User"
		WHERE
			id = $1`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)

	err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role, &us.Password)
	if err != nil {
		return user.User{}, err
	}
	return us, nil
}

func (r *repository) SelectAuthUser(ctx context.Context, email string, pass string) (us user.User, err error) {
	query := `
		SELECT
			id, email, full_name, phone, image, role, password
		FROM
			public."User"
		WHERE
			email = $1 AND password = $2`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, email, pass)

	err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role, &us.Password)
	if err != nil {
		return user.User{}, err
	}
	return us, nil
}

// GetUsers implements Repository.
func (r *repository) SelectUsers(ctx context.Context) (users []user.User, err error) {
	query := `
	SELECT
		id, email, full_name, phone, image, role, password
	FROM
		public."User"`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var us user.User
	for rows.Next() {
		err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role, &us.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, us)
	}
	return users, nil
}

func (r *repository) SelectUsersBySort(ctx context.Context, usr *user.User) (users []user.User, err error) {
	query := `
	SELECT
		id, email, full_name, phone, image, role, password
	FROM
		public."User"
	WHERE full_name ILIKE $1 OR email ILIKE $2 OR phone ILIKE $3 OR role = $4	
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, &usr.FullName, &usr.Email, &usr.Phone, &usr.Role)
	if err != nil {
		return nil, err
	}

	var us user.User
	for rows.Next() {
		err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role, &us.Password)
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
		email = $1, full_name = $2, phone = $3, role = $4
	WHERE
		id = $5
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, &user.Email, &user.FullName, &user.Phone, &user.Role, &user.Password, &user.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Обновлен", fmt.Sprintf("%s c id=%d / %s", "пользователь", user.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

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

	r.logger.LogEvents("Удален", fmt.Sprintf("%s c id=%d / %s", "пользователь", id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) SelectRole(ctx context.Context) (role []string, err error) {
	query := `
	SELECT name FROM public."Role"	`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var rl string

	for rows.Next() {
		err = rows.Scan(&rl)
		if err != nil {
			return nil, err
		}

		role = append(role, rl)
	}
	return role, nil
}

func (r *repository) FindUser(ctx context.Context, find string) (users []user.User, err error) {
	query := `
	SELECT
		id, email, full_name, phone, image, role, password
	FROM
		public."User"
	WHERE full_name ILIKE $1 OR email ILIKE $1 OR phone ILIKE $1 OR role = $1	
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	find = "%" + find + "%"

	rows, err := r.client.Query(ctx, query, find)
	if err != nil {
		return nil, err
	}

	var us user.User
	for rows.Next() {
		err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role, &us.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, us)
	}
	return users, nil
}
