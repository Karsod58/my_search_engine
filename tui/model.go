package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/Karsod58/search_engine/search"
)

type Model struct {
	input       textinput.Model
	spinner     spinner.Model
	searcher    *search.Searcher
	results     []search.Result
	suggestions []string
	loading     bool
	err         error
}

type searchFinishedMsg []search.Result

func New(searcher *search.Searcher) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter search query..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return Model{
		input:       ti,
		spinner:     sp,
		searcher:    searcher,
		results:     nil,
		suggestions: nil,
		loading:     false,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func searchCmd(s *search.Searcher, query string) tea.Cmd {
	return func() tea.Msg {
		res := s.Search(query, 5)
		return searchFinishedMsg(res)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			if m.input.Value() == "" || m.loading {
				return m, nil
			}

			m.loading = true
			m.results = nil
			m.suggestions = nil

			return m, tea.Batch(
				m.spinner.Tick,
				searchCmd(m.searcher, m.input.Value()),
			)

		case "ctrl+c", "esc":
			return m, tea.Quit

		default:
			query := m.input.Value()
			m.suggestions = m.searcher.Suggest(query)
		}

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case searchFinishedMsg:
		m.results = msg
		m.loading = false
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {

	title := "🔎 Go Search Engine\n\n"
	input := m.input.View() + "\n\n"

	if m.loading {
		return title + input + "Searching " + m.spinner.View() + "\n"
	}

	out := title + input

	if len(m.suggestions) > 0 {
		out += "Suggestions:\n"
		for _, s := range m.suggestions {
			out += "  " + s + "\n"
		}
		out += "\n"
	}

	if len(m.results) > 0 {
		out += "Results:\n\n"
		for i, r := range m.results {
			out += fmt.Sprintf(
				"%d. %s (%.4f)\n   %s\n\n",
				i+1,
				r.DocID,
				r.Score,
				r.Snippet,
			)
		}
	} else {
		out += "No results yet.\n"
	}

	out += "\n(Enter = search • Esc = quit)\n"
	return out
}