package screens

import (
	"fmt"
	"io"
	"strings"

	"dazzle/internal/domain"
	"dazzle/internal/ui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// operationItem adapts domain.Operation to list.Item.
type operationItem struct {
	op domain.Operation
}

func (i operationItem) Title() string {
	return fmt.Sprintf("%s %s", i.op.Method, i.op.Path)
}

func (i operationItem) Description() string { return i.op.Summary }
func (i operationItem) FilterValue() string {
	return string(i.op.Method) + " " + i.op.Path + " " + i.op.Summary
}

// operationDelegate renders operations with colored HTTP methods.
type operationDelegate struct{}

func (d operationDelegate) Height() int                             { return 2 }
func (d operationDelegate) Spacing() int                            { return 1 }
func (d operationDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d operationDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	op, ok := item.(operationItem)
	if !ok {
		return
	}

	method := styles.Method(string(op.op.Method))
	path := op.op.Path
	summary := op.op.Summary

	isSelected := index == m.Index()

	var title string
	if isSelected {
		title = lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("> %s %s", method, path))
		summary = lipgloss.NewStyle().Foreground(styles.Subtext1).Render("  " + summary)
	} else {
		title = fmt.Sprintf("  %s %s", method, path)
		summary = lipgloss.NewStyle().Foreground(styles.Overlay1).Render("  " + summary)
	}

	if strings.TrimSpace(op.op.Summary) != "" {
		fmt.Fprintf(w, "%s\n%s", title, summary)
	} else {
		fmt.Fprintf(w, "%s\n", title)
	}
}

// OperationsScreen displays a filterable list of API operations.
type OperationsScreen struct {
	list   list.Model
	width  int
	height int
}

func NewOperationsScreen(spec *domain.Spec, opSvc domain.OperationService) *OperationsScreen {
	ops := opSvc.SortOperations(opSvc.ListOperations(spec))

	items := make([]list.Item, len(ops))
	for i, op := range ops {
		items[i] = operationItem{op: op}
	}

	title := spec.Info.Title
	if title == "" {
		title = "Endpoints"
	}

	l := list.New(items, operationDelegate{}, 0, 0)
	l.Title = title
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.Title

	return &OperationsScreen{list: l}
}

func (s *OperationsScreen) Name() string { return "operations" }

func (s *OperationsScreen) Init() tea.Cmd { return nil }

func (s *OperationsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s *OperationsScreen) View() string {
	return s.list.View()
}
