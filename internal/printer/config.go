package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Renan-M-Fernandes/taskemon/internal/config"
	"github.com/google/gousb"
)

type Config struct {
	Enabled    bool   `json:"enabled"`
	Transport  string `json:"transport"`
	VendorID   string `json:"vendorID"`
	ProductID  string `json:"productID"`
	Endpoint   int    `json:"endpoint"`
	DevicePath string `json:"devicePath"`
	CutCommand string `json:"cutCommand"`
	Layout     Layout `json:"layout"`
}

type Layout struct {
	PaperWidth         int         `json:"paperWidth"`
	CharsPerLine       int         `json:"charsPerLine"`
	CutCommand         string      `json:"cutCommand"`
	Header             Header      `json:"header"`
	Title              Title       `json:"title"`
	Description        Description `json:"description"`
	Tag                Tag         `json:"tag"`
	QR                 QR          `json:"qr"`
	ShinyHint          ShinyHint   `json:"shinyHint"`
	Footer             Footer      `json:"footer"`
	FeedLinesBeforeCut int         `json:"feedLinesBeforeCut"`
}

type Header struct {
	Enabled bool   `json:"enabled"`
	Bold    bool   `json:"bold"`
	Size    int    `json:"size"`
	Justify string `json:"justify"`
	Style   int    `json:"style"`
}
type Title struct {
	Bold                bool   `json:"bold"`
	Size                int    `json:"size"`
	Justify             string `json:"justify"`
	LinesAllwaysVisible int    `json:"linesAllwaysVisible"`
	MaxLines            int    `json:"maxLines"`
}

type Description struct {
	Enabled             bool   `json:"enabled"`
	Bold                bool   `json:"bold"`
	Size                int    `json:"size"`
	Justify             string `json:"justify"`
	LinesAllwaysVisible int    `json:"linesAllwaysVisible"`
	MaxLines            int    `json:"maxLines"`
}

type Tag struct {
	Enabled bool   `json:"enabled"`
	Bold    bool   `json:"bold"`
	Size    int    `json:"size"`
	Justify string `json:"justify"`
	Style   int    `json:"style"`
}

type QR struct {
	Enabled    bool   `json:"enabled"`
	Size       int    `json:"size"`
	Justify    string `json:"justify"`
	Correction string `json:"correction"`
}

type ShinyHint struct {
	Enabled bool   `json:"enabled"`
	Bold    bool   `json:"bold"`
	Size    int    `json:"size"`
	Justify string `json:"justify"`
}

type Footer struct {
	Enabled bool   `json:"enabled"`
	Bold    bool   `json:"bold"`
	Size    int    `json:"size"`
	Justify string `json:"justify"`
}

type USBPrinter struct {
	VendorID       gousb.ID
	ProductID      gousb.ID
	EndpointNumber int
	Layout         Layout
}

func NewFromConfig(cfg config.Config) (Printer, error) {
	if !cfg.Printer.Enabled {
		return NoopPrinter{}, nil
	}

	layout := Layout{
		PaperWidth:   cfg.Printer.Layout.PaperWidth,
		CharsPerLine: cfg.Printer.Layout.CharsPerLine,
		CutCommand:   cfg.Printer.CutCommand,
		Header: Header{
			Enabled: cfg.Printer.Layout.Header.Enabled,
			Bold:    cfg.Printer.Layout.Header.Bold,
			Size:    cfg.Printer.Layout.Header.Size,
			Justify: cfg.Printer.Layout.Header.Justify,
			Style:   cfg.Printer.Layout.Header.Style,
		},
		Title: Title{
			Bold:                cfg.Printer.Layout.Title.Bold,
			Size:                cfg.Printer.Layout.Title.Size,
			Justify:             cfg.Printer.Layout.Title.Justify,
			LinesAllwaysVisible: cfg.Printer.Layout.Title.LinesAllwaysVisible,
			MaxLines:            cfg.Printer.Layout.Title.MaxLines,
		},
		Description: Description{
			Enabled:             cfg.Printer.Layout.Description.Enabled,
			Bold:                cfg.Printer.Layout.Description.Bold,
			Size:                cfg.Printer.Layout.Description.Size,
			Justify:             cfg.Printer.Layout.Description.Justify,
			LinesAllwaysVisible: cfg.Printer.Layout.Description.LinesAllwaysVisible,
			MaxLines:            cfg.Printer.Layout.Description.MaxLines,
		},
		Tag: Tag{
			Enabled: cfg.Printer.Layout.Tag.Enabled,
			Bold:    cfg.Printer.Layout.Tag.Bold,
			Size:    cfg.Printer.Layout.Tag.Size,
			Justify: cfg.Printer.Layout.Tag.Justify,
			Style:   cfg.Printer.Layout.Tag.Style,
		},
		QR: QR{
			Enabled:    cfg.Printer.Layout.QR.Enabled,
			Size:       cfg.Printer.Layout.QR.Size,
			Justify:    cfg.Printer.Layout.QR.Justify,
			Correction: cfg.Printer.Layout.QR.Correction,
		},
		ShinyHint: ShinyHint{
			Enabled: cfg.Printer.Layout.ShinyHint.Enabled,
			Bold:    cfg.Printer.Layout.ShinyHint.Bold,
			Size:    cfg.Printer.Layout.ShinyHint.Size,
			Justify: cfg.Printer.Layout.ShinyHint.Justify,
		},
		Footer: Footer{
			Enabled: cfg.Printer.Layout.Footer.Enabled,
			Bold:    cfg.Printer.Layout.Footer.Bold,
			Size:    cfg.Printer.Layout.Footer.Size,
			Justify: cfg.Printer.Layout.Footer.Justify,
		},
		FeedLinesBeforeCut: cfg.Printer.Layout.FeedLinesBeforeCut,
	}

	switch cfg.Printer.Transport {
	case "usb":
		vendorID, err := parseHexID(cfg.Printer.VendorID)
		if err != nil {
			return nil, fmt.Errorf("new from config: parse printer vendor id: %w", err)
		}

		productID, err := parseHexID(cfg.Printer.ProductID)
		if err != nil {
			return nil, fmt.Errorf("new from config: parse printer product id: %w", err)
		}

		return NewUSBPrinter(vendorID, productID, cfg.Printer.Endpoint, layout), nil
	// Will add network printer and such later

	case "noop":
		return NoopPrinter{}, nil

	default:
		return nil, fmt.Errorf("new from config: unsupported printer transport %q", cfg.Printer.Transport)
	}
}

func parseHexID(value string) (uint16, error) {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "0x")
	value = strings.TrimPrefix(value, "0X")

	id, err := strconv.ParseUint(value, 16, 16)
	if err != nil {
		return 0, err
	}

	return uint16(id), nil
}
