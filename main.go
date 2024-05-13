package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oas "github.com/getkin/kin-openapi/openapi3"
)

var appStyle = lipgloss.NewStyle().Padding(0, 0)

type Method string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

type endpoint struct {
	path    string
	method  Method
	summary string
}

func (e endpoint) Title() string {
	return fmt.Sprintf("%s %s", e.method, e.path)
}
func (e endpoint) Description() string { return e.summary }
func (e endpoint) FilterValue() string { return e.path }

type model struct {
	endpoints list.Model
}

func createModel(doc *oas.T) model {
	// TODO: be good if we could validate the document, but my test doc isn't valid

	endpoints := make([]list.Item, 0)
	for path, pathItem := range doc.Paths.Map() {
		if endpointItem := pathItem.Get; endpointItem != nil {
			endpoints = append(endpoints, endpoint{path: path, summary: endpointItem.Summary, method: Get})
		}

		if endpointItem := pathItem.Post; endpointItem != nil {
			endpoints = append(endpoints, endpoint{path: path, summary: endpointItem.Summary, method: Post})
		}

		if endpointItem := pathItem.Put; endpointItem != nil {
			endpoints = append(endpoints, endpoint{path: path, summary: endpointItem.Summary, method: Put})
		}

		if endpointItem := pathItem.Patch; endpointItem != nil {
			endpoints = append(endpoints, endpoint{path: path, summary: endpointItem.Summary, method: Patch})
		}

		if endpointItem := pathItem.Delete; endpointItem != nil {
			endpoints = append(endpoints, endpoint{path: path, summary: endpointItem.Summary, method: Delete})
		}
	}

	methodOrder := map[Method]int{
		Get:    1,
		Post:   2,
		Put:    3,
		Patch:  4,
		Delete: 5,
	}
	sort.Slice(endpoints, func(i, j int) bool {
		// sort by endpoint path first then by method
		if endpoints[i].(endpoint).path == endpoints[j].(endpoint).path {
			return methodOrder[endpoints[i].(endpoint).method] < methodOrder[endpoints[j].(endpoint).method]
		}

		return endpoints[i].(endpoint).path < endpoints[j].(endpoint).path
	})

	endpointsList := list.New(endpoints, list.NewDefaultDelegate(), 0, 0)
	endpointsList.Title = "endpoints"
	endpointsList.SetShowPagination(false)

	return model{
		endpoints: endpointsList,
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
		m.endpoints.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.endpoints.FilterState() == list.Filtering {
			break
		}

		switch {
		}
	}

	newListModel, cmd := m.endpoints.Update(msg)
	m.endpoints = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.endpoints.View())
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

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
