<?php
/**
 * Office1789 - Plugin pour injection du script auto-dark-mode
 */

class office1789_darkmode extends rcube_plugin
{
    public $task = '.*';

    function init()
    {
        $this->add_hook('render_page', array($this, 'inject_darkmode_script'));
    }

    function inject_darkmode_script($args)
    {
        // Injecter le script auto-dark-mode.js dans toutes les pages
        $script = '<script src="/skins/elastic/auto-dark-mode.js"></script>';
        
        if (isset($args['content'])) {
            // Insérer avant la fermeture du body
            $args['content'] = str_replace('</body>', $script . '</body>', $args['content']);
        }
        
        return $args;
    }
}
