package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "ctrl+c" {
			return m, tea.Quit
		}
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
	if key == "ctrl+c" {
		return m, tea.Quit
	}
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
			m.Screen = ScreenScheduleEditor
			m.Cursor = 0
		case 1:
			m.Screen = ScreenSettings
			m.Cursor = 0
		case 2:
			m.Screen = ScreenInstall
			m.Cursor = 0
		case 3:
			m.Screen = ScreenAbout
			m.Cursor = 0
		case 4:
			return m, tea.Quit
		}
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

// ── Schedule Editor ──

func (m Model) handleEditorKeys(key string) (tea.Model, tea.Cmd) {
	if m.Editing {
		return m.handleEditInput(key)
	}

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(m.Blocks) { // +1 for "Add" handled in view
			m.Cursor++
			if m.Cursor > len(m.Blocks) {
				m.Cursor = len(m.Blocks)
			}
		}
	case "enter", " ":
		if m.Cursor == len(m.Blocks) {
			// Add new block — default time = last block + 1h or 00:00
			lastH, lastM := 0, 0
			if len(m.Blocks) > 0 {
				lastH, lastM = parseTime(m.Blocks[len(m.Blocks)-1].Time)
			}
			newH := lastH + 1
			if newH > 23 {
				newH = 23
			}
			m.Blocks = append(m.Blocks, ScheduleBlock{
				Time:  fmt.Sprintf("%02d:%02d", newH, lastM),
				Label: "New Block",
			})
			m.Editing = true
			m.EditBlock = len(m.Blocks) - 1
			m.EditFocus = FocusHours
			m.EditHours = newH
			m.EditMinutes = lastM
			m.EditLabel = "New Block"
			m.Cursor = len(m.Blocks) - 1
		} else if m.Cursor < len(m.Blocks) {
			// Edit existing block
			h, min := parseTime(m.Blocks[m.Cursor].Time)
			m.Editing = true
			m.EditBlock = m.Cursor
			m.EditFocus = FocusHours
			m.EditHours = h
			m.EditMinutes = min
			m.EditLabel = m.Blocks[m.Cursor].Label
		}
	case "d":
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
	case "esc":
		if m.DeleteConfirm {
			m.DeleteConfirm = false
		} else {
			saveSchedule(m.Blocks)
			m.Screen = m.parentOf(ScreenScheduleEditor)
			m.Cursor = 0
		}
	}
	return m, nil
}

func (m Model) handleEditInput(key string) (tea.Model, tea.Cmd) {
	switch m.EditFocus {
	case FocusHours:
		return m.handleEditHours(key)
	case FocusMinutes:
		return m.handleEditMinutes(key)
	case FocusLabel:
		return m.handleEditLabel(key)
	}
	return m, nil
}

func (m Model) handleEditHours(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		minH, _ := m.minBlockTime(m.EditBlock)
		if m.EditHours < 23 {
			m.EditHours++
		}
		if m.EditHours < minH {
			m.EditHours = minH
		}
	case "down", "j":
		minH, _ := m.minBlockTime(m.EditBlock)
		if m.EditHours > minH {
			m.EditHours--
		}
	case "tab", "right":
		m.EditFocus = FocusMinutes
	case "enter":
		m.saveEditedBlock()
	case "esc":
		m.Editing = false
	}
	return m, nil
}

func (m Model) handleEditMinutes(key string) (tea.Model, tea.Cmd) {
	// Clamp minutes: if hours == prev hours, minutes must be > prev minutes
	minH, minM := m.minBlockTime(m.EditBlock)
	minAllowed := 0
	if m.EditHours == minH {
		minAllowed = minM + 1
		if minAllowed > 59 {
			minAllowed = 0
			m.EditHours++
		}
	}

	switch key {
	case "up", "k":
		if m.EditMinutes < 59 {
			m.EditMinutes++
		}
		// Enforce minimum
		if m.EditHours == minH && m.EditMinutes <= minM {
			m.EditMinutes = minAllowed
		}
	case "down", "j":
		if m.EditMinutes > minAllowed {
			m.EditMinutes--
		}
		if m.EditMinutes < 0 {
			m.EditMinutes = 0
		}
	case "tab", "right":
		m.EditFocus = FocusLabel
	case "left":
		m.EditFocus = FocusHours
	case "enter":
		m.saveEditedBlock()
	case "esc":
		m.Editing = false
	}
	return m, nil
}

func (m Model) handleEditLabel(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "tab", "right":
		m.EditFocus = FocusHours
	case "left":
		m.EditFocus = FocusMinutes
	case "enter":
		m.saveEditedBlock()
	case "esc":
		m.Editing = false
	case "backspace":
		if len(m.EditLabel) > 0 {
			m.EditLabel = m.EditLabel[:len(m.EditLabel)-1]
		}
	default:
		if len(key) == 1 && len(m.EditLabel) < 40 {
			m.EditLabel += key
		}
	}
	return m, nil
}

func (m *Model) saveEditedBlock() {
	if m.EditBlock >= 0 && m.EditBlock < len(m.Blocks) {
		timeStr := fmt.Sprintf("%02d:%02d", m.EditHours, m.EditMinutes)
		m.Blocks[m.EditBlock] = ScheduleBlock{
			Time:  timeStr,
			Label: m.EditLabel,
		}
		saveSchedule(m.Blocks)
	}
	m.Editing = false
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
		m.Screen = m.parentOf(ScreenLanguageSelector)
		m.Cursor = 0
	case "esc":
		m.Screen = m.parentOf(ScreenLanguageSelector)
		m.Cursor = 0
	}
	return m, nil
}

// ── Settings ──

func (m Model) handleSettingsKeys(key string) (tea.Model, tea.Cmd) {
	settingsCount := 2 // Language, ElevenLabs Voice

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < settingsCount { // includes Back
			m.Cursor++
		}
	case "enter", " ":
		switch m.Cursor {
		case 0:
			// Language selector
			m.Screen = ScreenLanguageSelector
			m.Cursor = 0
		case 1:
			// ElevenLabs Voice — show MCP link
			// (no-op, just informational for now)
		case 2:
			// Back
			m.Screen = m.parentOf(ScreenSettings)
			m.Cursor = 0
		}
	case "esc":
		m.Screen = m.parentOf(ScreenSettings)
		m.Cursor = 0
	}
	return m, nil
}

// ── Install ──

func (m Model) handleInstallKeys(key string) (tea.Model, tea.Cmd) {
	options := buildInstallOptions(m)

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
		return m.handleInstallSelect(options[m.Cursor])
	case "esc":
		m.Screen = m.parentOf(ScreenInstall)
		m.Cursor = 0
	}
	return m, nil
}

func buildInstallOptions(m Model) []string {
	var opts []string
	if !m.IsInstalled {
		opts = append(opts, "Install EchoRoutine — systemd service + timer")
	} else {
		if m.IsEnabled {
			opts = append(opts, "Pause announcements — keep installed, disable timer")
		} else {
			opts = append(opts, "Resume announcements — enable timer")
		}
		opts = append(opts, "Fully uninstall — remove service and timer")
	}
	opts = append(opts, "← Back")
	return opts
}

func (m Model) handleInstallSelect(selected string) (tea.Model, tea.Cmd) {
	switch {
	case strings.Contains(selected, "Install"):
		msg, err := installService()
		if err != nil {
			m.StatusMsg = fmt.Sprintf("Install failed: %v\n%s", err, msg)
		} else {
			m.IsInstalled = true
			m.IsEnabled = true
			m.StatusMsg = "EchoRoutine installed! Announcements will start on next timer tick."
		}
	case strings.Contains(selected, "Pause"):
		err := toggleTimer(false)
		if err != nil {
			m.StatusMsg = fmt.Sprintf("Failed to pause: %v", err)
		} else {
			m.IsEnabled = false
			m.StatusMsg = "Announcements paused. Timer disabled."
		}
	case strings.Contains(selected, "Resume"):
		err := toggleTimer(true)
		if err != nil {
			m.StatusMsg = fmt.Sprintf("Failed to resume: %v", err)
		} else {
			m.IsEnabled = true
			m.StatusMsg = "Announcements resumed! Timer enabled."
		}
	case strings.Contains(selected, "uninstall"):
		toggleTimer(false)
		m.IsInstalled = false
		m.IsEnabled = false
		m.StatusMsg = "Uninstalled. Service files left in place — remove manually if needed."
	case strings.Contains(selected, "Back"):
		m.Screen = m.parentOf(ScreenInstall)
		m.Cursor = 0
	}
	return m, nil
}

// ── About ──

func (m Model) handleAboutKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc", "q", "enter", " ":
		m.Screen = m.parentOf(ScreenAbout)
		m.Cursor = 0
	}
	return m, nil
}

// ── Helpers ──

func parseTime(t string) (int, int) {
	parts := strings.Split(t, ":")
	h, _ := strconv.Atoi(parts[0])
	m, _ := 0, 0
	if len(parts) > 1 {
		m, _ = strconv.Atoi(parts[1])
	}
	if h < 0 {
		h = 0
	}
	if h > 23 {
		h = 23
	}
	if m < 0 {
		m = 0
	}
	if m > 59 {
		m = 59
	}
	return h, m
}
