package printer

import (
	"testing"

	"github.com/Renan-M-Fernandes/taskemon/internal/config"
)

func TestNewFromConfigReturnsNoopWhenDisabled(t *testing.T) {
	cfg := config.Default()
	cfg.Printer.Enabled = false
	cfg.Printer.Transport = "usb"

	got, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig returned error: %v", err)
	}

	if _, ok := got.(NoopPrinter); !ok {
		t.Fatalf("printer type mismatch: got %T, want NoopPrinter", got)
	}
}

func TestNewFromConfigReturnsNoopTransport(t *testing.T) {
	cfg := config.Default()
	cfg.Printer.Enabled = true
	cfg.Printer.Transport = "noop"

	got, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig returned error: %v", err)
	}

	if _, ok := got.(NoopPrinter); !ok {
		t.Fatalf("printer type mismatch: got %T, want NoopPrinter", got)
	}
}

func TestNewFromConfigBuildsUSBPrinter(t *testing.T) {
	cfg := config.Default()
	cfg.Printer.Enabled = true
	cfg.Printer.Transport = "usb"
	cfg.Printer.VendorID = "0x0418"
	cfg.Printer.ProductID = "5011"
	cfg.Printer.Endpoint = 2
	cfg.Printer.CutCommand = "partial"
	cfg.Printer.Layout.CharsPerLine = 64
	cfg.Printer.Layout.FeedLinesBeforeCut = 3
	cfg.Printer.Layout.QR.Size = 7
	cfg.Printer.Layout.QR.Correction = "Q"

	got, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig returned error: %v", err)
	}

	usb, ok := got.(USBPrinter)
	if !ok {
		t.Fatalf("printer type mismatch: got %T, want USBPrinter", got)
	}
	if uint16(usb.VendorID) != 0x0418 {
		t.Fatalf("vendor id mismatch: got %#x, want %#x", uint16(usb.VendorID), 0x0418)
	}
	if uint16(usb.ProductID) != 0x5011 {
		t.Fatalf("product id mismatch: got %#x, want %#x", uint16(usb.ProductID), 0x5011)
	}
	if usb.EndpointNumber != 2 {
		t.Fatalf("endpoint mismatch: got %d, want 2", usb.EndpointNumber)
	}
	if usb.Layout.CutCommand != "partial" {
		t.Fatalf("cut command mismatch: got %q, want partial", usb.Layout.CutCommand)
	}
	if usb.Layout.CharsPerLine != 64 {
		t.Fatalf("chars per line mismatch: got %d, want 64", usb.Layout.CharsPerLine)
	}
	if usb.Layout.FeedLinesBeforeCut != 3 {
		t.Fatalf("feed lines mismatch: got %d, want 3", usb.Layout.FeedLinesBeforeCut)
	}
	if usb.Layout.QR.Size != 7 {
		t.Fatalf("qr size mismatch: got %d, want 7", usb.Layout.QR.Size)
	}
	if usb.Layout.QR.Correction != "Q" {
		t.Fatalf("qr correction mismatch: got %q, want Q", usb.Layout.QR.Correction)
	}
}

func TestNewFromConfigRejectsInvalidIDs(t *testing.T) {
	tests := []struct {
		name      string
		vendorID  string
		productID string
	}{
		{name: "vendor", vendorID: "nope", productID: "0x5011"},
		{name: "product", vendorID: "0x0418", productID: "nope"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Default()
			cfg.Printer.Enabled = true
			cfg.Printer.Transport = "usb"
			cfg.Printer.VendorID = tt.vendorID
			cfg.Printer.ProductID = tt.productID

			_, err := NewFromConfig(cfg)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestNewFromConfigRejectsUnsupportedTransport(t *testing.T) {
	cfg := config.Default()
	cfg.Printer.Enabled = true
	cfg.Printer.Transport = "bluetooth"

	_, err := NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected unsupported transport error")
	}
}

func TestParseHexID(t *testing.T) {
	tests := []struct {
		input string
		want  uint16
	}{
		{input: "0x0418", want: 0x0418},
		{input: "0X5011", want: 0x5011},
		{input: "5011", want: 0x5011},
		{input: " 0418 ", want: 0x0418},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseHexID(tt.input)
			if err != nil {
				t.Fatalf("parseHexID returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("id mismatch: got %#x, want %#x", got, tt.want)
			}
		})
	}
}
