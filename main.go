package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

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

type operation struct {
	path    string
	method  Method
	summary string
	tags    []string
}

func (e operation) Title() string {
	return fmt.Sprintf("%s %s", e.method, e.path)
}

func (e operation) Description() string {
	return fmt.Sprintf("%s %s", e.summary, strings.Join(e.tags, ","))
}
func (e operation) FilterValue() string { return e.path }

func NewOperation(path string, method Method, operationItem *oas.Operation) operation {
	return operation{
		path:    path,
		summary: operationItem.Summary,
		method:  method,
		tags:    operationItem.Tags,
	}
}

type model struct {
	operations list.Model
}

func createModel(doc *oas.T) model {
	// TODO: be good if we could validate the document, but my test doc isn't valid

	operations := make([]list.Item, 0)
	for path, pathItem := range doc.Paths.Map() {
		if operationItem := pathItem.Get; operationItem != nil {
			operations = append(operations, NewOperation(path, Get, operationItem))
		}

		if operationItem := pathItem.Post; operationItem != nil {
			operations = append(operations, NewOperation(path, Post, operationItem))
		}

		if operationItem := pathItem.Put; operationItem != nil {
			operations = append(operations, NewOperation(path, Put, operationItem))
		}

		if operationItem := pathItem.Patch; operationItem != nil {
			operations = append(operations, NewOperation(path, Patch, operationItem))
		}

		if operationItem := pathItem.Delete; operationItem != nil {
			operations = append(operations, NewOperation(path, Delete, operationItem))
		}
	}

	methodOrder := map[Method]int{
		Get:    1,
		Post:   2,
		Put:    3,
		Patch:  4,
		Delete: 5,
	}
	sort.Slice(operations, func(i, j int) bool {
		// sort by operation path first then by method
		if operations[i].(operation).path == operations[j].(operation).path {
			return methodOrder[operations[i].(operation).method] < methodOrder[operations[j].(operation).method]
		}

		return operations[i].(operation).path < operations[j].(operation).path
	})

	operationsList := list.New(operations, list.NewDefaultDelegate(), 0, 0)
	operationsList.Title = "endpoints"
	operationsList.SetShowPagination(false)

	return model{
		operations: operationsList,
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
