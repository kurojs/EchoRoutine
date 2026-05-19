# Schedule Announcer

AI-powered voice and desktop notifications for your daily schedule blocks.

At each scheduled block change, the system:
1. Triggers an **AI agent** (OpenCode headless) that checks the current time and determines the active block
2. Generates a **unique motivational message** in your chosen language (fresh every time, never hardcoded)
3. Speaks it aloud via **ElevenLabs TTS** (realistic voice)
4. Shows a **KDE desktop notification**

Runs as a systemd **user timer** — starts automatically at boot, no terminal or manual action needed.

---

## How It Works

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────────┐
│  systemd timer   │────▶│  block-announcer  │────▶│  opencode run        │
│  (fires at each  │     │  (bash wrapper)   │     │  (headless AI agent) │
│  block time)     │     │                   │     │                     │
└─────────────────┘     └──────────────────┘     │  - checks time + day │
                                                  │  - reads schedule     │
                                                  │  - generates message  │
                                                  └─────────┬───────────┘
                                                            │
                                              ┌─────────────┴─────────────┐
                                              ▼                           ▼
                                   ┌──────────────────┐      ┌──────────────────┐
                                   │  ElevenLabs TTS   │      │  KDE notify-send │
                                   │  (voice output)   │      │  (popup noti)    │
                                   └──────────────────┘      └──────────────────┘
```

Key insight: **Nothing is hardcoded.** Every message is generated fresh by the AI at runtime based on the current time, day of week, and schedule context. No message pools, no templates.

---

## Prerequisites

### Required
- **KDE Plasma** (for `notify-send`)
- **OpenCode** — [opencode.ai](https://opencode.ai)
- **systemd** (user mode)

### Optional but recommended
- **elevenlabs-mcp-tts** — [github.com/kurojs/elevenlabs-mcp-tts](https://github.com/kurojs/elevenlabs-mcp-tts)
  - MCP server for ElevenLabs TTS voice output
  - Configure your API key in `~/.local/share/elevenlabs-mcp-tts/.env`:
    ```env
    ELEVENLABS_API_KEY=your_key_here
    ELEVENLABS_VOICE_ID=h3KZVBOooxHZiKRxnsdE
    ```
- **curl** and **ffplay** (for audio playback)

> **Without ElevenLabs:** the system still sends desktop notifications at each block change. Voice is optional.

---

## Installation

### Quick install

```bash
git clone https://github.com/kurojs/schedule-announcer.git
cd schedule-announcer
chmod +x install.sh
./install.sh
```

### What gets installed

| File | Destination |
|------|-------------|
| `bin/block-announcer` | `~/.local/bin/block-announcer` |
| `config/block-announcer.service` | `~/.config/systemd/user/block-announcer.service` |
| `config/block-announcer.timer` | `~/.config/systemd/user/block-announcer.timer` |
| `config/schedule.txt` | `~/.config/schedule-announcer/schedule.txt` |

### Enable the timer

```bash
systemctl --user daemon-reload
systemctl --user enable --now block-announcer.timer
```

### Enable boot autostart

```bash
sudo loginctl enable-linger $USER
```

### Verify

```bash
systemctl --user status block-announcer.timer
systemctl --user list-timers --all | grep block-announcer
```

---

## Configuration

### Schedule

Edit `~/.config/schedule-announcer/schedule.txt`:

```
08:00 - Anki kanji + Radiko
08:30 - Coding katas
09:30 - Japanese textbooks
...
```

The script parses all `HH:MM` entries automatically. **No need to update the timer** — it reads the schedule file fresh each time.

### Language

Create `~/.config/schedule-announcer/language.txt` with one of:

```
es    # Spanish (default)
en    # English
jp    # Japanese
pt    # Portuguese
fr    # French
de    # German
```

The AI agent generates messages in your chosen language.

### Voice

Edit `~/.local/share/elevenlabs-mcp-tts/.env`:

```env
ELEVENLABS_VOICE_ID=your_voice_id_here
```

---

## Testing

```bash
# Run the service directly
systemctl --user start block-announcer.service

# Check logs
journalctl --user -u block-announcer.service -n 20 --no-pager
```

---

## Fallback Behavior

If OpenCode is not installed or fails:
- A plain desktop notification is sent: *"⛩️ Schedule: HH:MM — ¡Hora del bloque!"*
- No AI-generated message, no TTS voice

If ElevenLabs TTS is not configured:
- The AI still generates the message and sends the notification
- Voice output is skipped

---

## File Structure

```
~/.config/schedule-announcer/
├── schedule.txt        # Your daily schedule
└── language.txt        # Language setting (es/en/jp/...)

~/.local/bin/
└── block-announcer     # The trigger script

~/.config/systemd/user/
├── block-announcer.service  # systemd oneshot service
└── block-announcer.timer    # systemd timer

~/.local/share/elevenlabs-mcp-tts/
└── .env                     # ElevenLabs API config (optional)
```

---

## Uninstall

```bash
systemctl --user stop block-announcer.timer
systemctl --user disable block-announcer.timer
rm -f ~/.config/systemd/user/block-announcer.*
rm -f ~/.local/bin/block-announcer
rm -rf ~/.config/schedule-announcer
systemctl --user daemon-reload
```

---

## Roadmap

- [x] AI-generated messages with OpenCode
- [x] ElevenLabs TTS integration
- [x] Multi-language support
- [x] Desktop notifications with fallback
- [ ] **schedule-tui** — Bubbletea TUI for:
  - Visual schedule editor
  - Language selector
  - ElevenLabs API key setup
  - One-click install & enable
- [ ] **AUR package** — `yay -S schedule-announcer`

---

## Built With

- **OpenCode** — AI agent runner
- **elevenlabs-mcp-tts** — ElevenLabs voice via MCP
- **systemd** — user timer scheduling
- **Bash** — lightweight trigger script
- **Go + Bubbletea** — (upcoming) TUI installer and editor

---

## License

MIT
