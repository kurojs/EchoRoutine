package tui

import (
	"time"

	"github.com/charmbracelet/bubbletea"
)

// Screen represents the current view
type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenScheduleEditor
	ScreenLanguageSelector
	ScreenSettings
	ScreenInstall
	ScreenAbout
)

// ScheduleBlock represents a single time block
type ScheduleBlock struct {
	Time  string
	Label string
}

type colorPulseMsg struct{}

// Model holds all TUI state
type Model struct {
	Screen     Screen
	PrevScreen Screen
	Width      int
	Height     int
	Cursor     int
	Scroll     int

	// Color pulse animation
	PulseTick int

	// Schedule data
	Blocks          []ScheduleBlock
	ScheduleEnabled bool

	// Language
	Languages    []string
	SelectedLang int

	// API Key
	APIKey string

	// Install status
	IsInstalled bool
	StatusMsg   string

	// Editor state
	EditorMode    bool   // true = editing label, false = navigating
	EditingLabel  string // temp buffer for label editing
	EditingTime   string // temp buffer for time editing
	EditField     int    // 0 = time, 1 = label
	AddMode       bool   // true = adding new block
	DeleteConfirm bool   // true = confirm deletion
	ConfirmIdx    int    // index to delete
}

func NewModel() (Model, error) {
	blocks := loadSchedule()
	langs := getAllLanguages()
	currentLang := loadLanguage()

	langIdx := 0
	for i, l := range langs {
		if l == currentLang {
			langIdx = i
			break
		}
	}

	return Model{
		Screen:          ScreenDashboard,
		Blocks:          blocks,
		Languages:       langs,
		SelectedLang:    langIdx,
		ScheduleEnabled: isTimerEnabled(),
		IsInstalled:     isServiceInstalled(),
		APIKey:          loadAPIKey(),
	}, nil
}

func (m Model) Init() tea.Cmd {
	return colorPulseCmd()
}

func colorPulseCmd() tea.Cmd {
	return tea.Tick(colorPulseInterval, func(t time.Time) tea.Msg {
		return colorPulseMsg{}
	})
}

func (m Model) VisibleBlocks() []ScheduleBlock {
	maxVisible := m.Height - 12
	if maxVisible < 4 {
		maxVisible = 4
	}
	if m.Scroll > len(m.Blocks)-maxVisible {
		m.Scroll = max(0, len(m.Blocks)-maxVisible)
	}
	if m.Scroll < 0 {
		m.Scroll = 0
	}
	end := m.Scroll + maxVisible
	if end > len(m.Blocks) {
		end = len(m.Blocks)
	}
	return m.Blocks[m.Scroll:end]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
