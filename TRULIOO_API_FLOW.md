# 🚀 Fluxo Completo de Chamadas para Trulioo

## 🔄 **Quando as Chamadas são Feitas**

### 1. **Trigger da Requisição**
```bash
curl 'http://localhost/process-kyc/01K7JK5M-JNQHQ35RN-18WBH0XTB'
```

### 2. **Condição de Busca (CORRIGIDA)**
```sql
-- ❌ ANTES (Problema identificado)
WHERE package_file_id = $1 AND complete_kyc = true

-- ✅ AGORA (Corrigido)
WHERE package_file_id = $1 AND complete_kyc = false
```

**Explicação:** O sistema agora busca corretamente registros **pendentes** (não processados) ao invés de registros já processados.

## 📡 **Sequência de Chamadas para Trulioo**

Para cada registro encontrado, o sistema faz **4 chamadas sequenciais**:

### **🟢 Step 1: truliooInit()**
**Endpoint:** `GET https://api.workflow.prod.trulioo.com/interpreter-v2/test/flow/{FLOW_ID}`
- **Objetivo:** Obter mapeamento de campos disponíveis
- **Retorna:** Lista de campos obrigatórios e opcionais
- **Headers:** Cookie para sessão

### **🟡 Step 2: truliooBodySubmit()**
**Endpoint:** `POST https://api.workflow.prod.trulioo.com/interpreter-v2/test/submit/{FLOW_ID}`
- **Objetivo:** Submeter dados do cliente para processamento
- **Payload:** Dados pessoais mapeados do Step 1
- **Retorna:** `x-hf-session` header (sessão única)
- **Headers:** Content-Type: application/json

### **🔵 Step 3: truliooGenerateBearerToken()**
**Endpoint:** `POST https://auth-api.trulioo.com/connect/token`
- **Objetivo:** Gerar token OAuth2 para autenticação
- **Payload:** client_credentials grant
- **Retorna:** Bearer token
- **Headers:** Content-Type: application/x-www-form-urlencoded

### **🟣 Step 4: truliooDetailsFromClient()**
**Endpoint:** `GET https://api.workflow.prod.trulioo.com/export/test/v2/query/client/{X_HF_SESSION}?includeFullServiceDetails=true`
- **Objetivo:** Obter resultados completos do processamento KYC
- **Headers:** Authorization: Bearer {token}
- **Retorna:** Dados completos incluindo matches e watchlist

## 🔐 **Token Isolation por Sessão**

Cada job agora tem tokens únicos:

```go
type TruliooSession struct {
    XHfSession  string // Único do Step 2
    BearerToken string // Único do Step 3
    JobID       string // Identificador único
}
```

## 📊 **Dados Capturados do retorno.json**

Após o processamento, o sistema armazena:

### **1. Person Match Status**
```json
{
  "nodeTitle": "Person Match",
  "match": true,
  "serviceStatus": "COMPLETED"
}
```

### **2. Status Geral**
```json
{
  "status": "ACCEPTED"
}
```

### **3. Detalhes de Match por Campo**
```json
{
  "DatasourceFields": [
    {"FieldName": "FirstSurName", "Status": "match"},
    {"FieldName": "Address1", "Status": "match"},
    {"FieldName": "socialservice", "Status": "match"}
  ]
}
```

### **4. Watchlist Results**
```json
{
  "watchlistResults": {
    "Advanced Watchlist": {
      "watchlistStatus": "No Hit",
      "wlHitsNumber": 0
    }
  }
}
```

## 💾 **Atualização no Banco de Dados**

Após sucesso, o registro é atualizado:

```sql
UPDATE document_records
SET
  match = $1,              -- FlowData como JSON
  complete_kyc = true,     -- Marca como processado
  response = $4,           -- Response completa como JSON
  updated_at = NOW()       -- Timestamp da atualização
WHERE id = $2
```

## 🚦 **Estados do Processamento**

1. **Inicial**: `complete_kyc = false` → Aguardando processamento
2. **Em Processamento**: Worker pega o registro
3. **Concluído**: `complete_kyc = true` → Processamento finalizado

## ⚡ **Concurrent Processing**

- **5 Workers** processam simultaneamente
- **Tokens isolados** por job
- **Rate limiting** de 10 req/s
- **Timeout** de 30s por job

**Agora o sistema processa corretamente apenas registros pendentes e faz todas as 4 chamadas para Trulioo em sequência, armazenando os resultados completos no banco de dados!**