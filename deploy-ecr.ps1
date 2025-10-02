# Script para build e push da imagem Docker para AWS ECR
# Usage: .\deploy-ecr.ps1 [tag]

param(
    [string]$Tag = "latest"
)

# Configurações
$ECR_REGISTRY = "677688170820.dkr.ecr.us-east-1.amazonaws.com"
$ECR_REPOSITORY = "despensa-digital"
$AWS_REGION = "us-east-1"
$IMAGE_NAME = "despensa-digital"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Despensa Digital - Deploy to AWS ECR" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Registry: $ECR_REGISTRY"
Write-Host "Repository: $ECR_REPOSITORY"
Write-Host "Tag: $Tag"
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

try {
    # 1. Fazer login no ECR
    Write-Host "📝 Step 1: Logging in to AWS ECR..." -ForegroundColor Yellow
    $ecrPassword = aws ecr get-login-password --region $AWS_REGION
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to get ECR password"
    }
    
    $ecrPassword | docker login --username AWS --password-stdin $ECR_REGISTRY
    if ($LASTEXITCODE -ne 0) {
        throw "Docker login failed"
    }
    Write-Host "✅ Login successful!" -ForegroundColor Green
    Write-Host ""

    # 2. Build da imagem
    Write-Host "🔨 Step 2: Building Docker image..." -ForegroundColor Yellow
    docker build -t "${IMAGE_NAME}:${Tag}" .
    if ($LASTEXITCODE -ne 0) {
        throw "Docker build failed"
    }
    Write-Host "✅ Build successful!" -ForegroundColor Green
    Write-Host ""

    # 3. Tag da imagem
    Write-Host "🏷️  Step 3: Tagging image..." -ForegroundColor Yellow
    docker tag "${IMAGE_NAME}:${Tag}" "${ECR_REGISTRY}/${ECR_REPOSITORY}:${Tag}"
    if ($LASTEXITCODE -ne 0) {
        throw "Docker tag failed"
    }
    Write-Host "✅ Tag successful!" -ForegroundColor Green
    Write-Host ""

    # 4. Push da imagem
    Write-Host "🚀 Step 4: Pushing image to ECR..." -ForegroundColor Yellow
    docker push "${ECR_REGISTRY}/${ECR_REPOSITORY}:${Tag}"
    if ($LASTEXITCODE -ne 0) {
        throw "Docker push failed"
    }
    Write-Host "✅ Push successful!" -ForegroundColor Green
    Write-Host ""

    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "✨ Deploy completed successfully!" -ForegroundColor Green
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Image: ${ECR_REGISTRY}/${ECR_REPOSITORY}:${Tag}"
    Write-Host ""
    Write-Host "To pull this image:" -ForegroundColor Yellow
    Write-Host "  docker pull ${ECR_REGISTRY}/${ECR_REPOSITORY}:${Tag}"
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Make sure:" -ForegroundColor Yellow
    Write-Host "  1. AWS CLI is installed and configured"
    Write-Host "  2. You have permissions to push to ECR"
    Write-Host "  3. Docker is running"
    Write-Host ""
    exit 1
}
