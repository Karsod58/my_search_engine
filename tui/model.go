package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/Karsod58/search_engine/search"
)

type Model struct {
	input        textinput.Model
	spinner      spinner.Model
	searcher     *search.Searcher
	results      []search.Result
	suggestions  []string
	loading      bool
	err          error
	semanticMode bool 
	expansionMode bool 
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

func searchSemanticCmd(s *search.Searcher, query string) tea.Cmd {
	return func() tea.Msg {
	
		res := s.SemanticSearch(query, 5, 0.7)
		return searchFinishedMsg(res)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
			case "e":

	m.expansionMode = !m.expansionMode
	return m, nil


		case "s":
			
			m.semanticMode = !m.semanticMode
			return m, nil

		case "enter":
			if m.input.Value() == "" || m.loading {
				return m, nil
			}

			m.loading = true
			m.results = nil
			m.suggestions = nil

		
			var searchFn tea.Cmd
	if m.expansionMode {
		searchFn = searchExpansionCmd(m.searcher, m.input.Value())
	} else if m.semanticMode {
		searchFn = searchSemanticCmd(m.searcher, m.input.Value())
	} else {
		searchFn = searchCmd(m.searcher, m.input.Value())
	}

			return m, tea.Batch(
				m.spinner.Tick,
				searchFn,
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
func searchExpansionCmd(s *search.Searcher, query string) tea.Cmd {
	return func() tea.Msg {
		res := s.SearchWithExpansion(query, 5, s.GetExpander())
		return searchFinishedMsg(res)
	}
}


func (m Model) View() string {

	mode := "Keyword"
	if m.semanticMode {
		mode = "🤖 Semantic (AI)"
	}else if m.expansionMode {
		mode = "🧠 AI Expansion"
	}
	title := fmt.Sprintf("🔎 Go Search Engine [Mode: %s]\n\n", mode)
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
		
		if m.results[0].Stats != nil {
			stats := m.results[0].Stats
			out += fmt.Sprintf("Found %d results in %s (searched %d documents, matched %d terms)\n\n",
				stats.TotalResults, stats.QueryTime, stats.DocsSearched, stats.TermsMatched)
		}

		// Show corrections if any
		if len(m.results[0].Corrections) > 0 {
			out += "Showing results for: "
			for orig, corrected := range m.results[0].Corrections {
				out += fmt.Sprintf("%s → %s  ", orig, corrected)
			}
			out += "\n\n"
		}

		out += "Results:\n\n"
		for i, r := range m.results {
			// Display title (or DocID if no title)
			displayTitle := r.Title
			if displayTitle == "" {
				displayTitle = r.DocID
			}

			out += fmt.Sprintf("%d. %s (%.4f)\n", i+1, displayTitle, r.Score)

			// Display URL if available
			if r.URL != "" {
				out += fmt.Sprintf("   %s\n", r.URL)
			}

			// Display snippet
			if r.Snippet != "" {
				out += fmt.Sprintf("   %s\n", r.Snippet)
			}

			out += "\n"
		}
	} else {
		out += "No results yet.\n"
	}
	if len(m.results) > 0 && m.expansionMode {
		if expanded, ok := m.results[0].Corrections["expanded"]; ok && expanded != "" {
			out += fmt.Sprintf("AI Expanded: %s\n", expanded)
		}
		if intent, ok := m.results[0].Corrections["intent"]; ok {
			out += fmt.Sprintf("Intent: %s\n\n", intent)
		}
	}


	out += "\n(Enter = search • S = toggle mode • Esc = quit)\n"
	return out
}