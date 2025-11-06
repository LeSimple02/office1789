<?php
// Simuler exactement ce que fait Roundcube
define('INSTALL_PATH', '/var/www/html/');
define('RCUBE_INSTALL_PATH', '/var/www/html/');
define('RCUBE_CONFIG_DIR', '/var/www/html/config/');

class TestConfig {
    private $prop = array();
    
    public function load() {
        // Charger defaults.inc.php
        echo "1. Chargement defaults.inc.php\n";
        $this->load_from_file('defaults.inc.php');
        echo "   Plugins après defaults: ";
        print_r($this->prop['plugins']);
        
        // Charger config.inc.php
        echo "\n2. Chargement config.inc.php\n";
        $this->load_from_file('config.inc.php');
        echo "   Plugins après config.inc.php: ";
        print_r($this->prop['plugins']);
    }
    
    public function load_from_file($file) {
        $fpath = INSTALL_PATH . 'config/' . $file;
        if (file_exists($fpath) && is_readable($fpath)) {
            echo "   → Include $file\n";
            ob_start();
            include($fpath);
            ob_end_clean();
            
            if (isset($config) && is_array($config)) {
                echo "   → \$config trouvé dans $file\n";
                echo "   → \$config['plugins'] = ";
                var_dump($config['plugins'] ?? 'non défini');
                $this->merge($config);
            } else {
                echo "   → \$config NOT FOUND dans $file\n";
            }
        }
        return true;
    }
    
    private function merge($config) {
        $this->prop = array_merge($this->prop, $config);
    }
    
    public function get($key) {
        return $this->prop[$key] ?? null;
    }
}

$test = new TestConfig();
$test->load();

echo "\n3. RÉSULTAT FINAL:\n";
echo "   Plugins: ";
print_r($test->get('plugins'));
