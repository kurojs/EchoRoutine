#!/usr/bin/env bash
set -euo pipefail

# Schedule Announcer — install script
# Copies files to their proper locations and sets permissions.

BINDIR="${HOME}/.local/bin"
CONFDIR="${HOME}/.config/block-announcer"
SERVICEDIR="${HOME}/.config/systemd/user"
REPO_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "📦 Installing Schedule Announcer..."

# Bin
mkdir -p "${BINDIR}"
cp "${REPO_DIR}/bin/block-announcer" "${BINDIR}/block-announcer"
chmod +x "${BINDIR}/block-announcer"
echo "  ✓ ${BINDIR}/block-announcer"

# Config
mkdir -p "${CONFDIR}"
if [ ! -f "${CONFDIR}/schedule.txt" ]; then
    cp "${REPO_DIR}/config/schedule.txt" "${CONFDIR}/schedule.txt"
    echo "  ✓ ${CONFDIR}/schedule.txt"
else
    echo "  ∼ ${CONFDIR}/schedule.txt (already exists, skipping)"
fi

# systemd service
mkdir -p "${SERVICEDIR}"
cp "${REPO_DIR}/config/block-announcer.service" "${SERVICEDIR}/block-announcer.service"
cp "${REPO_DIR}/config/block-announcer.timer" "${SERVICEDIR}/block-announcer.timer"
echo "  ✓ ${SERVICEDIR}/block-announcer.service"
echo "  ✓ ${SERVICEDIR}/block-announcer.timer"

echo ""
echo "✅ Installation complete!"
echo ""
echo "Next steps:"
echo "  1. Edit your schedule:  nano ${CONFDIR}/schedule.txt"
echo "  2. Enable the timer:"
echo "     systemctl --user daemon-reload"
echo "     systemctl --user enable --now block-announcer.timer"
echo "  3. Enable boot autostart:"
echo "     sudo loginctl enable-linger \${USER}"
echo ""
echo "Verify with: systemctl --user status block-announcer.timer"
