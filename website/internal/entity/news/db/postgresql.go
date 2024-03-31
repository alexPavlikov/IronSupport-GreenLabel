package news_db

import (
	"context"

	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/news"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) news.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) SelectNews(ctx context.Context) ([]news.News, error) {
	query := `
	SELECT 
		id, title, avatar, create_date, text, video_link, deleted, author
	FROM 
		public."News"
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var n []news.News

	for rows.Next() {
		var news news.News
		err = rows.Scan(&news.Id, &news.Title, &news.Avatar, &news.CreateDate, &news.Text, &news.VideoLink, &news.Deleted, &news.Author)
		if err != nil {
			return nil, err
		}
		n = append(n, news)
	}
	return n, nil
}

func (r *repository) SelectUnDeletedNews(ctx context.Context) ([]news.News, error) {
	query := `
	SELECT 
		id, title, avatar, create_date, text, video_link, deleted, author
	FROM 
		public."News"
	WHERE deleted = false
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var n []news.News

	for rows.Next() {
		var news news.News
		err = rows.Scan(&news.Id, &news.Title, &news.Avatar, &news.CreateDate, &news.Text, &news.VideoLink, &news.Deleted, &news.Author)
		if err != nil {
			return nil, err
		}
		n = append(n, news)
	}
	return n, nil
}

func (r *repository) SelectPostOfNews(ctx context.Context, id int) (news.News, error) {
	query := `
	SELECT 
		id, title, avatar, create_date, text, video_link, deleted, author
	FROM 
		public."News"
	WHERE id = $1
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)

	var news news.News
	err := rows.Scan(&news.Id, &news.Title, &news.Avatar, &news.CreateDate, &news.Text, &news.VideoLink, &news.Deleted, &news.Author)
	if err != nil {
		return news, err
	}

	return news, nil
}

func (r *repository) InsertNews(ctx context.Context, nw *news.News) error {
	query := `
	INSERT INTO 
		public."News" (title, avatar, create_date, text, video_link, author)
	VALUES 
		($1, $2, $3, $4, $5)
	RETURNING id
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, &nw.Title, &nw.Avatar, &nw.CreateDate, &nw.Text, &nw.VideoLink, &nw.Author)
	err := rows.Scan(&nw.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateNews(ctx context.Context, nw news.News) {
	query := `
	UPDATE 
		public."News" 
	SET 
		title = $1, avatar = $2, create_date = $3, text = $4, video_link = $5
	WHERE 
		id = $6
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, nw.Title, nw.Avatar, nw.CreateDate, nw.Text, nw.VideoLink)
}

func (r *repository) CloseNews(ctx context.Context, id int) {
	query := `
	UPDATE
		public."News"
	SET
		deleted = true
	WHERE
		id = $1
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, id)
}
