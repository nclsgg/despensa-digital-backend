# Deploy para AWS ECR - Despensa Digital

Este documento descreve como fazer deploy da aplicação backend para o AWS Elastic Container Registry (ECR).

## 📋 Pré-requisitos

1. **AWS CLI instalado e configurado**
   ```bash
   aws --version
   aws configure
   ```

2. **Docker instalado e rodando**
   ```bash
   docker --version
   ```

3. **Credenciais AWS configuradas** com permissões para:
   - `ecr:GetAuthorizationToken`
   - `ecr:BatchCheckLayerAvailability`
   - `ecr:GetDownloadUrlForLayer`
   - `ecr:BatchGetImage`
   - `ecr:PutImage`
   - `ecr:InitiateLayerUpload`
   - `ecr:UploadLayerPart`
   - `ecr:CompleteLayerUpload`

## 🚀 Como Fazer Deploy

### Linux / macOS

```bash
# Dar permissão de execução (apenas na primeira vez)
chmod +x deploy-ecr.sh

# Deploy com tag 'latest' (padrão)
./deploy-ecr.sh

# Deploy com tag customizada
./deploy-ecr.sh v1.0.0
```

### Windows (PowerShell)

```powershell
# Deploy com tag 'latest' (padrão)
.\deploy-ecr.ps1

# Deploy com tag customizada
.\deploy-ecr.ps1 -Tag v1.0.0
```

### Windows (Command Prompt / CMD)

```cmd
REM Deploy com tag 'latest' (padrão)
deploy-ecr.bat

REM Deploy com tag customizada
deploy-ecr.bat v1.0.0
```

## 📦 O que o Script Faz

O script automatiza 4 etapas:

1. **Login no ECR** - Autentica o Docker com o AWS ECR
2. **Build da Imagem** - Cria a imagem Docker localmente
3. **Tag da Imagem** - Adiciona a tag do ECR à imagem
4. **Push da Imagem** - Envia a imagem para o ECR

## 🔧 Configuração

As configurações estão no início de cada script:

```bash
ECR_REGISTRY="677688170820.dkr.ecr.us-east-1.amazonaws.com"
ECR_REPOSITORY="despensa-digital"
AWS_REGION="us-east-1"
IMAGE_NAME="despensa-digital"
```

Se precisar alterar alguma configuração, edite essas variáveis nos scripts.

## 📝 Deploy Manual (Alternativa)

Se preferir executar os comandos manualmente:

```bash
# 1. Login no ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 677688170820.dkr.ecr.us-east-1.amazonaws.com

# 2. Build da imagem
docker build -t despensa-digital .

# 3. Tag da imagem
docker tag despensa-digital:latest 677688170820.dkr.ecr.us-east-1.amazonaws.com/despensa-digital:latest

# 4. Push da imagem
docker push 677688170820.dkr.ecr.us-east-1.amazonaws.com/despensa-digital:latest
```

## 🏷️ Estratégia de Tags

Recomendações para tags:

- `latest` - Última versão estável
- `v1.0.0`, `v1.1.0` - Versões específicas (SemVer)
- `dev` - Versão de desenvolvimento
- `staging` - Versão de staging
- `prod` - Versão de produção

Exemplo de workflow:

```bash
# Desenvolvimento
./deploy-ecr.sh dev

# Staging
./deploy-ecr.sh staging

# Produção (com versão)
./deploy-ecr.sh v1.0.0
./deploy-ecr.sh latest
```

## 🐳 Pull da Imagem

Para baixar a imagem do ECR:

```bash
# Login no ECR (se ainda não estiver logado)
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 677688170820.dkr.ecr.us-east-1.amazonaws.com

# Pull da imagem
docker pull 677688170820.dkr.ecr.us-east-1.amazonaws.com/despensa-digital:latest

# Rodar o container
docker run -p 8080:8080 677688170820.dkr.ecr.us-east-1.amazonaws.com/despensa-digital:latest
```

## 🔍 Verificar Imagens no ECR

```bash
# Listar imagens do repositório
aws ecr list-images --repository-name despensa-digital --region us-east-1

# Descrever imagens
aws ecr describe-images --repository-name despensa-digital --region us-east-1
```

## ❗ Troubleshooting

### Erro: "no basic auth credentials"
- Certifique-se de que fez login no ECR
- Execute: `aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 677688170820.dkr.ecr.us-east-1.amazonaws.com`

### Erro: "An error occurred (AccessDeniedException)"
- Verifique se suas credenciais AWS têm as permissões necessárias
- Execute: `aws sts get-caller-identity` para verificar sua identidade

### Erro: "repository does not exist"
- Certifique-se de que o repositório existe no ECR
- Crie o repositório: `aws ecr create-repository --repository-name despensa-digital --region us-east-1`

### Docker build lento
- Use cache de build: `docker build --cache-from despensa-digital:latest -t despensa-digital .`
- Considere usar multi-stage builds (já implementado no Dockerfile)

## 🔐 Segurança

- **Nunca commite credenciais AWS** no repositório
- Use **IAM roles** em ambientes de produção
- Configure **lifecycle policies** no ECR para gerenciar imagens antigas
- Use **tags imutáveis** para versões de produção

## 📚 Recursos Adicionais

- [AWS ECR Documentation](https://docs.aws.amazon.com/ecr/)
- [Docker Documentation](https://docs.docker.com/)
- [AWS CLI Reference](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ecr/index.html)

## 📞 Suporte

Se encontrar problemas, verifique:
1. Logs do Docker: `docker logs <container-id>`
2. AWS CloudWatch Logs
3. Status do serviço ECR: https://status.aws.amazon.com/
