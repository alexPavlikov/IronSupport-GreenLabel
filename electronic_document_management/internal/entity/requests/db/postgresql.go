package requests_db

import (
	"context"
	"fmt"
	"time"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/client"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/contract"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/equipment"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/objects"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/requests"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/user"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
	"github.com/lib/pq"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

// InsertRequest implements Repository.
func (r *repository) InsertRequest(ctx context.Context, req *requests.Request) error {
	query := `
	INSERT INTO public."Request"
	(title, client, worker, client_object, equipment, contract, description, priority, start_date, end_date, status, name)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	RETURNING id
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))
	fmt.Println(req.Id, req.Title, req.Client.Id, req.Worker.Id, req.ClientObject.Object.Id, req.Equipment.Id, req.Contract.Id, req.Description, req.Priority, req.StartDate, req.EndDate, req.Status.Name, req.Name)

	req.StartDate = time.Now().Format("2006-01-02")

	rows := r.client.QueryRow(ctx, query, req.Title, req.Client.Id, req.Worker.Id, req.ClientObject.Object.Id, req.Equipment.Id, req.Contract.Id, req.Description, req.Priority, req.StartDate, req.EndDate, req.Status.Name, req.Name)
	err := rows.Scan(&req.Id)
	if err != nil {
		fmt.Println("error", err)
		return err
	}

	r.logger.LogEvents("Добавлена", fmt.Sprintf("%s c id=%d на сотрудника - %d / %s", "заявка", req.Id, req.Worker.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

// SelectRequest implements Repository.
func (r *repository) SelectRequest(ctx context.Context, id int) (req requests.Request, err error) {
	query := `
	SELECT 
		id, title, client, worker, client_object, equipment, contract, description, priority, start_date, end_date, files, status, name
	FROM 
		public."Request"
	WHERE 
		id = $1`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return requests.Request{}, err
	}

	var client, worker, client_object, equipment, contract int
	var status string

	for rows.Next() {
		err = rows.Scan(&req.Id, &req.Title, &client, &worker, &client_object, &equipment, &contract, &req.Description, &req.Priority, &req.StartDate, &req.EndDate, pq.Array(&req.Files), &status, &req.Name)
		if err != nil {
			return requests.Request{}, err
		}
		req.Client, err = r.getRequestClient(ctx, client)
		if err != nil {
			return requests.Request{}, err
		}
		req.Worker, err = r.GetRequestWorker(ctx, worker)
		if err != nil {
			return requests.Request{}, err
		}
		req.Contract, err = r.getRequestContract(ctx, contract)
		if err != nil {
			return requests.Request{}, err
		}
		req.Equipment, err = r.getRequestEquipment(ctx, equipment)
		if err != nil {
			return requests.Request{}, err
		}
		req.ClientObject, err = r.getRequestClientObject(ctx, client_object)
		if err != nil {
			return requests.Request{}, err
		}
		req.Status, err = r.getRequestStatus(ctx, status)
		if err != nil {
			return requests.Request{}, err
		}
	}
	return req, nil
}

// SelectRequests implements Repository.
func (r *repository) SelectRequests(ctx context.Context) (reqs []requests.Request, err error) {
	query := `
	SELECT 
		id, title, client, worker, client_object, equipment, contract, description, priority, start_date, end_date, files, status, name
	FROM 
		public."Request"`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var client, worker, client_object, equipment, contract int
	var status string

	for rows.Next() {
		var rt requests.Request
		err = rows.Scan(&rt.Id, &rt.Title, &client, &worker, &client_object, &equipment, &contract, &rt.Description, &rt.Priority, &rt.StartDate, &rt.EndDate, pq.Array(&rt.Files), &status, &rt.Name)
		if err != nil {
			return nil, err
		}

		rt.Client, err = r.getRequestClient(ctx, client)
		if err != nil {
			return nil, err
		}
		rt.Worker, err = r.GetRequestWorker(ctx, worker)
		if err != nil {
			return nil, err
		}
		rt.Contract, err = r.getRequestContract(ctx, contract)
		if err != nil {
			return nil, err
		}
		rt.Equipment, err = r.getRequestEquipment(ctx, equipment)
		if err != nil {
			return nil, err
		}
		rt.ClientObject, err = r.getRequestClientObject(ctx, client_object)
		if err != nil {
			return nil, err
		}
		rt.Status, err = r.getRequestStatus(ctx, status)
		if err != nil {
			return nil, err
		}

		reqs = append(reqs, rt)
	}
	return reqs, nil
}

func (r *repository) SelectRequestsBySort(ctx context.Context, req requests.Request) (rs []requests.Request, err error) {
	query := `
	SELECT 
		id, title, client, worker, client_object, equipment, contract, description, priority, start_date, end_date, files, status, name
	FROM 
		public."Request"
	WHERE client = $1 AND worker = $2 AND client_object = $3 AND equipment = $4 AND status = $5
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, req.Client.Id, req.Worker.Id, req.ClientObject.Id, req.Equipment.Id, req.Status.Name)
	if err != nil {
		return nil, err
	}

	var rt requests.Request
	var client, worker, client_object, equipment, contract int
	var status string

	for rows.Next() {
		err = rows.Scan(&rt.Id, &rt.Title, &client, &worker, &client_object, &equipment, &contract, &rt.Description, &rt.Priority, &rt.StartDate, &rt.EndDate, pq.Array(&rt.Files), &status, &rt.Name)
		if err != nil {
			return nil, err
		}

		rt.Client, err = r.getRequestClient(ctx, client)
		if err != nil {
			return nil, err
		}
		rt.Worker, err = r.GetRequestWorker(ctx, worker)
		if err != nil {
			return nil, err
		}
		rt.Contract, err = r.getRequestContract(ctx, contract)
		if err != nil {
			return nil, err
		}
		rt.Equipment, err = r.getRequestEquipment(ctx, equipment)
		if err != nil {
			return nil, err
		}
		rt.ClientObject, err = r.getRequestClientObject(ctx, client_object)
		if err != nil {
			return nil, err
		}
		rt.Status, err = r.getRequestStatus(ctx, status)
		if err != nil {
			return nil, err
		}

		rs = append(rs, rt)
	}
	return rs, nil
}

// UpdateRequest implements Repository.
func (r *repository) UpdateRequest(ctx context.Context, req *requests.Request) error {
	query := `
	UPDATE 
		public."Request"
	SET 
		title = $1, client = $2, worker = $3, client_object = $4, equipment = $5, contract = $6, description = $7, priority = $8, status = $9, name = $10
	WHERE
		id = $11`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, req.Title, req.Client.Id, req.Worker.Id, req.ClientObject.Object.Id, req.Equipment.Id, req.Contract.Id, req.Description, req.Priority, req.Status.Name, req.Name, req.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Изменена", fmt.Sprintf("%s c id=%d сотрудником - %d / %s", "заявка", req.Id, req.Worker.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

// DeleteRequest implements Repository.
func (r *repository) CloseRequest(ctx context.Context, status string, id int) error {
	query := `
	UPDATE 
		public."Request"
	SET 
		status = $1
	WHERE
		id = $2`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, status, id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Закрыта", fmt.Sprintf("%s c id=%d / %s", "заявка", id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func NewRepository(client dbClient.Client, logger *logging.Logger) requests.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

// func (r *repository) GetInsertData(ctx context.Context) (data requests.RequestInsertDate, err error) {
// 	Title Priority Client Worker ClientObject Equipment Contract Status
// 	query := `
// 	SELECT id, equipment, type, cost FROM public."Services"
// 	`
// }

// client, worker, client_object, equipment, contract
// status

func (r *repository) getRequestClient(ctx context.Context, id int) (cl client.Client, err error) {
	query := `
	SELECT 
		id, name, inn, kpp, ogrn, owner, phone, email, address, create_date, status 
	FROM 
		public."Client"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&cl.Id, &cl.Name, &cl.INN, &cl.KPP, &cl.OGRN, &cl.Owner, &cl.Phone, &cl.Email, &cl.Address, &cl.CreateDate, &cl.Status)
	if err != nil {
		return client.Client{}, err
	}

	return cl, nil
}
func (r *repository) GetRequestWorker(ctx context.Context, id int) (us user.User, err error) {
	query := `
	SELECT 
		id, email, full_name, phone, image, role 
	FROM 
		public."User"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role)
	if err != nil {
		return user.User{}, err
	}

	return us, nil
}
func (r *repository) getRequestContract(ctx context.Context, id int) (ct contract.Contract, err error) {
	query := `
	SELECT 
		id, name, client, start_date, end_date, amount, file, status 
	FROM 
		public."Contract"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&ct.Id, &ct.Name, &ct.Client.Id, &ct.DataStart, &ct.DataEnd, &ct.Amount, &ct.File, &ct.Status)
	if err != nil {
		return contract.Contract{}, err
	}

	return ct, nil
}
func (r *repository) getRequestEquipment(ctx context.Context, id int) (eq equipment.Equipment, err error) {
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
	err = rows.Scan(&eq.Id, &eq.Name, &eq.Type, &eq.Manufacture, &eq.Model, &eq.UniqueNumber, &eq.Contract, &eq.CreateDate)
	if err != nil {
		return equipment.Equipment{}, err
	}

	return eq, nil
}
func (r *repository) getRequestClientObject(ctx context.Context, id int) (cl client.ClientObject, err error) {
	query := `
	SELECT 
		id, object
	FROM 
		public."Client_objects"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&cl.Id, &cl.Object.Id)
	if err != nil {
		return client.ClientObject{}, err
	}

	cl.Object, err = r.getRequestObject(ctx, cl.Object.Id)
	if err != nil {
		return client.ClientObject{}, err
	}

	return cl, nil
}
func (r *repository) getRequestStatus(ctx context.Context, id string) (rs requests.ReqStatus, err error) {
	query := `
	SELECT 
		name, color
	FROM 
		public."Request_status"
	WHERE 
		name = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&rs.Name, &rs.Color)
	if err != nil {
		return requests.ReqStatus{}, err
	}

	return rs, nil
}

func (r *repository) getRequestObject(ctx context.Context, id int) (obj objects.Object, err error) {
	query := `
	SELECT 
		id, name, address, work_schedule
	FROM 
		public."Objects"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule)
	if err != nil {
		return objects.Object{}, err
	}

	return obj, nil
}

func (r *repository) GetRequestPriority(ctx context.Context) (pr []string, err error) {
	query := `
	SELECT 
		name 
	FROM 
		public."Request_priority"
	`

	var p string

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&p)
		if err != nil {
			return nil, err
		}

		pr = append(pr, p)
	}
	return pr, nil
}

func (r *repository) GetRequestStatus(ctx context.Context) (rs []requests.ReqStatus, err error) {
	query := `
	SELECT 
		name, color 
	FROM 
		public."Request_status"
	`

	var p requests.ReqStatus

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&p.Name, &p.Color)
		if err != nil {
			return nil, err
		}

		rs = append(rs, p)
	}
	return rs, nil
}

func (r *repository) InsertRequestAnswer(ctx context.Context, ra *requests.ReqAns) error {
	query := `
	INSERT INTO public."Request_answer" (request_id, worker_id, text) VALUES ($1, $2, $3) RETURNING id
	`
	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, ra.Request.Id, ra.Worker.Id, ra.Text)
	err := rows.Scan(&ra.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SelectRequestAnswer(ctx context.Context, rid int) (ra []requests.ReqAns, err error) {
	query := `
		SELECT id, request_id, worker_id, "text" FROM public."Request_answer" WHERE request_id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, rid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var req requests.ReqAns

		err = rows.Scan(&req.Id, &req.Request.Id, &req.Worker.Id, &req.Text)
		if err != nil {
			return nil, err
		}

		req.Request, err = r.SelectRequest(ctx, req.Request.Id)
		if err != nil {
			return nil, err
		}

		req.Worker, err = r.GetRequestWorker(ctx, req.Worker.Id)
		if err != nil {
			return nil, err
		}

		ra = append(ra, req)
	}
	return ra, nil
}

func (r *repository) GetRequestWorkerByEmail(ctx context.Context, email string) (us user.User, err error) {
	query := `
	SELECT 
		id, email, full_name, phone, image, role 
	FROM 
		public."User"
	WHERE 
		email = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, email)
	err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role)
	if err != nil {
		return user.User{}, err
	}

	return us, nil
}

func (r *repository) FindRequests(ctx context.Context, find string) (rs []requests.Request, err error) {
	query := `
	SELECT 
		"Request".id, "Request".title, "Request".client, "Request".worker, "Request".client_object, "Request".equipment, "Request".contract, "Request".description, "Request".priority, "Request".start_date, "Request".end_date, "Request".files, "Request".status, "Request".name
	FROM 
		public."Request"
		JOIN "Client" ON "Request".client = "Client".id
		JOIN "User" ON "Request".worker = "User".id
		JOIN "Client_objects" ON "Request".client_object = "Client_objects".id
		JOIN "Objects" ON "Client_objects".object = "Objects".id
		JOIN "Equipment" ON "Request".equipment = "Equipment".id
	WHERE "Client".name ILIKE $1 OR "Request".title ILIKE $1 OR "Request".name ILIKE $1 OR "User".full_name ILIKE $1 OR "Objects".name ILIKE $1 OR "Equipment".name ILIKE $1
	
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	find = "%" + find + "%"

	rows, err := r.client.Query(ctx, query, find)
	if err != nil {
		return nil, err
	}

	var rt requests.Request
	var client, worker, client_object, equipment, contract int
	var status string

	for rows.Next() {
		err = rows.Scan(&rt.Id, &rt.Title, &client, &worker, &client_object, &equipment, &contract, &rt.Description, &rt.Priority, &rt.StartDate, &rt.EndDate, pq.Array(&rt.Files), &status, &rt.Name)
		if err != nil {
			return nil, err
		}

		rt.Client, err = r.getRequestClient(ctx, client)
		if err != nil {
			return nil, err
		}
		rt.Worker, err = r.GetRequestWorker(ctx, worker)
		if err != nil {
			return nil, err
		}
		rt.Contract, err = r.getRequestContract(ctx, contract)
		if err != nil {
			return nil, err
		}
		rt.Equipment, err = r.getRequestEquipment(ctx, equipment)
		if err != nil {
			return nil, err
		}
		rt.ClientObject, err = r.getRequestClientObject(ctx, client_object)
		if err != nil {
			return nil, err
		}
		rt.Status, err = r.getRequestStatus(ctx, status)
		if err != nil {
			return nil, err
		}

		rs = append(rs, rt)
	}
	return rs, nil
}
