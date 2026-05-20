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

// EditorFocus for the time/label editor
type EditorFocus int

const (
	FocusHours  EditorFocus = iota
	FocusMinutes
	FocusLabel
)

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

	// Install status
	IsInstalled bool
	IsEnabled   bool
	StatusMsg   string

	// Editor state
	Editing     bool
	EditBlock   int
	EditFocus   EditorFocus
	EditHours   int
	EditMinutes int
	EditLabel   string
	AddMode     bool

	// Delete confirmation
	DeleteConfirm bool
	ConfirmIdx    int
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
		IsEnabled:       isTimerEnabled(),
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

// parentOf returns the screen to go back to
func (m Model) parentOf(screen Screen) Screen {
	switch screen {
	case ScreenScheduleEditor:
		return ScreenDashboard
	case ScreenLanguageSelector:
		return ScreenSettings
	case ScreenSettings:
		return ScreenDashboard
	case ScreenInstall:
		return ScreenDashboard
	case ScreenAbout:
		return ScreenDashboard
	default:
		return ScreenDashboard
	}
}

// minBlockTime returns the minimum time (hours, minutes) for a block index
func (m Model) minBlockTime(idx int) (int, int) {
	if idx <= 0 {
		return 0, 0
	}
	h, min := parseTime(m.Blocks[idx-1].Time)
	return h, min
}
