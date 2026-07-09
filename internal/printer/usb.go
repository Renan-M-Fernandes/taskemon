package printer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	escposlib "github.com/Renan-M-Fernandes/taskemon/internal/escpos"
	"github.com/google/gousb"
)

var ErrPrinterNotFound = errors.New("printer not found")

func NewUSBPrinter(
	vendorID uint16,
	productID uint16,
	endpointNumber int,
	layout Layout,

) USBPrinter {
	return USBPrinter{
		VendorID:       gousb.ID(vendorID),
		ProductID:      gousb.ID(productID),
		EndpointNumber: endpointNumber,
		Layout:         layout,
	}
}

func (p USBPrinter) PrintTicket(ctx context.Context, ticket Ticket) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("print ticket: ctx %w", err)
	}

	usbCtx := gousb.NewContext()
	defer usbCtx.Close()

	dev, err := usbCtx.OpenDeviceWithVIDPID(p.VendorID, p.ProductID)
	if err != nil {
		return fmt.Errorf("print ticket: open usb printer: %w", err)
	}

	if dev == nil {
		return fmt.Errorf(
			"print ticket: %w: vendor=%#04x product=%#04x",
			ErrPrinterNotFound,
			uint16(p.VendorID),
			uint16(p.ProductID),
		)
	}

	defer dev.Close()

	err = dev.SetAutoDetach(true)
	if err != nil {
		return fmt.Errorf("print ticket: enable usb auto detach: %w", err)
	}

	intf, done, err := dev.DefaultInterface()
	if err != nil {
		return printerHint(fmt.Errorf("print ticket: claim printer interface: %w", err))
	}

	defer done()

	endpoint, err := intf.OutEndpoint(p.EndpointNumber)
	if err != nil {
		return fmt.Errorf("print ticket: open printer out endpoint %d: %w", p.EndpointNumber, err)
	}

	escposPrinter := escposlib.New(endpoint)

	if err := RenderTicket(escposPrinter, ticket, p.Layout); err != nil {
		return fmt.Errorf("print ticket: render ticket: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	if err := escposPrinter.PrintAndCut(p.Layout.FeedLinesBeforeCut, p.Layout.CutCommand); err != nil {
		return fmt.Errorf("print ticket: print ticket: %w", err)
	}
	return nil
}

func printerHint(err error) error {
	msg := err.Error()

	if strings.Contains(msg, "resource busy") ||
		strings.Contains(msg, "code -6") ||
		strings.Contains(msg, "claim interface") {
		return fmt.Errorf(
			"%w\n\n%s",
			err,
			"printer is busy: Linux may have claimed it with usblp. "+
				"Try enabling USB auto-detach "+
				"or run: sudo modprobe -r usblp",
		)
	}

	if strings.Contains(msg, "access") ||
		strings.Contains(msg, "permission") {
		return fmt.Errorf(
			"%w\n\n%s",
			err,
			"printer permission denied: run the service with access to the USB device group",
		)
	}

	return err
}
