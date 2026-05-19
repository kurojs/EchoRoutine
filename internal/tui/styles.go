package tui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Kurox theme colors
const (
	colorBg          = "#06080F"
	colorText        = "#F3F6F9"
	colorTextMuted   = "#5C6170"
	colorPurpleLight = "#C4B5FD"
	colorPurple      = "#A78BFA"
	colorPurpleDark  = "#7C3AED"
	colorGreen       = "#86EFAC"
	colorGreenDark   = "#22C55E"
	colorBorder      = "#1E293B"
	colorBorderAct   = "#A78BFA"
	colorWarning     = "#FCD34D"
	colorError       = "#FCA5A5"
	colorSurface     = "#0D1117"
)

// Pulse animation interval
const colorPulseInterval = 80 * time.Millisecond

// PulseColors cycles through purple shades for animation
var PulseColors = []string{
	colorPurpleLight,
	colorPurple,
	colorPurpleDark,
	colorPurple,
}

func PulseColor(tick int) string {
	return PulseColors[tick%len(PulseColors)]
}

// ── Base styles ──

var BaseStyle = lipgloss.NewStyle().
	Background(lipgloss.Color(colorBg)).
	Foreground(lipgloss.Color(colorText))

var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(colorPurpleLight)).
	Padding(0, 1)

var SubtitleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorTextMuted))

var AccentStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorPurple))

var GreenStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorGreen))

var MutedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorTextMuted))

var ErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorError))

var WarningStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorWarning))

// ── Box/Border styles ──

var BoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(colorBorder)).
	Padding(1, 2).
	Width(60)

var ActiveBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(colorPurple)).
	Padding(1, 2).
	Width(60)

// ── Menu item styles ──

func MenuItem(label string, selected bool) string {
	if selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorPurpleLight)).
			Background(lipgloss.Color(colorSurface)).
			Padding(0, 2).
			Render("▸ " + label + " ◂")
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Padding(0, 2).
		Render("  " + label)
}

func MenuItemGreen(label string, selected bool) string {
	if selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorGreen)).
			Background(lipgloss.Color(colorSurface)).
			Padding(0, 2).
			Render("▸ " + label + " ◂")
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Padding(0, 2).
		Render("  " + label)
}

// ── Divider ──

var Divider = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorBorder)).
	Render("──────────────────────────────────")

// ── Key hint style ──

var KeyHintStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorTextMuted)).
	Padding(0, 1)

func KeyHint(key, desc string) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorPurple)).
		Bold(true).
		Render(key)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextMuted)).
		Render(desc)
	return keyStyle + " " + descStyle
}

// ── Status indicator ──

func StatusDot(active bool) string {
	if active {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorGreen)).
			Render("●")
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextMuted)).
		Render("○")
}

// ── Banner ──

var BannerLines = []string{
	`╔═══════════════════════════════════════════╗`,
	`║                                           ║`,
	`║     ███████╗ ██████╗██╗  ██╗ ██████╗     ║`,
	`║     ██╔════╝██╔════╝██║  ██║██╔═══██╗    ║`,
	`║     █████╗  ██║     ███████║██║   ██║    ║`,
	`║     ██╔══╝  ██║     ██╔══██║██║   ██║    ║`,
	`║     ███████╗╚██████╗██║  ██║╚██████╔╝    ║`,
	`║     ╚══════╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝     ║`,
	`║                                           ║`,
	`║          ██████╗ ██╗   ██╗                ║`,
	`║          ██╔══██╗██║   ██║                ║`,
	`║          █████╔╝██║   ██║                ║`,
	`║          ██╔══██╗██║   ██║                ║`,
	`║          ██║  ██║╚██████╔╝                ║`,
	`║          ╚═╝  ╚═╝ ╚═════╝                 ║`,
	`║                                           ║`,
	`╚═══════════════════════════════════════════╝`,
}

func RenderBanner(color string) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)
	lines := make([]string, len(BannerLines))
	for i, line := range BannerLines {
		lines[i] = style.Render(line)
	}
	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}
