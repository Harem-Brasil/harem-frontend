<?php
declare(strict_types=1);

$root = dirname(__DIR__);
if (!is_dir($root . DIRECTORY_SEPARATOR . '.git')) {
    exit(0);
}

$previous = getcwd();
if (@chdir($root) !== true) {
    exit(0);
}

@exec('git checkout HEAD -- composer.lock', $unused, $code);
@chdir($previous);

exit(0);
