<?php
/**
 * Force loading Office1789 plugins
 * Include this at the END of config.inc.php
 */

// This will be called after Roundcube loads config
// We force plugin loading directly
if (isset($RCMAIL) && is_object($RCMAIL)) {
    error_log('[FORCE_PLUGINS] Forcing plugin load...');
    $RCMAIL->plugins->load_plugin('office1789_sso', true);
    $RCMAIL->plugins->load_plugin('office1789_darkmode', true);
    $RCMAIL->plugins->load_plugin('archive', true);
    $RCMAIL->plugins->load_plugin('zipdownload', true);
    error_log('[FORCE_PLUGINS] Plugins loaded: ' . implode(', ', $RCMAIL->plugins->loaded_plugins()));
}
