@echo off
REM Script para build e push da imagem Docker para AWS ECR
REM Usage: deploy-ecr.bat [tag]

setlocal

REM Configurações
set ECR_REGISTRY=677688170820.dkr.ecr.us-east-1.amazonaws.com
set ECR_REPOSITORY=despensa-digital
set AWS_REGION=us-east-1
set IMAGE_NAME=despensa-digital

REM Tag padrão é 'latest', mas pode ser passado como argumento
set TAG=%1
if "%TAG%"=="" set TAG=latest

echo ==========================================
echo Despensa Digital - Deploy to AWS ECR
echo ==========================================
echo Registry: %ECR_REGISTRY%
echo Repository: %ECR_REPOSITORY%
echo Tag: %TAG%
echo ==========================================
echo.

REM 1. Fazer login no ECR
echo 📝 Step 1: Logging in to AWS ECR...
for /f "tokens=*" %%i in ('aws ecr get-login-password --region %AWS_REGION%') do set ECR_PASSWORD=%%i
echo %ECR_PASSWORD% | docker login --username AWS --password-stdin %ECR_REGISTRY%
if %errorlevel% neq 0 (
    echo ❌ Login failed. Make sure AWS CLI is configured.
    exit /b 1
)
echo ✅ Login successful!
echo.

REM 2. Build da imagem
echo 🔨 Step 2: Building Docker image...
docker build -t %IMAGE_NAME%:%TAG% .
if %errorlevel% neq 0 (
    echo ❌ Build failed.
    exit /b 1
)
echo ✅ Build successful!
echo.

REM 3. Tag da imagem
echo 🏷️  Step 3: Tagging image...
docker tag %IMAGE_NAME%:%TAG% %ECR_REGISTRY%/%ECR_REPOSITORY%:%TAG%
if %errorlevel% neq 0 (
    echo ❌ Tag failed.
    exit /b 1
)
echo ✅ Tag successful!
echo.

REM 4. Push da imagem
echo 🚀 Step 4: Pushing image to ECR...
docker push %ECR_REGISTRY%/%ECR_REPOSITORY%:%TAG%
if %errorlevel% neq 0 (
    echo ❌ Push failed.
    exit /b 1
)
echo ✅ Push successful!
echo.

echo ==========================================
echo ✨ Deploy completed successfully!
echo ==========================================
echo Image: %ECR_REGISTRY%/%ECR_REPOSITORY%:%TAG%
echo.
echo To pull this image:
echo   docker pull %ECR_REGISTRY%/%ECR_REPOSITORY%:%TAG%
echo.

endlocal
