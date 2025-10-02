#!/bin/bash

# Script para build e push da imagem Docker para AWS ECR
# Usage: ./deploy-ecr.sh [tag]

set -e

# Configurações
ECR_REGISTRY="677688170820.dkr.ecr.us-east-1.amazonaws.com"
ECR_REPOSITORY="despensa-digital"
AWS_REGION="us-east-1"
IMAGE_NAME="despensa-digital"

# Tag padrão é 'latest', mas pode ser passado como argumento
TAG="${1:-latest}"

echo "=========================================="
echo "Despensa Digital - Deploy to AWS ECR"
echo "=========================================="
echo "Registry: $ECR_REGISTRY"
echo "Repository: $ECR_REPOSITORY"
echo "Tag: $TAG"
echo "=========================================="
echo ""

# 1. Fazer login no ECR
echo "📝 Step 1: Logging in to AWS ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REGISTRY
if [ $? -eq 0 ]; then
    echo "✅ Login successful!"
else
    echo "❌ Login failed. Make sure AWS CLI is configured."
    exit 1
fi
echo ""

# 2. Build da imagem
echo "🔨 Step 2: Building Docker image..."
docker build -t $IMAGE_NAME:$TAG .
if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed."
    exit 1
fi
echo ""

# 3. Tag da imagem
echo "🏷️  Step 3: Tagging image..."
docker tag $IMAGE_NAME:$TAG $ECR_REGISTRY/$ECR_REPOSITORY:$TAG
if [ $? -eq 0 ]; then
    echo "✅ Tag successful!"
else
    echo "❌ Tag failed."
    exit 1
fi
echo ""

# 4. Push da imagem
echo "🚀 Step 4: Pushing image to ECR..."
docker push $ECR_REGISTRY/$ECR_REPOSITORY:$TAG
if [ $? -eq 0 ]; then
    echo "✅ Push successful!"
else
    echo "❌ Push failed."
    exit 1
fi
echo ""

echo "=========================================="
echo "✨ Deploy completed successfully!"
echo "=========================================="
echo "Image: $ECR_REGISTRY/$ECR_REPOSITORY:$TAG"
echo ""
echo "To pull this image:"
echo "  docker pull $ECR_REGISTRY/$ECR_REPOSITORY:$TAG"
echo ""
