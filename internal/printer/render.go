package printer

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	escposlib "github.com/Renan-M-Fernandes/taskemon/internal/escpos"
)

func RenderTicket(p *escposlib.Escpos, ticket Ticket, la Layout) error {
	_, err := p.Initialize()
	if err != nil {
		return fmt.Errorf("render ticket: Initialize %w", err)
	}

	if la.Header.Enabled {
		err = printHeader(ticket.UserID, ticket.DueAt, ticket.CreatedAt, ticket.TaskID, la.CharsPerLine, la.Header, p)
		if err != nil {
			return err
		}
	}

	err = printTitle(ticket.Title, la.CharsPerLine, la.Title, p)
	if err != nil {
		return err
	}

	if la.Description.Enabled {
		err = printDescription(ticket.Description, la.CharsPerLine, la.Description, p)
		if err != nil {
			return err
		}
	}

	if la.QR.Enabled {
		err = printQR(ticket.QRValue, la.QR, p)
		if err != nil {
			return err
		}
	}

	if la.ShinyHint.Enabled && ticket.Shiny {
		if IsShiny() {
			err = printShinyOdd(la.ShinyHint, p)
			if err != nil {
				return err
			}
		}
	}

	if la.Tag.Enabled {
		err = printTag(ticket.Tag, la.CharsPerLine, la.Tag, p)
		if err != nil {
			return err
		}
	}

	if la.Footer.Enabled {
		err = printFooter(la.Footer, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func printHeader(user string, dueAt *time.Time, dateTime time.Time, taskID, charsPerLine int, cfg Header, p *escposlib.Escpos) error {
	user = strings.ToUpper(user)
	date := dateTime.Format("2006-01-02")
	dueDate, header := "", ""
	if dueAt != nil {
		dueDate = dueAt.Format("2006-01-02")
	}
	ID := fmt.Sprintf("%03d", taskID)
	userNew := user + " - " + ID

	switch cfg.Style {
	case 1:
		date = dueDate
		header = alignLeftRight(userNew, date, (charsPerLine / cfg.Size))
	case 2:
		date = date + " - " + dueDate
		header = alignLeftRight(userNew, date, (charsPerLine / cfg.Size))
	case 3:
		header = alignLeftRight(userNew, date, (charsPerLine / cfg.Size))
	case 4:
		header = user + " / " + dueDate
	case 5:
		header = user + " / " + date
	default:
		date = dueDate
		header = alignLeftRight(userNew, date, (charsPerLine / cfg.Size))
	}

	p.Bold(cfg.Bold)
	p.Size(uint8(cfg.Size), uint8(cfg.Size))
	switch cfg.Justify {
	case "left":
		p.Justify(escposlib.JustifyLeft)
	case "right":
		p.Justify(escposlib.JustifyRight)
	case "center":
		p.Justify(escposlib.JustifyCenter)
	default:
		p.Justify(escposlib.JustifyCenter)
	}

	_, err := p.WriteLine(header)
	if err != nil {
		return fmt.Errorf("print header: write line %w", err)
	}
	return nil
}

func printTitle(title string, charsPerLine int, cfg Title, p *escposlib.Escpos) error {
	titleLines := SplitIntoLines(title, (charsPerLine / cfg.Size), cfg.MaxLines)

	p.Bold(cfg.Bold)
	p.Size(uint8(cfg.Size), uint8(cfg.Size))
	switch cfg.Justify {
	case "left":
		p.Justify(escposlib.JustifyLeft)
	case "right":
		p.Justify(escposlib.JustifyRight)
	case "center":
		p.Justify(escposlib.JustifyCenter)
	default:
		p.Justify(escposlib.JustifyCenter)
	}

	err := printAllwaysVisible(titleLines, cfg.LinesAllwaysVisible, cfg.MaxLines, p)
	if err != nil {
		return fmt.Errorf("print title: print allways visible %w", err)
	}

	_, err = p.LineFeed()
	if err != nil {
		return fmt.Errorf("print title: line feed %w", err)
	}
	return nil
}

func printDescription(description string, charsPerLine int, cfg Description, p *escposlib.Escpos) error {
	descriptionLines := SplitIntoLines(description, (charsPerLine / cfg.Size), cfg.MaxLines)

	p.Bold(cfg.Bold)
	p.Size(uint8(cfg.Size), uint8(cfg.Size))
	switch cfg.Justify {
	case "left":
		p.Justify(escposlib.JustifyLeft)
	case "right":
		p.Justify(escposlib.JustifyRight)
	case "center":
		p.Justify(escposlib.JustifyCenter)
	default:
		p.Justify(escposlib.JustifyLeft)
	}

	err := printAllwaysVisible(descriptionLines, cfg.LinesAllwaysVisible, cfg.MaxLines, p)
	if err != nil {
		return fmt.Errorf("print description: print allways visible %w", err)
	}
	return nil
}

func printAllwaysVisible(lines []string, allwaysVisible, maxLines int, p *escposlib.Escpos) error {
	if allwaysVisible == 0 {
		for i, line := range lines {
			if i < maxLines || maxLines == 0 {
				_, err := p.WriteLine(line)
				if err != nil {
					return err
				}
			}
		}
	} else {
		for i := 0; i < allwaysVisible; i++ {
			line := ""
			if i < len(lines) {
				line = lines[i]
			}

			if i < maxLines || maxLines == 0 {
				_, err := p.WriteLine(line)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func printQR(qrValue string, cfg QR, p *escposlib.Escpos) error {

	p.Size(uint8(cfg.Size), uint8(cfg.Size))
	switch cfg.Justify {
	case "left":
		p.Justify(escposlib.JustifyLeft)
	case "right":
		p.Justify(escposlib.JustifyRight)
	case "center":
		p.Justify(escposlib.JustifyCenter)
	default:
		p.Justify(escposlib.JustifyRight)
	}

	if qrValue != "" {
		_, err := p.QRCode(qrValue, true, uint8(cfg.Size), qrCorrectionLevel(cfg.Correction))
		if err != nil {
			return fmt.Errorf("print qr: qr Code %w", err)
		}
	}
	return nil
}

func qrCorrectionLevel(correction string) uint8 {
	switch strings.ToUpper(strings.TrimSpace(correction)) {
	case "L":
		return escposlib.QRCodeErrorCorrectionLevelL
	case "M":
		return escposlib.QRCodeErrorCorrectionLevelM
	case "Q":
		return escposlib.QRCodeErrorCorrectionLevelQ
	case "H":
		return escposlib.QRCodeErrorCorrectionLevelH
	default:
		return escposlib.QRCodeErrorCorrectionLevelH
	}
}

func printShinyOdd(cfg ShinyHint, p *escposlib.Escpos) error {
	p.Bold(cfg.Bold)
	p.Size(uint8(cfg.Size), uint8(cfg.Size))
	switch cfg.Justify {
	case "left":
		p.Justify(escposlib.JustifyLeft)
	case "right":
		p.Justify(escposlib.JustifyRight)
	case "center":
		p.Justify(escposlib.JustifyCenter)
	default:
		p.Justify(escposlib.JustifyCenter)
	}

	_, err := p.WriteLine("* This card was generated with a Shiny Charm *")
	if err != nil {
		return fmt.Errorf("print shiny odd: write line %w", err)
	}
	return nil
}

func printTag(t string, charsPerLine int, cfg Tag, p *escposlib.Escpos) error {
	tag := formatTag(t, (charsPerLine / cfg.Size))

	p.Bold(cfg.Bold)
	p.Size(uint8(cfg.Size), uint8(cfg.Size))
	switch cfg.Justify {
	case "left":
		p.Justify(escposlib.JustifyLeft)
	case "right":
		p.Justify(escposlib.JustifyRight)
	case "center":
		p.Justify(escposlib.JustifyCenter)
	default:
		p.Justify(escposlib.JustifyCenter)
	}

	_, err := p.WriteLine(tag)
	if err != nil {
		return fmt.Errorf("print tag: write line %w", err)
	}
	return nil
}

func printFooter(cfg Footer, p *escposlib.Escpos) error {
	footer := footerMessage()

	p.Bold(cfg.Bold)
	p.Size(uint8(cfg.Size), uint8(cfg.Size))
	switch cfg.Justify {
	case "left":
		p.Justify(escposlib.JustifyLeft)
	case "right":
		p.Justify(escposlib.JustifyRight)
	case "center":
		p.Justify(escposlib.JustifyCenter)
	default:
		p.Justify(escposlib.JustifyCenter)
	}

	_, err := p.WriteLine(footer)
	if err != nil {
		return fmt.Errorf("print footer: write line %w", err)
	}
	return nil
}

func alignLeftRight(left, right string, charsPerLine int) string {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	if charsPerLine <= 0 {
		charsPerLine = 48
	}

	spaceCount := charsPerLine - len(left) - len(right)

	if spaceCount <= 1 {
		return left + " " + right
	}

	return left + strings.Repeat(" ", spaceCount) + right
}

func formatTag(tag string, charsPerLine int) string {
	if tag != "" {
		tag = "#" + tag
		tag = strings.TrimSpace(tag)
	}

	spaceCount := (charsPerLine - len(tag)) / 2

	if spaceCount <= 1 {
		return tag
	}

	return strings.Repeat("-", spaceCount) + tag + strings.Repeat("-", spaceCount)
}

func SplitIntoLines(text string, charsPerLine, maxLines int) []string {
	if charsPerLine <= 0 {
		return nil
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	var current strings.Builder

	appendCurrent := func() bool {
		if current.Len() == 0 {
			return false
		}

		if maxLines > 0 && len(lines) >= maxLines {
			return true
		}

		lines = append(lines, current.String())
		current.Reset()
		return maxLines > 0 && len(lines) >= maxLines
	}

	for _, word := range words {
		if current.Len() == 0 {
			current.WriteString(word)
			continue
		}

		if current.Len()+1+len(word) <= charsPerLine {
			current.WriteByte(' ')
			current.WriteString(word)
			continue
		}

		if appendCurrent() {
			return lines
		}
		current.WriteString(word)
	}

	appendCurrent()

	return lines
}

func footerMessage() string {
	footers := [...]string{
		"Waiting for a trainer.",
		"A wild task appeared!",
		"Ready when you are.",
		"Don't leave me in the PC.",
		"Still needs a hero.",
		"Procrastination used Protect.",
		"The next badge starts here.",
		"Adventure awaits.",
		"One quest at a time.",
		"Time to grind.",
		"Your future self is watching.",
		"Roll for initiative.",
		"Main quest available.",
		"Side quest accepted?",
		"Not very effective... yet.",
		"This won't finish itself.",
		"The Pokédex believes in you.",
		"Your party is waiting.",
		"A Poké Ball won't solve this one.",
		"Procrastination used Protect.",
		"Loading motivation...",
		"Procrastination used Protect.",
		"The timer is ticking.",
		"Procrastination used Protect.",
		"One more step to level up.",
		"Challenge accepted?",
		"A rare reward may be inside.",
		"Someone has to do it.",
		"The backlog grew restless.",
		"Choose your next move.",
		"(^. .^)/ A wild cat appeared",
	}

	easterEggs := [...]string{
		"Ah shit, here we go again.",
		"All we had to do was follow the damn train, CJ!",
		"The cake is a lie.",
		"Do a barrel roll!",
		"LEEROOOOY JENKINS!",
		"War. War never changes.",
		"Hey, you. You're finally awake.",
		"Would you kindly...",
		"Praise the Sun!",
		"You died... just kidding.",
		"Nothing is true. Everything is permitted.",
		"Snake? Snake?! SNAAAAAKE!",
		"Finish him!",
		"Hadouken!",
		"Press F to pay respects.",
		"Wake up, Samurai. We've got tasks to finish.",
		"The princess is in another castle.",
		"A NEW HAND TOUCHES THE BEACON!",
		"You must construct additional Pylons.",
		"Job's done!",
		"Zug zug.",
		"The cake wasn't a lie... this time.",
		"Would you intercept me? I'd intercept me.",
		"Mission passed! Respect +",
		"GG. Next quest!",
	}

	if rand.IntN(100) < 20 {
		return easterEggs[rand.IntN(len(easterEggs))]
	}
	return footers[rand.IntN(len(footers))]
}

func IsShiny() bool {
	return rand.IntN(10) == 0
}
