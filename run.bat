@echo off
setlocal EnableExtensions EnableDelayedExpansion
title BATIQA-AI Server
cd /d "%~dp0"

echo.
echo  ============================================
echo    BATIQA AI Guest Assistant - Launcher
echo  ============================================
echo.

rem -- Cek Go --
where go >nul 2>nul
if errorlevel 1 (
    echo  [ERROR] Go tidak ditemukan di PATH. Install Go 1.24+ dulu.
    pause
    exit /b 1
)

rem -- Auto-buat .env dari contoh jika belum ada --
if not exist .env (
    copy /y .env.example .env >nul
    echo  [INFO] File .env dibuat dari .env.example
)

rem -- Cek MySQL, nyalakan via Docker jika perlu --
set PORT_OK=0
powershell -NoProfile -Command "if(Test-NetConnection -ComputerName 127.0.0.1 -Port 3306 -InformationLevel Quiet -WarningAction SilentlyContinue){exit 0}else{exit 1}" >nul 2>nul
if not errorlevel 1 set PORT_OK=1

if "!PORT_OK!"=="0" (
    where docker >nul 2>nul
    if not errorlevel 1 (
        echo  [INFO] Menyalakan MySQL via Docker Compose...
        docker-compose up -d >nul 2>nul || docker compose up -d >nul 2>nul
        echo  [INFO] Menunggu MySQL siap.
        for /l %%i in (1,1,20) do (
            if "!PORT_OK!"=="0" (
                timeout /t 3 /nobreak >nul
                powershell -NoProfile -Command "if(Test-NetConnection -ComputerName 127.0.0.1 -Port 3306 -InformationLevel Quiet -WarningAction SilentlyContinue){exit 0}else{exit 1}" >nul 2>nul
                if not errorlevel 1 set PORT_OK=1
            )
        )
    ) else (
        echo  [WARN] MySQL tidak jalan dan Docker tidak ada.
        echo  [WARN] Server tetap jalan tanpa database ^(mode degraded^).
    )
)
if "!PORT_OK!"=="1" echo  [OK] MySQL siap di localhost:3306

rem -- Build --
echo  [INFO] Build aplikasi...
go build -o bin\api.exe .\cmd\api
if errorlevel 1 (
    echo  [ERROR] Build gagal. Perbaiki error di atas lalu jalankan lagi.
    pause
    exit /b 1
)

rem -- Migrasi database (jika DB tersedia) --
if "!PORT_OK!"=="1" (
    go build -o bin\migrate.exe .\cmd\migrate
    bin\migrate.exe
    if errorlevel 1 (
        echo  [WARN] Migrasi gagal - cek kredensial DB di file .env
    ) else (
        echo  [OK] Migrasi + seed selesai
    )
)

rem -- Matikan instance lama jika masih jalan --
taskkill /f /im api.exe >nul 2>nul

echo.
echo  ============================================
echo    Server : http://localhost:8080/
echo    Tamu   : http://localhost:8080/
echo    Staff  : http://localhost:8080/staff/login.html
echo    Login  : admin@batiqa.com / batiqa123
echo.
echo    Tekan Ctrl+C untuk berhenti
echo  ============================================
echo.
bin\api.exe
pause
