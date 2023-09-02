chcp 65001
@echo off
:loop
@echo off&amp;color 0A
cls
echo,
echo 请选择要编译的系统环境：
echo,
echo 1. Windows_amd64
echo 2. linux_amd64

set/p action=请选择:
if %action% == 1 goto build_Windows_amd64
if %action% == 2 goto build_linux_amd64

:build_Windows_amd64
echo 编译Windows版本64位
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o evn_api/target/evn_api.exe evn_api/main.go
go build -o evn_article/target/evn_article.exe evn_article/main.go
go build -o evn_other/target/evn_other.exe evn_other/main.go
go build -o evn_user/target/evn_user.exe evn_user/main.go
go build -o evn_video/target/evn_video.exe evn_video/main.go
go build -o evn_ws/target/evn_ws.exe evn_ws/main.go
:build_linux_amd64
echo 编译Linux版本64位
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o evn_api/target/evn_api evn_api/main.go
go build -o evn_article/target/evn_article evn_article/main.go
go build -o evn_other/target/evn_other evn_other/main.go
go build -o evn_user/target/evn_user evn_user/main.go
go build -o evn_video/target/evn_video evn_video/main.go
go build -o evn_ws/target/evn_ws evn_ws/main.go