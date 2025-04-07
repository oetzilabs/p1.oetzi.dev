package broker

type Broker struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Servers []string `json:"servers"`
}

func NewBroker(id, name, url string) Broker {

	return Broker{
		ID:      id,
		Name:    name,
		URL:     url,
		Servers: []string{},
	}
}

func (b *Broker) UpdateServer(server string) error {
	// TODO: implement
	// we gotta check if the server is online/available and has good metrics
	return nil
}

func (b *Broker) Update() error {
	// TODO: implement
	for _, server := range b.Servers {
		// TODO: implement
		server := server
		err := b.UpdateServer(server)
		if err != nil {
			return err
		}
	}
	return nil
}
