package contract

type Contract struct {
	Id        int
	Name      string
	Client    Client
	DataStart string
	DataEnd   string
	Amount    int
	File      string
	Status    bool
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
