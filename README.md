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

# LICENSE

Copyright 2017 TSUBOUCHI, Yuuki <yuki.tsubo@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License"): you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
