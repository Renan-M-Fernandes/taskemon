package printer

import "time"

type Ticket struct {
	TaskID      int
	UserID      string
	Title       string
	Description string
	Tag         string
	DueAt       *time.Time
	CreatedAt   time.Time
	Shiny       bool
	QRValue     string
}
