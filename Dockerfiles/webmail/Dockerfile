FROM roundcube/roundcubemail:1.6.11-fpm-alpine

# ─── Elastic2025 required patches ──────────────────────────────────────────
# Based on https://github.com/bijanbina/BCube/commit/acf9dcd53b791e8ab6bc30f8f2a5386f3fcec67f

# 1. app.js — fix widescreen column order (subject before date)
RUN APP_JS=$(find /var/www/html -name 'app.js' | head -1) && \
    if [ -n "$APP_JS" ]; then \
        sed -i "s/{className: 'subject', cells: \\['fromto', 'date', 'size', 'status', 'subject'\\]}/{className: 'subject', cells: ['fromto', 'subject', 'date', 'size', 'status']}/" "$APP_JS" && \
        echo "Patched $APP_JS"; \
    else \
        echo "WARNING: app.js not found, skipping"; \
    fi

# 2. rcmail_action.php — put treetoggle div BEFORE the folder link
COPY patch_rcmail_action.php /tmp/patch_rcmail_action.php
RUN php /tmp/patch_rcmail_action.php && rm /tmp/patch_rcmail_action.php
