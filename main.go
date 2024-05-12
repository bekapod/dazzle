package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oas "github.com/getkin/kin-openapi/openapi3"
)

var (
	appStyle = lipgloss.NewStyle().Padding(0, 0)

	listItemStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("241"))
	listItemSelectedStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("170"))
)

type path struct {
	path string
}

func (p path) Title() string       { return p.path }
func (p path) Description() string { return "" }
func (p path) FilterValue() string { return p.path }

type pathDelegate struct{}

func (d pathDelegate) Height() int                             { return 1 }
func (d pathDelegate) Spacing() int                            { return 0 }
func (d pathDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d pathDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(path)
	if !ok {
		return
	}

	str := i.Title()

	fn := listItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return listItemSelectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	paths list.Model
}

func createModel(doc *oas.T) model {
	// TODO: be good if we could validate the document, but my test doc isn't valid

	paths := make([]list.Item, 0)
	for p := range doc.Paths.Map() {
		paths = append(paths, path{path: p})
	}
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].(path).path < paths[j].(path).path
	})

	pathsList := list.New(paths, pathDelegate{}, 0, 0)
	pathsList.Title = "Paths"
	pathsList.SetShowStatusBar(false)
	pathsList.SetFilteringEnabled(false)
	pathsList.SetShowPagination(false)

	return model{
		paths: pathsList,
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
		m.paths.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.paths.FilterState() == list.Filtering {
			break
		}

		switch {
		}
	}

	newListModel, cmd := m.paths.Update(msg)
	m.paths = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.paths.View())
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
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
