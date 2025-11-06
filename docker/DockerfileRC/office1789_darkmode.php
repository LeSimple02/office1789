<?php
/**
 * Plugin Office1789 Dark Mode
 * Injecte auto-dark-mode.js dans toutes les pages Roundcube
 */
class office1789_darkmode extends rcube_plugin
{
    public $task = '.*';

    function init()
    {
        $this->add_hook('render_page', array($this, 'inject_dark_mode_script'));
    }

    function inject_dark_mode_script($args)
    {
        // Injecter auto-dark-mode.js avant </body>
        $script = '<script src="/skins/elastic/auto-dark-mode.js"></script>';
        $args['content'] = str_replace('</body>', $script . '</body>', $args['content']);
        
        return $args;
    }
}
