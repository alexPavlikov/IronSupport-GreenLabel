package guest

import "context"

type Repository interface {
	SelectGuests(ctx context.Context) (guests []Guests, err error)
	SelectGuestByColumn(ctx context.Context, column string, value interface{}) (g Guests, err error)
	SelectOrganization(ctx context.Context, inn int) (org Organization, err error)
	InsertOrganization(ctx context.Context, org Organization) error
	InsertGuest(ctx context.Context, gst *Guests) error
	UpdateOrganization(ctx context.Context, ogr Organization) error
	UpdateGuest(ctx context.Context, gst Guests) error
	BannedGuest(ctx context.Context, id int) error
	CheckAuthGuest(ctx context.Context, email string, password string) (g Guests, err error)
	SelectTrustCompany(ctx context.Context) (tc []TrustCompany, err error)
}
