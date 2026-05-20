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

	// Schedule status
	status := lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Render(
		fmt.Sprintf("\n  %s  Schedule:  %s  (%d blocks)",
			StatusDot(m.IsEnabled),
			statusText(m.IsEnabled),
			len(m.Blocks),
		),
	)

	// Language
	langDisplay := m.Languages[m.SelectedLang]
	langShort := strings.SplitN(langDisplay, " - ", 2)[0]
	langStatus := MutedStyle.Render(fmt.Sprintf("  %s  Language:  %s", StatusDot(true), langShort))

	// MCP Voice status (check if elevenlabs-mcp-tts is configured)
	mcpStatus := MutedStyle.Render(fmt.Sprintf("  %s  Voice:     %s", StatusDot(isMCPConfigured()), mcpDisplay()))

	// Next block (based on real current time)
	var nextBlockStr string
	if m.IsEnabled && len(m.Blocks) > 0 {
		if b, ok := getNextBlock(m.Blocks); ok {
			nextBlockStr = GreenStyle.Render(fmt.Sprintf("\n  ◉  Next:  %s — %s", b.Time, b.Label))
		} else {
			nextBlockStr = MutedStyle.Render("\n  ◉  No upcoming blocks")
		}
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

	subtitle := AccentStyle.Render("\n  Your AI-powered daily routine voice  ")

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
		Divider,
		status,
		langStatus,
		mcpStatus,
		nextBlockStr,
		"",
		Divider,
		"",
		menu,
		"",
		Divider,
		"",
		hints,
	)

	return lipgloss.NewStyle().Padding(0, 2).Render(content)
}

// ── Schedule Editor ──

func (m Model) renderEditor() string {
	var lines []string

	title := TitleStyle.Render("Edit Schedule")
	lines = append(lines, title, "")

	if m.DeleteConfirm {
		confirmStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorWarning)).Bold(true)
		lines = append(lines, confirmStyle.Render("  Delete this block? (y/n)"))
		lines = append(lines, "")
	}

	// Block list
	for i, block := range m.Blocks {
		selected := i == m.Cursor && !m.Editing
		prefix := "  "
		if selected {
			prefix = "▸ "
		}

		var line string
		if m.Editing && i == m.EditBlock {
			// Show editing state with time picker
			hoursStr := fmt.Sprintf("%02d", m.EditHours)
			minStr := fmt.Sprintf("%02d", m.EditMinutes)

			var hoursDisplay, minDisplay, labelDisplay string

			switch m.EditFocus {
			case FocusHours:
				hoursDisplay = lipgloss.NewStyle().
					Foreground(lipgloss.Color(colorPurpleLight)).
					Background(lipgloss.Color(colorSurface)).
					Bold(true).
					Render(hoursStr)
				minDisplay = MutedStyle.Render(minStr)
			case FocusMinutes:
				hoursDisplay = AccentStyle.Render(hoursStr)
				minDisplay = lipgloss.NewStyle().
					Foreground(lipgloss.Color(colorPurpleLight)).
					Background(lipgloss.Color(colorSurface)).
					Bold(true).
					Render(minStr)
			default:
				hoursDisplay = AccentStyle.Render(hoursStr)
				minDisplay = AccentStyle.Render(minStr)
			}

			if m.EditFocus == FocusLabel {
				labelDisplay = lipgloss.NewStyle().
					Foreground(lipgloss.Color(colorPurpleLight)).
					Background(lipgloss.Color(colorSurface)).
					Bold(true).
					Render(m.EditLabel + "▌")
			} else {
				labelDisplay = MutedStyle.Render(m.EditLabel)
			}

			badge := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorTextMuted)).
				Render("[editing]")

			line = fmt.Sprintf("%s %s:%s  %s  %s", prefix, hoursDisplay, minDisplay, labelDisplay, badge)
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
	addSelected := m.Cursor == len(m.Blocks) && !m.Editing
	lines = append(lines, MenuItemGreen("+ Add Block", addSelected))

	lines = append(lines, "", Divider, "")

	// Key hints
	if m.Editing {
		var hints string
		switch m.EditFocus {
		case FocusHours:
			hints = lipgloss.JoinHorizontal(lipgloss.Left,
				KeyHint("↑↓", "Hour  "),
				KeyHint("Tab", "Minutes  "),
				KeyHint("Enter", "Save  "),
				KeyHint("Esc", "Cancel"),
			)
		case FocusMinutes:
			hints = lipgloss.JoinHorizontal(lipgloss.Left,
				KeyHint("↑↓", "Minute  "),
				KeyHint("Tab", "Label  "),
				KeyHint("Enter", "Save  "),
				KeyHint("Esc", "Cancel"),
			)
		case FocusLabel:
			hints = lipgloss.JoinHorizontal(lipgloss.Left,
				KeyHint("Type", "Label text  "),
				KeyHint("Enter", "Save  "),
				KeyHint("Esc", "Cancel"),
			)
		}
		lines = append(lines, hints)
	} else {
		hints := lipgloss.JoinHorizontal(lipgloss.Left,
			KeyHint("↑↓", "Navigate  "),
			KeyHint("Enter", "Edit  "),
			KeyHint("d", "Delete  "),
			KeyHint("Esc", "Back"),
		)
		lines = append(lines, hints)
	}

	lines = append(lines, "")
	status := fmt.Sprintf("  %d blocks — times shown in 24h format", len(m.Blocks))
	lines = append(lines, MutedStyle.Render(status))

	return lipgloss.NewStyle().Padding(0, 2).Render(strings.Join(lines, "\n"))
}

// ── Language Selector ──

func (m Model) renderLanguageSelector() string {
	var lines []string

	lines = append(lines, TitleStyle.Render("Select Language"))
	lines = append(lines, "")

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

	mcpOk := isMCPConfigured()

	settings := []struct {
		name  string
		value string
		desc  string
	}{
		{"Language", langShort, "Voice language for announcements"},
		{"ElevenLabs Voice", mcpDisplay(), "Voice engine (via OpenCode MCP)"},
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

		// Show description on selected or always for MCP
		descStr := ""
		if !mcpOk && i == 1 {
			descStr = "\n     " + WarningStyle.Render("Run: git clone https://github.com/kurojs/elevenlabs-mcp-tts")
		} else if selected {
			descStr = "\n     " + MutedStyle.Render(s.desc)
		}

		lines = append(lines, fmt.Sprintf("%s%s: %s%s", prefix, nameStyle, valueStyle, descStr))
	}

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

	// Status section with descriptions
	if m.IsInstalled {
		lines = append(lines, GreenStyle.Render("  ◉ Service installed"))
		if m.IsEnabled {
			lines = append(lines, GreenStyle.Render("  ◉ Timer active — announcements run on schedule"))
		} else {
			lines = append(lines, WarningStyle.Render("  ◉ Timer paused — no announcements until enabled"))
		}
	} else {
		lines = append(lines, MutedStyle.Render("  ◉ Service not installed"))
		lines = append(lines, MutedStyle.Render("    Install to enable automatic daily announcements"))
	}
	lines = append(lines, "")

	// Status message
	if m.StatusMsg != "" {
		msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen)).Padding(0, 2)
		if strings.Contains(m.StatusMsg, "fail") || strings.Contains(m.StatusMsg, "Failed") {
			msgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorWarning)).Padding(0, 2)
		}
		lines = append(lines, msgStyle.Render(m.StatusMsg))
		lines = append(lines, "")
	}

	// Options
	options := buildInstallOptions(m)
	for i, opt := range options {
		lines = append(lines, MenuItem(opt, m.Cursor == i))
	}

	lines = append(lines, "", Divider, "")
	lines = append(lines, MutedStyle.Render("  EchoRoutine uses systemd --user timers"))
	lines = append(lines, MutedStyle.Render("  Run at boot: loginctl enable-linger $USER"))
	lines = append(lines, "")

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
  Each time block gets announced via ElevenLabs TTS.

  Built for KDE Plasma | Arch Linux
  Part of the kurojs ecosystem.

  github.com/kurojs/EchoRoutine

 ── Dependencies ──

  OpenCode (headless AI runner)
    → github.com/opencode-ai/opencode

  elevenlabs-mcp-tts (voice engine)
    → github.com/kurojs/elevenlabs-mcp-tts

 ── Quick start ──

  echoroutine    Launch this TUI
  q / Ctrl+C     Quit
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
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen)).Render("Active")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextMuted)).Render("Paused")
}

func mcpDisplay() string {
	if isMCPConfigured() {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen)).Render("✓ Configured")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorWarning)).Render("Not configured")
}
