package api

type WebSocketDataUpdate struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type WebSocketMessageData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
