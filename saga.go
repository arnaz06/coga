package coga

// list of coga event
var (
	EventStart    = "start"
	EventRollback = "rollback"
)

// Message is the struct represent message data.
type Message struct {
	ID      string      `json:"id"`
	Service string      `json:"service"`
	Event   string      `json:"event"`
	Data    interface{} `json:"data"`
}
