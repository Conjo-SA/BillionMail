FROM roundcube/roundcubemail:1.6.11-fpm-alpine

# ─── Elastic2025 required patches ──────────────────────────────────────────
# Based on https://github.com/bijanbina/BCube/commit/acf9dcd53b791e8ab6bc30f8f2a5386f3fcec67f
#
# 1. app.js — fix widescreen column order (subject before date)
RUN sed -i \
    "s/{className: 'subject', cells: \['fromto', 'date', 'size', 'status', 'subject'\]}/{className: 'subject', cells: ['fromto', 'subject', 'date', 'size', 'status']}/" \
    /var/www/html/program/js/app.js

# 2. rcmail_action.php — put treetoggle div BEFORE the folder link so the
#    left-sidebar CSS selectors in Elastic2025 work correctly.
COPY patch_rcmail_action.php /tmp/patch_rcmail_action.php
RUN php /tmp/patch_rcmail_action.php && rm /tmp/patch_rcmail_action.php
