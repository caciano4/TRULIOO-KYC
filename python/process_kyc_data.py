#!/usr/bin/env python3
"""
Script completo para processar dados KYC e gerar PDFs
Uso: python process_kyc_data.py [options]
"""

import argparse
import os
import sys
from pathlib import Path
from convert_user_data import process_relatori_kyc_file

def main():
    parser = argparse.ArgumentParser(description='Process KYC data and generate PDFs')
    parser.add_argument('--input', '-i', default='relatori_kyc',
                       help='Input file with user data (default: relatori_kyc)')
    parser.add_argument('--response', '-r',
                       help='Response JSON file for PDF generation')
    parser.add_argument('--pdf', '-p', action='store_true',
                       help='Generate PDF after processing')
    parser.add_argument('--output', '-o',
                       help='Output PDF filename')

    args = parser.parse_args()

    print("🔄 KYC Data Processor")
    print("=" * 40)

    # Processar dados do arquivo
    print(f"📝 Processing input file: {args.input}")
    process_relatori_kyc_file(args.input)

    # Verificar se request.json foi gerado
    if not os.path.exists('request.json'):
        print("❌ request.json not generated. Check your input data.")
        return 1

    # Gerar PDF se solicitado
    if args.pdf:
        if not args.response:
            print("❌ Response file required for PDF generation")
            print("💡 Use --response to specify response JSON file")
            return 1

        if not os.path.exists(args.response):
            print(f"❌ Response file not found: {args.response}")
            return 1

        print(f"📄 Generating PDF...")

        # Determinar nome do arquivo de saída
        output_file = args.output or "kyc_report.pdf"

        # Executar gerador de PDF
        try:
            from professional_pdf_generator import ProfessionalKYCPDFGenerator
            generator = ProfessionalKYCPDFGenerator()
            result = generator.generate_pdf('request.json', args.response, output_file)
            print(f"✅ PDF generated: {result}")
        except ImportError:
            print("❌ PDF generator not available")
            print("💡 Make sure professional_pdf_generator.py is in the same directory")
            return 1
        except Exception as e:
            print(f"❌ Error generating PDF: {e}")
            return 1

    print("✅ Process completed successfully!")
    return 0

if __name__ == "__main__":
    sys.exit(main())