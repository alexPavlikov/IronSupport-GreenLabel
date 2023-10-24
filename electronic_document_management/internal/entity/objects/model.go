package objects

type Object struct {
	Id             int
	Name           string
	Address        string
	WorkSchedule   string
	Client         Client
	ClientObjectId int
}
type Client struct {
	Id         int
	Name       string
	INN        string
	KPP        string
	OGRN       string
	Owner      string
	Phone      string
	Email      string
	Address    string
	CreateDate string
	Status     bool
}
