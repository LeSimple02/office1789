#!/bin/sh
# Writes runtime config file for the SPA so we don't need to rebuild on env changes
set -e

# Default values if not provided
: ${BACKEND_URL:="https://backend.office1789.com"}
: ${MAIL_URL:="https://mail.office1789.com"}
: ${CHAT_URL:="https://chat.office1789.com"}
: ${DOCS_URL:="https://docs.office1789.com"}

CONFIG_FILE="/usr/share/nginx/html/app-config.js"
cat > "$CONFIG_FILE" <<-EOF
window.__APP_CONFIG__ = {
  "VITE_APP_API": "${BACKEND_URL}",
  "VITE_APP_MAIL_URL": "${MAIL_URL}",
  "VITE_APP_CHAT_URL": "${CHAT_URL}",
  "VITE_APP_DOCS_URL": "${DOCS_URL}"
};
EOF

# Ensure file permissions
chmod 644 "$CONFIG_FILE"

# Execute the passed command (nginx)
exec "$@"
