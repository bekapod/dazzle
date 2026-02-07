package screens

import (
	"fmt"
	"sort"
	"strings"

	"dazzle/internal/domain"
	"dazzle/internal/ui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// DetailPanel renders operation metadata in a scrollable viewport.
type DetailPanel struct {
	viewport viewport.Model
	op       *domain.Operation
	width    int
	height   int
}

func NewDetailPanel(width, height int) *DetailPanel {
	vp := viewport.New(max(1, width-1), height)
	return &DetailPanel{
		viewport: vp,
		width:    width,
		height:   height,
	}
}

func (d *DetailPanel) SetOperation(op domain.Operation) {
	d.op = &op
	d.viewport.SetContent(d.renderContent())
	d.viewport.GotoTop()
}

func (d *DetailPanel) Clear() {
	d.op = nil
	d.viewport.SetContent("")
	d.viewport.GotoTop()
}

func (d *DetailPanel) SetSize(width, height int) {
	d.width = width
	d.height = height
	d.viewport.Width = max(1, width-1) // reserve 1 col for scrollbar
	d.viewport.Height = height
	if d.op != nil {
		d.viewport.SetContent(d.renderContent())
	}
}

func (d *DetailPanel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return cmd
}

func (d *DetailPanel) View() string {
	if d.op == nil {
		return styles.Muted.Render("Select an operation to view details")
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, d.viewport.View(), d.renderScrollbar())
}

func (d *DetailPanel) renderScrollbar() string {
	total := d.viewport.TotalLineCount()
	visible := d.viewport.VisibleLineCount()
	showThumb := total > visible
	thumbH := 0
	thumbTop := 0
	if showThumb {
		thumbH = max(1, d.height*visible/total)
		thumbTop = int(float64(d.height-thumbH) * d.viewport.ScrollPercent())
	}

	var b strings.Builder
	for i := range d.height {
		if i > 0 {
			b.WriteByte('\n')
		}
		if showThumb && i >= thumbTop && i < thumbTop+thumbH {
			b.WriteString(lipgloss.NewStyle().Foreground(styles.Overlay1).Render("┃"))
		} else {
			b.WriteString(" ")
		}
	}
	return b.String()
}

func (d *DetailPanel) renderContent() string {
	if d.op == nil {
		return ""
	}

	var b strings.Builder
	op := d.op

	method := styles.Method(string(op.Method))
	path := lipgloss.NewStyle().Bold(true).Render(op.Path)
	b.WriteString(method + " " + path)
	b.WriteString("\n")

	if op.Summary != "" {
		b.WriteString("\n")
		b.WriteString(op.Summary)
		b.WriteString("\n")
	}

	if op.Description != "" {
		b.WriteString("\n")
		b.WriteString(renderMarkdown(op.Description, max(1, d.viewport.Width-2)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(sectionHeader("Parameters"))
	if len(op.Parameters) == 0 {
		b.WriteString(styles.Muted.Render("  None"))
		b.WriteString("\n")
	} else {
		for _, p := range op.Parameters {
			b.WriteString(renderParameter(p, d.viewport.Width))
		}
	}

	b.WriteString("\n")
	b.WriteString(sectionHeader("Request Body"))
	if op.RequestBody == nil {
		b.WriteString(styles.Muted.Render("  None"))
		b.WriteString("\n")
	} else {
		b.WriteString(renderRequestBody(op.RequestBody, d.viewport.Width))
	}

	b.WriteString("\n")
	b.WriteString(sectionHeader("Responses"))
	if len(op.Responses) == 0 {
		b.WriteString(styles.Muted.Render("  None"))
		b.WriteString("\n")
	} else {
		b.WriteString(renderResponses(op.Responses))
	}

	return b.String()
}

func sectionHeader(title string) string {
	return styles.Title.Render(title) + "\n"
}

func renderParameter(p domain.Parameter, width int) string {
	var parts []string
	parts = append(parts, lipgloss.NewStyle().Bold(true).Render(p.Name))
	parts = append(parts, styles.Muted.Render(string(p.In)))
	if p.Schema != nil {
		if t := renderSchemaType(p.Schema); t != "" {
			parts = append(parts, t)
		}
	}
	if p.Required {
		parts = append(parts, lipgloss.NewStyle().Foreground(styles.Red).Render("required"))
	}

	line := "  " + strings.Join(parts, "  ")
	if p.Description != "" {
		desc := renderMarkdown(p.Description, max(1, width-6))
		line += "\n" + indent(desc, "    ")
	}
	return line + "\n"
}

func renderRequestBody(rb *domain.RequestBody, width int) string {
	var b strings.Builder
	if rb.Required {
		b.WriteString("  " + lipgloss.NewStyle().Foreground(styles.Red).Render("required") + "\n")
	}
	if rb.Description != "" {
		desc := renderMarkdown(rb.Description, max(1, width-4))
		b.WriteString(indent(desc, "  ") + "\n")
	}

	for _, contentType := range sortedKeys(rb.Content) {
		mt := rb.Content[contentType]
		b.WriteString("  " + lipgloss.NewStyle().Foreground(styles.Blue).Render(contentType) + "\n")
		if mt.Schema != nil {
			b.WriteString(renderSchemaProperties(mt.Schema, "    "))
		}
	}
	return b.String()
}

func renderResponses(responses map[string]domain.Response) string {
	var b strings.Builder
	for _, code := range sortedKeys(responses) {
		resp := responses[code]
		status := lipgloss.NewStyle().Bold(true).Foreground(styles.Green).Render(code)
		b.WriteString("  " + status)
		if resp.Description != "" {
			b.WriteString("  " + resp.Description)
		}
		b.WriteString("\n")

		for _, contentType := range sortedKeys(resp.Content) {
			mt := resp.Content[contentType]
			b.WriteString("    " + lipgloss.NewStyle().Foreground(styles.Blue).Render(contentType) + "\n")
			if mt.Schema != nil {
				b.WriteString(renderSchemaProperties(mt.Schema, "      "))
			}
		}
	}
	return b.String()
}

// renderSchemaType formats a schema type for inline display.
func renderSchemaType(s *domain.Schema) string {
	t := string(s.Type)
	if t == "" {
		return ""
	}
	if s.Type == domain.SchemaTypeArray && s.Items != nil {
		return fmt.Sprintf("array[%s]", string(s.Items.Type))
	}
	if s.Format != "" {
		return fmt.Sprintf("%s (%s)", t, s.Format)
	}
	return t
}

func renderSchemaProperties(s *domain.Schema, indent string) string {
	if s.Type == domain.SchemaTypeArray && s.Items != nil {
		if len(s.Items.Properties) > 0 {
			return indent + "array[object]:\n" + renderObjectProperties(s.Items, indent+"  ")
		}
		return indent + renderSchemaType(s) + "\n"
	}
	if len(s.Properties) == 0 {
		return indent + renderSchemaType(s) + "\n"
	}
	return renderObjectProperties(s, indent)
}

func renderObjectProperties(s *domain.Schema, indent string) string {
	names := sortedKeys(s.Properties)

	requiredSet := make(map[string]struct{}, len(s.Required))
	for _, r := range s.Required {
		requiredSet[r] = struct{}{}
	}

	var b strings.Builder
	for _, name := range names {
		prop := s.Properties[name]
		line := indent + lipgloss.NewStyle().Bold(true).Render(name) + ": " + renderSchemaType(prop)
		if _, ok := requiredSet[name]; ok {
			line += "  " + lipgloss.NewStyle().Foreground(styles.Red).Render("required")
		}
		b.WriteString(line + "\n")
	}
	return b.String()
}

// indent prepends prefix to every line of text.
func indent(text, prefix string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

// renderMarkdown renders text as terminal markdown using glamour. If rendering
// fails, the original text is returned as-is.
func renderMarkdown(text string, width int) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(styles.MarkdownStyle()),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return text
	}
	out, err := r.Render(text)
	if err != nil {
		return text
	}
	return strings.TrimRight(out, "\n")
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
