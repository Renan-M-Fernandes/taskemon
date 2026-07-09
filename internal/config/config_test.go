package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	if cfg.Server.Host != "0.0.0.0" {
		t.Fatalf("server host mismatch: got %q, want %q", cfg.Server.Host, "0.0.0.0")
	}
	if cfg.Server.Port != 8080 {
		t.Fatalf("server port mismatch: got %d, want 8080", cfg.Server.Port)
	}
	if cfg.Database.Path != "./database/taskemon.db" {
		t.Fatalf("database path mismatch: got %q", cfg.Database.Path)
	}
	if cfg.Printer.Transport != "noop" {
		t.Fatalf("printer transport mismatch: got %q, want noop", cfg.Printer.Transport)
	}
	if cfg.Printer.QRMode != "pokemon_placeholder" {
		t.Fatalf("printer qr mode mismatch: got %q, want pokemon_placeholder", cfg.Printer.QRMode)
	}
	if cfg.Printer.Layout.CharsPerLine != 48 {
		t.Fatalf("chars per line mismatch: got %d, want 48", cfg.Printer.Layout.CharsPerLine)
	}
	if cfg.Printer.Layout.Description.LinesAllwaysVisible != 0 {
		t.Fatalf("description lines mismatch: got %d, want 0", cfg.Printer.Layout.Description.LinesAllwaysVisible)
	}
}

func TestLoadMissingFileReturnsDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Printer.Transport != "noop" {
		t.Fatalf("printer transport mismatch: got %q, want noop", cfg.Printer.Transport)
	}
	if cfg.Printer.QRMode != "pokemon_placeholder" {
		t.Fatalf("printer qr mode mismatch: got %q, want pokemon_placeholder", cfg.Printer.QRMode)
	}
}

func TestLoadOverridesAndNormalizesEnums(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	data := []byte(`{
		"database": {
			"path": "./Database/Taskemon.DB"
		},
		"printer": {
			"enabled": true,
			"transport": "USB",
			"vendorID": "0X0418",
			"productID": "0X5011",
			"devicePath": "/dev/USB/lp0",
			"cutCommand": "FULL",
			"qrMode": "TASK_COMPLETION",
			"baseURL": "http://localhost:8080/",
			"layout": {
				"header": { "justify": "RIGHT" },
				"title": { "justify": "LEFT" },
				"description": { "justify": "CENTER" },
				"tag": { "justify": "CENTER" },
				"qr": { "justify": "RIGHT" },
				"shinyHint": { "justify": "CENTER" },
				"footer": { "justify": "LEFT" }
			}
		}
	}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Printer.Transport != "usb" {
		t.Fatalf("printer transport mismatch: got %q, want usb", cfg.Printer.Transport)
	}
	if cfg.Printer.CutCommand != "full" {
		t.Fatalf("cut command mismatch: got %q, want full", cfg.Printer.CutCommand)
	}
	if cfg.Printer.QRMode != "task_completion" {
		t.Fatalf("qr mode mismatch: got %q, want task_completion", cfg.Printer.QRMode)
	}
	if cfg.Printer.Layout.Header.Justify != "right" {
		t.Fatalf("header justify mismatch: got %q, want right", cfg.Printer.Layout.Header.Justify)
	}
	if cfg.Printer.Layout.Title.Justify != "left" {
		t.Fatalf("title justify mismatch: got %q, want left", cfg.Printer.Layout.Title.Justify)
	}

	if cfg.Database.Path != "./Database/Taskemon.DB" {
		t.Fatalf("database path should preserve case: got %q", cfg.Database.Path)
	}
	if cfg.Printer.DevicePath != "/dev/USB/lp0" {
		t.Fatalf("device path should preserve case: got %q", cfg.Printer.DevicePath)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected invalid json error")
	}
}
