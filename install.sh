#!/usr/bin/env bash
set -euo pipefail

# Schedule Announcer — install script
# Copies files to their proper locations and sets permissions.

BINDIR="${HOME}/.local/bin"
CONFDIR="${HOME}/.config/schedule-announcer"
SERVICEDIR="${HOME}/.config/systemd/user"
REPO_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "📦 Installing Schedule Announcer..."
echo ""

# --- Prerequisite validation ---
echo "🔍 Checking prerequisites..."
MISSING=""

check_cmd() {
    if ! command -v "$1" &>/dev/null; then
        MISSING="${MISSING}  ❌ $1 ($2)\n"
    else
        echo "  ✅ $1"
    fi
}

check_cmd opencode "AI agent headless runner — https://opencode.ai"
check_cmd notify-send "KDE desktop notifications (libnotify)"
check_cmd systemctl "systemd user mode"

if systemctl --user list-units --type=timer &>/dev/null 2>&1; then
    echo "  ✅ systemd (user mode)"
else
    MISSING="${MISSING}  ❌ systemd user mode (timer no disponible)\n"
fi

# Check for MCP server config (optional but recommended)
MCP_DIR="${HOME}/.local/share/elevenlabs-mcp-tts"
if [ -f "${MCP_DIR}/.env" ] || [ -d "${MCP_DIR}" ]; then
    echo "  ✅ elevenlabs-mcp-tts (MCP server found)"
else
    echo "  ⚠️  elevenlabs-mcp-tts no detectado — https://github.com/kurojs/elevenlabs-mcp-tts"
    echo "     (opcional: sin TTS, solo notificaciones de escritorio)"
fi

if [ -n "${MISSING}" ]; then
    echo ""
    echo "❌ Prerequisites missing:"
    printf "${MISSING}"
    echo ""
    echo "Install missing dependencies and re-run this script."
    exit 1
fi

echo ""

# --- Install files ---
mkdir -p "${BINDIR}"
mkdir -p "${CONFDIR}"
mkdir -p "${SERVICEDIR}"

# Bin
cp "${REPO_DIR}/bin/block-announcer" "${BINDIR}/block-announcer"
chmod +x "${BINDIR}/block-announcer"
echo "  ✓ ${BINDIR}/block-announcer"

# Config
if [ ! -f "${CONFDIR}/schedule.txt" ]; then
    cp "${REPO_DIR}/config/schedule.txt" "${CONFDIR}/schedule.txt"
    echo "  ✓ ${CONFDIR}/schedule.txt"
else
    echo "  ∼ ${CONFDIR}/schedule.txt (already exists, skipping)"
fi

# systemd service
cp "${REPO_DIR}/config/block-announcer.service" "${SERVICEDIR}/block-announcer.service"
cp "${REPO_DIR}/config/block-announcer.timer" "${SERVICEDIR}/block-announcer.timer"
echo "  ✓ ${SERVICEDIR}/block-announcer.service"
echo "  ✓ ${SERVICEDIR}/block-announcer.timer"

echo ""
echo "✅ Installation complete!"
echo ""
echo "Next steps:"
echo "  1. Edit your schedule:  nano ${CONFDIR}/schedule.txt"
echo "  2. (Optional) Set language:  echo 'en' > ${CONFDIR}/language.txt"
echo "     Supported: es, en, jp, pt, fr, de"
echo "  3. Enable the timer:"
echo "     systemctl --user daemon-reload"
echo "     systemctl --user enable --now block-announcer.timer"
echo "  4. Enable boot autostart:"
echo "     sudo loginctl enable-linger \${USER}"
echo ""
echo "Verify with: systemctl --user status block-announcer.timer"
echo ""
echo "Need the elevenlabs-mcp-tts server for voice?"
echo "  → https://github.com/kurojs/elevenlabs-mcp-tts"
