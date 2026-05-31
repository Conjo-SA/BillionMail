<?php
/**
 * Applies the Elastic2025 sidebar fix to rcmail_action.php.
 * Moves the treetoggle div BEFORE the folder link (<a>) inside render_folder_tree_html().
 */

// Find rcmail_action.php dynamically
$found = trim(shell_exec('find /var/www/html -name rcmail_action.php 2>/dev/null | head -1'));
if (empty($found)) {
    echo "WARNING: rcmail_action.php not found – patch skipped.\n";
    exit(0);
}
$file = $found;
echo "Patching $file\n";

$src = file_get_contents($file);

// Regex: match html::a(...) . (treetoggle ...) and swap order
$pattern = '/(html::a\(\$link_attrib,\s*\$html_name\))\s*\.\s*(\(!empty\(\$folder\[.folders.\]\)\s*\?\s*html::div\([^)]+\)\s*:\s*\'\'\))/s';

if (!preg_match($pattern, $src)) {
    echo "WARNING: pattern not found in rcmail_action.php – patch skipped.\n";
    exit(0);
}

$new = preg_replace($pattern, '$2' . "\n                . " . '$1', $src, 1);
file_put_contents($file, $new);
echo "rcmail_action.php patched successfully.\n";
