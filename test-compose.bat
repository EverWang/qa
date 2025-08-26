@echo off
chcp 65001 >nul
echo Testing modified Docker Compose configuration...
echo.

echo Switching to WSL environment...
wsl -d Ubuntu -e bash -c "cd /mnt/d/GITVIEW/qa && echo 'Current directory:' && pwd"

echo.
echo Stopping existing containers...
wsl -d Ubuntu -e bash -c "cd /mnt/d/GITVIEW/qa && docker compose -f docker-compose.local.yml down"

echo.
echo Starting containers with new configuration...
wsl -d Ubuntu -e bash -c "cd /mnt/d/GITVIEW/qa && docker compose -f docker-compose.local.yml up -d"

echo.
echo Checking container status...
wsl -d Ubuntu -e bash -c "cd /mnt/d/GITVIEW/qa && docker compose -f docker-compose.local.yml ps"

echo.
echo Test completed!
echo Access URLs:
echo - Backend API: http://localhost:8080
echo - Admin Web: http://localhost:3000
echo - Miniprogram: http://localhost:3001
echo - Nginx Proxy: http://localhost:8081
echo.
pause