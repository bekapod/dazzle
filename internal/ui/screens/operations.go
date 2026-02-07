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

type panelFocus int

const (
	focusList panelFocus = iota
	focusDetail
)

// OperationsScreen displays a split-pane view: filterable operation list
// on the left, operation detail on the right.
type OperationsScreen struct {
	list   list.Model
	detail *DetailPanel
	focus  panelFocus
	lastID string
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

	s := &OperationsScreen{
		list:   l,
		detail: NewDetailPanel(0, 0),
	}

	s.syncDetail()
	return s
}

func (s *OperationsScreen) Name() string { return "operations" }

func (s *OperationsScreen) Init() tea.Cmd { return nil }

func (s *OperationsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.layoutPanels()
		return s, nil

	case tea.KeyMsg:
		if msg.String() == "tab" {
			s.toggleFocus()
			return s, nil
		}
		// Detail panel has no text input, so 'q' always means quit.
		if s.focus == focusDetail && msg.String() == "q" {
			return s, tea.Quit
		}
		// Route key messages based on which panel is focused.
		var cmd tea.Cmd
		if s.focus == focusDetail {
			cmd = s.detail.Update(msg)
		} else {
			s.list, cmd = s.list.Update(msg)
			s.syncDetail()
		}
		return s, cmd

	case tea.MouseMsg:
		// Route mouse scroll to whichever panel the cursor is over.
		var cmd tea.Cmd
		if s.panelAt(msg.X) == focusDetail {
			cmd = s.detail.Update(msg)
		} else {
			s.list, cmd = s.list.Update(msg)
			s.syncDetail()
		}
		return s, cmd
	}

	// Non-key messages (e.g. FilterMatchesMsg) always go to the list so it
	// can process async results regardless of which panel is focused.
	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	s.syncDetail()
	return s, cmd
}

func (s *OperationsScreen) View() string {
	if s.width == 0 {
		return ""
	}

	listWidth := s.listWidth()
	detailWidth := s.width - listWidth
	contentH := max(1, s.height-2)
	listContentW := max(1, listWidth-2)
	detailContentW := max(1, detailWidth-2) // border (2); padding is inside Width

	var activeBorder, inactiveBorder lipgloss.Style
	activeBorder = lipgloss.NewStyle().
		Height(contentH).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Blue)
	inactiveBorder = lipgloss.NewStyle().
		Height(contentH).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Surface1)

	var listBorder, detailBorder lipgloss.Style
	if s.focus == focusList {
		listBorder = activeBorder.Width(listContentW)
		detailBorder = inactiveBorder.Width(detailContentW).PaddingLeft(1).PaddingRight(1)
	} else {
		listBorder = inactiveBorder.Width(listContentW)
		detailBorder = activeBorder.Width(detailContentW).PaddingLeft(1).PaddingRight(1)
	}

	listView := listBorder.Render(s.list.View())
	detailView := detailBorder.Render(s.detail.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, listView, detailView)
}

func (s *OperationsScreen) listWidth() int {
	return s.width / 3
}

func (s *OperationsScreen) layoutPanels() {
	listWidth := s.listWidth()
	detailWidth := s.width - listWidth
	contentH := max(1, s.height-2)

	// Account for border (1 char each side), clamped to avoid negative sizes.
	// Detail panel also has 1 char horizontal padding on each side.
	s.list.SetSize(max(1, listWidth-2), contentH)
	s.detail.SetSize(max(1, detailWidth-4), contentH)
}

// panelAt returns which panel occupies the given x coordinate.
func (s *OperationsScreen) panelAt(x int) panelFocus {
	if x >= s.listWidth() {
		return focusDetail
	}
	return focusList
}

func (s *OperationsScreen) toggleFocus() {
	if s.focus == focusList {
		s.focus = focusDetail
	} else {
		s.focus = focusList
	}
}

func (s *OperationsScreen) syncDetail() {
	item, ok := s.list.SelectedItem().(operationItem)
	if !ok {
		s.lastID = ""
		s.detail.Clear()
		return
	}
	if item.op.ID == s.lastID {
		return
	}
	s.lastID = item.op.ID
	s.detail.SetOperation(item.op)
}
