package contract_db

import (
	"context"
	"fmt"
	"time"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/contract"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) contract.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertContract(ctx context.Context, contract *contract.Contract) error {
	query := `
	INSERT INTO public."Contract" 
		(name, client, start_date, end_date, amount, file, status)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, &contract.Name, &contract.Client.Id, &contract.DataStart, &contract.DataEnd, &contract.Amount, &contract.File, &contract.Status)

	err := rows.Scan(&contract.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Добавлен", fmt.Sprintf("%s c id=%d / %s", "контракт", contract.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) SelectContract(ctx context.Context, id int) (contract contract.Contract, err error) {
	query := `
	SELECT 
		id, name, client, start_date, end_date, amount, file, status
	FROM 
		public."Contract"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)

	err = rows.Scan(&contract.Id, &contract.Name, &contract.Client.Id, &contract.DataStart, &contract.DataEnd, &contract.Amount, &contract.File, &contract.Status)
	if err != nil {
		return contract, err
	}

	contract.Client, err = r.getClientObject(ctx, contract.Client.Id)
	if err != nil {
		return contract, err
	}

	return contract, nil
}

func (r *repository) SelectContractsBySort(ctx context.Context, ct contract.Contract) (cts []contract.Contract, err error) {
	query := `
	SELECT 
		"Contract".id, "Contract".name, "Contract".client, "Contract".start_date, "Contract".end_date, "Contract".amount, "Contract".file, "Contract".status, 
		"Client".name
	FROM 
		public."Contract"
	JOIN "Client" ON "Client".id = "Contract".client	
	WHERE 
		"Client".name ILIKE $1 OR "Contract".start_date ILIKE $2 OR "Contract".end_date ILIKE $3 OR "Contract".status = $4;
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	var contract contract.Contract

	if ct.Client.Name != "" {
		ct.Client.Name = "%" + ct.Client.Name + "%"
	}
	if ct.DataStart != "" {
		ct.DataStart = "%" + ct.DataStart + "%"
	}
	if ct.DataEnd != "" {
		ct.DataEnd = "%" + ct.DataEnd + "%"
	}

	rows, err := r.client.Query(ctx, query, ct.Client.Name, ct.DataStart, ct.DataEnd, ct.Status)
	if err != nil {
		return cts, err
	}
	for rows.Next() {
		err = rows.Scan(&contract.Id, &contract.Name, &contract.Client.Id, &contract.DataStart, &contract.DataEnd, &contract.Amount, &contract.File, &contract.Status, &contract.Client.Name)
		if err != nil {
			return nil, err
		}
		contract.Client, err = r.getClientObject(ctx, contract.Client.Id)
		if err != nil {
			return nil, err
		}

		cts = append(cts, contract)
	}

	return cts, nil
}

func (r *repository) SelectContracts(ctx context.Context) (contracts []contract.Contract, err error) {
	query := `
	SELECT 
		id, name, client, start_date, end_date, amount, file, status
	FROM 
		public."Contract"
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	var clientId int

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var contract contract.Contract
	for rows.Next() {
		err = rows.Scan(&contract.Id, &contract.Name, &clientId, &contract.DataStart, &contract.DataEnd, &contract.Amount, &contract.File, &contract.Status)
		if err != nil {
			return nil, err
		}
		fmt.Println(clientId, "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		contract.Client, err = r.getClientObject(ctx, clientId)
		if err != nil {
			return nil, err
		}

		contracts = append(contracts, contract)
	}
	return contracts, nil
}

func (r *repository) UpdateContract(ctx context.Context, contract *contract.Contract) error {
	query := `
	UPDATE 
		public."Contract"
	SET 
		name = $1, client = $2, start_date = $3, end_date = $4, amount = $5, status = $6
	WHERE 
		id = $7
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	fmt.Println("llllllllllll", contract.Name, contract.Client.Id, contract.DataStart, contract.DataEnd, contract.Amount, contract.Status, contract.Id)

	_ = r.client.QueryRow(ctx, query, contract.Name, contract.Client.Id, contract.DataStart, contract.DataEnd, contract.Amount, contract.Status, contract.Id)

	r.logger.LogEvents("Изменен", fmt.Sprintf("%s c id=%d / %s", "контракт", contract.Id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) CloseContract(ctx context.Context, id int) error {
	query := `
	UPDATE 
		public."Contract"
	SET 
		status = "false"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Закрыт", fmt.Sprintf("%s c id=%d / %s", "контракт", id, fmt.Sprint(time.Now().Format("15:04 2006-01-02"))))

	return nil
}

func (r *repository) getClientObject(ctx context.Context, id int) (cl contract.Client, err error) {
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
		return contract.Client{}, err
	}

	return cl, nil
}

//-----

func (r *repository) SelectClients(ctx context.Context) (clnts []contract.Client, err error) {
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
	var cl contract.Client
	for rows.Next() {
		err = rows.Scan(&cl.Id, &cl.Name, &cl.INN, &cl.KPP, &cl.OGRN, &cl.Owner, &cl.Phone, &cl.Email, &cl.Address, &cl.CreateDate, &cl.Status)
		if err != nil {
			return nil, err
		}
		clnts = append(clnts, cl)
	}
	return clnts, nil
}

func (r *repository) FindContract(ctx context.Context, text string) (cts []contract.Contract, err error) {
	query := `
	SELECT 
		"Contract".id, "Contract".name, "Contract".client, "Contract".start_date, "Contract".end_date, "Contract".amount, "Contract".file, "Contract".status, "Client".name
	FROM 
		public."Contract"
	JOIN "Client" ON "Client".id = "Contract".client	
	WHERE 
		"Client".name ILIKE $1 OR "Contract".name ILIKE $1;
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	var contract contract.Contract

	text = "%" + text + "%"

	rows, err := r.client.Query(ctx, query, text)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&contract.Id, &contract.Name, &contract.Client.Id, &contract.DataStart, &contract.DataEnd, &contract.Amount, &contract.File, &contract.Status, &contract.Client.Name)
		if err != nil {
			return nil, err
		}
		contract.Client, err = r.getClientObject(ctx, contract.Client.Id)
		if err != nil {
			return nil, err
		}

		cts = append(cts, contract)
	}

	return cts, nil
}
