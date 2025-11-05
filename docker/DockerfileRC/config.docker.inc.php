<?php
  $config['db_dsnw'] = 'sqlite:////var/roundcube/db/sqlite.db?mode=0646';
  $config['db_dsnr'] = '';
  $config['imap_host'] = 'localhost:143';
  $config['smtp_host'] = 'localhost:587';
  $config['temp_dir'] = '/tmp/roundcube-temp';
  $config['skin'] = 'elastic';
  $config['plugins'] = array_filter(array_unique(array_merge($config['plugins'], ['office1789_sso', 'office1789_darkmode', 'archive', 'zipdownload'])));
