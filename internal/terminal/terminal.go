package terminal

import (
	"fmt"
	"strings"
	"unicode/utf8"

	catppuccin "github.com/catppuccin/go"
	"github.com/tcondeixa/resume/internal/resume"
)

const (
	escChar = "\x1b"
	csi     = escChar + "["
	FG      = 38
	BG      = 48

	invertColors = csi + "7m"
	resetColors  = csi + "m"
	bold         = csi + "1m"
	unbold       = csi + "22m"
	blink        = csi + "5m"
	unblink      = csi + "25m"
	newLine      = "\r\n"

	outerSpace       = "   "
	innerSpace       = "    "
	leftFrame        = outerSpace + "┃"
	leftMargin       = leftFrame + innerSpace
	rightFrame       = "┃" + outerSpace
	rightMargin      = innerSpace + rightFrame
	leftUpperCorner  = outerSpace + "┏"
	rightUpperCorner = "┓" + outerSpace
	leftLowerCorner  = outerSpace + "┗"
	rightLowerCorner = "┛" + outerSpace
)

var nonTextSpace = utf8.RuneCountInString(leftMargin + rightMargin)

type Terminal struct {
	Width  int
	Height int
	Term   string
	Theme  catppuccin.Flavor
}

func (t *Terminal) Render(r *resume.Resume) (string, error) {
	builder := strings.Builder{}
	builder.Grow(200)

	emptyLine := t.emptyLine()
	separatorLine := t.separatorLine()

	err := t.firstLine(&builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine)
	if err != nil {
		return "", err
	}

	err = t.Header(r, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + separatorLine)
	if err != nil {
		return "", err
	}

	err = t.Skills(r, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + separatorLine)
	if err != nil {
		return "", err
	}

	err = t.Experience(r, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + separatorLine)
	if err != nil {
		return "", err
	}

	err = t.Education(r, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + separatorLine)
	if err != nil {
		return "", err
	}

	err = t.Certifications(r, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + emptyLine)
	if err != nil {
		return "", err
	}

	err = t.lastLine(&builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(newLine)
	if err != nil {
		return "", err
	}

	err = t.Footer(&builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(newLine)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (t *Terminal) Footer(builder *strings.Builder) error {
	text := "  by tcondeixa: https://github.com/tcondeixa/resume" + rightMargin
	pad := t.paddingStr(text)
	text = Blink("  by tcondeixa: ") + FgColor("https://github.com/tcondeixa/resume", t.Theme.Blue())
	_, err := builder.WriteString(pad + text + strings.Repeat(" ", utf8.RuneCountInString(rightMargin)) + newLine)
	if err != nil {
		return err
	}

	return nil
}

func (t *Terminal) Header(r *resume.Resume, builder *strings.Builder) error {
	emptyLine := t.emptyLine()
	pad := t.paddingStr(leftMargin + r.Header.Name + rightMargin)
	_, err := builder.WriteString(
		leftMargin + FgColor(Blink(strings.ToUpper(r.Header.Name)), t.Theme.Green()) + pad + rightMargin + newLine,
	)
	if err != nil {
		return err
	}

	_, err = builder.WriteString(emptyLine)
	if err != nil {
		return err
	}

	for _, line := range splitTextLines(r.Header.Summary, t.Width-nonTextSpace) {
		pad := t.paddingStr(leftMargin + line + rightMargin)
		_, err = builder.WriteString(leftMargin + line + pad + rightMargin + newLine)
		if err != nil {
			return err
		}
	}

	_, err = builder.WriteString(emptyLine)
	if err != nil {
		return err
	}

	location := " " + r.Header.Location
	pad = t.paddingStr(leftMargin + location + rightMargin)
	location = FgColor(Bold(" "), t.Theme.Red()) + r.Header.Location
	_, err = builder.WriteString(leftMargin + location + pad + rightMargin + newLine)
	if err != nil {
		return err
	}

	for _, link := range r.Header.Links {
		pad := t.paddingStr(leftMargin + link.Icon + " " + link.URL + rightMargin)
		_, err = builder.WriteString(
			leftMargin + link.Icon + " " + FgColor(link.URL, t.Theme.Blue()) + pad + rightMargin + newLine,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Terminal) Skills(r *resume.Resume, builder *strings.Builder) error {
	emptyLine := t.emptyLine()
	err := t.Title("Skills", builder)
	if err != nil {
		return err
	}

	_, err = builder.WriteString(emptyLine)
	if err != nil {
		return err
	}

	for _, skill := range r.Skills {
		aSkills := strings.Join(skill.Skills, " | ")

		text := skill.Area + ": " + aSkills
		for _, line := range splitTextLines(text, t.Width-nonTextSpace) {
			pad := t.paddingStr(leftMargin + line + rightMargin)
			line = strings.ReplaceAll(line, skill.Area, Bold(skill.Area))
			_, err := builder.WriteString(leftMargin + line + pad + rightMargin + newLine)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Terminal) Experience(r *resume.Resume, builder *strings.Builder) error {
	emptyLine := t.emptyLine()
	err := t.Title("Experience", builder)
	if err != nil {
		return err
	}

	previousCompany := ""
	for _, exp := range r.Experience {
		_, err = builder.WriteString(emptyLine)
		if err != nil {
			return err
		}

		if exp.Company != previousCompany {
			pad := t.paddingStr(leftMargin + exp.Company + rightMargin)
			_, err := builder.WriteString(
				leftMargin + FgColor(Bold(exp.Company), t.Theme.Lavender()) + pad + rightMargin + newLine,
			)
			if err != nil {
				return err
			}
			previousCompany = exp.Company
		}

		text := exp.Title + "      " + exp.StartDate + " - " + exp.EndDate
		pad := t.paddingStr(leftMargin + text + rightMargin)
		if pad != "" {
			_, err = builder.WriteString(
				leftMargin + Bold(exp.Title) + pad + Bold("      "+exp.StartDate+" - "+exp.EndDate) + rightMargin + newLine,
			)
			if err != nil {
				return err
			}
		} else {
			pad = t.paddingStr(leftMargin + exp.Title + rightMargin)
			_, err = builder.WriteString(
				leftMargin + Bold(exp.Title) + pad + rightMargin + newLine,
			)
			if err != nil {
				return err
			}

			text := exp.StartDate + " - " + exp.EndDate
			pad := t.paddingStr(leftMargin + text + rightMargin)
			_, err = builder.WriteString(
				leftMargin + Bold(text) + pad + rightMargin + newLine,
			)
			if err != nil {
				return err
			}
		}

		_, err = builder.WriteString(emptyLine)
		if err != nil {
			return err
		}

		for _, line := range splitTextLines(exp.Summary, t.Width-nonTextSpace) {
			pad := t.paddingStr(leftMargin + line + rightMargin)
			_, err = builder.WriteString(leftMargin + line + pad + rightMargin + newLine)
			if err != nil {
				return err
			}
		}

		_, err = builder.WriteString(emptyLine)
		if err != nil {
			return err
		}

		maxTextLineLen := t.Width - utf8.RuneCountInString(leftMargin+rightMargin)
		for _, highlight := range exp.Highlights {
			highlight = "■ " + highlight
			for _, line := range splitTextLines(highlight, maxTextLineLen) {
				pad := t.paddingStr(leftMargin + line + rightMargin)
				line = strings.ReplaceAll(line, "■", FgColor("■", t.Theme.Lavender()))
				_, err = builder.WriteString(leftMargin + line + pad + rightMargin + newLine)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t *Terminal) Education(r *resume.Resume, builder *strings.Builder) error {
	emptyLine := t.emptyLine()
	err := t.Title("Education", builder)
	if err != nil {
		return err
	}

	for _, edu := range r.Education {
		_, err = builder.WriteString(emptyLine)
		if err != nil {
			return err
		}

		pad := t.paddingStr(leftMargin + edu.Institution + rightMargin)
		_, err := builder.WriteString(leftMargin + Bold(edu.Institution) + pad + rightMargin + newLine)
		if err != nil {
			return err
		}

		_, err = builder.WriteString(emptyLine)
		if err != nil {
			return err
		}

		for _, achievement := range edu.Achievements {
			for _, line := range splitTextLines(achievement, t.Width-nonTextSpace) {
				pad := t.paddingStr(leftMargin + line + rightMargin)
				_, err := builder.WriteString(leftMargin + line + pad + rightMargin + newLine)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t *Terminal) Certifications(r *resume.Resume, builder *strings.Builder) error {
	emptyLine := t.emptyLine()
	err := t.Title("Certifications", builder)
	if err != nil {
		return err
	}

	_, err = builder.WriteString(emptyLine)
	if err != nil {
		return err
	}

	for _, cert := range r.Certifications {
		text := fmt.Sprintf("%s %s %s (%s)", cert.Link.Icon, cert.Name, cert.Authority, "Issued "+cert.Issued)
		for _, line := range splitTextLines(text, t.Width-nonTextSpace) {
			pad := t.paddingStr(leftMargin + line + rightMargin)
			line = strings.ReplaceAll(line, cert.Authority, Bold(cert.Authority))
			line = strings.ReplaceAll(line, cert.Name, Link(cert.Name, cert.Link.URL, t.Theme))
			_, err := builder.WriteString(leftMargin + line + pad + rightMargin + newLine)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Terminal) paddingStr(text string) string {
	if utf8.RuneCountInString(text) >= t.Width {
		return ""
	}

	return strings.Repeat(" ", t.Width-utf8.RuneCountInString(text))
}

func (t *Terminal) firstLine(builder *strings.Builder) error {
	_, err := builder.WriteString(
		leftUpperCorner + strings.Repeat(
			"━",
			t.Width-utf8.RuneCountInString(leftUpperCorner+rightUpperCorner),
		) + rightUpperCorner + newLine,
	)

	return err
}

func (t *Terminal) lastLine(builder *strings.Builder) error {
	_, err := builder.WriteString(
		leftLowerCorner +
			strings.Repeat(
				"━",
				t.Width-utf8.RuneCountInString(leftLowerCorner+rightLowerCorner),
			) + rightLowerCorner + newLine,
	)
	return err
}

func (t *Terminal) emptyLine() string {
	return leftFrame + strings.Repeat(" ", t.Width-utf8.RuneCountInString(leftFrame+rightFrame)) + rightFrame + newLine
}

func (t *Terminal) separatorLine() string {
	return leftMargin + strings.Repeat(
		"━",
		t.Width-utf8.RuneCountInString(leftMargin+rightMargin),
	) + rightMargin + newLine
}

func FgColor(text string, color catppuccin.Color) string {
	return Color(text, color, FG)
}

func BgColor(text string, color catppuccin.Color) string {
	return Color(text, color, BG)
}

func Color(text string, color catppuccin.Color, ground int) string {
	return fmt.Sprintf(
		"%s%d;2;%d;%d;%dm%s%s",
		csi,
		ground,
		color.RGB[0],
		color.RGB[1],
		color.RGB[2],
		text,
		resetColors,
	)
}

func Link(text string, url string, theme catppuccin.Flavor) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, FgColor(text, theme.Blue()))
}

func (t *Terminal) Title(title string, builder *strings.Builder) error {
	text := leftMargin + title + rightMargin
	_, err := builder.WriteString(
		leftMargin + FgColor(Bold(title), t.Theme.Red()) + t.paddingStr(text) + rightMargin + newLine,
	)
	return err
}

func Blink(text string) string {
	return fmt.Sprintf("%s%s%s", blink, text, unblink)
}

func Bold(text string) string {
	return fmt.Sprintf("%s%s%s", bold, text, unbold)
}

func splitTextLines(text string, maxTextSize int) []string {
	runes := []rune(text)
	output := []string{}
	i := 0

	for i < len(runes) {
		end := min(i+maxTextSize, len(runes))

		if end == len(runes) {
			output = append(output, string(runes[i:end]))
			break
		}

		// Find last space within the chunk
		originalEnd := end
		for end > i && runes[end] != ' ' {
			end--
		}

		// If no space found, use original end (force break)
		if end == i {
			end = originalEnd
		}

		output = append(output, string(runes[i:end]))

		// Skip the space if we broke on one
		if end < len(runes) && runes[end] == ' ' {
			end++
		}

		i = end
	}

	return output
}
