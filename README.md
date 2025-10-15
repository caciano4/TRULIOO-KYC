# KORE Trulioo KYC

Sistema de processamento KYC (Know Your Customer) integrado com Trulioo para verificação de identidade.

## 📋 Sobre o Projeto

Este projeto é uma aplicação Go que gerencia processos de KYC através da integração com a API do Trulioo. O sistema permite upload de documentos, processamento de verificações de identidade e armazenamento de resultados.

## 🚀 Tecnologias

- **Go 1.23.4** - Linguagem principal
- **PostgreSQL** - Banco de dados
- **Docker & Docker Compose** - Containerização
- **Air** - Hot reload para desenvolvimento
- **PgAdmin** - Interface de administração do banco

## 📁 Estrutura do Projeto

```
├── cmd/                    # Arquivo principal da aplicação
├── config/                 # Configurações (database, logs, env)
├── controllers/            # Controladores HTTP
├── middleware/             # Middlewares da aplicação
├── models/                 # Modelos de dados
├── resources/              # Resources/handlers
├── routes/                 # Definição de rotas
├── services/               # Lógica de negócio
├── migrations/             # Migrações do banco de dados
├── utils/                  # Utilitários
├── validations/            # Validações
├── views/                  # Templates/views
├── static/                 # Arquivos estáticos
├── db/                     # Dados do PostgreSQL
├── BRUNO_API/              # Coleção de APIs para testes
└── docker-compose.yml      # Configuração Docker
```

## ⚙️ Pré-requisitos

- Docker e Docker Compose
- Go 1.23.4+ (para desenvolvimento local)
- Make (opcional, para usar comandos do Makefile)

## 🔧 Configuração e Instalação

### Início Rápido com Make

```bash
# Inicia todo o projeto (setup + docker + migrações)
make start

# Ou passo a passo:
make setup        # Configura ambiente
make docker-up    # Inicia serviços Docker
make migrate-up   # Executa migrações
```

### Configuração Manual

1. **Clone o repositório**
   ```bash
   git clone <repository-url>
   cd KORE-TRULIOO
   ```

2. **Configure as variáveis de ambiente**
   ```bash
   cp .env.example .env  # Se o arquivo existir
   # Edite o .env com suas configurações
   ```

   **⚠️ Configuração do Trulioo (Importante):**
   ```env
   # Credenciais Trulioo
   CLIENT_ID=seu_client_id_aqui
   CLIENT_SECRET=seu_client_secret_aqui
   FLOW_ID=seu_flow_id_aqui
   X_HF_SESSION=sua_session_aqui

   # Ambiente Trulioo (test/prod)
   TRULIOO_ENV=test    # Para ambiente de teste
   # TRULIOO_ENV=prod  # Para ambiente de produção
   ```

   **Como alternar entre ambientes:**
   - Para **TESTE**: `TRULIOO_ENV=test` (padrão)
   - Para **PRODUÇÃO**: `TRULIOO_ENV=prod`

3. **Inicie os serviços**
   ```bash
   docker-compose up -d
   ```

4. **Execute as migrações**
   ```bash
   MIGRATION_ACTION=up docker-compose up migrate
   ```

## 🏃‍♂️ Como Executar

### Desenvolvimento com Hot Reload
```bash
make dev
# ou
air
```

### Produção
```bash
make build
make run
```

### Docker
```bash
make docker-up
```

## 📝 Comandos Úteis

```bash
# Ver todos os comandos disponíveis
make help

# Desenvolvimento
make dev          # Inicia com hot reload
make build        # Compila a aplicação
make run          # Executa aplicação compilada
make test         # Executa testes

# Docker
make docker-up    # Inicia todos os serviços
make docker-down  # Para todos os serviços
make logs         # Mostra logs da aplicação

# Banco de dados
make migrate-up   # Executa migrações
make migrate-down # Reverte migrações

# Limpeza
make clean        # Remove arquivos temporários
```

## 🌐 Endpoints de Acesso

- **Aplicação**: http://localhost:80
- **PgAdmin**: http://localhost:8080
- **PostgreSQL**: localhost:5432

## 🔄 Funcionalidades Principais

### Processamento KYC
- Upload de documentos
- Integração com API Trulioo
- Processamento em 4 etapas
- Armazenamento de respostas
- Logs detalhados de requisições

### Gestão de Registros
- CRUD de registros KYC
- Armazenamento de documentos
- Controle de status de verificação
- Histórico de tentativas

### Campos Suportados
- Driver License Number
- Driver License Version Number
- Voter ID
- Passport Number
- Dados pessoais e endereço

## 🗃️ Banco de Dados

O projeto usa PostgreSQL com as seguintes tabelas principais:
- `records` - Dados dos registros KYC
- `document_records` - Documentos associados
- Migrações automáticas via Docker

## 🧪 Testes

Para executar testes:
```bash
make test
# ou
go test ./...
```

## 📊 API Testing

Use a coleção Bruno API localizada em `BRUNO_API/` para testar os endpoints.

## 🐛 Troubleshooting

### Problemas Comuns

1. **Erro de conexão com banco**
   - Verifique se o PostgreSQL está rodando: `make logs`
   - Confirme as variáveis de ambiente no `.env`

2. **Migrações falhando**
   - Pare os containers: `make docker-down`
   - Remova volumes: `docker-compose down -v`
   - Inicie novamente: `make start`

3. **Aplicação não inicia**
   - Verifique logs: `make logs`
   - Confirme se todas as dependências estão instaladas

## 📈 Melhorias Implementadas

### Últimas Atualizações
- ✅ Novos campos no banco (DriverLicense, Passport, VoterID, etc.)
- ✅ Armazenamento de respostas das etapas 1-4
- ✅ Sistema de logs melhorado
- ✅ Tratamento de erros aprimorado
- ✅ Correção do sistema de Bearer Token
- ✅ Estrutura de rotas otimizada

### Próximas Melhorias
- [ ] Processamento assíncrono de requests KYC
- [ ] Contador de falhas por registro
- [ ] Interface para download de respostas
- [ ] Botão para requisição manual de KYC
- [ ] Coluna para identificar requests falhadas

## 🤝 Contribuição

1. Faça fork do projeto
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Push para a branch
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob licença [inserir licença aqui].

---

**Desenvolvido por:** KORE Team
**Última atualização:** Outubro 2024