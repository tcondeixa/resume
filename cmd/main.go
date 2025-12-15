package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/tcondeixa/resume/internal/resume"
	"github.com/tcondeixa/resume/internal/terminal"
	"golang.org/x/term"

	"github.com/yaml/go-yaml"
)

func buildResume() (*resume.Resume, error) {
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
	resume, err := buildResume()
	if err != nil {
		log.Errorf("Could not build the resume: %v", err)
		os.Exit(1)
	}

	outputFormat := flag.String("format", "terminal", "output format")
	flag.Parse()

	switch *outputFormat {
	case "terminal":
		term, err := getTerminalInfo()
		if err != nil {
			log.Errorf("Could not get terminal info: %v", err)
			os.Exit(1)
		}

		output, err := term.Render(resume)
		if err != nil {
			log.Errorf("Could not get the resume: %v", err)
			os.Exit(1)
		}

		fmt.Print(output)
		return
	}
}

func getTerminalInfo() (*terminal.Terminal, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}

	return &terminal.Terminal{
		Term:   os.Getenv("TERM"),
		Width:  width,
		Height: height,
	}, nil
}
