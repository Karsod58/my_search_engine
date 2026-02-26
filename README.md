# Search Engine

A lightweight, in-memory search engine built in Go that implements TF-IDF (Term Frequency-Inverse Document Frequency) ranking with an inverted index. Features an interactive terminal UI powered by Bubble Tea.

## Features

- **Inverted Index**: Efficient document retrieval using an inverted index data structure
- **TF-IDF Ranking**: Ranks search results based on term frequency and inverse document frequency
- **Text Processing**: Tokenization, stopword filtering, and normalization
- **Interactive TUI**: Terminal-based user interface for real-time search
- **Top-K Results**: Returns the most relevant documents for any query

## Architecture

The project is organized into modular packages:

- `documents/` - Document structure and management
- `index/` - Inverted index implementation with TF-IDF scoring
- `processor/` - Text processing (tokenization, stopword removal)
- `search/` - Search logic and result ranking
- `tui/` - Terminal user interface using Bubble Tea

## How It Works

1. **Indexing**: Documents are processed to extract terms and their frequencies
2. **TF Calculation**: Term frequency is computed as `count / total_terms`
3. **IDF Calculation**: Inverse document frequency is computed as `log(N / df)` where N is total documents and df is document frequency
4. **Search**: Query terms are processed and matched against the index
5. **Ranking**: Documents are scored using TF-IDF and top-K results are returned

## Installation

```bash
go mod download
```

## Usage

Run the search engine:

```bash
go run main.go
```

This launches an interactive terminal interface where you can enter search queries and view ranked results.

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal output

## Project Structure

```
.
├── documents/       # Document data structures
├── index/          # Inverted index implementation
├── processor/      # Text processing and tokenization
├── search/         # Search and ranking logic
├── tui/            # Terminal user interface
└── main.go         # Application entry point
```

## Technical Details

### Inverted Index

An inverted index maps each term to the documents that contain it along with their TF scores. This enables fast lookups when searching for terms.

### Text Processing

- Converts text to lowercase
- Removes non-alphanumeric characters (except spaces)
- Filters common stopwords
- Calculates term frequencies

### Ranking Algorithm

Documents are ranked using the TF-IDF score:
```
score = Σ (TF(term, doc) × IDF(term))
```

Where the sum is over all query terms present in the document.


