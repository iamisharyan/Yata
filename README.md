# Yata

A fast, keyboard-driven **terminal file manager written in Go**, built on the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework. Yata gives you a Miller-columns layout (parent · current · preview), fuzzy file search, file-type icons, and a built-in file preview, all without leaving your terminal.

## Features

- **Three-pane layout**: parent directory, current directory and a live preview of the selected file or folder
- **Vim-style navigation**: `h` `j` `k` `l`, `gg`, `G`, plus arrow keys
- **Fuzzy finder**: press `/` to search every file and folder under the current directory
- **File preview**: text files are previewed instantly (first 200 lines); binary files are detected and skipped
- **Fullscreen preview**: expand the preview to read longer files, with keyboard and mouse scrolling
- **File operations**: copy, cut, paste and rename without leaving the app
- **Open in your editor**: opens files in `$EDITOR` (falls back to `nvim`)
- **Hidden files toggle**: show or hide dotfiles with one key
- **File-type icons**: colored icons via [go-devicons](https://github.com/epilande/go-devicons)
- **Mouse support**: click to select, scroll wheel to scroll the preview
- **Status bar**: shows the current path, selected item and permissions

## Requirements

- [Go](https://go.dev/dl/) 1.26 or newer (see `go.mod`)
- A terminal with true-color support
- A [Nerd Font](https://www.nerdfonts.com/) for icons to render correctly (the fuzzy finder prompt also uses a Nerd Font glyph)
- Optional: `nvim` or any editor set in `$EDITOR`

## Installation

```bash
git clone https://github.com/iamisharyan/Yata.git
cd Yata
go build -o yata .
```

Run it from any directory; Yata opens in the directory you launch it from:

```bash
./yata
```

Or run without building:

```bash
go run .
```

## Keybindings

### Navigation

| Key | Action |
| --- | --- |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `l` / `→` / `Enter` | Enter the selected directory |
| `h` / `←` | Go to the parent directory |
| `gg` | Jump to the top |
| `G` | Jump to the bottom |
| `w` / `Ctrl+d` | Jump 10 entries down |
| `o` / `Ctrl+u` | Jump 10 entries up |
| `.` | Toggle hidden files |

### Search

| Key | Action |
| --- | --- |
| `/` | Open the fuzzy finder |
| `↓` / `Ctrl+l` | Next result |
| `↑` / `Ctrl+k` | Previous result |
| `Enter` / `=` | Jump to the result (folders are entered; files open in fullscreen preview) |
| `Esc` | Close the finder |

### Preview

| Key | Action |
| --- | --- |
| `f` | Toggle fullscreen preview for the selected file |
| `Esc` | Leave fullscreen preview |
| `alt+1` / `alt+2` | Scroll down / up one line |
| `0` / `9` | Scroll half a page down / up |
| Mouse wheel | Scroll the preview |

### File operations

| Key | Action |
| --- | --- |
| `c` | Copy the selected item |
| `x` | Cut the selected item |
| `p` | Paste into the current directory |
| `r` | Rename (type the new name, `Enter` to confirm, `Esc` to cancel) |
| `n` | Open the selected file (see below) |

### General

| Key | Action |
| --- | --- |
| `q` / `Ctrl+c` | Quit |

## Opening files

Pressing `n` on a file opens it depending on its type:

- **Text and code files**: opened in `$EDITOR`, defaulting to `nvim`
- **Images** (`png`, `jpg`, `jpeg`, `gif`, `bmp`, `webp`, `tiff`, `ico`) and **documents** (`doc`, `docx`, `docm`, `dot`, `dotx`, `dotm`, `rtf`, `odt`, `pdf`): opened with the system default application using `cmd /c start`, so this currently works on **Windows only**

## Project structure

```
.
├── main.go              # Entry point, starts the Bubble Tea program
└── internal/
    ├── tui/             # Model, update loop and view rendering
    ├── fs/              # Directory reading, sorting and fuzzy search
    ├── ops/             # Copy, move and rename operations
    ├── preview/         # Text preview and binary detection
    ├── icons/           # File-type icons
    └── theme/           # Colors and pane styles
```

## Built with

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles): viewport and text input components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss): styling and layout
- [Chroma](https://github.com/alecthomas/chroma): syntax highlighting
- [sahilm/fuzzy](https://github.com/sahilm/fuzzy): fuzzy matching
- [go-devicons](https://github.com/epilande/go-devicons): file icons
- [fileflow](https://github.com/spf13/fileflow): safe file copy, move and rename

## Known limitations

- Copying whole directories is not supported yet (cutting and moving a directory works within the same filesystem)
- Opening images and documents with the system default app is Windows-only for now
- Deleting and creating files or folders are not implemented yet

## Contributing

Issues and pull requests are welcome. Fork the repo, create a feature branch, and open a PR describing your change.

## License

No license file has been added to this repository yet. Add one (for example MIT) and update this section.
