#!/usr/bin/env bash
# Bring a fresh Ubuntu VPS up to the point where prepbot can run.
# Run once, as a sudo-capable user. Idempotent where practical.
set -euo pipefail

echo "==> 0. Verify Telegram + Anthropic are reachable from this region (spec §15)"
curl -sI https://api.telegram.org  | head -1
curl -sI https://api.anthropic.com | head -1
echo "    If either hung or failed, this region blocks them — recreate the instance elsewhere."

echo "==> 1. Packages"
sudo apt-get update
sudo apt-get install -y postgresql ufw fail2ban unattended-upgrades rsync

echo "==> 2. Postgres (apt build binds to 127.0.0.1 by default)"
sudo -u postgres psql -tc "SELECT 1 FROM pg_roles WHERE rolname='prepbot'" | grep -q 1 || \
  sudo -u postgres psql -c "CREATE USER prepbot WITH PASSWORD 'CHANGE_ME';"
sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname='prepbot'" | grep -q 1 || \
  sudo -u postgres psql -c "CREATE DATABASE prepbot OWNER prepbot;"

echo "==> 3. System user + directories"
id -u prepbot >/dev/null 2>&1 || sudo useradd --system --home /opt/prepbot --shell /usr/sbin/nologin prepbot
sudo mkdir -p /opt/prepbot /var/backups
sudo chown -R prepbot:prepbot /opt/prepbot /var/backups

echo "==> 4. SSH hardening — disable password auth"
sudo sed -i 's/^#*PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo systemctl restart ssh

echo "==> 5. Firewall — OpenSSH only (long polling needs no inbound port)"
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow OpenSSH
sudo ufw --force enable

echo "==> 6. Daily backup cron (03:00, keeps the progress history off the box weekly by hand)"
echo '0 3 * * * prepbot pg_dump -Fc prepbot > /var/backups/prepbot-$(date +\%F).dump' | \
  sudo tee /etc/cron.d/prepbot-backup >/dev/null
sudo chmod 0644 /etc/cron.d/prepbot-backup

cat <<'NEXT'

==> Done. Remaining manual steps:
  1. Put secrets in /etc/prepbot.env (mode 600):
       sudo install -m 600 /dev/null /etc/prepbot.env
       sudo nano /etc/prepbot.env            # fill from .env.example
  2. Install the service unit:
       sudo cp deploy/prepbot.service /etc/systemd/system/
       sudo systemctl daemon-reload
       sudo systemctl enable --now prepbot
  3. Ship the binary + data:  make deploy REMOTE_HOST=<ip>
  4. Watch it:                journalctl -u prepbot -f
NEXT
