package services_db

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/services"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) services.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertServices(ctx context.Context, sr *services.Services) error {
	query := `
	INSERT INTO 
		public."Services" (equipment, type, cost)
	VALUES 
		($1, $2, $3)
	RETURNING 
		id
	`

	r.logger.Tracef("Query : %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, &sr.Equipment, &sr.Type, &sr.Cost)
	err := rows.Scan(&sr.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Добавлена", fmt.Sprintf("%s c id=:%d", "услуга", &sr.Id))

	return nil
}

func (r *repository) SelectServices(ctx context.Context) (srvs []services.Services, err error) {
	query := `
	SELECT 
		id, equipment, type, cost
	FROM 
		public."Services"
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var src services.Services

	for rows.Next() {
		err = rows.Scan(&src.Id, &src.Equipment, &src.Type, &src.Cost)
		if err != nil {
			return nil, err
		}
		srvs = append(srvs, src)
		fmt.Println(srvs)
	}
	return srvs, nil
}

func (r *repository) SelectService(ctx context.Context, id int) (srv services.Services, err error) {
	query := `
	SELECT 
		id, equipment, type, cost 
	FROM 
		public."Services" 
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return services.Services{}, err
	}

	err = rows.Scan(&srv.Id, &srv.Equipment, &srv.Type, &srv.Cost)
	if err != nil {
		return services.Services{}, err
	}
	return srv, nil
}

func (r *repository) UpdateServices(ctx context.Context, srv *services.Services) error {
	query := `
	UPDATE 
		public."Services" 
	SET 
		equipment = $1, type = $2, cost = $3 
	WHERE 
		id = $4 
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, &srv.Equipment, &srv.Type, &srv.Cost, &srv.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Изменена", fmt.Sprintf("%s c id=:%d", "услуга", &srv.Id))

	return nil
}

func (r *repository) DeleteServices(ctx context.Context, id int) error {
	query := `
	DELETE FROM 
		public."Sevices" 
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Удалена", fmt.Sprintf("%s c id=:%d", "услуга", id))

	return nil
}
