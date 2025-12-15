package terminal

import (
	"fmt"
	"strings"
	"unicode/utf8"

	catppuccin "github.com/catppuccin/go"
	"github.com/tcondeixa/resume/internal/resume"
)

const (
	bellChar = "\x07"
	escChar  = "\x1b"
	csi      = escChar + "["
	FG       = 38
	BG       = 48

	invertColors = csi + "7m"
	resetColors  = csi + "m"
	bold         = csi + "1m"
	unbold       = csi + "22m"
	blink        = csi + "5m"
	unblink      = csi + "25m"
	newLine      = "\r\n"

	leftFrame        = "  ┃"
	leftMargin       = leftFrame + "    "
	rightFrame       = "┃  "
	rightMargin      = "    " + rightFrame
	leftUpperCorner  = "  ┏"
	rightUpperCorner = "┓  "
	leftLowerCorner  = "  ┗"
	rightLowerCorner = "┛  "
)

type Terminal struct {
	Width  int
	Height int
	Term   string
}

func (t *Terminal) Render(r *resume.Resume) (string, error) {
	theme := catppuccin.Mocha
	builder := strings.Builder{}
	builder.Grow(100)

	emptyLine := leftFrame + strings.Repeat(" ", t.Width-len(leftFrame)-len(rightFrame)) + rightFrame + newLine
	separatorLine := leftMargin + strings.Repeat("━", t.Width-len(leftMargin)-len(rightMargin)) + rightMargin + newLine

	_, err := builder.WriteString(
		leftUpperCorner + strings.Repeat(
			"━",
			t.Width-len(leftUpperCorner)-len(rightUpperCorner),
		) + rightUpperCorner + newLine,
	)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine)
	if err != nil {
		return "", err
	}

	err = Header(r, theme, t.Width, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + separatorLine)
	if err != nil {
		return "", err
	}

	err = Skills(r, theme, t.Width, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + separatorLine)
	if err != nil {
		return "", err
	}

	err = Experience(r, theme, t.Width, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + separatorLine)
	if err != nil {
		return "", err
	}

	err = Education(r, theme, t.Width, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + separatorLine)
	if err != nil {
		return "", err
	}

	err = Certifications(r, theme, t.Width, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(emptyLine + emptyLine)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(
		leftLowerCorner +
			strings.Repeat(
				"━",
				t.Width-len(leftLowerCorner)-len(rightLowerCorner),
			) + rightLowerCorner + newLine,
	)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(newLine)
	if err != nil {
		return "", err
	}

	err = Footer(theme, t.Width, &builder)
	if err != nil {
		return "", err
	}

	_, err = builder.WriteString(newLine)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func Footer(theme catppuccin.Flavor, width int, builder *strings.Builder) error {
	text := "  by tcondeixa: https://github.com/tcondeixa/resume"
	pad := paddingStr(width, utf8.RuneCountInString(text)+6)
	text = Blink("  by tcondeixa: ") + FgColor("https://github.com/tcondeixa/resume", theme.Blue())
	_, err := builder.WriteString(pad + text + "      " + newLine)
	if err != nil {
		return err
	}

	return nil
}

func Header(r *resume.Resume, theme catppuccin.Flavor, width int, builder *strings.Builder) error {
	emptyLine := leftFrame + strings.Repeat(" ", width-len(leftFrame)-len(rightFrame)) + rightFrame + newLine
	pad := paddingStr(width, len(leftMargin)+len(rightMargin)+len(r.Header.Name))
	_, err := builder.WriteString(
		leftMargin + FgColor(Blink(strings.ToUpper(r.Header.Name)), theme.Green()) + pad + rightMargin + newLine,
	)
	if err != nil {
		return err
	}

	_, err = builder.WriteString(emptyLine)
	if err != nil {
		return err
	}

	for _, line := range splitTextLines(r.Header.Summary, width-len(leftMargin)-len(rightMargin)) {
		pad := paddingStr(width, len(leftMargin)+len(rightMargin)+len(line))
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
	pad = paddingStr(width, len(leftMargin)+len(rightMargin)+utf8.RuneCountInString(location))
	location = FgColor(Bold(" "), theme.Red()) + r.Header.Location
	_, err = builder.WriteString(leftMargin + location + pad + rightMargin + newLine)
	if err != nil {
		return err
	}

	for _, link := range r.Header.Links {
		pad := paddingStr(width, len(leftMargin)+len(rightMargin)+utf8.RuneCountInString(link.Icon+" "+link.URL))
		_, err = builder.WriteString(
			leftMargin + link.Icon + " " + FgColor(link.URL, theme.Blue()) + pad + rightMargin + newLine,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func Skills(r *resume.Resume, theme catppuccin.Flavor, width int, builder *strings.Builder) error {
	emptyLine := leftFrame + strings.Repeat(" ", width-len(leftFrame)-len(rightFrame)) + rightFrame + newLine
	_, err := builder.WriteString(leftMargin + Title(theme, width, "Skills") + rightMargin + newLine)
	if err != nil {
		return err
	}

	_, err = builder.WriteString(emptyLine)
	if err != nil {
		return err
	}

	for _, skill := range r.Skills {
		aSkills := strings.Join(skill.Skills, " | ")
		pad := paddingStr(width, len(leftMargin)+len(rightMargin)+len(skill.Area)+2+len(aSkills))
		_, err := builder.WriteString(leftMargin + Bold(skill.Area) + ": " + aSkills + pad + rightMargin + newLine)
		if err != nil {
			return err
		}
	}

	return nil
}

func Experience(r *resume.Resume, theme catppuccin.Flavor, width int, builder *strings.Builder) error {
	emptyLine := leftFrame + strings.Repeat(" ", width-len(leftFrame)-len(rightFrame)) + rightFrame + newLine
	_, err := builder.WriteString(leftMargin + Title(theme, width, "Experience") + rightMargin + newLine)
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
			pad := paddingStr(width, len(leftMargin)+len(rightMargin)+len(exp.Company))
			_, err := builder.WriteString(
				leftMargin + FgColor(Bold(exp.Company), theme.Lavender()) + pad + rightMargin + newLine,
			)
			if err != nil {
				return err
			}
			previousCompany = exp.Company
		}

		pad := paddingStr(width, len(leftMargin)+len(rightMargin)+len(exp.Title)+len(exp.StartDate)+3+len(exp.EndDate))
		_, err = builder.WriteString(
			leftMargin + Bold(exp.Title) + pad + Bold(exp.StartDate+" - "+exp.EndDate) + rightMargin + newLine,
		)
		if err != nil {
			return err
		}

		_, err = builder.WriteString(emptyLine)
		if err != nil {
			return err
		}

		for _, line := range splitTextLines(exp.Summary, width-len(leftMargin)-len(rightMargin)) {
			pad := paddingStr(width, len(leftMargin)+len(rightMargin)+len(line))
			_, err = builder.WriteString(leftMargin + line + pad + rightMargin + newLine)
			if err != nil {
				return err
			}
		}

		_, err = builder.WriteString(emptyLine)
		if err != nil {
			return err
		}

		for _, highlight := range exp.Highlights {
			highlight = "■ " + highlight
			for i, line := range splitTextLines(highlight, width-len(leftMargin)-len(rightMargin)) {
				pad := paddingStr(width, len(leftMargin)+len(rightMargin)+utf8.RuneCountInString(line))
				if i == 0 {
					line = FgColor("■", theme.Lavender()) + line[1:]
				}
				_, err = builder.WriteString(leftMargin + line + pad + rightMargin + newLine)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func Education(r *resume.Resume, theme catppuccin.Flavor, width int, builder *strings.Builder) error {
	emptyLine := leftFrame + strings.Repeat(" ", width-len(leftFrame)-len(rightFrame)) + rightFrame + newLine
	_, err := builder.WriteString(leftMargin + Title(theme, width, "Education") + rightMargin + newLine)
	if err != nil {
		return err
	}

	for _, edu := range r.Education {
		_, err = builder.WriteString(emptyLine)
		if err != nil {
			return err
		}

		pad := paddingStr(width, len(leftMargin)+len(rightMargin)+len(edu.Institution))
		_, err := builder.WriteString(leftMargin + Bold(edu.Institution) + pad + rightMargin + newLine)
		if err != nil {
			return err
		}

		_, err = builder.WriteString(emptyLine)
		if err != nil {
			return err
		}

		for _, achievement := range edu.Achievements {
			pad := paddingStr(width, len(leftMargin)+len(rightMargin)+len(achievement))
			_, err := builder.WriteString(leftMargin + achievement + pad + rightMargin + newLine)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Certifications(r *resume.Resume, theme catppuccin.Flavor, width int, builder *strings.Builder) error {
	emptyLine := leftFrame + strings.Repeat(" ", width-len(leftFrame)-len(rightFrame)) + rightFrame + newLine
	_, err := builder.WriteString(leftMargin + Title(theme, width, "Certifications") + rightMargin + newLine)
	if err != nil {
		return err
	}

	_, err = builder.WriteString(emptyLine)
	if err != nil {
		return err
	}

	for _, cert := range r.Certifications {
		text := fmt.Sprintf("%s %s %s (%s)", cert.Link.Icon, cert.Name, cert.Authority, "Issued "+cert.Issued)
		pad := paddingStr(width, len(leftMargin)+len(rightMargin)+utf8.RuneCountInString(text))
		text = fmt.Sprintf(
			"%s %s %s (%s)",
			cert.Link.Icon,
			Link(cert.Name, cert.Link.URL, theme),
			Bold(cert.Authority),
			"Issued "+cert.Issued,
		)
		_, err := builder.WriteString(leftMargin + text + pad + rightMargin + newLine)
		if err != nil {
			return err
		}
	}

	return nil
}

func paddingStr(width int, usedLen int) string {
	if usedLen >= width {
		return ""
	}

	return strings.Repeat(" ", width-usedLen)
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

func Title(theme catppuccin.Flavor, width int, title string) string {
	pad := paddingStr(width, len(leftMargin)+len(rightMargin)+len(title))
	return FgColor(Bold(title), theme.Red()) + pad
}

func Blink(text string) string {
	return fmt.Sprintf("%s%s%s", blink, text, unblink)
}

func Bold(text string) string {
	return fmt.Sprintf("%s%s%s", bold, text, unbold)
}

func splitTextLines(text string, lineLen int) []string {
	output := []string{}
	i := 0
	for i < len(text) {
		end := min(i+lineLen, len(text))
		if end == len(text) {
			output = append(output, text[i:end])
			break
		}
		for text[end] != ' ' && end < len(text) {
			end--
		}
		output = append(output, text[i:end])
		i += end + 1 - i
	}

	return output
}
