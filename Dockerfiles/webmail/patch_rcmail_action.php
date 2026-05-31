<?php
/**
 * Applies the Elastic2025 sidebar fix to rcmail_action.php.
 * Moves the treetoggle div BEFORE the folder link (<a>) inside render_folder_tree_html().
 *
 * Original (RoundCube 1.6.x):
 *   html::a($link_attrib, $html_name)
 *   . (!empty($folder['folders']) ? html::div('treetoggle ...', '&nbsp;') : '')
 *
 * After patch:
 *   (!empty($folder['folders']) ? html::div('treetoggle ...', '&nbsp;') : '')
 *   . html::a($link_attrib, $html_name)
 */

$file = '/var/www/html/program/include/rcmail_action.php';
$src  = file_get_contents($file);

// The exact pattern from RoundCube 1.6.11 (may span multiple lines)
$search = <<<'SEARCH'
                html::a($link_attrib, $html_name)
                . (!empty($folder['folders']) ? html::div(
                    'treetoggle ' . ($collapsed ? 'collapsed' : 'expanded'), '&nbsp;') : '')
SEARCH;

$replace = <<<'REPLACE'
                (!empty($folder['folders']) ? html::div(
                    'treetoggle ' . ($collapsed ? 'collapsed' : 'expanded'), '&nbsp;') : '')
                . html::a($link_attrib, $html_name)
REPLACE;

if (strpos($src, $search) === false) {
    // Try a more lenient match (whitespace-normalised)
    $pattern = '/html::a\(\$link_attrib,\s*\$html_name\)\s*\.\s*\(!empty\(\$folder\[.folders.\]\)\s*\?\s*html::div\(\s*.treetoggle\s*.\s*\.\s*\(\$collapsed\s*\?\s*.collapsed.\s*:\s*.expanded.\),\s*.&nbsp;.\)\s*:\s*\'\'\)/s';
    $repl = "(!empty(\$folder['folders']) ? html::div(\n                    'treetoggle ' . (\$collapsed ? 'collapsed' : 'expanded'), '&nbsp;') : '')\n                . html::a(\$link_attrib, \$html_name)";
    $new = preg_replace($pattern, $repl, $src, 1, $count);
    if ($count === 0) {
        echo "WARNING: rcmail_action.php pattern not found – patch skipped.\n";
        exit(0);
    }
    file_put_contents($file, $new);
    echo "rcmail_action.php patched (regex match).\n";
} else {
    file_put_contents($file, str_replace($search, $replace, $src));
    echo "rcmail_action.php patched (exact match).\n";
}
