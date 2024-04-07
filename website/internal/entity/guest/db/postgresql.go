package guest_db

import (
	"context"
	"fmt"

	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) guest.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) SelectGuests(ctx context.Context) (guests []guest.Guests, err error) {

	query := `
		SELECT 
			id, email, firstname, lastname, patronymic, phone, password, age, organization, card_number, card_date, card_cvv, card_bank, banned
		FROM 
			public."Guest"
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var g guest.Guests

		err = rows.Scan(&g.Id, &g.Email, &g.Firstname, &g.Lastname, &g.Patronymic, &g.Phone, &g.Password, &g.Age, &g.Organization.INN, &g.SaveCard.Number, &g.SaveCard.Date, &g.SaveCard.CVV, &g.SaveCard.Bank, &g.Banned)
		if err != nil {
			return nil, err
		}

		g.Organization, err = r.SelectOrganization(ctx, g.Organization.INN)
		if err != nil {
			return nil, err
		}

		guests = append(guests, g)
	}
	return guests, nil
}

func (r *repository) SelectGuestByColumn(ctx context.Context, column string, value interface{}) (g guest.Guests, err error) {

	query := `
		SELECT 
			id, email, firstname, lastname, patronymic, phone, password, age, organization, card_number, card_date, card_cvv, card_bank, banned
		FROM 
			public."Guest"
		WHERE $1 = $2
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, column, value)

	err = rows.Scan(&g.Id, &g.Email, &g.Firstname, &g.Lastname, &g.Patronymic, &g.Phone, &g.Password, &g.Age, &g.Organization.INN, &g.SaveCard.Number, &g.SaveCard.Date, &g.SaveCard.CVV, &g.SaveCard.Bank, &g.Banned)
	if err != nil {
		return guest.Guests{}, err
	}

	g.Organization, err = r.SelectOrganization(ctx, g.Organization.INN)
	if err != nil {
		return guest.Guests{}, err
	}

	return g, nil
}

func (r *repository) SelectOrganization(ctx context.Context, inn int) (org guest.Organization, err error) {
	query := `
		SELECT 
			inn, kpp, ogrn, name, fullname, city, country
		FROM 
			public."Organization"
		WHERE inn = $1
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, inn)

	err = rows.Scan(&org.INN, &org.KPP, &org.OGRN, &org.Name, &org.Fullname, &org.City, &org.Country)
	if err != nil {
		return guest.Organization{}, err
	}

	return org, nil
}

func (r *repository) InsertOrganization(ctx context.Context, org guest.Organization) error {
	query := `
	INSERT INTO 
		public."Ogranization" (inn, kpp, ogrn, name, fullname, city, country)
	VALUEST 
		($1, $2, $3, $4, $5, $6, $7)
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, org.INN, org.KPP, org.OGRN, org.Name, org.Fullname, org.City, org.Country)

	return nil
}

func (r *repository) InsertGuest(ctx context.Context, gst *guest.Guests) error {
	query := `
	INSERT INTO 
		public."Guest" (email, firstname, lastname, patronymic, phone, password, age, organization, card_number, card_date, card_cvv, card_bank)
	VALUEST 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id	
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	row := r.client.QueryRow(ctx, query, &gst.Email, &gst.Firstname, &gst.Lastname, &gst.Patronymic, &gst.Phone, &gst.Password, &gst.Age, &gst.Organization.INN, &gst.SaveCard.Number, &gst.SaveCard.Date, &gst.SaveCard.CVV, &gst.SaveCard.Bank)

	err := row.Scan(&gst.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateOrganization(ctx context.Context, ogr guest.Organization) error {
	query := `
	UPDATE INTO 
		public."Organization"
	SET
		kpp = $1, ogrn = $2, name = $3, fullname = $4, city = $5, country = $6
	WHERE 
		inn = $7
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, ogr.KPP, ogr.OGRN, ogr.Name, ogr.Fullname, ogr.City, ogr.Country, ogr.INN)

	return nil
}

func (r *repository) UpdateGuest(ctx context.Context, gst guest.Guests) error {
	query := `
	UPDATE INTO
		public."Guest"
	SET 
		email = $1, firstname = $2, lastname = $3, patronymic = $4, phone = $5, password = $6, age = $7, organization = $8, card_number = $9, card_date = $10, card_cvv = $11, card_bank = $12
	WHERE
		id = $13
	`
	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, gst.Email, gst.Firstname, gst.Lastname, gst.Patronymic, gst.Phone, gst.Password, gst.Age, gst.Organization.INN, gst.SaveCard.Number, gst.SaveCard.Date, gst.SaveCard.CVV, gst.SaveCard.Bank, gst.Id)

	return nil
}

func (r *repository) BannedGuest(ctx context.Context, id int) error {
	query := `
	UPDATE INTO
		public."Guest"
	SET 
		banned = true
	WHERE
		id = $1
	`
	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query, id)

	return nil
}

func (r *repository) CheckAuthGuest(ctx context.Context, email string, password string) (g guest.Guests, err error) { //work
	query := `
		SELECT 
			id, email, firstname, lastname, patronymic, phone, password, age, organization, card_number, card_date, card_cvv, card_bank, banned
		FROM 
			public."Guest"
		WHERE email = $1 AND password = $2 AND banned = false
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, email, password)

	err = rows.Scan(&g.Id, &g.Email, &g.Firstname, &g.Lastname, &g.Patronymic, &g.Phone, &g.Password, &g.Age, &g.Organization.INN, &g.SaveCard.Number, &g.SaveCard.Date, &g.SaveCard.CVV, &g.SaveCard.Bank, &g.Banned)
	if err != nil {
		return guest.Guests{}, err
	}

	g.Organization, err = r.SelectOrganization(ctx, g.Organization.INN)
	if err != nil {
		return guest.Guests{}, err
	}

	return g, nil
}

func (r *repository) SelectTrustCompany(ctx context.Context) (tc []guest.TrustCompany, err error) {
	query := `
	SELECT name, description, logo FROM public."TrustCompany"
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
			fmt.Println(err)
			return nil, err
		}

		tc = append(tc, t)
	}

	return tc, nil
}
