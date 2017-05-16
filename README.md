mkr-check
=========

A CLI to run check commands in mackerel-agent.conf.

# Usage

```
[plugin.checks.nginx-log-age]
command = "/usr/local/bin/check-file-age --warning-age=1200 --critical-age=2400 --file=/var/log/nginx/access.log
notification_interval = 600
max_check_attempts = 2

[plugin.checks.nginx-procs]
command = "/usr/local/bin/check-procs --pattern='^nginx' --critical-under=1  --warn-over=5"
notification_interval = 60
max_check_attempts = 3
```

```shell
$ mkr-check
 OK: Procs OK: Found 1 matching processes; cmd /^nginx/
 OK: FileAge OK: /var/log/nginx/access.log is 2 seconds old (16:57:47) and 6291934 bytes.
```
