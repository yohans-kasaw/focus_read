# Focus Read

Focus Read is a terminal-based speed reading tool designed to help you read EPUB books and text content efficiently and without distractions. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework, it provides a clean, focused interface for consuming content.

## Features

-   **EPUB Support**: Parse and read EPUB files directly in your terminal.
-   **Text File Support**: Read plain text files.
-   **Clipboard Integration**: Instantly read text from your clipboard using the `--paste` flag.
    -   Automatically cleans up whitespace.
    -   Intelligently splits pasted text into sentences for better flow.
-   **Progress Tracking**: Automatically saves your reading progress.
    -   Resume exactly where you left off when you reopen a book.
    -   Launch without arguments to continue the last opened book.
-   **Distraction-Free UI**: A minimal interface that displays text chunk by chunk to maintain focus.

## Installation

### Prerequisites

-   [Go](https://go.dev/) (version 1.24 or higher)

### Build from Source

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/focus_read.git
    cd focus_read
    ```

2.  Build the application:
    ```bash
    go build -o focusRead .
    ```

## Usage

### Reading a Book

To read an EPUB or text file, simply pass the file path as an argument:

```bash
./focusRead path/to/your/book.epub
```

### Resuming Reading

If you run the application without any arguments, it will automatically resume the last book you were reading:

```bash
./focusRead
```

### Reading from Clipboard

To read text currently in your clipboard (useful for articles or long posts):

```bash
./focusRead --paste
```
*Note: The pasted content is automatically saved to the `./paste` directory for future reference.*

### Controls

| Key | Action |
| :--- | :--- |
| `Space` / `Right Arrow` | Next text chunk |
| `Left Arrow` | Previous text chunk |
| `q` | Quit the application |

## Project Structure

-   `main.go`: Entry point of the application.
-   `cli/`: Contains the Bubble Tea TUI implementation.
-   `epub/`: Custom EPUB parsing logic.
-   `cache/`: Stores debug output of parsed text chunks.
-   `paste/`: Stores content saved from clipboard paste mode.
-   `history.json`: Tracks reading progress for each file.

## License

[MIT](LICENSE)
