## Server

- [frp](https://github.com/fatedier/frp)

- [frp for OpenWRT](https://github.com/kuoruan/openwrt-frp)

## OpenWRT 安装 frpc

如果编译的 OpenWRT 没有包含 frp，那个可以下载编译好的 ipk 文件安装，参考[openwrt-frp](https://github.com/kuoruan/openwrt-frp/blob/master/README.md)

## 配置

- frps.ini

```bash
$ cat /etc/frp/frps.ini
[common]
bind_addr = 0.0.0.0
bind_port = 6000
bind_udp_port = 6001
kcp_bind_port = 6000
dashboard_addr = 0.0.0.0
dashboard_port = 7500
dashboard_user = admin
dashboard_pwd = tooyoungtoosimple
log_file = /var/log/frps.log
log_level = warn
log_max_days = 5
detailed_errors_to_clinet = true
authentication_method = token
token = tooyoungtoosimple
```

## 运行

```bash
$ nohup /usr/bin/frps -c /etc/frp/frps.ini &
```

## Client

Windows 以运行 frpc 转发 rdp 为例

## frpc 配置文件

```bash
[common]
server_addr = 10.5.47.31
server_port = 6000
token = tooyoungtoosimple
protocol = kcp

[desktop-rdp]
type = tcp
local_ip = 127.0.0.1
local_port = 3389
remote_port = 23389
```

## 创建 frpc 运行脚本 `fprc.vbs`

```bash
set ws=WScript.CreateObject("WScript.Shell")
ws.Run "C:\Users\dell\desktop\frp\frpc.exe -c C:\Users\dell\desktop\frp\frpc_kcp.ini",0
```

## 创建 frpc 进程监控脚本 `check_frpc.bat`

```bash
@echo off
title frpc_check
color fc
set frpc=frpc.exe
:START_CHECK

tasklist | findstr "%frpc%" || goto STARTPRO
@echo frpc is running

ping -n 10 127.0.0.1 > nul
goto START_CHECK

:STARTPRO
@echo start frpc...
C:\Users\dell\desktop\frp\frpc.exe -c C:\Users\dell\desktop\frp\frpc_kcp.ini

ping -n 10 127.0.0.1 > nul
```

## 创建 frpc 进程监控脚本的运行脚本 `check_frpc.vbs`

```bash
set ws=WScript.CreateObject("WScript.Shell")
ws.Run "C:\Users\dell\desktop\frp\check_frpc.bat /start",0
```

## 将脚本放到对应的目录

运行 `shell:startup`, 将 `fprc.vbs`， `check_frpc.vbs` 放至所打开的目录，`check_frpc.bat` 放至 `check_frpc.vbs` 所定义的目录中

## Ref

- [简单高效的内网穿透](https://lighti.me/4807.html)


