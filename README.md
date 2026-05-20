<p align="center">
  <img src="https://img.shields.io/badge/status-active-success?color=%2386EFAC&style=flat-square" alt="Status">
  <img src="https://img.shields.io/badge/platform-linux-%23A78BFA&style=flat-square" alt="Platform">
  <img src="https://img.shields.io/badge/kde-plasma-%23C4B5FD&style=flat-square" alt="KDE">
  <img src="https://img.shields.io/badge/license-MIT-%2386EFAC&style=flat-square" alt="License">
</p>

<h1 align="center" style="color: #C4B5FD; font-weight: 700;">EchoRoutine</h1>

<p align="center" style="color: #A78BFA; font-size: 1.15em;">
  <em>Your AI-powered daily routine voice assistant</em>
</p>

<p align="center">
  Each block of your day gets announced via ElevenLabs TTS with custom AI motivation.
  Because your schedule deserves a voice.
</p>

<br>

<table align="center">
  <tr>
    <td width="33%"><img src="https://i.imgur.com/XeCCdQm.png" alt="Dashboard" width="100%"></td>
    <td width="33%"><img src="https://i.imgur.com/xiLPHBW.png" alt="Schedule Editor" width="100%"></td>
    <td width="33%"><img src="https://i.imgur.com/Mz8RHpw.png" alt="Language Selector" width="100%"></td>
  </tr>
  <tr align="center">
    <td style="color: #5C6170;"><small>Dashboard animated banner</small></td>
    <td style="color: #5C6170;"><small>Schedule editor</small></td>
    <td style="color: #5C6170;"><small>Language selector</small></td>
  </tr>
</table>

<br>

---

## Features

- **AI Voice Announcements** — each block transition triggers a unique ElevenLabs TTS message
- **Visual Schedule Editor** — TUI with intuitive time picker (arrow keys, Tab to cycle fields)
- **100+ Languages** — extensive voice language selector for ElevenLabs
- **Kurox Theme** — purple (`#A78BFA`) x green (`#86EFAC`) on dark terminal
- **systemd Timer Integration** — runs headless via OpenCode, survives reboots
- **Desktop Notifications** — fallback `notify-send` when TTS is unavailable
- **TUI + Headless** — configure visually or edit files directly
- **Auto-resume** — announces the last block on boot if you missed it

---

## Installation

### From AUR (recommended)

```bash
# with your favorite AUR helper
yay -S echoroutine

# or manually
git clone https://aur.archlinux.org/echoroutine.git
cd echoroutine
makepkg -si
```

Then enable the timer:

```bash
systemctl --user daemon-reload
systemctl --user enable --now block-announcer.timer
loginctl enable-linger $USER
```

### From source

| Dependency | Why | Install |
|------------|-----|---------|
| [OpenCode](https://github.com/anomalyco/opencode) | Headless AI runner for announcements | [anomalyco/opencode](https://github.com/anomalyco/opencode) |
| [elevenlabs-mcp-tts](https://github.com/kurojs/elevenlabs-mcp-tts) | ElevenLabs TTS voice engine | `git clone` + setup |
| `notify-send` | Desktop notifications (libnotify) | `sudo pacman -S libnotify` |
| `systemd` | Timer service (user mode) | Built-in on Arch |
| `go` | Build the TUI | `sudo pacman -S go` |

```bash
git clone https://github.com/kurojs/schedule-announcer
cd schedule-announcer
./install.sh
```

### Enable at boot (any method)

```bash
systemctl --user daemon-reload
systemctl --user enable --now block-announcer.timer
loginctl enable-linger $USER
```

---

## Usage

Run `echoroutine` to launch the TUI configuration dashboard:

```
╔══════════════════════════════════════╗
║           EchoRoutine               ║    animated color pulse
║    Your AI-powered daily voice       ║
╚══════════════════════════════════════╝

  ●  Schedule:  Active  (7 blocks)
  ○  Language:  English
  ●  Voice:     Configured

  ◉  Next:  14:00 — Nihongo study

  ▸ Edit Schedule
    Settings
    Install / Manage
    About
    Exit
```

### TUI Screens

| Screen | What you can do |
|--------|----------------|
| **Dashboard** | See status at a glance, navigate to other screens |
| **Edit Schedule** | Add/edit/delete blocks. Arrows change hour/minute, Tab cycles fields |
| **Settings** | Pick language, check ElevenLabs MCP status |
| **Install / Manage** | Install service, pause/resume announcements |
| **About** | Dependencies, version info |

### Schedule Format

Edit `~/.config/echoroutine/schedule.txt` directly or use the TUI:

```
# EchoRoutine — Daily Schedule
# Format: HH:MM Label

07:00 Morning Routine
09:00 Deep Work
12:00 Lunch Break
14:00 Nihongo Study
16:00 Project Work
18:00 Exercise
20:00 Evening Wind-down
```

Times are in **24h format**. Blocks are announced in order. No overlapping times.

### Language

Pick from 100+ ElevenLabs-supported languages in the TUI Settings, or:

```bash
echo "Japanese - ja-JP" > ~/.config/echoroutine/language.txt
```

Any language string ElevenLabs supports works — English, Japanese, Korean, French, and more.

---

## Architecture

```
                   +------------------+
                   |   systemd timer  |    runs every minute
                   |  (headless mode) |
                   +--------+---------+
                            |
                   +--------v---------+
                   |  block-announcer |    bash trigger script
                   |  (Bash + flock)  |
                   +--------+---------+
                            |
              +-------------+-------------+
              |             |             |
     +--------v---+  +------v------+     |
     |  OpenCode  |  | notify-send |     |
     |  + MCP     |  | (fallback)  |     |
     +--------+---+  +-------------+     |
              |                          |
     +--------v---+              +-------v-------+
     | ElevenLabs |              |  EchoRoutine  |
     |  TTS Voice |              |  TUI (config) |
     +------------+              +---------------+
```

---

## Dependencies

### [OpenCode](https://github.com/anomalyco/opencode)
The headless AI runner that generates daily motivation. EchoRoutine injects your current block into OpenCode's context and gets back a unique, contextual announcement — not the same canned message every day.

### [elevenlabs-mcp-tts](https://github.com/kurojs/elevenlabs-mcp-tts)
Kuro's MCP server that bridges OpenCode with ElevenLabs text-to-speech. Turns AI text into natural-sounding voice in any supported language.

---

## Tips

- **Try it**: Run `block-announcer` manually to hear your first announcement
- **Edit schedule**: The TUI editor is the easiest way, but `schedule.txt` is plain text
- **Debug timer**: `systemctl --user status block-announcer.timer`
- **Logs**: `journalctl --user -u block-announcer.service -f`

---

## License

MIT (c) [kurojs](https://github.com/kurojs)
