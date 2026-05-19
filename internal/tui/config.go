package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// configDir returns the config directory relative to the home/install location
func configDir() string {
	// Try the repo's config dir first (development), then ~/.config/echoroutine
	candidates := []string{
		filepath.Join(execDir(), "..", "..", "config"),
		filepath.Join(os.Getenv("HOME"), ".config", "echoroutine"),
		filepath.Join(os.Getenv("HOME"), ".config", "schedule-announcer"),
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	// Default to repo config
	repoConfig := filepath.Join(execDir(), "..", "..", "config")
	return repoConfig
}

func execDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

// ── Schedule ──

func loadSchedule() []ScheduleBlock {
	schedulePath := filepath.Join(configDir(), "schedule.txt")
	file, err := os.Open(schedulePath)
	if err != nil {
		// Return defaults
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
	defer file.Close()

	var blocks []ScheduleBlock
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) >= 1 {
			time := parts[0]
			label := ""
			if len(parts) >= 2 {
				label = parts[1]
			}
			blocks = append(blocks, ScheduleBlock{Time: time, Label: label})
		}
	}
	return blocks
}

func saveSchedule(blocks []ScheduleBlock) error {
	schedulePath := filepath.Join(configDir(), "schedule.txt")
	var lines []string
	lines = append(lines, "# EchoRoutine Daily Schedule")
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
	return os.WriteFile(langPath, []byte(lang+"\n"), 0644)
}

// ── API Key ──

func loadAPIKey() string {
	key := os.Getenv("ELEVENLABS_API_KEY")
	if key != "" {
		return "configured"
	}

	keyPath := filepath.Join(configDir(), "elevenlabs.key")
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return ""
	}
	key = strings.TrimSpace(string(data))
	if key != "" {
		return "configured"
	}
	return ""
}

func saveAPIKey(key string) error {
	keyPath := filepath.Join(configDir(), "elevenlabs.key")
	return os.WriteFile(keyPath, []byte(key+"\n"), 0600)
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
	servicePath := filepath.Join(os.Getenv("HOME"), ".config", "systemd", "user", "block-announcer.service")
	_, err := os.Stat(servicePath)
	return err == nil
}

func installService() (string, error) {
	installSh := filepath.Join(execDir(), "..", "..", "install.sh")
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
	return cmd.Run()
}

// ── Next block info ──

func getNextBlock(blocks []ScheduleBlock) (ScheduleBlock, string) {
	if len(blocks) == 0 {
		return ScheduleBlock{}, ""
	}

	// Simple approach: just return the first block for now
	// In a full impl we'd parse times and compare with current time
	return blocks[0], "scheduled"
}
