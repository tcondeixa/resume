package main

import (
	"errors"
	"flag"
	"os"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/tcondeixa/resume/internal/resume"
	"github.com/tcondeixa/resume/internal/terminal"

	"github.com/yaml/go-yaml"
)

func buildResume(filePath string) (*resume.Resume, error) {
	content, err := os.ReadFile(filePath)
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
	address := flag.String("address", "0.0.0.0:22", "address to listen on [host:port]")
	filePath := flag.String("file-path", "./resume.yaml", "path to the resume YAML file")

	flag.Parse()
	resume, err := buildResume(*filePath)
	if err != nil {
		log.Errorf("Could not build the resume: %v", err)
		os.Exit(1)
	}

	srv, err := wish.NewServer(
		wish.WithAddress(*address),
		wish.WithHostKeyPath(".ssh/resume_ed25519"),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(session ssh.Session) {
					term, err := getPTYInfo(session)
					if err != nil {
						log.Errorf("Could not get PTY info: %v", err)
						wish.Error(session, "Could not get terminal info")
					}

					formatedResume, err := term.Render(resume)
					if err != nil {
						log.Errorf("Could not get the resume: %v", err)
						wish.Error(session, "Could not generate the resume")
					} else {
						wish.Print(session, formatedResume)
					}
					next(session)
				}
			},

			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	log.Info("Starting SSH server", "address", address)
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not start server", "error", err)
	}
}

func getPTYInfo(session ssh.Session) (*terminal.Terminal, error) {
	pty, _, _ := session.Pty()
	return &terminal.Terminal{
		Term:   pty.Term,
		Width:  pty.Window.Width,
		Height: pty.Window.Height,
	}, nil
}
