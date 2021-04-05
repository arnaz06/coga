package coga

// Message is the struct represent message data.
type Message struct {
	ID      string `json:"id"`
	Service string `json:"service"`
	Event   string `json:"event"`
	Data    []byte `json:"data"`
}
