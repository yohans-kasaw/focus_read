from textual.app import App
from textual.widgets import Static


class FocusReadApp(App):

    def compose(self):
        yield Static("Hello, world!")

    def on_key(self, event):
        if event.key == "q":
            self.exit()

app = FocusReadApp()
app.run()

