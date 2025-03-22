import curses
import sys
import subprocess
from datetime import datetime
import pyperclip

def get_clipboard_text():
    result = pyperclip.paste()
    return result

def main(stdscr):
    curses.curs_set(0)  # Hide the cursor
    stdscr.nodelay(1)   # Non-blocking input
    stdscr.timeout(100) # Refresh every 100 milliseconds

    # Initialize color pair for green text
    curses.start_color()
    curses.init_pair(1, curses.COLOR_GREEN, curses.COLOR_BLACK)

    text = get_clipboard_text().split()
    if not text:
        return

    idx = 0
    y, x = stdscr.getmaxyx()
    x_center = x // 2

    while idx < len(text):
        stdscr.clear()
        word = text[idx]
        x_position = x_center - (len(word) // 2)
        
        line_y = (y // 2)
        for i, char in enumerate(word):
            if i < len(word) // 2:
                stdscr.addstr(line_y, x_position + i, char, curses.color_pair(1) | curses.A_BOLD)
            else:
                stdscr.addstr(line_y, x_position + i, char)

        stdscr.refresh()

        key = stdscr.getch()
        if key == curses.KEY_LEFT and idx > 0:
            idx -= 1
        elif key == curses.KEY_RIGHT:
            if idx == len(text) - 1:
                break
            idx += 1
        elif key == ord('q'):
            break

    save_text_to_file(text)
    sys.exit()

def save_text_to_file(text):
    current_time = datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
    first_few_words = "_".join(text[:5])
    filename = f"text_{current_time}_{first_few_words}.txt"
    with open(filename, 'w') as file:
        file.write(" ".join(text))
    print(f"Text saved to {filename}")

if __name__ == "__main__":
    curses.wrapper(main)
