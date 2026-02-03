package main

import (
	"fmt"
	"os"

	"dazzle/internal/application"
	"dazzle/internal/infrastructure"
	"dazzle/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func usage() {
	fmt.Println("Usage: dazzle <URL>")
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "warning: failed to close debug log:", err)
		}
	}()

	if len(os.Args) != 2 {
		usage()

		fmt.Fprintln(os.Stderr, "\nPlease provide a valid OpenAPI document")
		os.Exit(1)
	}

	docSource := os.Args[1]

	// Initialize repository (infrastructure layer)
	repo, err := infrastructure.NewOpenAPIRepository(docSource)
	if err != nil {
		fmt.Println("Error loading OpenAPI document:", err)
		os.Exit(1)
	}

	// Initialize service (application layer)
	service := application.NewOperationService(repo)

	// Initialize UI (presentation layer)
	appModel, err := ui.NewAppModel(service)
	if err != nil {
		fmt.Println("Error creating application:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(appModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}