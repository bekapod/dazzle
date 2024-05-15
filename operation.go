package main

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oas "github.com/getkin/kin-openapi/openapi3"
	"github.com/muesli/reflow/truncate"
)

const (
	ellipsis = "…"
)

const (
	Get    = "GET"
	Post   = "POST"
	Put    = "PUT"
	Patch  = "PATCH"
	Delete = "DELETE"
)

type operation struct {
	path    string
	method  string
	summary string
}

func (o operation) Path() string {
	return o.path
}

func (o operation) Summary() string {
	return o.summary
}

func (o operation) Method() string {
	return o.method
}

func (o operation) FilterValue() string { return o.path }

func NewOperation(path string, method string, operationItem *oas.Operation) operation {
	return operation{
		path:    path,
		summary: operationItem.Summary,
		method:  method,
	}
}

func NewOperationList(doc *oas.T) list.Model {
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

	methodOrder := map[string]int{
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

	return newOperationListModel(operations)
}

func newOperationListModel(operations []list.Item) list.Model {
	styles := NewOperationListStyles()

	filterInput := textinput.New()
	filterInput.Prompt = "Filter: "
	filterInput.PromptStyle = styles.FilterPrompt
	filterInput.Cursor.Style = styles.FilterCursor
	filterInput.CharLimit = 64
	filterInput.Focus()

	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = styles.ActivePaginationDot.String()
	p.InactiveDot = styles.InactivePaginationDot.String()

	h := help.New()
	h.Styles = NewHelpStyles()

	m := list.Model{
		KeyMap:                list.DefaultKeyMap(),
		Filter:                list.DefaultFilter,
		Styles:                styles,
		Title:                 "Endpoints",
		FilterInput:           filterInput,
		StatusMessageLifetime: time.Second,
		Paginator:             p,
		Help:                  h,
	}

	m.SetDelegate(NewOperationDelegate())
	m.SetItems(operations)
	m.SetWidth(0)
	m.SetHeight(0)
	m.SetShowTitle(true)
	m.SetShowFilter(true)
	m.SetShowStatusBar(true)
	m.SetShowPagination(true)
	m.SetShowHelp(true)
	m.SetStatusBarItemName("item", "items")
	m.SetFilteringEnabled(true)

	return m
}

type Operation interface {
	list.Item
	Path() string
	Summary() string
	Method() string
}

type OperationDelegate struct {
	Styles  OperationItemStyles
	height  int
	spacing int
}

func NewOperationDelegate() OperationDelegate {
	return OperationDelegate{
		Styles:  NewOperationItemStyles(),
		height:  2,
		spacing: 1,
	}
}

func (d OperationDelegate) Height() int {
	return d.height
}

func (d OperationDelegate) Spacing() int {
	return d.spacing
}

func (d OperationDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d OperationDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		itemStyle             lipgloss.Style
		methodColour          lipgloss.AdaptiveColor
		path, summary, method string
		matchedRunes          []int
		s                     = &d.Styles
	)
	palette := NewPalette()

	if i, ok := item.(Operation); ok {
		path = i.Path()
		summary = i.Summary()
		method = i.Method()
		switch i.Method() {
		case "GET":
			methodColour = palette.Methods.Get
		case "POST":
			methodColour = palette.Methods.Post
		case "PUT":
			methodColour = palette.Methods.Put
		case "PATCH":
			methodColour = palette.Methods.Patch
		case "DELETE":
			methodColour = palette.Methods.Delete
		}
	} else {
		return
	}

	if m.Width() <= 0 {
		return
	}

	pathWidth := uint(m.Width() - s.NormalPath.GetPaddingLeft() - s.NormalPath.GetPaddingRight())
	path = truncate.StringWithTail(path, pathWidth, ellipsis)
	summaryWidth := uint(m.Width() - s.NormalSummary.GetPaddingLeft() - s.NormalSummary.GetPaddingRight())
	summary = truncate.StringWithTail(summary, summaryWidth, ellipsis)

	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if isFiltered && index < len(m.VisibleItems()) {
		matchedRunes = m.MatchesForItem(index)
	}

	if emptyFilter {
		itemStyle = s.DimmedItem
		path = s.DimmedPath.Render(path)
		summary = s.DimmedSummary.Render(summary)
		method = s.DimmedMethod.Foreground(methodColour).Render(method)
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			unmatched := s.SelectedPath.Inline(true)
			matched := unmatched.Copy().Inherit(s.FilterMatch)
			path = lipgloss.StyleRunes(path, matchedRunes, matched, unmatched)
		}
		itemStyle = s.SelectedItem
		path = s.SelectedPath.Render(path)
		summary = s.SelectedSummary.Render(summary)
		method = s.SelectedMethod.Render(method)
	} else {
		if isFiltered {
			unmatched := s.NormalPath.Inline(true)
			matched := unmatched.Copy().Inherit(s.FilterMatch)
			path = lipgloss.StyleRunes(path, matchedRunes, matched, unmatched)
		}
		itemStyle = s.NormalItem
		path = s.NormalPath.Render(path)
		summary = s.NormalSummary.Render(summary)
		method = s.NormalMethod.Foreground(methodColour).Render(method)
	}

	content := fmt.Sprintf("%s%s\n%s", method, path, summary)
	fmt.Fprintf(w, "%s", itemStyle.Render(content))
}
