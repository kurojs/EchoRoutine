package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	content := ""
	switch m.Screen {
	case ScreenDashboard:
		content = m.renderDashboard()
	case ScreenScheduleEditor:
		content = m.renderEditor()
	case ScreenLanguageSelector:
		content = m.renderLanguageSelector()
	case ScreenSettings:
		content = m.renderSettings()
	case ScreenInstall:
		content = m.renderInstall()
	case ScreenAbout:
		content = m.renderAbout()
	}

	w := m.Width
	h := m.Height
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 40
	}
	return BaseStyle.Render(lipgloss.Place(
		w, h,
		lipgloss.Center, lipgloss.Center,
		content,
	))
}

// ── Dashboard ──

func (m Model) renderDashboard() string {
	pulseColor := PulseColor(m.PulseTick)

	banner := RenderBanner(pulseColor)

	nextBlock, _ := getNextBlock(m.Blocks)

	// Status
	status := lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Render(
		fmt.Sprintf("\n  %s  Schedule:  %s  (%d blocks)",
			StatusDot(m.ScheduleEnabled),
			statusText(m.ScheduleEnabled),
			len(m.Blocks),
		),
	)

	langDisplay := m.Languages[m.SelectedLang]
	langShort := strings.SplitN(langDisplay, " - ", 2)[0]

	langStatus := MutedStyle.Render(fmt.Sprintf("  %s  Language:  %s", StatusDot(true), langShort))

	apiStatus := MutedStyle.Render(fmt.Sprintf("  %s  API Key:   %s", StatusDot(m.APIKey != ""), apiKeyDisplay(m.APIKey)))

	// Next block info
	var nextBlockStr string
	if m.ScheduleEnabled && len(m.Blocks) > 0 {
		nextBlockStr = GreenStyle.Render(fmt.Sprintf("\n  ◉  Next:  %s — %s", nextBlock.Time, nextBlock.Label))
	} else {
		nextBlockStr = MutedStyle.Render("\n  ◉  No upcoming blocks")
	}

	// Menu
	menuItems := []string{
		"Edit Schedule",
		"Settings",
		"Install / Manage",
		"About",
		"Exit",
	}

	var menuLines []string
	for i, item := range menuItems {
		menuLines = append(menuLines, MenuItem(item, m.Cursor == i))
	}
	menu := lipgloss.JoinVertical(lipgloss.Left, menuLines...)

	// Subtitle
	subtitle := AccentStyle.Render("\n  Your AI-powered daily routine voice  ")

	// Divider
	div := Divider

	// Key hints
	hints := lipgloss.JoinHorizontal(lipgloss.Left,
		KeyHint("↑↓", "Navigate  "),
		KeyHint("Enter", "Select  "),
		KeyHint("q", "Quit"),
	)

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		banner,
		subtitle,
		"",
		div,
		status,
		langStatus,
		apiStatus,
		nextBlockStr,
		"",
		div,
		"",
		menu,
		"",
		div,
		"",
		hints,
	)

	return lipgloss.NewStyle().Padding(0, 2).Render(content)
}

// ── Schedule Editor ──

func (m Model) renderEditor() string {
	var lines []string

	title := TitleStyle.Render("Schedule Editor")
	lines = append(lines, title, "")

	if m.DeleteConfirm {
		confirmStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorWarning)).
			Bold(true)
		lines = append(lines, confirmStyle.Render("  Delete this block? (y/n)"))
		lines = append(lines, "")
	}

	// Block list
	for i, block := range m.Blocks {
		selected := i == m.Cursor && !m.EditorMode
		prefix := "  "
		if selected {
			prefix = "▸ "
		}

		var line string
		if m.EditorMode && i == m.Cursor {
			// Show editing state
			timeDisplay := m.EditingTime
			if m.EditField == 0 {
				timeDisplay = lipgloss.NewStyle().
					Foreground(lipgloss.Color(colorPurpleLight)).
					Background(lipgloss.Color(colorSurface)).
					Underline(true).
					Render(m.EditingTime)
			} else {
				timeDisplay = MutedStyle.Render(m.EditingTime)
			}

			labelDisplay := m.EditingLabel
			if m.EditField == 1 {
				labelDisplay = lipgloss.NewStyle().
					Foreground(lipgloss.Color(colorPurpleLight)).
					Background(lipgloss.Color(colorSurface)).
					Underline(true).
					Render(m.EditingLabel)
			} else {
				labelDisplay = MutedStyle.Render(m.EditingLabel)
			}

			line = fmt.Sprintf("%s %s  %s  [Tab to switch, Enter to save]",
				prefix,
				timeDisplay,
				labelDisplay,
			)
		} else {
			style := MutedStyle
			if selected {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPurpleLight))
			}
			line = fmt.Sprintf("%s %s  %s", prefix, block.Time, style.Render(block.Label))
		}
		lines = append(lines, line)
	}

	// Add block option
	addSelected := m.Cursor == len(m.Blocks) && !m.EditorMode
	lines = append(lines, MenuItemGreen("+ Add Block", addSelected))

	lines = append(lines, "", Divider, "")

	// Key hints
	if m.EditorMode {
		hints := lipgloss.JoinHorizontal(lipgloss.Left,
			KeyHint("Tab", "Switch field  "),
			KeyHint("Enter", "Save  "),
			KeyHint("Esc", "Cancel"),
		)
		lines = append(lines, hints)
	} else {
		hints := lipgloss.JoinHorizontal(lipgloss.Left,
			KeyHint("↑↓", "Navigate  "),
			KeyHint("Enter", "Edit  "),
			KeyHint("d", "Delete  "),
			KeyHint("t", "Toggle  "),
			KeyHint("Esc", "Back"),
		)
		lines = append(lines, hints)
	}

	lines = append(lines, "")
	status := fmt.Sprintf("  Blocks: %d | Status: %s", len(m.Blocks), statusText(m.ScheduleEnabled))
	lines = append(lines, MutedStyle.Render(status))

	return lipgloss.NewStyle().Padding(0, 2).Render(strings.Join(lines, "\n"))
}

// ── Language Selector ──

func (m Model) renderLanguageSelector() string {
	var lines []string

	lines = append(lines, TitleStyle.Render("Select Language"))
	lines = append(lines, "")

	// Scrollable language list
	visibleLanguages := m.Height - 8
	if visibleLanguages < 5 {
		visibleLanguages = 5
	}

	start := 0
	if m.Cursor >= visibleLanguages {
		start = m.Cursor - visibleLanguages + 1
	}
	end := start + visibleLanguages
	if end > len(m.Languages) {
		end = len(m.Languages)
	}

	for i := start; i < end; i++ {
		selected := i == m.Cursor
		lang := m.Languages[i]
		isCurrent := i == m.SelectedLang

		prefix := "  "
		if selected {
			prefix = "▸ "
		}

		suffix := ""
		if isCurrent {
			suffix = lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen)).Render("  ✓")
		}

		style := MutedStyle
		if selected {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPurpleLight))
		}
		if isCurrent && !selected {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen))
		}

		lines = append(lines, fmt.Sprintf("%s%s%s", prefix, style.Render(lang), suffix))
	}

	lines = append(lines, "", Divider, "")
	hints := lipgloss.JoinHorizontal(lipgloss.Left,
		KeyHint("↑↓", "Navigate  "),
		KeyHint("Enter", "Select  "),
		KeyHint("Esc", "Back"),
	)
	lines = append(lines, hints)

	return lipgloss.NewStyle().Padding(0, 2).Render(strings.Join(lines, "\n"))
}

// ── Settings ──

func (m Model) renderSettings() string {
	var lines []string

	lines = append(lines, TitleStyle.Render("Settings"))
	lines = append(lines, "")

	langDisplay := m.Languages[m.SelectedLang]
	langShort := strings.SplitN(langDisplay, " - ", 2)[0]

	settings := []struct {
		name  string
		value string
	}{
		{"Language", langShort},
		{"API Key", apiKeyDisplay(m.APIKey)},
	}

	for i, s := range settings {
		selected := m.Cursor == i
		prefix := "  "
		if selected {
			prefix = "▸ "
		}

		style := MutedStyle
		if selected {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPurpleLight))
		}

		nameStyle := style.Render(s.name)
		valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Render(s.value)
		lines = append(lines, fmt.Sprintf("%s%s: %s", prefix, nameStyle, valueStyle))
	}

	// Back option
	lines = append(lines, "")
	backSelected := m.Cursor == 2
	lines = append(lines, MenuItem("← Back", backSelected))

	lines = append(lines, "", Divider, "")
	hints := lipgloss.JoinHorizontal(lipgloss.Left,
		KeyHint("↑↓", "Navigate  "),
		KeyHint("Enter", "Select  "),
		KeyHint("Esc", "Back"),
	)
	lines = append(lines, hints)

	return lipgloss.NewStyle().Padding(0, 2).Render(strings.Join(lines, "\n"))
}

// ── Install ──

func (m Model) renderInstall() string {
	var lines []string

	lines = append(lines, TitleStyle.Render("Install / Manage"))
	lines = append(lines, "")

	// Status section
	if m.IsInstalled {
		lines = append(lines, GreenStyle.Render(fmt.Sprintf("  ◉ Service: Installed")))
		if m.ScheduleEnabled {
			lines = append(lines, GreenStyle.Render(fmt.Sprintf("  ◉ Timer:  Enabled")))
		} else {
			lines = append(lines, WarningStyle.Render(fmt.Sprintf("  ◉ Timer:  Disabled")))
		}
	} else {
		lines = append(lines, MutedStyle.Render("  ◉ Service: Not installed"))
	}
	lines = append(lines, "")

	// Status message
	if m.StatusMsg != "" {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorGreen)).
			Padding(0, 2).
			Render(m.StatusMsg))
		lines = append(lines, "")
	}

	// Menu
	options := []string{}
	if !m.IsInstalled {
		options = append(options, "Install Service")
	} else {
		if m.ScheduleEnabled {
			options = append(options, "Disable Timer")
		} else {
			options = append(options, "Enable Timer")
		}
	}
	options = append(options, "← Back")

	for i, opt := range options {
		lines = append(lines, MenuItem(opt, m.Cursor == i))
	}

	lines = append(lines, "", Divider, "")
	hints := lipgloss.JoinHorizontal(lipgloss.Left,
		KeyHint("↑↓", "Navigate  "),
		KeyHint("Enter", "Execute  "),
		KeyHint("Esc", "Back"),
	)
	lines = append(lines, hints)

	return lipgloss.NewStyle().Padding(0, 2).Render(strings.Join(lines, "\n"))
}

// ── About ──

func (m Model) renderAbout() string {
	title := TitleStyle.Render("About EchoRoutine")

	desc := fmt.Sprintf(`
  EchoRoutine %s

  Your AI-powered daily routine voice assistant.
  Built with ♥ for KDE Plasma | Arch Linux

  Each block of your day gets announced via
  ElevenLabs TTS with custom AI motivation.

  EchoRoutine is part of the kurojs ecosystem.
  github.com/kurojs/schedule-announcer

 ── Tech Stack ──

  Trigger: Bash + systemd timer
  Voice:  ElevenLabs TTS (via MCP)
  TUI:    Go + Bubbletea
  Theme:  Kurox (purple × green)

 ── Commands ──

  echoroutine    Launch this TUI
  Ctrl+C / q     Quit
`, lipgloss.NewStyle().Foreground(lipgloss.Color(colorPurpleLight)).Render("v1.0.0"))

	content := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextMuted)).
		Render(desc)

	hints := lipgloss.JoinHorizontal(lipgloss.Left,
		KeyHint("Esc", "Back  "),
		KeyHint("q", "Quit"),
	)

	return lipgloss.JoinVertical(lipgloss.Center,
		"",
		title,
		content,
		"",
		Divider,
		"",
		hints,
	)
}

// ── Helpers ──

func statusText(enabled bool) string {
	if enabled {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen)).Render("Enabled")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextMuted)).Render("Disabled")
}

func apiKeyDisplay(key string) string {
	if key == "" {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorWarning)).Render("Not set")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen)).Render("✓ Configured")
}
