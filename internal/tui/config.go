package tui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ── Config path detection ──

func configDir() string {
	// Priority:
	// 1. ~/.config/echoroutine/ (new standard)
	// 2. Repo relative (dev mode — looks for config/ next to bin/)
	// 3. ~/.config/schedule-announcer/ (legacy)

	home := os.Getenv("HOME")

	candidates := []string{
		filepath.Join(home, ".config", "echoroutine"),
		filepath.Join(devConfigDir()),
		filepath.Join(home, ".config", "schedule-announcer"),
	}

	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}

	// If none exist, default to ~/.config/echoroutine/ (create later if needed)
	echoroutineDir := filepath.Join(home, ".config", "echoroutine")
	os.MkdirAll(echoroutineDir, 0755)
	return echoroutineDir
}

func devConfigDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	// Walk up from bin/ to find config/
	dir := filepath.Dir(exe)
	for i := 0; i < 4; i++ {
		cfg := filepath.Join(dir, "config")
		if info, err := os.Stat(cfg); err == nil && info.IsDir() {
			return cfg
		}
		dir = filepath.Dir(dir)
	}
	return ""
}

// ── MCP / ElevenLabs detection ──

func isMCPConfigured() bool {
	// Method 1: Check if elevenlabs-mcp-tts source directory exists
	mcpDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "elevenlabs-mcp-tts")
	if info, err := os.Stat(mcpDir); err == nil && info.IsDir() {
		return true
	}

	// Method 2: Check OpenCode config for elevenlabs-tts MCP entry
	opencodeConfig := filepath.Join(os.Getenv("HOME"), ".config", "opencode", "opencode.jsonc")
	if data, err := os.ReadFile(opencodeConfig); err == nil {
		content := string(data)
		if strings.Contains(content, "elevenlabs-tts") || strings.Contains(content, "elevenlabs") {
			return true
		}
	}

	// Method 3: Check opencode.json
	opencodeConfig2 := filepath.Join(os.Getenv("HOME"), ".config", "opencode", "opencode.json")
	if data, err := os.ReadFile(opencodeConfig2); err == nil {
		var cfg struct {
			MCP map[string]interface{} `json:"mcp"`
		}
		if json.Unmarshal(data, &cfg) == nil {
			for name := range cfg.MCP {
				if strings.Contains(name, "elevenlabs") {
					return true
				}
			}
		}
		// Also check raw string
		content := string(data)
		if strings.Contains(content, "elevenlabs-tts") || strings.Contains(content, "elevenlabs") {
			return true
		}
	}

	return false
}

// ── Schedule ──

func loadSchedule() []ScheduleBlock {
	schedulePath := filepath.Join(configDir(), "schedule.txt")
	file, err := os.Open(schedulePath)
	if err != nil {
		return defaultSchedule()
	}
	defer file.Close()

	var blocks []ScheduleBlock
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) >= 1 && len(parts[0]) == 5 && parts[0][2] == ':' {
			time := parts[0]
			label := ""
			if len(parts) >= 2 {
				label = strings.TrimSpace(parts[1])
			}
			blocks = append(blocks, ScheduleBlock{Time: time, Label: label})
		}
	}
	if len(blocks) == 0 {
		return defaultSchedule()
	}
	return blocks
}

func defaultSchedule() []ScheduleBlock {
	return []ScheduleBlock{
		{Time: "07:00", Label: "Morning Routine"},
		{Time: "09:00", Label: "Deep Work"},
		{Time: "12:00", Label: "Lunch Break"},
		{Time: "14:00", Label: "日本語の勉強"},
		{Time: "16:00", Label: "Project Work"},
		{Time: "18:00", Label: "Exercise"},
		{Time: "20:00", Label: "Evening Wind-down"},
	}
}

func saveSchedule(blocks []ScheduleBlock) error {
	schedulePath := filepath.Join(configDir(), "schedule.txt")
	// Ensure config dir exists
	os.MkdirAll(filepath.Dir(schedulePath), 0755)

	var lines []string
	lines = append(lines, "# EchoRoutine — Daily Schedule")
	lines = append(lines, "# Format: HH:MM Label")
	lines = append(lines, "")
	for _, b := range blocks {
		lines = append(lines, fmt.Sprintf("%s %s", b.Time, b.Label))
	}
	return os.WriteFile(schedulePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

// ── Language ──

func loadLanguage() string {
	langPath := filepath.Join(configDir(), "language.txt")
	data, err := os.ReadFile(langPath)
	if err != nil {
		return "English (US) - en-US"
	}
	return strings.TrimSpace(string(data))
}

func saveLanguage(lang string) error {
	langPath := filepath.Join(configDir(), "language.txt")
	os.MkdirAll(filepath.Dir(langPath), 0755)
	return os.WriteFile(langPath, []byte(lang+"\n"), 0644)
}

// ── Systemd status ──

func isTimerEnabled() bool {
	cmd := exec.Command("systemctl", "--user", "is-enabled", "block-announcer.timer")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "enabled"
}

func isServiceInstalled() bool {
	home := os.Getenv("HOME")
	candidates := []string{
		filepath.Join(home, ".config", "systemd", "user", "block-announcer.service"),
		"/usr/lib/systemd/user/block-announcer.service",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

func installService() (string, error) {
	// Use install.sh — try repo path first, then relative to binary
	home := os.Getenv("HOME")

	candidates := []string{
		filepath.Join(devConfigDir(), "..", "install.sh"),
		filepath.Join(devConfigDir(), "..", "..", "install.sh"),
		filepath.Join(home, ".local", "bin", "install.sh"),
	}

	var installSh string
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			installSh = c
			break
		}
	}

	if installSh == "" {
		return "", fmt.Errorf("install.sh not found — re-clone the repository or run install.sh manually")
	}

	cmd := exec.Command("bash", installSh)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("install failed: %w", err)
	}
	return string(out), nil
}

func toggleTimer(enable bool) error {
	action := "enable"
	if !enable {
		action = "disable"
	}
	cmd := exec.Command("systemctl", "--user", action, "block-announcer.timer")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("systemctl %s failed: %w\n%s", action, err, string(out))
	}
	return nil
}

// ── Next block ──

func getNextBlock(blocks []ScheduleBlock) (ScheduleBlock, bool) {
	if len(blocks) == 0 {
		return ScheduleBlock{}, false
	}

	now := time.Now()
	currentMin := now.Hour()*60 + now.Minute()

	for _, b := range blocks {
		h, m := parseTime(b.Time)
		blockMin := h*60 + m
		if blockMin > currentMin {
			return b, true
		}
	}
	// All passed for today — wrap to first block tomorrow
	return blocks[0], true
}
