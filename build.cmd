@echo off
go build ^
-ldflags "-w -s" ^
-trimpath ^
-o auth.exe ^
-v ^
./cmd/