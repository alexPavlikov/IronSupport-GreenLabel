package services_db

import (
	"context"
	"fmt"
	"time"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/equipment"
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

	r.logger.LogEvents("Добавлена", fmt.Sprintf("%s c id=%d / %s", "услуга", sr.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) SelectServices(ctx context.Context) (srvs []services.Services, err error) {
	query := `
	SELECT 
		s.id, s.equipment, s.type, s.cost, eq.Name 
	FROM 
		public."Services" s
	JOIN 
		"Equipment" eq ON s.equipment = eq.Id
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var src services.Services

	for rows.Next() {
		err = rows.Scan(&src.Id, &src.Equipment, &src.Type, &src.Cost, &src.EquipmentStructure.Name)
		if err != nil {
			return nil, err
		}

		srvs = append(srvs, src)
	}
	return srvs, nil
}

func (r *repository) SelectService(ctx context.Context, id int) (srv services.Services, err error) {
	query := `
	SELECT 
		s.id, s.equipment, s.type, s.cost, eq.Name 
	FROM 
		public."Services" s
	JOIN 
		"Equipment" eq ON s.equipment = eq.Id
	WHERE 
		s.id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)

	err = rows.Scan(&srv.Id, &srv.Equipment, &srv.Type, &srv.Cost, &srv.EquipmentStructure.Name)
	if err != nil {
		fmt.Println(err)
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

	r.logger.LogEvents("Изменена", fmt.Sprintf("%s c id=%d / %s", "услуга", srv.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

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

	r.logger.LogEvents("Удалена", fmt.Sprintf("%s c id=%d / %s", "услуга", id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) SelectServiceType(ctx context.Context) (types []string, err error) {
	query := `
	SELECT name FROM public."Services_type"
	`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var t string

	for rows.Next() {
		err = rows.Scan(&t)
		if err != nil {
			return nil, err
		}

		types = append(types, t)
	}
	return types, nil
}

func (r *repository) InsertServicesType(ctx context.Context, types string) error {
	query := `
	INSERT INTO 
		public."Services_type" (name) 
	VALUES 
		($1)
	`

	fmt.Println(types)

	_ = r.client.QueryRow(ctx, query, types)
	return nil
}

func (r *repository) SelectAllEquipment(ctx context.Context) (eqs []equipment.Equipment, err error) {
	query := `
		SELECT 
			id, name, type, manufacturer, model, unique_number, contract, create_date
		FROM 
			public."Equipment"
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var eq equipment.Equipment

	for rows.Next() {
		err = rows.Scan(&eq.Id, &eq.Name, &eq.Type, &eq.Manufacture, &eq.Model, &eq.UniqueNumber, &eq.Contract, &eq.CreateDate)
		if err != nil {
			return nil, err
		}

		eqs = append(eqs, eq)
	}
	return eqs, nil
}

func (r *repository) SelectServicesBySort(ctx context.Context, srv *services.Services) (srvc []services.Services, err error) {
	query := `
	SELECT 
		s.id, s.equipment, s.type, s.cost, eq.Name 
	FROM 
		public."Services" s
	JOIN 
		"Equipment" eq ON s.equipment = eq.Id
	WHERE s.equipment = $1 OR s.type = $2
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, &srv.Equipment, &srv.Type)
	if err != nil {
		return nil, err
	}

	var s services.Services

	for rows.Next() {
		err = rows.Scan(&s.Id, &s.Equipment, &s.Type, &s.Cost, &s.EquipmentStructure.Name)
		if err != nil {
			return nil, err
		}

		srvc = append(srvc, s)
	}
	return srvc, nil
}

func (r *repository) FindService(ctx context.Context, find string) (srvc []services.Services, err error) {
	query := `
	SELECT 
		s.id, s.equipment, s.type, s.cost, eq.name 
	FROM 
		public."Services" s
	JOIN 
		"Equipment" eq ON s.equipment = eq.Id
	WHERE eq.name ILIKE $1 OR s.type ILIKE $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	find = "%" + find + "%"

	rows, err := r.client.Query(ctx, query, find)
	if err != nil {
		return nil, err
	}

	var s services.Services

	for rows.Next() {
		err = rows.Scan(&s.Id, &s.Equipment, &s.Type, &s.Cost, &s.EquipmentStructure.Name)
		if err != nil {
			return nil, err
		}

		srvc = append(srvc, s)
	}
	return srvc, nil
}
