package equipment

type Equipment struct {
	Id           int
	Name         string
	Type         string
	Manufacture  string
	Model        string
	UniqueNumber string
	Contract     string
	CreateDate   string
}

type SortEq struct {
	Type        []string
	Manufacture []string
	Model       []string
}
