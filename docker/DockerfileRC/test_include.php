<?php
$config = array();
include("/var/www/html/config/config.inc.php");
echo "Plugins après include config.inc.php:\n";
print_r($config['plugins']);
