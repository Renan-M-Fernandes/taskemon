package printer

import "context"

type Printer interface {
	PrintTicket(ctx context.Context, ticket Ticket) error
}
