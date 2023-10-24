package client_db

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/client"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) client.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertClient(ctx context.Context, clnt *client.Client) error {
	query := `
	INSERT INTO public."Client" 
		(name, inn, kpp, ogrn, owner, phone, email, address, create_date, status)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, &clnt.Name, &clnt.INN, &clnt.KPP, &clnt.OGRN, &clnt.Owner, &clnt.Phone, &clnt.Email, &clnt.Address, &clnt.CreateDate, &clnt.Status)
	if err != nil {
		return err
	}

	err = rows.Scan(&clnt.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Добавлен", fmt.Sprintf("%s c id=:%d", "клиент", &clnt.Id))

	return nil
}

func (r *repository) SelectClient(ctx context.Context, id int) (cl client.Client, err error) {
	query := `
	SELECT 
		id, name, inn, kpp, ogrn, owner, phone, email, address, create_date, status
	FROM 
		public."Client"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return client.Client{}, err
	}

	err = rows.Scan(&cl.Id, &cl.Name, &cl.INN, &cl.KPP, &cl.OGRN, &cl.Owner, &cl.Phone, &cl.Email, &cl.Address, &cl.CreateDate, &cl.Status)
	if err != nil {
		return client.Client{}, err
	}
	return cl, nil
}

func (r *repository) SelectClients(ctx context.Context) (clnts []client.Client, err error) {
	query := `
	SELECT 
		id, name, inn, kpp, ogrn, owner, phone, email, address, create_date, status
	FROM 
		public."Client"
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var cl client.Client
	for rows.Next() {
		err = rows.Scan(&cl.Id, &cl.Name, &cl.INN, &cl.KPP, &cl.OGRN, &cl.Owner, &cl.Phone, &cl.Email, &cl.Address, &cl.CreateDate, &cl.Status)
		if err != nil {
			return nil, err
		}
		clnts = append(clnts, cl)
	}
	return clnts, nil
}

func (r *repository) UpdateClient(ctx context.Context, cl *client.Client) error {
	query := `
	UPDATE 
		public."Client"
	SET 
		name = $1, inn = $2, kpp = $3, ogrn = $4, owner = $5, phone = $6, email = $7, address = $8, create_date = $9
	WHERE 
		id = $11
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, &cl.Name, &cl.INN, &cl.KPP, &cl.OGRN, &cl.Owner, &cl.Phone, &cl.Email, &cl.Address, &cl.CreateDate, &cl.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Обновлен", fmt.Sprintf("%s c id=:%d", "клиент", &cl.Id))

	return nil
}

func (r *repository) DeleteClient(ctx context.Context, id int) error {
	query := `
	UPDATE 
		public."Client"
	SET 
		status = $1
	WHERE 
		id = $2
	`
	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, false, id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Удален", fmt.Sprintf("%s c id=:%d", "клиент", id))

	return nil
}
