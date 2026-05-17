# Schedule Announcer

AI-powered voice and desktop notifications for your daily schedule blocks.

At each scheduled block change, the system:
1. Triggers an **AI agent** (OpenCode headless) that checks the current time and determines the active block
2. Generates a **unique motivational message** in Spanish (fresh every time, never hardcoded)
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
- **GitHub Copilot** or other OpenCode-compatible AI provider
- **curl** and **ffplay** (for ElevenLabs TTS audio playback)
- **systemd** (user mode)

### Optional but recommended
- **ElevenLabs API key** configured as an OpenCode MCP server for voice output

---

## Installation

### 1. Clone this repository

```bash
git clone https://github.com/kurojs/schedule-announcer.git
cd schedule-announcer
```

### 2. Run the install script

```bash
chmod +x install.sh
./install.sh
```

This copies:
- `bin/block-announcer` → `~/.local/bin/block-announcer`
- `config/block-announcer.service` → `~/.config/systemd/user/block-announcer.service`
- `config/block-announcer.timer` → `~/.config/systemd/user/block-announcer.timer`
- `config/schedule.txt` → `~/.config/block-announcer/schedule.txt`

### 3. Customize your schedule

Edit `~/.config/block-announcer/schedule.txt` with your own daily schedule.
The format is:

```
08:00 - Activity name (description)
08:30 - Another activity
...
```

### 4. Enable the timer

```bash
systemctl --user daemon-reload
systemctl --user enable block-announcer.timer
systemctl --user start block-announcer.timer
```

### 5. Enable lingering (for boot autostart)

For the timer to start automatically when the computer boots (without needing to log into the desktop first):

```bash
sudo loginctl enable-linger $USER
```

After this, reboot or log out/in and the timer will be active perpetually.

### 6. Verify it's working

```bash
systemctl --user status block-announcer.timer
systemctl --user list-timers --all | grep block-announcer
```

You should see the next scheduled trigger time.

---

## Testing

To test the system at any time (even between blocks):

```bash
# Run the service directly (will only fire if current time matches a block)
systemctl --user start block-announcer.service

# Or check the logs
journalctl --user -u block-announcer.service -n 20 --no-pager
```

---

## How the AI Generates Messages

The system does NOT use hardcoded messages. Instead:

1. The systemd timer fires at each block time (e.g., 08:00, 09:30, etc.)
2. A bash script launches OpenCode in headless mode with the schedule file as context
3. The AI agent:
   - Checks the actual system time with `date`
   - Determines which schedule block is active
   - Generates a unique motivational message in Spanish
   - Calls ElevenLabs TTS to speak it
   - Sends a KDE notification
4. The session ends

This means **every message is unique** — generated fresh each time based on the AI's creativity.

---

## Customization

### Voice

To change the ElevenLabs voice, edit `~/.local/share/elevenlabs-mcp-tts/.env`:

```env
ELEVENLABS_VOICE_ID=your_voice_id_here
```

### Language

The system language can be changed by forking this repo and updating the prompt instructions in `bin/block-announcer`. Change the phrase "en ESPAÑOL" to your preferred language.

### Schedule

Your schedule is at `~/.config/block-announcer/schedule.txt`. Edit freely — the AI reads it fresh on each trigger.

---

## File Structure

```
~/.config/block-announcer/
├── schedule.txt        # Your daily schedule

~/.local/bin/
└── block-announcer     # The trigger script

~/.config/systemd/user/
├── block-announcer.service  # systemd oneshot service
└── block-announcer.timer    # systemd timer (11 daily triggers)
```

---

## Uninstall

```bash
systemctl --user stop block-announcer.timer
systemctl --user disable block-announcer.timer
rm -f ~/.config/systemd/user/block-announcer.*
rm -f ~/.local/bin/block-announcer
rm -rf ~/.config/block-announcer
systemctl --user daemon-reload
```

---

## License

MIT
