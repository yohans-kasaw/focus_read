from textual.app import App, ComposeResult
from textual.containers import Container, Horizontal, Vertical
from textual.widgets import Static
from textual.binding import Binding
from textual.reactive import reactive
import pyperclip
import json
import os

PREVIEW_WINDOW_SIZE = 30  


class WordDisplay(Static):
    """A widget to display the current word with enhanced formatting."""

    DEFAULT_CSS = """
    WordDisplay {
        width: 100%;
        margin: 12 0 0 0;
        content-align: center middle;
    }
    
    #word-preview {
        text-align: center;
    }
    """

    def __init__(self, word: str = "", *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.word = word

    def compose(self) -> ComposeResult:
        yield Static(self.word, id="word-preview")

    def update_word(self, word: str) -> None:
        """Update the displayed word."""
        self.word = word
        current_word_widget = self.query_one("#word-preview", Static)
        current_word_widget.update(self.word)


class Progress(Static):
    """Display reading progress."""

    DEFAULT_CSS = """
    Progress {
        dock: top;
        padding: 0 1;
        height: 1;
        background: $panel;
        color: $text-muted;
    }
    """

    def compose(self) -> ComposeResult:
        yield Static("", id="progress-display")

    def update_progress(self, current: int, total: int) -> None:
        """Update the progress display."""
        progress = (current / total) * 100 if total > 0 else 0
        progress_text = f"Progress: {progress:.1f}%"
        progress_widget = self.query_one("#progress-display", Static)
        progress_widget.update(progress_text)


class TextPreview(Static):
    """A widget to display text preview with current word highlight."""

    DEFAULT_CSS = """
    TextPreview {
        width: 70%;
        content-align: center middle;
    }
    
    #preview-text {
        text-align: center;
        color: $text-muted;
    }
    """

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.full_text = []  
        self.display_start = 0 

    def compose(self) -> ComposeResult:
        yield Static("", id="preview-text")

    def set_text(self, text: str) -> None:
        """Set the full text for preview."""
        self.full_text = [word for word in text.split() if word.strip()]
        self.display_start = 0
        self.update_preview(0)

    def update_preview(self, current_index: int) -> None:
        """Update the preview text with highlighted current word."""
        if not self.full_text:
            return

        if current_index >= self.display_start + PREVIEW_WINDOW_SIZE or current_index < self.display_start:
            self.display_start = current_index

        display_end = min(self.display_start + PREVIEW_WINDOW_SIZE, len(self.full_text))
        display_words = self.full_text[self.display_start : display_end]

        preview_parts = []
        for i, word in enumerate(display_words):
            actual_index = i + self.display_start
            if actual_index == current_index:
                preview_parts.append(f"[bold green]{word}[/]")
            else:
                preview_parts.append(f"[dim]{word}[/]")

        preview_text = " ".join(preview_parts)
        preview_widget = self.query_one("#preview-text", Static)
        preview_widget.update(preview_text)


class FocusReadApp(App):
    """A modern, feature-rich speed reading application."""

    BINDINGS = [
        Binding("left", "prev_word", "Previous"),
        Binding("right", "next_word", "Next"),
        Binding("q", "quit", "Quit"),
    ]

    current_index = reactive(0)
    words = reactive([])

    def __init__(self):
        super().__init__()
        self.word_display = WordDisplay()
        self.progress = Progress()
        self.text_preview = TextPreview()

    DEFAULT_CSS = """
    #main-container {
        height: 100%;
        width: 100%;
        align: center middle;
    }
    
    #word-container {
        align: center middle;
        height: 60%;
    }

    #text-preview-container {
        width: 100%;
        height: 40%;
        align: center bottom;
    }
    """

    def compose(self) -> ComposeResult:
        """Create child widgets for the app."""
        yield self.progress
        yield Container(
            Vertical(
                Vertical(self.word_display, id="word-container"),
                Vertical(self.text_preview, id="text-preview-container"),
                id="main-container",
            )
        )

    def on_mount(self) -> None:
        """Initialization after mounting."""
        self.load_clipboard_text()
        self.update_display()

    def load_clipboard_text(self) -> None:
        """Load text from clipboard."""
        text = pyperclip.paste()

        if text.strip() == "":
            text = "No text found in clipboard"
        self.words = [word for word in text.split() if word.strip()]
        self.current_index = 0
        self.text_preview.set_text(text)

    def update_display(self) -> None:
        """Update the display with current word and statistics."""
        if not self.words:
            return

        word = self.words[self.current_index]

        orp_index = len(word) // 3
        formatted_word = f"[bold yellow]{word[:orp_index]}[/][white]{word[orp_index:]}[/]"
        self.word_display.update_word(formatted_word)

        self.progress.update_progress(self.current_index + 1, len(self.words))

        self.text_preview.update_preview(self.current_index)

    def action_next_word(self) -> None:
        """Move to the next word."""
        if self.current_index < len(self.words) - 1:
            self.current_index += 1
            self.update_display()

    def action_prev_word(self) -> None:
        """Move to the previous word."""
        if self.current_index > 0:
            self.current_index -= 1
            self.update_display()


if __name__ == "__main__":
    app = FocusReadApp()
    app.run()
