# Resume SSH Server

A Go application that serves your resume over SSH. Connect via SSH to view a beautifully formatted resume in your terminal.

## Overview

This project creates an SSH server that displays your resume when users connect.
The resume is defined in a YAML file and rendered according to the client's terminal capabilities.

## Features

- SSH server for serving resumes
- YAML-based resume configuration
- Terminal-aware formatting
- Configurable host and port
- Host key authentication

## Installation

```bash
make
```

## Usage

### Basic Usage

```bash
./resume-server
```

This starts the server on `0.0.0.0:22` using `./resume.yaml` as the resume file.

### Command Line Options

```bash
./resume --address "localhost:2222" --file-path "/path/to/my-resume.yaml"
```

**Flags:**

- `-address`: Server address (default: `0.0.0.0:22`)
- `-file-path`: Path to resume YAML file (default: `./resume.yaml`)

### Connecting to Your Resume

Once the server is running, users can view your resume by connecting via SSH:

```bash
ssh your-server-address
```

## Resume Configuration

Create a `resume.yaml` file with your resume data.
The exact structure depends on your `internal/resume` package implementation.

## Requirements

- Go 1.19+
- SSH host key (automatically generated at `.ssh/resume_ed25519`)

## Dependencies

- `github.com/charmbracelet/log`
- `github.com/charmbracelet/ssh`
- `github.com/charmbracelet/wish`
- `github.com/yaml/go-yaml`
