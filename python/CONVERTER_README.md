# 🔄 KYC Data Converter

Converte dados de usuários em formato tabular para JSON compatible com a API Trulioo e gera PDFs automaticamente.

## Arquivos Principais

- `convert_user_data.py` - Converte dados tabulares para JSON
- `process_kyc_data.py` - Script completo com geração de PDF
- `relatori_kyc` - Arquivo de entrada com dados dos usuários

## Formato do Arquivo de Entrada

O arquivo `relatori_kyc` deve conter uma linha por usuário com campos separados por TAB:

```
TA_Responsible	Type	Email	UserID	FirstName	MiddleName	LastName	DOB	Phone	Address	City	Postal	State	Country	SSN
```

### Exemplo:
```
Marie-Lou	Merge	missinglimbranch@mac.com	648486	Russell		Martin II	1962-07-23	3252421782	3838 N FM 644	Loraine	79532	TX	US	452211429
Wei	Transfer	john.doe@email.com	123456	John	Michael	Doe	1985-03-15	5551234567	123 Main St	New York	10001	NY	US	123456789
```

## Como Usar

### 1. Conversão Simples (apenas JSON)

```bash
# Ativar ambiente virtual
source venv/bin/activate

# Converter dados para JSON
python convert_user_data.py
```

**Saída:**
- `request.json` - Para um usuário
- `request_1.json`, `request_2.json`, etc. - Para múltiplos usuários
- `request_combined.json` - Todos os usuários em um array

### 2. Conversão + Geração de PDF

```bash
# Ativar ambiente virtual
source venv/bin/activate

# Converter e gerar PDF (precisa de response.json)
python process_kyc_data.py --pdf --response response.json --output relatorio.pdf
```

### 3. Opções Avançadas

```bash
# Especificar arquivo de entrada diferente
python convert_user_data.py --input meus_dados.txt

# Gerar PDF com nome específico
python process_kyc_data.py --pdf --response response.json --output "Russell_Martin_Report.pdf"
```

## Campos Suportados

| Campo Original | Mapeamento JSON | Obrigatório |
|---|---|---|
| TA_Responsible | Comentário | Não |
| Type | Comentário | Não |
| Email | Comentário | Não |
| UserID | Client Reference | Não |
| FirstName | `6716b75a1287d277472c8d82` | Sim |
| MiddleName | `674dcb7ce686813288f9045f` | Não |
| LastName | `6716b75a1287d277472c8d83` | Sim |
| DOB | `6716b75a1287d277472c8d84` | Sim |
| Phone | Comentário | Não |
| Address | `6716b75a1287d277472c8d86` | Sim |
| City | `6716b75a1287d277472c8d8c` | Sim |
| Postal | `6716b75a1287d277472c8d87` | Sim |
| State | `6716b75a1287d277472c8d88` | Sim |
| Country | `6716b75a1287d277472c8d81` | Sim |
| SSN | `6744facf99661447b4b58ff7` | Não |

## Formato de Saída JSON

```json
{
    // TA: Marie-Lou
    // Type: Merge
    // Email: missinglimbranch@mac.com
    // UserId: 648486
    // Name: Russell Martin II
    // DOB: 1962-07-23
    // Phone: 3252421782
    // Address: 3838 N FM 644, Loraine, TX, 79532, US
    // SSN: 452211429
    "67228aef1e5e2108d84020a2": "manual-2025-10-15-648486",
    "6716b75a1287d277472c8d82": "Russell",
    "6716b75a1287d277472c8d83": "Martin II",
    "6716b75a1287d277472c8d84": "1962-07-23",
    "6716b75a1287d277472c8d86": "3838 N FM 644",
    "6716b75a1287d277472c8d8c": "Loraine",
    "6716b75a1287d277472c8d88": "TX",
    "6716b75a1287d277472c8d87": "79532",
    "6716b75a1287d277472c8d81": "US",
    "6744facf99661447b4b58ff7": "452211429"
}
```

## Recursos

✅ **Conversão automática de formatos de data**
✅ **Geração de Client Reference ID único**
✅ **Comentários informativos no JSON**
✅ **Suporte a múltiplos usuários**
✅ **Integração com gerador de PDF**
✅ **Validação de campos obrigatórios**
✅ **Suporte a campos opcionais**

## Troubleshooting

### Problema: "Nenhum usuário encontrado"
- Verifique se o arquivo `relatori_kyc` existe
- Confirme que os dados estão separados por TAB
- Certifique-se de que há pelo menos os campos obrigatórios

### Problema: "Campos malformados"
- Verifique se todos os campos estão presentes
- Use TAB como separador (não espaços)
- Confirme o formato da data (YYYY-MM-DD ou DD-MM-YYYY)

### Problema: "PDF não gerado"
- Verifique se o arquivo `response.json` existe
- Confirme que o ambiente virtual está ativo
- Certifique-se de que todas as dependências estão instaladas

## Dependências

```bash
pip install reportlab python-dateutil
```