package util

import (
	"time"

	"github.com/jordic/goics"
)

// ClassEvent represent info of a class
type ClassEvent struct {
	StartTime, EndTime                   time.Time
	Name, Location, Teacher, Description string
}

// ICSEvent aims to encode an event into ics format
type ICSEvent struct {
	component goics.Componenter
}

// EmitICal is a must for ICalEmitter
func (evt *ICSEvent) EmitICal() goics.Componenter {
	return evt.component
}
