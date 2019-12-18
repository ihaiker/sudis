#!/usr/bin/env bash

mkdir -p /etc/sudis/programs
cp bin/sudis /usr/local/bin/sudis
cp conf/sudis.toml.example /etc/sudis/sudis.toml
cp conf/logs.toml.example /etc/sudis/logs.toml

cp cmds/sudis-master.service /lib/systemd/system/
cp cmds/sudis-server.service /lib/systemd/system/
cp cmds/sudis-single.service /lib/systemd/system/

chmod +x /lib/systemd/system/sudis-*

echo "你可以执行下列命令来设置开机启动"
echo "systemctl enable sudis-master.service"
echo "systemctl enable sudis-server.service"
echo "systemctl enable sudis-single.service"
