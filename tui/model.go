package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

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
	summaryMode bool
	chatMode bool
	chatHistory []ChatEntry
	chatAnswer string
	chatSources []string
}
type ChatEntry struct{
	Question string
	Answer string
	Sources []string
}
type chatFinishedMsg struct{
	answer string
	sources []string
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
func chatCmd(s *search.Searcher, question string) tea.Cmd {
	return func() tea.Msg {
		answer, sources, err := s.ChatWithDocs(question, 5)
		if err != nil {
			return chatFinishedMsg{answer: "Error: " + err.Error(), sources: nil}
		}
		return chatFinishedMsg{answer: answer, sources: sources}
	}
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
		var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "c":
			// Toggle chat mode
			m.chatMode = !m.chatMode
			if m.chatMode {
				// Clear other modes
				m.semanticMode = false
				m.expansionMode = false
				m.summaryMode = false
				m.results = nil
			}
			return m, nil

		case "ctrl+l":
			// Clear chat history
			if m.chatMode && m.searcher.GetRAGChat() != nil {
				m.searcher.GetRAGChat().ClearHistory()
				m.chatHistory = []ChatEntry{}
				m.chatAnswer = ""
				m.chatSources = nil
			}
			return m, nil

		case "enter":
			if m.input.Value() == "" || m.loading {
				return m, nil
			}

			m.loading = true

			if m.chatMode {
				// Chat mode
				return m, tea.Batch(
					m.spinner.Tick,
					chatCmd(m.searcher, m.input.Value()),
				)
			} else {
				// Search mode
				m.results = nil
				m.suggestions = nil

				var searchFn tea.Cmd
				if m.summaryMode {
					searchFn = searchSummaryCmd(m.searcher, m.input.Value())
				} else if m.expansionMode {
					searchFn = searchExpansionCmd(m.searcher, m.input.Value())
				} else if m.semanticMode {
					searchFn = searchSemanticCmd(m.searcher, m.input.Value())
				} else {
					searchFn = searchCmd(m.searcher, m.input.Value())
				}

				return m, tea.Batch(m.spinner.Tick, searchFn)
			}

		case "ctrl+c", "esc":
			return m, tea.Quit

		default:
			if !m.chatMode {
				query := m.input.Value()
				m.suggestions = m.searcher.Suggest(query)
			}
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

	case chatFinishedMsg:
		m.chatAnswer = msg.answer
		m.chatSources = msg.sources
		m.chatHistory = append(m.chatHistory, ChatEntry{
			Question: m.input.Value(),
			Answer:   msg.answer,
			Sources:  msg.sources,
		})
		m.input.SetValue("")
		m.loading = false
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
func searchSummaryCmd(s *search.Searcher, query string) tea.Cmd {
	return func() tea.Msg {
		res := s.SearchWithSummary(query, 5)
		return searchFinishedMsg(res)
	}
}

func searchExpansionCmd(s *search.Searcher, query string) tea.Cmd {
	return func() tea.Msg {
		res := s.SearchWithExpansion(query, 5, s.GetExpander())
		return searchFinishedMsg(res)
	}
}


func (m Model) View() string {
	if m.chatMode {
		return m.renderChatView()
	}
	return m.renderSearchView()
}

func (m Model) renderChatView() string {
	out := "💬 RAG Chat Mode - Ask questions about your documents\n\n"

	
	if len(m.chatHistory) > 0 {
		
		start := 0
		if len(m.chatHistory) > 3 {
			start = len(m.chatHistory) - 3
		}

		for _, entry := range m.chatHistory[start:] {
			out += fmt.Sprintf("You: %s\n", entry.Question)
			out += fmt.Sprintf("AI: %s\n", entry.Answer)
			if len(entry.Sources) > 0 {
				out += fmt.Sprintf("   Sources: %s\n", strings.Join(entry.Sources, ", "))
			}
			out += "\n"
		}
	}

	
	if m.loading {
		out += "Thinking " + m.spinner.View() + "\n"
	}

	out += m.input.View() + "\n\n"

	if len(m.chatHistory) == 0 && !m.loading {
		out += "Start a conversation! Ask anything about your indexed documents.\n"
	}

	out += "\n(Enter = ask • C = toggle mode • Ctrl+L = clear history • Esc = quit)\n"
	return out
}

func (m Model) renderSearchView() string {
	mode := "Keyword"
	if m.summaryMode {
		mode = "📝 AI Summary"
	} else if m.expansionMode {
		mode = "🧠 AI Expansion"
	} else if m.semanticMode {
		mode = "🤖 Semantic"
	}

	title := fmt.Sprintf("🔎 Go Search Engine [Mode: %s]\n\n", mode)
	input := m.input.View() + "\n\n"

	if m.loading {
		return title + input + "Searching " + m.spinner.View() + "\n"
	}

	out := title + input

	

	out += "\n(Enter = search • S = semantic • E = expansion • A = summary • C = chat • Esc = quit)\n"
	return out
}