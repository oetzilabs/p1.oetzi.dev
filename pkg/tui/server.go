package tui

type Server struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func NewServer(id string, name string, url string) Server {
	return Server{
		Id:   id,
		Name: name,
		Url:  url,
	}
}
