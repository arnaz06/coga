package ebus

import (
	"github.com/arnaz06/coga"
)

type Handler interface {
	Handle(coga.SystemEvent)
}
