package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oas "github.com/getkin/kin-openapi/openapi3"
)

var appStyle = lipgloss.NewStyle().Padding(0, 0)

type model struct {
	operations list.Model
}

func createModel(doc *oas.T) model {
	// TODO: be good if we could validate the document, but my test doc isn't valid

	return model{
		operations: NewOperationList(doc),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.operations.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.operations.FilterState() == list.Filtering {
			break
		}

		switch {
		}
	}

	newListModel, cmd := m.operations.Update(msg)
	m.operations = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.operations.View())
}

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
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		usage()

		fmt.Fprintln(os.Stderr, "\nPlease provide a valid OpenAPI document")
		os.Exit(1)
	}

	raw_url := os.Args[1]
	parsed_url, err := url.ParseRequestURI(raw_url)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		os.Exit(1)
	}

	ctx := context.Background()
	loader := &oas.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromURI(parsed_url)
	if err != nil {
		fmt.Println("Error loading OpenAPI document:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(createModel(doc), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
