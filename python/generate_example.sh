#!/bin/bash
# Script para gerar exemplo de PDF KYC

set -e

echo "🚀 KYC PDF Generator - Exemplo"
echo "==============================="

# Verificar se estamos no diretório correto
if [ ! -f "simple_pdf_generator.py" ]; then
    echo "❌ Erro: Execute este script no diretório python/"
    exit 1
fi

# Ativar ambiente virtual se existir
if [ -d "venv" ]; then
    echo "🔄 Ativando ambiente virtual..."
    source venv/bin/activate
else
    echo "⚠️  Ambiente virtual não encontrado. Executando com Python global..."
fi

# Verificar se os arquivos de exemplo existem
if [ ! -f "request.json" ] || [ ! -f "response.json" ]; then
    echo "❌ Erro: Arquivos request.json ou response.json não encontrados"
    echo "   Certifique-se de que os arquivos de exemplo estão no diretório python/"
    exit 1
fi

# Gerar PDF profissional
echo "📄 Gerando PDF profissional de exemplo..."
python professional_pdf_generator.py request.json response.json professional_example_report.pdf

# Verificar se o PDF foi criado
if [ -f "professional_example_report.pdf" ]; then
    echo "✅ Sucesso! PDF profissional gerado: professional_example_report.pdf"
    echo "📁 Tamanho do arquivo: $(du -h professional_example_report.pdf | cut -f1)"
    echo ""
    echo "🎨 Este PDF tem o visual profissional que corresponde ao template HTML!"
    echo ""
    echo "Para visualizar o PDF:"
    echo "  macOS: open professional_example_report.pdf"
    echo "  Linux: xdg-open professional_example_report.pdf"
    echo "  Windows: start professional_example_report.pdf"
else
    echo "❌ Erro: PDF não foi criado"
    exit 1
fi