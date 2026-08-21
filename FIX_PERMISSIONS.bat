@echo off
:: FIX_PERMISSIONS.bat - Jalankan sebagai Administrator (Right-click -> Run as administrator)
:: Memberikan write access ke D:\AI ASSISTANT BATIQA agar Phase 1 bisa di-sync ke lokasi asli

echo Fixing permissions for D:\AI ASSISTANT BATIQA ...
icacls "D:\AI ASSISTANT BATIQA" /grant Users:(OI)(CI)F /T
if %errorlevel% equ 0 (
  echo SUCCESS: Permissions fixed.
  echo Sekarang copy file dari C:\Users\%USERNAME%\BATIQA-AI ke D:\AI ASSISTANT BATIQA
  echo Contoh: robocopy "C:\Users\%USERNAME%\BATIQA-AI" "D:\AI ASSISTANT BATIQA" /E /XD .git bin
) else (
  echo FAILED: Jalankan file ini sebagai Administrator!
)
pause
