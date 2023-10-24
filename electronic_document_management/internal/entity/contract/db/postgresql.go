package contract_db

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/contract"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/pkg/utils"
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
	RETURNIMG id
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, &contract.Name, &contract.Client.Id, &contract.DataStart, &contract.DataEnd, &contract.Amount, &contract.File, &contract.Status)

	err := rows.Scan(contract.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Добавлен", fmt.Sprintf("%s c id=:%d", "контракт", &contract.Id))

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

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return contract, err
	}
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
		name = $1, client = $2, start_date = $3, end_date = $4, amount = $5, file = $6, status = $7
	WHERE 
		id = $8
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, &contract.Name, &contract.Client.Id, &contract.DataStart, &contract.DataEnd, &contract.Amount, &contract.File, &contract.Status, &contract.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Изменен", fmt.Sprintf("%s c id=:%d", "контракт", &contract.Id))

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

	r.logger.LogEvents("Закрыт", fmt.Sprintf("%s c id=:%d", "контракт", id))

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
