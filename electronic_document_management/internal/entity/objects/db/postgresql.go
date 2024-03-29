package objects_db

import (
	"context"
	"fmt"
	"time"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/objects"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) objects.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertObject(ctx context.Context, obj *objects.Object) error {
	query := `
	INSERT INTO 
		public."Objects" (name, address, work_schedule)
	VALUES 
		($1, $2, $3)
	RETURNING 
		id
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, &obj.Name, &obj.Address, &obj.WorkSchedule)

	err := rows.Scan(&obj.Id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("TEST", obj.Id, obj.Client.Id)

	err = r.InsertClientObject(ctx, obj.Client.Id, obj.Id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	r.logger.LogEvents("Добавлен", fmt.Sprintf("%s c id=%d / %s", "объект", obj.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) SelectObject(ctx context.Context, id int) (obj objects.Object, err error) {
	query := `
	SELECT 
		"Objects".id, "Objects".name, "Objects".address, "Objects".work_schedule, "Client_objects".id, "Client".Name, "Client".Id  
	FROM 
		public."Objects" 
	JOIN "Client_objects" ON "Objects".id = "Client_objects".object 
	JOIN "Client" ON "Client".id = "Client_objects".client  
	WHERE 
		"Objects".id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)

	err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule, &obj.ClientObjectId, &obj.Client.Name, &obj.Client.Id)
	if err != nil {
		return objects.Object{}, err
	}

	obj.Client, err = r.SelectClientById(ctx, obj.Client.Id)
	if err != nil {
		return objects.Object{}, err
	}
	return obj, nil
}

func (r *repository) SelectObjects(ctx context.Context) (objs []objects.Object, err error) {
	query := `
	SELECT 
		"Objects".id, "Objects".name, "Objects".address, "Objects".work_schedule, "Client_objects".id, "Client".Name, "Client".Id  
	FROM 
		public."Objects" 
	JOIN "Client_objects" ON "Objects".id = "Client_objects".object 
	JOIN "Client" ON "Client".id = "Client_objects".client 
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var obj objects.Object

	for rows.Next() {

		err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule, &obj.ClientObjectId, &obj.Client.Name, &obj.Client.Id)
		if err != nil {
			return nil, err
		}

		obj.Client, err = r.SelectClientById(ctx, obj.Client.Id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}

	return objs, nil
}

func (r *repository) SelectObjectBySorted(ctx context.Context, ob *objects.Object) (obs []objects.Object, err error) {
	query := `
	SELECT 
		"Objects".id, "Objects".name, "Objects".address, "Objects".work_schedule, "Client_objects".id, "Client".Name, "Client".Id  
	FROM 
		public."Objects" 
	JOIN "Client_objects" ON "Objects".id = "Client_objects".object 
	JOIN "Client" ON "Client".id = "Client_objects".client
	WHERE "Objects".name ILIKE $1 OR "Client".Id = $2 
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	if ob.Name != "" {
		ob.Name = "%" + ob.Name + "%"
	}

	fmt.Println(ob.Name, ob.Client.Id)

	rows, err := r.client.Query(ctx, query, ob.Name, ob.Client.Id)
	if err != nil {
		return nil, err
	}

	var obj objects.Object

	for rows.Next() {

		err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule, &obj.ClientObjectId, &obj.Client.Name, &obj.Client.Id)
		if err != nil {
			return nil, err
		}

		obj.Client, err = r.SelectClientById(ctx, obj.Client.Id)
		if err != nil {
			return nil, err
		}

		fmt.Println(obj)

		obs = append(obs, obj)
	}

	fmt.Println(obs)

	return obs, nil
}

func (r *repository) UpdateObject(ctx context.Context, obj *objects.Object) error {
	query := `
	UPDATE INTO 
		public."Objects" 
	SET 
		name = $1, address = $2, work_schedule = $3
	WHERE 
		id = $4
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, &obj.Name, &obj.Address, &obj.WorkSchedule, &obj.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Обновлен", fmt.Sprintf("%s c id=%d / %s", "объект", obj.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) UpdateClientObject(ctx context.Context, obj *objects.Object) error {
	query := `
	UPDATE INTO 
		public."Client_object" 
	SET 
		client = $1
	WHERE 
		id = $2
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, &obj.Client.Id, &obj.ClientObjectId)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Обновлен", fmt.Sprintf("%s c id=%d / %s", "объект", obj.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) DeleteObject(ctx context.Context, id int) error {
	query := `
	DELETE 
	FROM 
		public."Objects" 
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Удален", fmt.Sprintf("%s c id=%d / %s", "объект", id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) SelectClient(ctx context.Context) (clnts []objects.Client, err error) {
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
	var cl objects.Client
	for rows.Next() {
		err = rows.Scan(&cl.Id, &cl.Name, &cl.INN, &cl.KPP, &cl.OGRN, &cl.Owner, &cl.Phone, &cl.Email, &cl.Address, &cl.CreateDate, &cl.Status)
		if err != nil {
			return nil, err
		}
		clnts = append(clnts, cl)
	}
	return clnts, nil
}

func (r *repository) SelectClientById(ctx context.Context, id int) (clnts objects.Client, err error) {
	query := `
	SELECT 
		id, name, inn, kpp, ogrn, owner, phone, email, address, create_date, status
	FROM 
		public."Client"
	WHERE id = $1
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)

	err = rows.Scan(&clnts.Id, &clnts.Name, &clnts.INN, &clnts.KPP, &clnts.OGRN, &clnts.Owner, &clnts.Phone, &clnts.Email, &clnts.Address, &clnts.CreateDate, &clnts.Status)
	if err != nil {
		return clnts, err
	}
	return clnts, nil
}

func (r *repository) InsertClientObject(ctx context.Context, client int, obj int) error {
	query := `
	INSERT INTO public."Client_objects" (client, object)
	VALUES 
		($1, $2)
	RETURNING id
	`

	var id int

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, client, obj)

	err := rows.Scan(&id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Добавлен", fmt.Sprintf("%s c id=%d / %s", "объект", id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) FindObject(ctx context.Context, find string) (obs []objects.Object, err error) {
	query := `
	SELECT 
		"Objects".id, "Objects".name, "Objects".address, "Objects".work_schedule, "Client_objects".id, "Client".Name, "Client".Id  
	FROM 
		public."Objects" 
	JOIN "Client_objects" ON "Objects".id = "Client_objects".object 
	JOIN "Client" ON "Client".id = "Client_objects".client
	WHERE "Objects".name ILIKE $1 OR "Client".Name ILIKE $1 OR "Objects".address ILIKE $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	find = "%" + find + "%"

	rows, err := r.client.Query(ctx, query, find)
	if err != nil {
		return nil, err
	}

	var obj objects.Object

	for rows.Next() {

		err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule, &obj.ClientObjectId, &obj.Client.Name, &obj.Client.Id)
		if err != nil {
			return nil, err
		}

		obj.Client, err = r.SelectClientById(ctx, obj.Client.Id)
		if err != nil {
			return nil, err
		}

		obs = append(obs, obj)
	}

	fmt.Println(obs)

	return obs, nil
}
