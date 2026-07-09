package printer

import (
	"bytes"
	"strings"
	"testing"
	"time"

	escposlib "github.com/Renan-M-Fernandes/taskemon/internal/escpos"
)

func TestSplitIntoLinesWrapsByWordAndMaxLines(t *testing.T) {
	got := SplitIntoLines("milk paper towel coffee beans", 12, 2)
	want := []string{"milk paper", "towel coffee"}

	if len(got) != len(want) {
		t.Fatalf("line count mismatch: got %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d mismatch: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFormatTag(t *testing.T) {
	got := formatTag("work", 17)
	want := "------#work------"

	if got != want {
		t.Fatalf("tag mismatch: got %q, want %q", got, want)
	}
}

func TestQRCorrectionLevel(t *testing.T) {
	tests := []struct {
		input string
		want  uint8
	}{
		{input: "L", want: escposlib.QRCodeErrorCorrectionLevelL},
		{input: "m", want: escposlib.QRCodeErrorCorrectionLevelM},
		{input: "Q", want: escposlib.QRCodeErrorCorrectionLevelQ},
		{input: "H", want: escposlib.QRCodeErrorCorrectionLevelH},
		{input: "invalid", want: escposlib.QRCodeErrorCorrectionLevelH},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := qrCorrectionLevel(tt.input)
			if got != tt.want {
				t.Fatalf("correction mismatch: got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestRenderTicketWritesNativeTicketContent(t *testing.T) {
	var buf bytes.Buffer
	p := escposlib.New(&buf)

	layout := Layout{
		CharsPerLine: 48,
		Header: Header{
			Enabled: true,
			Bold:    true,
			Size:    1,
			Justify: "center",
			Style:   3,
		},
		Title: Title{
			Bold:                true,
			Size:                2,
			Justify:             "left",
			LinesAllwaysVisible: 2,
			MaxLines:            2,
		},
		Description: Description{
			Enabled:             true,
			Size:                1,
			Justify:             "left",
			LinesAllwaysVisible: 5,
			MaxLines:            5,
		},
		Tag: Tag{
			Enabled: true,
			Bold:    true,
			Size:    1,
			Justify: "center",
		},
		QR: QR{
			Enabled:    true,
			Size:       4,
			Justify:    "right",
			Correction: "H",
		},
		ShinyHint: ShinyHint{Enabled: false},
		Footer:    Footer{Enabled: false},
	}

	err := RenderTicket(p, Ticket{
		TaskID:      7,
		UserID:      "ash",
		Title:       "Buy groceries today",
		Description: "Milk paper towels coffee beans and rice",
		Tag:         "home",
		CreatedAt:   time.Date(2026, 7, 8, 10, 0, 0, 0, time.UTC),
		QRValue:     "taskemon://tasks/ash/7/complete",
	}, layout)
	if err != nil {
		t.Fatalf("RenderTicket returned error: %v", err)
	}
	if err := p.Print(); err != nil {
		t.Fatalf("Print returned error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"ASH - 007", "2026-07-08", "Buy groceries", "today", "Milk paper towels", "#home", "taskemon://tasks/ash/7/complete"} {
		if !strings.Contains(out, want) {
			t.Fatalf("rendered output does not contain %q: %q", want, out)
		}
	}
}
