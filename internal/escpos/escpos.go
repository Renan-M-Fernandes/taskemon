// SPDX-FileCopyrightText: 2025 James Pond <james@cipher.host>
// SPDX-FileCopyrightText: 2026 Renan Fernandes
//
// SPDX-License-Identifier: EUPL-1.2

// Package escpos is a small vendored ESC/POS implementation based on
// github.com/hennedo/escpos, patched for generic POS-80 thermal printers.
package escpos

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"math"
)

type Style struct {
	Bold       bool
	Width      uint8
	Height     uint8
	Reverse    bool
	Underline  uint8
	UpsideDown bool
	Rotate     bool
	Justify    uint8
}

const (
	JustifyLeft   uint8 = 0
	JustifyCenter uint8 = 1
	JustifyRight  uint8 = 2

	QRCodeErrorCorrectionLevelL uint8 = 48
	QRCodeErrorCorrectionLevelM uint8 = 49
	QRCodeErrorCorrectionLevelQ uint8 = 50
	QRCodeErrorCorrectionLevelH uint8 = 51

	esc byte = 0x1B
	gs  byte = 0x1D
	fs  byte = 0x1C
)

type Escpos struct {
	dst   *bufio.Writer
	Style Style
	err   error
}

func New(dst io.Writer) *Escpos {
	return &Escpos{
		dst: bufio.NewWriter(dst),
		Style: Style{
			Width:  1,
			Height: 1,
		},
	}
}

func (e *Escpos) Print() error {
	if e.err != nil {
		return e.err
	}

	return e.dst.Flush()
}

func (e *Escpos) remember(err error) {
	if err != nil && e.err == nil {
		e.err = err
	}
}

func (e *Escpos) PrintAndCut(linesBeforeCut int, cutCommand string) error {
	if linesBeforeCut != 0 {
		if _, err := e.LineFeedD(uint8(linesBeforeCut)); err != nil {
			return fmt.Errorf("feed before cut: %w", err)
		}
	}

	if _, err := e.CutCommand(cutCommand); err != nil {
		return fmt.Errorf("cut: %w", err)
	}

	return e.Print()
}

func (e *Escpos) WriteRaw(data []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}

	if len(data) == 0 {
		return 0, nil
	}

	n, err := e.dst.Write(data)
	if err != nil {
		e.remember(err)
	}

	return n, err
}

func (e *Escpos) Write(data string) (int, error) {
	return e.WriteRaw([]byte(data))
}

func (e *Escpos) WriteLine(data string) (int, error) {
	return e.WriteRaw([]byte(data + "\n"))
}

func (e *Escpos) WriteGBK(data string) (int, error) {
	return 0, fmt.Errorf("WriteGBK is not available in this vendored build")
}

// WriteWEU writes a string using a small CP850 conversion table for common
// Western European / Portuguese characters. Unsupported non-ASCII runes are
// replaced with '?'.
func (e *Escpos) WriteWEU(data string) (int, error) {
	return e.Write(encodeCP850(data))
}

func (e *Escpos) Bold(p bool) *Escpos {
	e.Style.Bold = p
	_, err := e.WriteRaw([]byte{esc, 'E', boolToByte(p)})
	e.remember(err)
	return e
}

func (e *Escpos) Reverse(p bool) *Escpos {
	e.Style.Reverse = p
	_, err := e.WriteRaw([]byte{gs, 'B', boolToByte(p)})
	e.remember(err)
	return e
}

func (e *Escpos) Rotate(p bool) *Escpos {
	e.Style.Rotate = p
	_, err := e.WriteRaw([]byte{esc, 'V', boolToByte(p)})
	e.remember(err)
	return e
}

func (e *Escpos) UpsideDown(p bool) *Escpos {
	e.Style.UpsideDown = p
	_, err := e.WriteRaw([]byte{esc, '{', boolToByte(p)})
	e.remember(err)
	return e
}

func (e *Escpos) Justify(p uint8) *Escpos {
	p = clampJustify(p)
	e.Style.Justify = p
	_, err := e.WriteRaw([]byte{esc, 'a', p})
	e.remember(err)
	return e
}

func (e *Escpos) Underline(p uint8) *Escpos {
	p = clampUnderline(p)
	e.Style.Underline = p
	_, err := e.WriteRaw([]byte{esc, '-', p})
	e.remember(err)
	return e
}

func (e *Escpos) Size(width, height uint8) *Escpos {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if width > 8 {
		width = 8
	}
	if height > 8 {
		height = 8
	}

	e.Style.Width = width
	e.Style.Height = height

	_, err := e.WriteRaw([]byte{gs, '!', sizeByte(width, height)})
	e.remember(err)

	return e
}

func (e *Escpos) HRIPosition(p uint8) (int, error) {
	if p > 3 {
		p = 0
	}
	return e.WriteRaw([]byte{gs, 'H', p})
}

func (e *Escpos) HRIFont(p bool) (int, error)        { return e.WriteRaw([]byte{gs, 'f', boolToByte(p)}) }
func (e *Escpos) BarcodeHeight(p uint8) (int, error) { return e.WriteRaw([]byte{gs, 'h', p}) }

func (e *Escpos) BarcodeWidth(p uint8) (int, error) {
	if p < 2 {
		p = 2
	}
	if p > 6 {
		p = 6
	}
	return e.WriteRaw([]byte{gs, 'w', p})
}

func (e *Escpos) UPCA(code string) (int, error)  { return e.barcode(0, code, 11, 12) }
func (e *Escpos) UPCE(code string) (int, error)  { return e.barcode(1, code, 11, 12) }
func (e *Escpos) EAN13(code string) (int, error) { return e.barcode(2, code, 12, 13) }
func (e *Escpos) EAN8(code string) (int, error)  { return e.barcode(3, code, 7, 8) }

func (e *Escpos) barcode(kind byte, code string, minLen, maxLen int) (int, error) {
	if len(code) < minLen || len(code) > maxLen {
		return 0, fmt.Errorf("code length must be between %d and %d", minLen, maxLen)
	}
	if !onlyDigits(code) {
		return 0, fmt.Errorf("code can only contain numerical characters")
	}
	byteCode := append([]byte(code), 0)
	return e.WriteRaw(append([]byte{gs, 'k', kind}, byteCode...))
}

func (e *Escpos) QRCode(code string, model bool, size uint8, correctionLevel uint8) (int, error) {
	if len(code) > 7089 {
		return 0, fmt.Errorf("qr code data too long")
	}
	if size < 1 {
		size = 1
	}
	if size > 16 {
		size = 16
	}
	if correctionLevel < QRCodeErrorCorrectionLevelL || correctionLevel > QRCodeErrorCorrectionLevelH {
		correctionLevel = QRCodeErrorCorrectionLevelM
	}

	m := byte(49)
	if model {
		m = 50
	}
	if _, err := e.WriteRaw([]byte{gs, '(', 'k', 4, 0, 49, 65, m, 0}); err != nil {
		return 0, err
	}
	if _, err := e.WriteRaw([]byte{gs, '(', 'k', 3, 0, 49, 67, size}); err != nil {
		return 0, err
	}
	if _, err := e.WriteRaw([]byte{gs, '(', 'k', 3, 0, 49, 69, correctionLevel}); err != nil {
		return 0, err
	}

	codeLength := len(code) + 3
	pH := byte(int(math.Floor(float64(codeLength) / 256)))
	pL := byte(codeLength - 256*int(pH))

	written, err := e.WriteRaw(append([]byte{gs, '(', 'k', pL, pH, 49, 80, 48}, []byte(code)...))
	if err != nil {
		return written, err
	}
	_, err = e.WriteRaw([]byte{gs, '(', 'k', 3, 0, 49, 81, 48})
	return written, err
}

func (e *Escpos) PrintImage(img image.Image) (int, error) {
	return e.PrintImageBitImage(img)
}

func (e *Escpos) PrintNVBitImage(p uint8, mode uint8) (int, error) {
	if p == 0 {
		return 0, fmt.Errorf("start index of nv bit images starts at 1")
	}
	if mode > 3 {
		return 0, fmt.Errorf("mode only supports values from 0 to 3")
	}
	return e.WriteRaw([]byte{fs, 'd', p, mode})
}

func (e *Escpos) LineFeed() (int, error)           { return e.WriteRaw([]byte{'\n'}) }
func (e *Escpos) LineFeedD(p uint8) (int, error)   { return e.WriteRaw([]byte{esc, 'd', p}) }
func (e *Escpos) DefaultLineSpacing() (int, error) { return e.WriteRaw([]byte{esc, '2'}) }
func (e *Escpos) LineSpacing(p uint8) (int, error) { return e.WriteRaw([]byte{esc, '3', p}) }

func (e *Escpos) MotionUnits(x, y uint8) (int, error) { return e.WriteRaw([]byte{gs, 'P', x, y}) }

func (e *Escpos) Cut() (int, error)                { return e.WriteRaw([]byte{gs, 'V', 0x00}) }
func (e *Escpos) PartialCut() (int, error)         { return e.WriteRaw([]byte{gs, 'V', 1}) }
func (e *Escpos) CutFeed(lines uint8) (int, error) { return e.WriteRaw([]byte{gs, 'V', 65, lines}) }

func (e *Escpos) CutCommand(command string) (int, error) {
	switch command {
	case "", "gs_v_0", "full":
		return e.Cut()
	case "partial":
		return e.PartialCut()
	case "feed":
		return e.CutFeed(0)
	default:
		return 0, fmt.Errorf("unsupported cut command %q", command)
	}
}

func (e *Escpos) Initialize() (int, error) {
	e.Style = Style{
		Width:   1,
		Height:  1,
		Justify: JustifyLeft,
	}

	return e.WriteRaw([]byte{esc, '@'})
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func sizeByte(width, height uint8) byte {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if width > 8 {
		width = 8
	}
	if height > 8 {
		height = 8
	}
	return ((width - 1) << 4) | (height - 1)
}

func clampJustify(p uint8) uint8 {
	if p > 2 {
		return 0
	}
	return p
}

func clampUnderline(p uint8) uint8 {
	if p > 2 {
		return 0
	}
	return p
}

func onlyDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func encodeCP850(s string) string {
	out := make([]byte, 0, len(s))
	for _, r := range s {
		if r >= 0x20 && r <= 0x7e || r == '\n' || r == '\r' || r == '\t' {
			out = append(out, byte(r))
			continue
		}
		if b, ok := cp850[r]; ok {
			out = append(out, b)
		} else {
			out = append(out, '?')
		}
	}
	return string(out)
}

var cp850 = map[rune]byte{
	'Ç': 0x80, 'ü': 0x81, 'é': 0x82, 'â': 0x83, 'ä': 0x84, 'à': 0x85,
	'å': 0x86, 'ç': 0x87, 'ê': 0x88, 'ë': 0x89, 'è': 0x8a, 'ï': 0x8b,
	'î': 0x8c, 'ì': 0x8d, 'Ä': 0x8e, 'Å': 0x8f, 'É': 0x90, 'æ': 0x91,
	'Æ': 0x92, 'ô': 0x93, 'ö': 0x94, 'ò': 0x95, 'û': 0x96, 'ù': 0x97,
	'ÿ': 0x98, 'Ö': 0x99, 'Ü': 0x9a, 'ø': 0x9b, '£': 0x9c, 'Ø': 0x9d,
	'á': 0xa0, 'í': 0xa1, 'ó': 0xa2, 'ú': 0xa3, 'ñ': 0xa4, 'Ñ': 0xa5,
	'ª': 0xa6, 'º': 0xa7, '¿': 0xa8, '®': 0xa9, '¬': 0xaa, '½': 0xab,
	'¼': 0xac, '¡': 0xad, '«': 0xae, '»': 0xaf, 'Á': 0xb5, 'Â': 0xb6,
	'À': 0xb7, '©': 0xb8, 'ã': 0xc6, 'Ã': 0xc7, '¤': 0xcf, 'ð': 0xd0,
	'Ð': 0xd1, 'Ê': 0xd2, 'Ë': 0xd3, 'È': 0xd4, 'Í': 0xd6, 'Î': 0xd7,
	'Ï': 0xd8, 'Ó': 0xe0, 'ß': 0xe1, 'Ô': 0xe2, 'Ò': 0xe3, 'õ': 0xe4,
	'Õ': 0xe5, 'µ': 0xe6, 'þ': 0xe7, 'Þ': 0xe8, 'Ú': 0xe9, 'Û': 0xea,
	'Ù': 0xeb, 'ý': 0xec, 'Ý': 0xed, '¯': 0xee, '´': 0xef, '±': 0xf1,
	'¾': 0xf3, '¶': 0xf4, '§': 0xf5, '÷': 0xf6, '¸': 0xf7, '°': 0xf8,
	'¨': 0xf9, '·': 0xfa, '¹': 0xfb, '³': 0xfc, '²': 0xfd,
}
