# Terminal Wordle

I just wanted to play more wordle...

### How to Run:
Make sure to have [Go](https://go.dev/doc/install) installed.

#### Set up:
1. `git clone git@github.com:koutaroyumiba/terminal-wordle.git`
2. `cd terminal-wordle`

#### Running the app:
Run `./wordle` from the command line (make sure you are in the `terminal-wordle` directory)

For development:
From the `terminal-wordle` directory, run:
1. `go mod tidy`
2. `go run .`

### Notes:
- stats are saved in `stats.json` in root

### Logs:
- 12 October 2025
    - omg terminal UI is so cool
    - i was so locked in today (peak procrastination)
- 14 October 2025
    - wordle bots are also so cool
    - now we can find how many words we are guessing from
- 08 February 2026
    - refreshened the UI a little bit
    - changed the format structure to be more clean

🗂️ Project Map (Recommended)
wordle/
├── cmd/
│   └── wordle/
│       └── main.go          ← Program entrypoint (very small)
├── internal/
│   ├── tui/
│   │   ├── model.go         ← Bubble Tea model + Update()
│   │   ├── view.go          ← View() + layout (windows)
│   │   ├── styles.go        ← lipgloss styles
│   │   └── keyboard.go     ← keyboard rendering logic
│   ├── game/
│   │   ├── game.go          ← GameState, rules, validation
│   │   ├── cell.go          ← Cell, CellState, helpers
│   │   ├── stats.go         ← statistics + tracking
│   │   └── game_test.go
│   └── bot/
│       ├── bot.go           ← solver logic
│       ├── analysis.go      ← word filtering, heuristics
│       └── bot_test.go
├── assets/
│   └── words.txt            ← word lists, etc.
├── go.mod
├── go.sum
└── README.md

main
 ↓
tui
 ↓
game ← bot

Rules:

tui can import game and bot

bot can import game (for types)

game must NOT import tui or bot

main imports only tui

This keeps your core logic UI-agnostic, which is exactly what Go encourages.

🔌 What Each Layer Is Responsible For
cmd/wordle/main.go

Purpose: Wiring only

func main() {
    p := tea.NewProgram(tui.InitialModel())
    if err := p.Start(); err != nil {
        log.Fatal(err)
    }
}


🚫 No logic
🚫 No styling
🚫 No state

internal/game/

The “truth” of the application

Handles:

Word validation

Guess application

Cell state transitions

Win/loss detection

Statistics tracking

Should NOT know:

That a TUI exists

That colors exist

That a keyboard exists

💡 You should be able to reuse game in:

Web app

GUI app

API server

Tests only

internal/bot/

Pure logic / intelligence

Handles:

Candidate word filtering

Entropy / heuristics

Suggestions

Should:

Be deterministic

Be test-heavy

Have zero UI code

internal/tui/

Presentation + interaction only

Split it further:

model.go

Bubble Tea model

Init, Update

State transitions

view.go

Layout composition

Window composition:

game window

stats window

message window

lipgloss.JoinHorizontal / Vertical

styles.go

All lipgloss styles in one place

No logic

keyboard.go

Rendering keyboard

Mapping known letters → colors

This keeps your View() readable and sane.

🪟 Your “Window” Layout (Mentally)

Think of your UI as three independent panes:

┌───────────────────────────────┬───────────────┐
│                               │               │
│  Game Board                   │  Statistics   │
│  + Keyboard                   │               │
│                               │               │
└───────────────────────────────┴───────────────┘
┌───────────────────────────────────────────────┐
│ Message / Errors                               │
└───────────────────────────────────────────────┘


In code, this becomes:

main := lipgloss.JoinHorizontal(
    lipgloss.Top,
    gameView,
    statsView,
)

return lipgloss.JoinVertical(
    lipgloss.Left,
    main,
    messageView,
)


This is idiomatic Bubble Tea.

🧪 Testing Strategy

game: heavy tests

bot: heavy tests

tui: minimal tests (mostly visual/manual)

internal/game/game_test.go
internal/bot/bot_test.go

🚀 Why This Structure Scales Well

With this layout, you can later:

✅ Add a web UI
✅ Add a solver-only CLI
✅ Add multiplayer
✅ Replace TUI without touching logic

That’s real engineering separation.

🔑 Rule of Thumb (Memorize This)

If you delete internal/tui, your game should still work.
If you delete cmd, your logic should still compile.

If that’s true — your structure is solid.

If you want next:

I can refactor your existing files into this structure

Or help you split your current tui.go into model/view/styles cleanly

Or show how to use interfaces to decouple bot ↔ game

Just tell me what you want to tackle next.
