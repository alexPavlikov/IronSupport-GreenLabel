package admin_db

import (
	"context"
	"fmt"

	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/admin"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/news"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"
	site "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website"
	"github.com/lib/pq"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) admin.Repository {
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
		fmt.Println(err)
		return nil, err
	}

	var n []news.News

	for rows.Next() {
		var news news.News
		err = rows.Scan(&news.Id, &news.Title, &news.Avatar, &news.CreateDate, &news.Text, &news.VideoLink, &news.Deleted, &news.Author)
		if err != nil {
			fmt.Println(err)
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
		($1, $2, $3, $4, $5, $6)
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
		title = $1, avatar = $2, text = $3, video_link = $4, deleted = $5
	WHERE 
		id = $6
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, nw.Title, nw.Avatar, nw.Text, nw.VideoLink, nw.Deleted, nw.Id)
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

func (r *repository) SelectProducts(ctx context.Context) (pr []product.Product, err error) {
	query := `
	SELECT id, article, name, full_name, waight, unit_of_meas, remains, price, category, discount, on_the_way
	FROM public."Product"
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p product.Product
		err = rows.Scan(&p.Id, &p.Article, &p.Name, &p.FullName, &p.Waight, &p.UnitOfMeasurement, &p.Remains, &p.Price, &p.Category.Name, &p.Discount.Percent, &p.OnTheWay)
		if err != nil {
			return nil, err
		}

		//add discount
		//add category

		pr = append(pr, p)
	}
	return pr, nil
}

func (r *repository) SelectTrustCompany(ctx context.Context) (tc []guest.TrustCompany, err error) {
	query := `
	SELECT 
		name, description, logo 
	FROM 
		public."TrustCompany"
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		var t guest.TrustCompany
		err = rows.Scan(&t.Name, &t.Description, &t.Logo)
		if err != nil {
			return nil, err
		}

		tc = append(tc, t)
	}

	return tc, nil
}

func (r *repository) SelectTrustCompanyByName(ctx context.Context, name string) (t guest.TrustCompany, err error) {
	query := `
	SELECT 
		name, description, logo 
	FROM 
		public."TrustCompany" 
	WHERE 
		name = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, name)

	err = rows.Scan(&t.Name, &t.Description, &t.Logo)
	if err != nil {
		return guest.TrustCompany{}, err
	}

	return t, nil
}

func (r *repository) UpdateTrustCompany(ctx context.Context, tc guest.TrustCompany) error {
	query := `
	UPDATE 
		public."TrustCompany" 
	SET 
		name = $1, description = $2, logo = $3
	WHERE 
		name = $4
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	r.client.QueryRow(ctx, query, tc.Name, tc.Description, tc.Logo, tc.Name)

	return nil
}

func (r *repository) InsertTrustCompany(ctx context.Context, tc guest.TrustCompany) error {
	query := `
	INSERT INTO public."TrustCompany" (name, description, logo) VALUES ($1, $2, $3)
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, tc.Name, tc.Description, tc.Logo)

	return nil
}

func (r *repository) SelectVacancy(ctx context.Context) (vc []site.Vacancy, err error) {
	query := `
	SELECT 
		name, options, active 
	FROM 
		public."Vacancy"
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var v site.Vacancy
		err = rows.Scan(&v.Name, pq.Array(&v.Options), &v.Active)
		if err != nil {
			return nil, err
		}

		vc = append(vc, v)
	}

	return vc, nil
}

func (r *repository) SelectVacancyByName(ctx context.Context, name string) (v site.Vacancy, err error) {
	query := `
	SELECT 
		name, options, active 
	FROM 
		public."Vacancy"
	WHERE 
		name = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, name)

	err = rows.Scan(&v.Name, pq.Array(&v.Options), &v.Active)
	if err != nil {
		return site.Vacancy{}, err
	}

	return v, nil
}

func (r *repository) InsertVacancy(ctx context.Context, v site.Vacancy) error {
	query := `
	INSERT INTO public."Vacancy" (name, options, active) VALUES($1, $2, $3)
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, v.Name, pq.Array(v.Options), v.Active)

	return nil
}

func (r *repository) UpdateVacancy(ctx context.Context, v site.Vacancy, name string) error {
	query := `
	UPDATE 
		public."Vacancy" 
	SET 
		name = $1, options = array_append('', $2), active = $3
	WHERE 
		name = $4
	`

	query2 := fmt.Sprintf(`
	UPDATE 
	public."Vacancy" 
	SET 
		name = '%s', options = '%v', active = %v
	WHERE 
		name = '%s
	`, v.Name, v.Options, v.Active, name)

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query2)

	return nil
}

func (r *repository) SelectAllSubscriber(ctx context.Context) (emails []string, err error) {
	query := `
	SELECT email FROM public."Subscribers"
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var e string
		err = rows.Scan(&e)
		if err != nil {
			return nil, err
		}

		emails = append(emails, e)
	}
	return emails, nil
}
