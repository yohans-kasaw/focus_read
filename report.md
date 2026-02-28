# Focus Read - Project Report

## Overview
**Focus Read** is a terminal-based application designed for reading EPUB books and text files directly from the command line. It emphasizes a distraction-free reading experience ("focus mode") by presenting text in a clean, centered interface. The application also supports a unique "Paste Mode" for quickly reading content from the system clipboard.

## Features
- **EPUB Support:** Parses and renders standard EPUB e-book files.
- **Plain Text Support:** Reads standard text files.
- **Paste Mode:** Automatically grabs content from the system clipboard, sanitizes it, and presents it for reading.
- **Progress Tracking:** Automatically saves and restores reading progress for each file.
- **Terminal User Interface (TUI):** A minimalist, keyboard-driven interface built for focus.

## Architecture

The project is written in **Go** and follows a modular architecture, separating the CLI interface, core logic, and data handling.

### 1. Core Components

*   **Entry Point (`main.go`):** 
    -   Handles command-line arguments (e.g., `-paste`).
    -   Initializes the `ProgressStore` to manage reading history.
    -   Routes execution to either `handlePasteMode` or `handleNormalMode` based on user input.
    -   Launches the CLI and saves progress upon exit.

*   **CLI Module (`cli/`):**
    -   Built using the **Bubble Tea** framework (`github.com/charmbracelet/bubbletea`).
    -   **Model:** Manages the state, including the list of text segments and the current reading index.
    -   **View:** Renders the text using **Lipgloss** for styling, ensuring content is centered and readable.
    -   **Update:** Handles keyboard navigation (Left/Right arrows, Space, 'q' to quit).

*   **EPUB Module (`epub/`):**
    -   Responsible for parsing EPUB files (which are ZIP archives).
    -   **Parsing Logic:**
        -   Extracts the `container.xml` to locate the root file (OPF).
        -   Parses the Package Document (OPF) to understand the book structure (Manifest, Spine).
        -   Parses the Table of Contents (NCX).
        -   Extracts and normalizes text from HTML content files using `golang.org/x/net/html`.
    -   Flattens the hierarchical structure of an EPUB into a linear sequence of text segments for the reader.

### 2. Data Management

*   **Progress Store (`readingprogress.go`):**
    -   Persists reading progress in a `history.json` file located in the `cache/` directory.
    -   Tracks the file path and the last read index for multiple books.
    
*   **Helper Functions (`helpers.go`):**
    -   **Paste Handling:** Reads clipboard content, sanitizes filenames, and splits text into sentence-like segments for better readability.
    -   **File Validation:** Ensures paths are valid files and not directories.
    -   **Debugging:** Includes a `debugWriteTexts` function to dump parsed content to JSON files for inspection.

### 3. File Structure

- **`main.go`**: Application entry point.
- **`cli/`**: Contains the TUI logic and styling.
- **`epub/`**: Contains EPUB parsing and data structures.
- **`cache/`**: Stores reading history and debug files.
- **`paste/`**: Stores temporary files created from clipboard content.
- **`helpers.go`**: Utility functions for file I/O, clipboard access, and text processing.
- **`readingprogress.go`**: Logic for loading and saving reading progress.

## Dependencies

- **UI/UX:** `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`
- **System:** `github.com/atotto/clipboard` (Clipboard access)
- **HTML Parsing:** `golang.org/x/net/html`
- **Utilities:** `github.com/google/uuid`
