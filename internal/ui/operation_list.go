package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"dazzle/internal/domain"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

const ellipsis = "…"

// OperationListItem adapts domain.Operation to list.Item
type OperationListItem struct {
	operation domain.Operation
}

func NewOperationListItem(op domain.Operation) OperationListItem {
	return OperationListItem{operation: op}
}

func (item OperationListItem) FilterValue() string {
	return item.operation.Path
}

func (item OperationListItem) Path() string {
	return item.operation.Path
}

func (item OperationListItem) Summary() string {
	return item.operation.Summary
}

func (item OperationListItem) Method() string {
	return string(item.operation.Method)
}

// OperationList manages the list of operations
type OperationList struct {
	model   list.Model
	service domain.OperationService
}

func NewOperationList(service domain.OperationService) (*OperationList, error) {
	operations, err := service.ListOperations()
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, len(operations))
	for i, op := range operations {
		items[i] = NewOperationListItem(op)
	}

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

	listModel := list.Model{
		KeyMap:                list.DefaultKeyMap(),
		Filter:                list.DefaultFilter,
		Styles:                styles,
		Title:                 "Endpoints",
		FilterInput:           filterInput,
		StatusMessageLifetime: time.Second,
		Paginator:             p,
		Help:                  h,
	}

	listModel.SetDelegate(NewOperationDelegate())
	listModel.SetItems(items)
	listModel.SetWidth(0)
	listModel.SetHeight(0)
	listModel.SetShowTitle(true)
	listModel.SetShowFilter(true)
	listModel.SetShowStatusBar(true)
	listModel.SetShowPagination(true)
	listModel.SetShowHelp(true)
	listModel.SetStatusBarItemName("item", "items")
	listModel.SetFilteringEnabled(true)

	return &OperationList{
		model:   listModel,
		service: service,
	}, nil
}

func (ol *OperationList) Update(msg tea.Msg) (list.Model, tea.Cmd) {
	newListModel, cmd := ol.model.Update(msg)
	ol.model = newListModel
	return newListModel, cmd
}

func (ol *OperationList) View() string {
	return ol.model.View()
}

func (ol *OperationList) SetSize(width, height int) {
	ol.model.SetSize(width, height)
}

func (ol *OperationList) FilterState() list.FilterState {
	return ol.model.FilterState()
}

// OperationDelegate handles rendering of individual operation items
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

	operationItem, ok := item.(OperationListItem)
	if !ok {
		return
	}

	path = operationItem.Path()
	summary = operationItem.Summary()
	method = operationItem.Method()

	switch method {
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

	if m.Width() <= 0 {
		return
	}

	itemHSize, _ := s.NormalItem.GetFrameSize()
	methodHSize, _ := s.NormalMethod.GetFrameSize()
	pathWidth := uint(m.Width() - itemHSize - methodHSize)
	path = truncate.StringWithTail(method+path, pathWidth, ellipsis)
	path = strings.Replace(path, method, "", 1)
	summaryWidth := uint(m.Width() - itemHSize)
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
			matched := unmatched.Inherit(s.FilterMatch)
			path = lipgloss.StyleRunes(path, matchedRunes, matched, unmatched)
		}
		itemStyle = s.SelectedItem
		path = s.SelectedPath.Render(path)
		summary = s.SelectedSummary.Render(summary)
		method = s.SelectedMethod.Render(method)
	} else {
		if isFiltered {
			unmatched := s.NormalPath.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			path = lipgloss.StyleRunes(path, matchedRunes, matched, unmatched)
		}
		itemStyle = s.NormalItem
		path = s.NormalPath.Render(path)
		summary = s.NormalSummary.Render(summary)
		method = s.NormalMethod.Foreground(methodColour).Render(method)
	}

	content := fmt.Sprintf("%s%s\n%s", method, path, summary)
	_, err := fmt.Fprintf(w, "%s", itemStyle.Render(content))
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
}