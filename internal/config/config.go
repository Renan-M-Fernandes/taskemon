package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Printer  PrinterConfig  `json:"printer"`
}

type ServerConfig struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	CORSOrigins []string `json:"corsOrigins"`
}

type DatabaseConfig struct {
	Path string `json:"path"`
}

type PrinterConfig struct {
	Enabled    bool   `json:"enabled"`
	Transport  string `json:"transport"`
	VendorID   string `json:"vendorID"`
	ProductID  string `json:"productID"`
	Endpoint   int    `json:"endpoint"`
	DevicePath string `json:"devicePath"`
	CutCommand string `json:"cutCommand"`
	QRMode     string `json:"qrMode"`
	BaseURL    string `json:"baseURL"`
	Layout     Layout `json:"layout"`
}

type Layout struct {
	PaperWidth         int         `json:"paperWidth"`
	CharsPerLine       int         `json:"charsPerLine"`
	Header             header      `json:"header"`
	Title              title       `json:"title"`
	Description        description `json:"description"`
	Tag                tag         `json:"tag"`
	QR                 qr          `json:"qr"`
	ShinyHint          ShinyHint   `json:"shinyHint"`
	Footer             footer      `json:"footer"`
	FeedLinesBeforeCut int         `json:"feedLinesBeforeCut"`
}

type header struct {
	Enabled bool   `json:"enabled"`
	Bold    bool   `json:"bold"`
	Size    int    `json:"size"`
	Justify string `json:"justify"`
	Style   int    `json:"style"`
}
type title struct {
	Bold                bool   `json:"bold"`
	Size                int    `json:"size"`
	Justify             string `json:"justify"`
	LinesAllwaysVisible int    `json:"linesAllwaysVisible"`
	MaxLines            int    `json:"maxLines"`
}

type description struct {
	Enabled             bool   `json:"enabled"`
	Bold                bool   `json:"bold"`
	Size                int    `json:"size"`
	Justify             string `json:"justify"`
	LinesAllwaysVisible int    `json:"linesAllwaysVisible"`
	MaxLines            int    `json:"maxLines"`
}

type tag struct {
	Enabled bool   `json:"enabled"`
	Bold    bool   `json:"bold"`
	Size    int    `json:"size"`
	Justify string `json:"justify"`
	Style   int    `json:"style"`
}

type qr struct {
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

type footer struct {
	Enabled bool   `json:"enabled"`
	Bold    bool   `json:"bold"`
	Size    int    `json:"size"`
	Justify string `json:"justify"`
}

func Load(path string) (Config, error) {
	cfg := Default()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}

		return Config{}, fmt.Errorf("open config: %w", err)
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}

	cfg.Printer.Transport = strings.ToLower(cfg.Printer.Transport)
	cfg.Printer.CutCommand = strings.ToLower(cfg.Printer.CutCommand)
	cfg.Printer.QRMode = strings.ToLower(cfg.Printer.QRMode)
	cfg.Printer.Layout.Header.Justify = strings.ToLower(cfg.Printer.Layout.Header.Justify)
	cfg.Printer.Layout.Title.Justify = strings.ToLower(cfg.Printer.Layout.Title.Justify)
	cfg.Printer.Layout.Description.Justify = strings.ToLower(cfg.Printer.Layout.Description.Justify)
	cfg.Printer.Layout.Tag.Justify = strings.ToLower(cfg.Printer.Layout.Tag.Justify)
	cfg.Printer.Layout.QR.Justify = strings.ToLower(cfg.Printer.Layout.QR.Justify)
	cfg.Printer.Layout.ShinyHint.Justify = strings.ToLower(cfg.Printer.Layout.ShinyHint.Justify)
	cfg.Printer.Layout.Footer.Justify = strings.ToLower(cfg.Printer.Layout.Footer.Justify)

	return cfg, nil
}

func Default() Config {
	return Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
			CORSOrigins: []string{
				"http://localhost:8123",
			},
		},
		Database: DatabaseConfig{
			Path: "./database/taskemon.db",
		},
		Printer: PrinterConfig{
			Enabled:    false,
			Transport:  "noop",
			VendorID:   "0x0418",
			ProductID:  "0x5011",
			DevicePath: "/dev/usb/lp0",
			Endpoint:   1,
			CutCommand: "gs_v_0",
			QRMode:     "pokemon_placeholder",
			Layout: Layout{
				PaperWidth:   80,
				CharsPerLine: 48,

				Header: header{
					Enabled: true,
					Bold:    true,
					Size:    1,
					Justify: "center",
					Style:   1,
				},

				Title: title{
					Bold:                true,
					Size:                2,
					Justify:             "left",
					LinesAllwaysVisible: 1,
					MaxLines:            2,
				},

				Description: description{
					Enabled:             true,
					Bold:                false,
					Size:                1,
					Justify:             "left",
					LinesAllwaysVisible: 0,
					MaxLines:            0,
				},

				Tag: tag{
					Enabled: true,
					Bold:    true,
					Size:    1,
					Justify: "center",
					Style:   1,
				},

				QR: qr{
					Enabled:    true,
					Size:       4,
					Justify:    "right",
					Correction: "H",
				},

				ShinyHint: ShinyHint{
					Enabled: true,
					Bold:    false,
					Size:    1,
					Justify: "center",
				},

				Footer: footer{
					Enabled: true,
					Bold:    true,
					Size:    1,
					Justify: "center",
				},

				FeedLinesBeforeCut: 1,
			},
		},
	}
}
