package post_db

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/abs/internal/entity/post"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) post.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertAbs(ctx context.Context, abs *post.Post) error {
	query := `
	INSERT INTO public."Abs" 
		(user_id, deadline, text, color, xcord, ycord)
	VALUES 
		($1, $2, $3, $4, $5, $6)
	RETURNING 
		id
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	row := r.client.QueryRow(ctx, query, &abs.User, &abs.Deadline, &abs.Text, &abs.Color, &abs.PosX, &abs.PosY)
	err := row.Scan(&abs.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SelectAllAbs(ctx context.Context) (abs []post.Post, err error) {
	query := `
	SELECT 
		id, user_id, deadline, text, color, xcord, ycord
	FROM 
		public."Abs"
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var a post.Post

	for rows.Next() {
		err = rows.Scan(&a.Id, &a.User, &a.Deadline, &a.Text, &a.Color, &a.PosX, &a.PosY)
		if err != nil {
			return nil, err
		}

		abs = append(abs, a)
	}
	return abs, nil
}

func (r *repository) SelectAbs(ctx context.Context, id int) (abs post.Post, err error) {
	query := `
	SELECT 
		id, user_id, deadline, text, color, xcord, ycord
	FROM 
		public."Abs"
	WHERE 
		id = $1
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	row := r.client.QueryRow(ctx, query, id)
	err = row.Scan(&abs.Id, &abs.User, &abs.Deadline, &abs.Text, &abs.Color, &abs.PosX, &abs.PosY)
	if err != nil {
		return post.Post{}, nil
	}
	return abs, nil
}

func (r *repository) SelectUserAbs(ctx context.Context, user_id int) (abs []post.Post, err error) {
	query := `
	SELECT 
		id, user_id, deadline, text, color, xcord, ycord
	FROM 
		public."Abs"
	WHERE 
		user_id = $1
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	var a post.Post

	for rows.Next() {
		err = rows.Scan(&a.Id, &a.User, &a.Deadline, &a.Text, &a.Color, &a.PosX, &a.PosY)
		if err != nil {
			return nil, err
		}

		abs = append(abs, a)
	}
	return abs, nil
}

func (r *repository) UpdateUserAbs(ctx context.Context, abs post.Post) error {
	query := `
	UPDATE 
		public."Abs" 
	SET
		deadline = $1, text = $2, color = $3, xcord = $4, ycord = $5
	WHERE 
		id = $6
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, abs.Deadline, abs.Text, abs.Color, abs.PosX, abs.PosY, abs.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateCordAbs(ctx context.Context, x, y, id int) error {
	query := `
	UPDATE 
		public."Abs" 
	SET
		xcord = $1, ycord = $2
	WHERE 
		id = $3
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, x, y, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateTextAbs(ctx context.Context, text string, id int) error {
	query := `
	UPDATE 
		public."Abs" 
	SET
		text = $1
	WHERE 
		id = $2
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, text, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteUserAbs(ctx context.Context, abs int) error {
	query := `
	DELETE 
		FROM public."Abs"
	WHERE 
		id = $1
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Exec(ctx, query, abs)
	if err != nil {
		return err
	}
	return nil
}
