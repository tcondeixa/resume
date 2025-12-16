package main

import (
	"flag"
	"fmt"
	"os"

	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/log"
	"github.com/tcondeixa/resume/internal/resume"
	"github.com/tcondeixa/resume/internal/terminal"
	"golang.org/x/term"

	"github.com/yaml/go-yaml"
)

func loadResume() (*resume.Resume, error) {
	content, err := os.ReadFile("resume.yaml")
	if err != nil {
		return nil, err
	}

	var resume resume.Resume
	err = yaml.Unmarshal(content, &resume)
	if err != nil {
		return nil, err
	}

	return &resume, nil
}

func main() {
	outputFormat := flag.String("output-format", "terminal", "output format")
	flag.Parse()

	resume, err := loadResume()
	if err != nil {
		log.Fatalf("Could not load the resume: %v", err)
	}

	switch *outputFormat {
	case "terminal":
		term, err := getTerminalInfo()
		if err != nil {
			log.Fatalf("Could not get terminal info: %v", err)
		}

		output, err := term.Render(resume)
		if err != nil {
			log.Fatalf("Could not render the resume: %v", err)
		}

		fmt.Print(output)
	default:
		log.Fatalf("Unknown output format: %s", *outputFormat)
	}
}

func getTerminalInfo() (*terminal.Terminal, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}

	return &terminal.Terminal{
		Theme:  catppuccin.Mocha,
		Term:   os.Getenv("TERM"),
		Width:  width,
		Height: height,
	}, nil
}
