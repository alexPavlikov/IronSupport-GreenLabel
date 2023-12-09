package yandexdisk

type (
	Config struct {
		AuthToken   string
		ProjectList []Project
	}
)

func readConfig() {
	config.AuthToken = "a6f1c482280a473aa1f210999f85f643"
	config.ProjectList = []Project{
		{
			"projectA", "docker", "C:/Users/admin/Desktop/ФИЗТЕХ/Санин/1-16.doc", "", "C:/Users/admin/Desktop/ФИЗТЕХ/Санин/1-16.doc",
		}, {
			"defly", "docker", "", "defly_postgres_1", "",
		},
	}
}
