package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		return m.handleKeyPress(key)
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case colorPulseMsg:
		m.PulseTick++
		return m, colorPulseCmd()
	}
	return m, nil
}

func (m Model) handleKeyPress(key string) (tea.Model, tea.Cmd) {
	switch m.Screen {
	case ScreenDashboard:
		return m.handleDashboardKeys(key)
	case ScreenScheduleEditor:
		return m.handleEditorKeys(key)
	case ScreenLanguageSelector:
		return m.handleLanguageKeys(key)
	case ScreenSettings:
		return m.handleSettingsKeys(key)
	case ScreenInstall:
		return m.handleInstallKeys(key)
	case ScreenAbout:
		return m.handleAboutKeys(key)
	}
	return m, nil
}

// ── Dashboard ──

func (m Model) handleDashboardKeys(key string) (tea.Model, tea.Cmd) {
	options := []string{"Edit Schedule", "Settings", "Install / Manage", "About", "Exit"}

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
		}
	case "enter", " ":
		switch m.Cursor {
		case 0:
			m.PrevScreen = m.Screen
			m.Screen = ScreenScheduleEditor
			m.Cursor = 0
		case 1:
			m.PrevScreen = m.Screen
			m.Screen = ScreenSettings
			m.Cursor = 0
		case 2:
			m.PrevScreen = m.Screen
			m.Screen = ScreenInstall
			m.Cursor = 0
		case 3:
			m.PrevScreen = m.Screen
			m.Screen = ScreenAbout
			m.Cursor = 0
		case 4:
			return m, tea.Quit
		}
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		return m, tea.Quit
	}
	return m, nil
}

// ── Schedule Editor ──

func (m Model) handleEditorKeys(key string) (tea.Model, tea.Cmd) {
	if m.EditorMode {
		return m.handleEditorInput(key)
	}

	// Navigation mode
	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(m.Blocks)+1 { // +1 for "Add block" option
			m.Cursor++
		}
	case "enter", " ":
		if m.Cursor == len(m.Blocks) {
			// Add new block
			m.Blocks = append(m.Blocks, ScheduleBlock{Time: "00:00", Label: "New Block"})
			m.EditorMode = true
			m.EditField = 0
			m.EditingTime = "00:00"
			m.EditingLabel = "New Block"
			m.Cursor = len(m.Blocks) - 1
		} else if m.Cursor < len(m.Blocks) {
			// Edit existing block
			m.EditorMode = true
			m.EditField = 0
			m.EditingTime = m.Blocks[m.Cursor].Time
			m.EditingLabel = m.Blocks[m.Cursor].Label
		}
	case "d":
		// Delete block
		if m.Cursor < len(m.Blocks) {
			m.DeleteConfirm = true
			m.ConfirmIdx = m.Cursor
		}
	case "y":
		if m.DeleteConfirm {
			m.Blocks = append(m.Blocks[:m.ConfirmIdx], m.Blocks[m.ConfirmIdx+1:]...)
			m.DeleteConfirm = false
			if m.Cursor >= len(m.Blocks) {
				m.Cursor = len(m.Blocks) - 1
			}
		}
	case "n":
		m.DeleteConfirm = false
	case "t":
		// Toggle schedule on/off
		m.ScheduleEnabled = !m.ScheduleEnabled
	case "esc":
		if m.DeleteConfirm {
			m.DeleteConfirm = false
		} else {
			saveSchedule(m.Blocks)
			m.Screen = m.PrevScreen
			m.Cursor = 0
		}
	}
	return m, nil
}

func (m Model) handleEditorInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "tab":
		// Switch between time and label
		m.EditField = (m.EditField + 1) % 2
	case "enter":
		// Save edit
		if m.Cursor < len(m.Blocks) {
			m.Blocks[m.Cursor] = ScheduleBlock{
				Time:  m.EditingTime,
				Label: m.EditingLabel,
			}
		}
		m.EditorMode = false
	case "esc":
		m.EditorMode = false
	case "backspace":
		if m.EditField == 0 && len(m.EditingTime) > 0 {
			m.EditingTime = m.EditingTime[:len(m.EditingTime)-1]
		} else if m.EditField == 1 && len(m.EditingLabel) > 0 {
			m.EditingLabel = m.EditingLabel[:len(m.EditingLabel)-1]
		}
	default:
		if len(key) == 1 {
			if m.EditField == 0 && len(m.EditingTime) < 5 {
				// Only allow digits and : for time
				if strings.ContainsAny(key, "0123456789:") {
					m.EditingTime += key
				}
			} else if m.EditField == 1 && len(m.EditingLabel) < 40 {
				m.EditingLabel += key
			}
		}
	}
	return m, nil
}

// ── Language Selector ──

func (m Model) handleLanguageKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(m.Languages)-1 {
			m.Cursor++
		}
	case "enter", " ":
		m.SelectedLang = m.Cursor
		lang := m.Languages[m.SelectedLang]
		saveLanguage(lang)
		m.Screen = m.PrevScreen
		m.Cursor = 0
	case "esc":
		m.Screen = m.PrevScreen
		m.Cursor = 0
	}
	return m, nil
}

// ── Settings ──

func (m Model) handleSettingsKeys(key string) (tea.Model, tea.Cmd) {
	settingsCount := 2

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < settingsCount {
			m.Cursor++
		}
	case "enter", " ":
		switch m.Cursor {
		case 0:
			// Language selector
			m.PrevScreen = m.Screen
			m.Screen = ScreenLanguageSelector
			m.Cursor = 0
		case 1:
			// API Key
			m.Screen = ScreenSettings
			// For now just show status
		case 2:
			// Back
			m.Screen = m.PrevScreen
			m.Cursor = 0
		}
	case "esc":
		m.Screen = m.PrevScreen
		m.Cursor = 0
	}
	return m, nil
}

// ── Install ──

func (m Model) handleInstallKeys(key string) (tea.Model, tea.Cmd) {
	installed := m.IsInstalled
	enabled := m.ScheduleEnabled

	options := []string{}
	if !installed {
		options = append(options, "Install Service")
	} else {
		if enabled {
			options = append(options, "Disable Timer")
		} else {
			options = append(options, "Enable Timer")
		}
		options = append(options, "Uninstall Service")
	}
	options = append(options, "Back")

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
		}
	case "enter", " ":
		selected := options[m.Cursor]
		switch selected {
		case "Install Service":
			msg, err := installService()
			if err != nil {
				m.StatusMsg = fmt.Sprintf("Install failed: %v\n%s", err, msg)
			} else {
				m.IsInstalled = true
				m.StatusMsg = "EchoRoutine installed successfully!"
				toggleTimer(true)
				m.ScheduleEnabled = true
			}
		case "Enable Timer":
			err := toggleTimer(true)
			if err != nil {
				m.StatusMsg = fmt.Sprintf("Failed to enable: %v", err)
			} else {
				m.ScheduleEnabled = true
				m.StatusMsg = "Timer enabled!"
			}
		case "Disable Timer":
			err := toggleTimer(false)
			if err != nil {
				m.StatusMsg = fmt.Sprintf("Failed to disable: %v", err)
			} else {
				m.ScheduleEnabled = false
				m.StatusMsg = "Timer disabled."
			}
		case "Uninstall Service":
			toggleTimer(false)
			m.IsInstalled = false
			m.ScheduleEnabled = false
			m.StatusMsg = "Timer disabled. Remove service files manually if needed."
		case "Back":
			m.Screen = m.PrevScreen
			m.Cursor = 0
		}
	case "esc":
		m.Screen = m.PrevScreen
		m.Cursor = 0
	}
	return m, nil
}

// ── About ──

func (m Model) handleAboutKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc", "q", "enter", " ":
		m.Screen = m.PrevScreen
		m.Cursor = 0
	}
	return m, nil
}
