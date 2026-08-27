<?php
// Wrapper config qui charge config.docker puis config.local (priorité à local)

// Charger config.docker.inc.php EN PREMIER
if (file_exists(__DIR__ . '/config.docker.inc.php')) {
    include(__DIR__ . '/config.docker.inc.php');
}

// Puis charger config.local.inc.php QUI ÉCRASE (priorité maximale)
if (file_exists(__DIR__ . '/config.local.inc.php')) {
    include(__DIR__ . '/config.local.inc.php');
}
