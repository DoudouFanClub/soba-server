@echo off
:: Launch PowerShell and run Ubuntu, then run redis-cli in the same session
start powershell.exe -NoExit -Command "ubuntu run bash -c 'redis-cli'"

:: Launch Command Prompt and run mongod
start cmd /k "mongod"

:: Wait a moment to ensure both Services start
timeout /t 5 /nobreak

:: Wait for the user to press a key before closing the script
pause