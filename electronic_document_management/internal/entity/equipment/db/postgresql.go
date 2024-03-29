package equipment_db

import (
	"context"
	"fmt"
	"time"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/equipment"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) equipment.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertEquipment(ctx context.Context, eq *equipment.Equipment) error {
	query := `
	INSERT INTO 
		public."Equipment" (name, type, manufacturer, model, unique_number, contract, create_date)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING 
		id
	`
	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	fmt.Println(eq.Model)

	rows := r.client.QueryRow(ctx, query, &eq.Name, &eq.Type, &eq.Manufacture, &eq.Model, &eq.UniqueNumber, &eq.Contract, &eq.CreateDate)
	err := rows.Scan(&eq.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Добавлено", fmt.Sprintf("%s c id=%d / %s", "оборудование", eq.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) SelectEquipment(ctx context.Context, id int) (eq equipment.Equipment, err error) {
	query := `
		SELECT 
			id, name, type, manufacturer, model, unique_number, contract, create_date
		FROM 
			public."Equipment"
		WHERE 
			id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	if err != nil {
		return equipment.Equipment{}, err
	}

	err = rows.Scan(&eq.Id, &eq.Name, &eq.Type, &eq.Manufacture, &eq.Model, &eq.UniqueNumber, &eq.Contract, &eq.CreateDate)
	if err != nil {
		return equipment.Equipment{}, err
	}
	return eq, nil
}

func (r *repository) SelectEquipments(ctx context.Context) (eqs []equipment.Equipment, err error) {
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

func (r *repository) SelectEquipmentsBySort(ctx context.Context, eq *equipment.Equipment) (eqs []equipment.Equipment, err error) {
	query := `
		SELECT 
			id, name, type, manufacturer, model, unique_number, contract, create_date
		FROM 
			public."Equipment"
		WHERE name ILIKE $1 OR type = $2 OR manufacturer = $3 OR model = $4 OR unique_number ILIKE $5
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, &eq.Name, &eq.Type, &eq.Manufacture, &eq.Model, &eq.UniqueNumber)
	if err != nil {
		return nil, err
	}

	var e equipment.Equipment

	for rows.Next() {
		err = rows.Scan(&e.Id, &e.Name, &e.Type, &e.Manufacture, &e.Model, &e.UniqueNumber, &e.Contract, &e.CreateDate)
		if err != nil {
			return nil, err
		}

		eqs = append(eqs, e)
	}
	return eqs, nil
}

func (r *repository) UpdateEquipment(ctx context.Context, eq *equipment.Equipment) error {
	query := `
	UPDATE 
		public."Equipment" 
	SET
		name = $1, type = $2, manufacturer = $3, model = $4, unique_number = $5
	WHERE
		id = $6
		`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, &eq.Name, &eq.Type, &eq.Manufacture, &eq.Model, &eq.UniqueNumber, &eq.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Обновлено", fmt.Sprintf("%s c id=%d / %s", "оборудование", eq.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) DeleteEquipment(ctx context.Context, id int) error {
	query := `
	DELETE FROM 
		public."Equipment" 
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Удалено", fmt.Sprintf("%s c id=%d / %s", "оборудование", id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) SelectAllType(ctx context.Context) (types []string, err error) {
	query := `
	SELECT DISTINCT ON (type) type
	FROM public."Equipment"
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

func (r *repository) SelectAllManufacture(ctx context.Context) (manufacturers []string, err error) {

	query := `
	SELECT DISTINCT ON (manufacturer) manufacturer
	FROM public."Equipment"
	`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var m string

	for rows.Next() {
		err = rows.Scan(&m)
		if err != nil {
			return nil, err
		}

		manufacturers = append(manufacturers, m)
	}
	return manufacturers, nil
}

func (r *repository) SelectAllModel(ctx context.Context) (models []string, err error) {

	query := `
	SELECT DISTINCT ON (model) model
	FROM public."Equipment"
	`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var m string

	for rows.Next() {
		err = rows.Scan(&m)
		if err != nil {
			return nil, err
		}

		models = append(models, m)
	}
	return models, nil
}

func (r *repository) FindEquipment(ctx context.Context, find string) (eqs []equipment.Equipment, err error) {
	query := `
		SELECT 
			id, name, type, manufacturer, model, unique_number, contract, create_date
		FROM 
			public."Equipment"
		WHERE name ILIKE $1 OR type ILIKE $1 OR manufacturer ILIKE $1 OR model ILIKE $1 OR unique_number ILIKE $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	find = "%" + find + "%"

	rows, err := r.client.Query(ctx, query, find)
	if err != nil {
		return nil, err
	}

	var e equipment.Equipment

	for rows.Next() {
		err = rows.Scan(&e.Id, &e.Name, &e.Type, &e.Manufacture, &e.Model, &e.UniqueNumber, &e.Contract, &e.CreateDate)
		if err != nil {
			return nil, err
		}

		eqs = append(eqs, e)
	}
	return eqs, nil
}
