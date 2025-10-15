# 🔄 Guia de Alternância de Ambientes Trulioo

Este documento explica como alternar entre os ambientes de TESTE e PRODUÇÃO do Trulioo.

## 📋 Configuração Atual

O sistema agora suporta alternância dinâmica entre ambientes através da variável `TRULIOO_ENV`.

### URLs por Ambiente

#### 🧪 TESTE (test)
- **Rota 1** (Init): `https://api.workflow.prod.trulioo.com/interpreter-v2/test/flow/:flowId`
- **Rota 2** (Submit): `https://api.workflow.prod.trulioo.com/interpreter-v2/test/submit/:flowId`
- **Rota 3** (Auth): `https://auth-api.trulioo.com/connect/token` (mesmo para ambos)
- **Rota 4** (Query): `https://api.workflow.prod.trulioo.com/export/test/v2/query/client/:session`

#### 🚀 PRODUÇÃO (prod)
- **Rota 1** (Init): `https://api.workflow.prod.trulioo.com/interpreter-v2/flow/:flowId`
- **Rota 2** (Submit): `https://api.workflow.prod.trulioo.com/interpreter-v2/submit/:flowId`
- **Rota 3** (Auth): `https://auth-api.trulioo.com/connect/token` (mesmo para ambos)
- **Rota 4** (Query): `https://api.workflow.prod.trulioo.com/export/v2/query/client/:session`

## 🔧 Como Alternar Ambientes

### Método 1: Editar arquivo .env

1. **Para TESTE (padrão):**
   ```env
   TRULIOO_ENV=test
   ```

2. **Para PRODUÇÃO:**
   ```env
   TRULIOO_ENV=prod
   ```

3. **Reiniciar aplicação:**
   ```bash
   docker-compose restart app
   ```

### Método 2: Variável de ambiente temporária

1. **Para TESTE:**
   ```bash
   TRULIOO_ENV=test docker-compose up app
   ```

2. **Para PRODUÇÃO:**
   ```bash
   TRULIOO_ENV=prod docker-compose up app
   ```

### Método 3: Docker Compose override

1. **Criar docker-compose.override.yml:**
   ```yaml
   version: '3.8'
   services:
     app:
       environment:
         - TRULIOO_ENV=prod  # ou test
   ```

2. **Iniciar normalmente:**
   ```bash
   docker-compose up -d
   ```

## 📊 Verificação do Ambiente

O sistema registra nos logs qual ambiente está sendo usado:
```
Using Trulioo environment: test
```

### Verificar logs:
```bash
make logs
# ou
docker-compose logs app
```

## ⚠️ Importante

- **TESTE**: Use durante desenvolvimento e testes
- **PRODUÇÃO**: Use apenas quando tiver certeza e dados reais
- Sempre verifique os logs para confirmar o ambiente ativo
- As credenciais podem ser diferentes entre ambientes

## 🔍 Troubleshooting

### Problema: Ambiente não está mudando
**Solução:**
1. Verifique se a variável está correta no .env
2. Reinicie a aplicação: `docker-compose restart app`
3. Verifique os logs: `docker-compose logs app`

### Problema: Erro de autenticação
**Solução:**
1. Verifique se as credenciais estão corretas para o ambiente
2. Confirme o FLOW_ID correto para o ambiente
3. Verifique o X_HF_SESSION se necessário