# KYC PDF Report Generator

Este projeto Python gera relatórios PDF de KYC a partir de dados JSON de request e response da API Trulioo.

## Três Versões Disponíveis

### 1. **professional_pdf_generator.py** (⭐ RECOMENDADO ⭐)
- **Visual Profissional**: Corresponde exatamente ao template HTML existente
- **Estrutura Correta**: Seções organizadas como Application Details, Screening Details, Processing Steps
- **Styling Avançado**: Cores, fontes e layout idênticos ao template
- **Números Grandes**: Screening results com números de destaque (32px)
- **Status Colorido**: ACCEPTED (verde) / REJECTED (vermelho)
- **Tabelas Profissionais**: Com bordas, cores alternadas e alinhamento perfeito

### 2. **simple_pdf_generator.py** (Básico)
- **Fácil instalação**: Usa apenas ReportLab
- **PDF simples**: Layout básico mas funcional
- **Todas as informações**: Dados completos mas sem styling avançado

### 3. **pdf_generator.py** (HTML Template - Avançado)
- **Template HTML**: Usa o template existente do sistema
- **Mais dependências**: Requer WeasyPrint e bibliotecas do sistema
- **Estilo HTML**: Converte HTML diretamente para PDF

## Requisitos

- Python 3.8+
- Virtual environment (recomendado)

## Instalação

1. Navegue até o diretório do projeto:
```bash
cd python/
```

2. Crie um ambiente virtual:
```bash
python3 -m venv venv
source venv/bin/activate  # No Windows: venv\Scripts\activate
```

3. Instale as dependências:
```bash
pip install -r requirements.txt
```

## Uso

### Gerador Profissional (⭐ RECOMENDADO ⭐)

```bash
# Ativar ambiente virtual
source venv/bin/activate

# Gerar PDF profissional com os arquivos de exemplo
python professional_pdf_generator.py request.json response.json

# Especificar arquivo de saída
python professional_pdf_generator.py request.json response.json ryan_nelson_professional.pdf

# Usar o script de exemplo (mais fácil)
./generate_example.sh
```

### Gerador Simples (Básico)

```bash
# Gerar PDF básico
python simple_pdf_generator.py request.json response.json basic_report.pdf
```

### Gerador com Template HTML

```bash
# Instalar dependências do sistema (macOS)
brew install pango glib libffi gobject-introspection

# Gerar PDF usando template HTML
python pdf_generator.py request.json response.json template_report.pdf
```

### Como Módulo Python

```python
from pdf_generator import KYCPDFGenerator

# Inicializar gerador
generator = KYCPDFGenerator()

# Gerar PDF
output_file = generator.generate_pdf(
    request_file="request.json",
    response_file="response.json",
    output_file="report.pdf"
)

print(f"PDF gerado: {output_file}")
```

## Estrutura dos Arquivos

### request.json
Arquivo JSON com os dados da requisição enviada para Trulioo:
```json
{
    "67228aef1e5e2108d84020a2": "manual-2025-09-11-645687",
    "6716b75a1287d277472c8d82": "Ryan",
    "6716b75a1287d277472c8d83": "Nelson",
    // ... outros campos
}
```

### response.json
Arquivo JSON com a resposta da API Trulioo:
```json
{
    "id": "68e774cb2d00003b009cbcf1",
    "status": "REJECTED",
    "flowData": {
        // ... dados do fluxo de verificação
    }
    // ... outros campos
}
```

## Funcionalidades

- ✅ **Parsing Inteligente**: Extrai automaticamente informações dos JSONs da Trulioo
- ✅ **Template HTML**: Usa o template existente do sistema (`templates/kyc_report.html`)
- ✅ **Informações Completas**: Inclui dados pessoais, endereço, status de verificação
- ✅ **Screening Details**: Mostra hits de Watchlist, Adverse Media e PEP
- ✅ **Processing Steps**: Lista todas as etapas de verificação executadas
- ✅ **Match Status**: Indica quais dados foram verificados com sucesso
- ✅ **Styling**: Mantém toda a formatação e cores do template original

## Dados Extraídos

O gerador extrai e processa:

### Informações Pessoais
- Nome completo (first, middle, last name)
- Data de nascimento
- Endereço completo
- ID Nacional/SSN
- Cliente Reference ID

### Resultados de Screening
- **Watchlist Hits**: Número de matches em listas de observação
- **Adverse Media Hits**: Número de matches em mídia adversa
- **PEP Hits**: Número de matches em Pessoas Politicamente Expostas

### Status de Processamento
- Status geral (ACCEPTED/REJECTED)
- Etapas de processamento executadas
- Status de match para cada verificação

## Exemplo de Saída

O PDF gerado inclui:

1. **Application Details**: Dados pessoais com status de match
2. **Screening Details**: Contadores de hits e detalhes de processamento
3. **Processing Steps**: Tabela com todas as etapas executadas
4. **Overall Status**: Status final da verificação

## Tratamento de Erros

O script inclui tratamento robusto de erros para:
- Arquivos JSON malformados
- Campos ausentes nos dados
- Problemas de template
- Erros de geração de PDF

## Dependências

- **weasyprint**: Geração de PDF a partir de HTML/CSS
- **jinja2**: Engine de templates
- **json-datetime**: Manipulação de datas em JSON
- **python-dateutil**: Utilities para datas

## Personalização

### Template HTML
O template está localizado em `../templates/kyc_report.html` e pode ser personalizado conforme necessário.

### Styling
O CSS é aplicado automaticamente para otimizar a impressão em PDF formato A4.

### Campos Adicionais
Para adicionar novos campos, edite a função `parse_trulioo_response()` no arquivo `pdf_generator.py`.