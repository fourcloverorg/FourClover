@echo off

REM generate build folder
if exist build rmdir /s /q build
mkdir build

REM generate builds for all platforms

REM set environment variables
set APP_NAME=fourclover
set APP_VERSION=v010

set CGO_ENABLED=0

set GOOS=windows
set GOARCH=386
go build -o build/%APP_NAME%_%APP_VERSION%_%GOOS%_%GOARCH%.exe

set GOOS=windows
set GOARCH=amd64
go build -o build/%APP_NAME%_%APP_VERSION%_%GOOS%_%GOARCH%.exe

set GOOS=linux
set GOARCH=386
go build -o build/%APP_NAME%_%APP_VERSION%_%GOOS%_%GOARCH%

set GOOS=linux
set GOARCH=amd64
go build -o build/%APP_NAME%_%APP_VERSION%_%GOOS%_%GOARCH%

set GOOS=linux
set GOARCH=arm
go build -o build/%APP_NAME%_%APP_VERSION%_%GOOS%_%GOARCH%

set GOOS=linux
set GOARCH=arm64
go build -o build/%APP_NAME%_%APP_VERSION%_%GOOS%_%GOARCH%

set GOOS=darwin
set GOARCH=amd64
go build -o build/%APP_NAME%_%APP_VERSION%_%GOOS%_%GOARCH%


REM generate sha256 checksums for all files in build folder and save to checksums.txt
REM eg. e4f32c364d596f67c6c119718f48df85bc1cc2e2a29c78e81849763caf5dce8d, FOURCLOVER_v010_windows_amd64.exe
for /f "delims=" %%a in ('dir /b /a-d "build"') do (
  for /f "delims=" %%b in ('certutil -hashfile "build\%%a" SHA256') do (
    echo %%b, %%a >> build\checksums.txt
  )
)

REM remove "CertUtil: -hashfile command completed successfully." lines from checksums.txt
for /f "delims=" %%a in ('findstr /v "CertUtil: -hashfile command completed successfully." "build\checksums.txt"') do (
  echo %%a >> build\checksums.txt.tmp
)
move /y build\checksums.txt.tmp build\checksums.txt

REM remove "SHA256 hash of file build\%%a:, %%a" lines from checksums.txt
for /f "delims=" %%a in ('findstr /v "SHA256 hash of file build\%%a:, %%a" "build\checksums.txt"') do (
  echo %%a >> build\checksums.txt.tmp
)
move /y build\checksums.txt.tmp build\checksums.txt

echo Checksums saved to checksums.txt.

echo Builds generated successfully!
