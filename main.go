package main

import (
	"context"
	"fmt"
	"os"

	"dazzle/internal/application"
	"dazzle/internal/infrastructure/openapi"
	"dazzle/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 {
		fmt.Println("dazzle — spec-aware API explorer")
		fmt.Println()
		fmt.Println("Usage: dazzle <spec-file-or-url>")
		return fmt.Errorf("expected exactly one argument")
	}

	source := os.Args[1]

	if f := os.Getenv("DAZZLE_DEBUG"); f != "" {
		logFile, err := tea.LogToFile(f, "dazzle")
		if err != nil {
			return fmt.Errorf("debug log: %w", err)
		}
		defer logFile.Close()
	}

	repo := openapi.NewRepository()
	specSvc := application.NewSpecService(repo)
	opSvc := application.NewOperationService()

	app := ui.NewAppModel(context.Background(), specSvc, opSvc, source)

	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
