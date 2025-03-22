from textual.app import App, ComposeResult
from textual.containers import Container, Horizontal, Vertical
from textual.widgets import Static 
from textual.binding import Binding
from textual.reactive import reactive
import pyperclip
import json
import os

class WordDisplay(Static):
    """A widget to display the current word with enhanced formatting."""
    
    DEFAULT_CSS = """
    WordDisplay {
        width: 100%;
        content-align: center middle;
    }
    
    #current-word {
        text-align: center;
    }
    """
    
    def __init__(self, word: str = "", *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.word = word
        
    def compose(self) -> ComposeResult:
        yield Static(self.word, id="current-word")

    def update_word(self, word: str) -> None:
        """Update the displayed word."""
        self.word = word
        current_word_widget = self.query_one("#current-word", Static)
        current_word_widget.update(self.word)

class Stats(Static):
    """Display reading statistics."""
    
    def compose(self) -> ComposeResult:
        yield Static("", id="stats-display")
    def update_stats(self, current: int, total: int) -> None:
        """Update the statistics display."""
        progress = (current / total) * 100 if total > 0 else 0
        stats = f"Progress: ({progress:.1f}%)"
        stats_widget = self.query_one("#stats-display", Static)
        stats_widget.update(stats)

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
        self.stats = Stats()
        
    DEFAULT_CSS = """
    #main-container {
        height: 100%;
        align: center middle;
    }
    """

    def compose(self) -> ComposeResult:
        """Create child widgets for the app."""
        yield Container(
            Vertical(
                self.stats,
                self.word_display,
                id="main-container"
            )
        )
        
    def on_mount(self) -> None:
        """Initialization after mounting."""
        self.load_clipboard_text()
        self.update_display()
        
    def load_clipboard_text(self) -> None:
        """Load text from clipboard."""
        text = pyperclip.paste()
        ## if text is empty line with no charactors
        if text.strip() == "":
            text = "No text found in clipboard"
        self.words = [word for word in text.split() if word.strip()]
        self.current_index = 0
        
    def update_display(self) -> None:
        """Update the display with current word and statistics."""
        if not self.words:
            return
            
        word = self.words[self.current_index]
        # Enhanced word display with ORP (Optimal Recognition Point)
        orp_index = len(word) // 3
        formatted_word = f"[green red]{word[:orp_index]}[/][dim]{word[orp_index:]}[/]"
        self.word_display.update_word(formatted_word)
        
        # Update progress bar and stats
        progress = (self.current_index / len(self.words)) * 100
        self.stats.update_stats(self.current_index + 1, len(self.words))
        
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
