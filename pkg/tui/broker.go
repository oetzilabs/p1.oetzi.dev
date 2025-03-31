package tui

type Broker struct {
	Servers []Server `json:"servers"`
}
type BrokerData struct {
	Data Broker `json:"data"`
}
