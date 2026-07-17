@echo off
:: pqai Windows installer (CMD / Command Prompt)
:: Usage: curl -fsSL https://raw.githubusercontent.com/noaa/pqai_cli/main/install.bat -o "%TEMP%\pqai-install.bat" && "%TEMP%\pqai-install.bat"
::
:: Requires: Windows 10 build 1803+ (curl and tar are built in)

setlocal EnableDelayedExpansion

set REPO=noaa/pqai_cli
set BINARY=pqai
set ASSET=pqai-windows-amd64
set INST_DIR=%LOCALAPPDATA%\pqai

echo.
echo ^>^> Fetching latest release...

:: Query the latest tag from the GitHub API
curl -fsSL "https://api.github.com/repos/%REPO%/releases/latest" -o "%TEMP%\pqai-release.json"
if errorlevel 1 (
    echo ERROR: Failed to contact GitHub API. Check your internet connection.
    exit /b 1
)

:: Parse tag_name (pure CMD, no PowerShell needed)
for /f "tokens=2 delims=:, " %%A in (
    'findstr /i "tag_name" "%TEMP%\pqai-release.json"'
) do (
    set TAG=%%~A
    set TAG=!TAG:"=!
    goto :got_tag
)
:got_tag
del "%TEMP%\pqai-release.json" 2>nul

if "!TAG!"=="" (
    echo ERROR: Could not determine latest release tag.
    echo Please visit: https://github.com/%REPO%/releases
    exit /b 1
)

set ZIP_NAME=%ASSET%.zip
set URL=https://github.com/%REPO%/releases/download/!TAG!/%ZIP_NAME%
set TMP_DIR=%TEMP%\pqai-install-%RANDOM%

echo ^>^> Installing %BINARY% !TAG! for Windows/amd64
echo.

mkdir "%TMP_DIR%" 2>nul

echo ^>^> Downloading %URL%
curl -fsSL "%URL%" -o "%TMP_DIR%\%ZIP_NAME%"
if errorlevel 1 (
    echo ERROR: Download failed. URL: %URL%
    rmdir /s /q "%TMP_DIR%" 2>nul
    exit /b 1
)

echo ^>^> Extracting...
tar -xf "%TMP_DIR%\%ZIP_NAME%" -C "%TMP_DIR%"
if errorlevel 1 (
    echo ERROR: Extraction failed.
    rmdir /s /q "%TMP_DIR%" 2>nul
    exit /b 1
)

:: The archive contains a plain "pqai.exe" (no platform suffix)
if not exist "%INST_DIR%" mkdir "%INST_DIR%"
copy /y "%TMP_DIR%\%BINARY%.exe" "%INST_DIR%\%BINARY%.exe" >nul
if errorlevel 1 (
    echo ERROR: Could not copy binary to %INST_DIR%
    rmdir /s /q "%TMP_DIR%" 2>nul
    exit /b 1
)

rmdir /s /q "%TMP_DIR%" 2>nul

echo.
echo ^>^> Installed: %INST_DIR%\%BINARY%.exe
echo.

:: Add install dir to user PATH if not already present
echo !PATH! | findstr /i /c:"%INST_DIR%" >nul 2>&1
if errorlevel 1 (
    for /f "tokens=2*" %%A in (
        'reg query "HKCU\Environment" /v PATH 2^>nul'
    ) do set USER_PATH=%%B

    if "!USER_PATH!"=="" (
        reg add "HKCU\Environment" /v PATH /t REG_EXPAND_SZ /d "%INST_DIR%" /f >nul
    ) else (
        reg add "HKCU\Environment" /v PATH /t REG_EXPAND_SZ /d "!USER_PATH!;%INST_DIR%" /f >nul
    )

    echo ^>^> Added %INST_DIR% to your PATH.
    echo ^>^> Please open a NEW Command Prompt window to use pqai.
) else (
    echo ^>^> PATH already contains %INST_DIR%
)

echo.
echo ^>^> Done! Open a new Command Prompt and run: %BINARY% help
echo.

endlocal
