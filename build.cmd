@echo off
go build ^
-ldflags "-w -s" ^
-gcflags "all=-N -l" ^
-trimpath ^
-o auth.exe ^
-v ^
./cmd/