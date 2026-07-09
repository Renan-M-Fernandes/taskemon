package printer

import "context"

type NoopPrinter struct{}

func (NoopPrinter) PrintTicket(ctx context.Context, ticket Ticket) error {
	return nil
}
