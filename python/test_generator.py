#!/usr/bin/env python3
"""
Test script for KYC PDF Generator
"""

import os
import sys
from pathlib import Path
from pdf_generator import KYCPDFGenerator


def test_pdf_generation():
    """Test PDF generation with sample files"""

    current_dir = Path(__file__).parent
    request_file = current_dir / "request.json"
    response_file = current_dir / "response.json"

    # Check if files exist
    if not request_file.exists():
        print(f"❌ Request file not found: {request_file}")
        return False

    if not response_file.exists():
        print(f"❌ Response file not found: {response_file}")
        return False

    try:
        print("🔄 Initializing PDF generator...")
        generator = KYCPDFGenerator()

        print("🔄 Generating PDF report...")
        output_file = generator.generate_pdf(
            request_file=str(request_file),
            response_file=str(response_file),
            output_file="test_report.pdf"
        )

        # Check if PDF was created
        if Path(output_file).exists():
            file_size = Path(output_file).stat().st_size
            print(f"✅ Success! PDF generated: {output_file}")
            print(f"📄 File size: {file_size:,} bytes")
            return True
        else:
            print(f"❌ PDF file was not created: {output_file}")
            return False

    except Exception as e:
        print(f"❌ Error generating PDF: {e}")
        return False


def test_data_parsing():
    """Test data parsing functionality"""

    current_dir = Path(__file__).parent
    response_file = current_dir / "response.json"

    if not response_file.exists():
        print(f"❌ Response file not found: {response_file}")
        return False

    try:
        print("🔄 Testing data parsing...")
        generator = KYCPDFGenerator()

        # Load and parse response
        response_data = generator.load_json_file(str(response_file))
        parsed_data = generator.parse_trulioo_response(response_data)

        # Print extracted data
        print("\n📊 Extracted Data:")
        print(f"  Client Name: {parsed_data['ClientName']}")
        print(f"  Application ID: {parsed_data['ApplicationID']}")
        print(f"  Status: {parsed_data['Status']}")
        print(f"  Full Name: {parsed_data['FullName']}")
        print(f"  Address: {parsed_data['FullAddress']}")
        print(f"  Date of Birth: {parsed_data['DateOfBirth']}")
        print(f"  National ID: {parsed_data['NationalID']}")
        print(f"  Watchlist Hits: {parsed_data['WatchlistHits']}")
        print(f"  Adverse Media Hits: {parsed_data['AdverseMediaHits']}")
        print(f"  PEP Hits: {parsed_data['PEPHits']}")
        print(f"  Processing Steps: {len(parsed_data['ProcessingSteps'])}")

        print("✅ Data parsing successful!")
        return True

    except Exception as e:
        print(f"❌ Error parsing data: {e}")
        return False


def main():
    """Run all tests"""
    print("🧪 KYC PDF Generator - Test Suite")
    print("=" * 50)

    # Test 1: Data parsing
    print("\n1. Testing data parsing...")
    parsing_success = test_data_parsing()

    # Test 2: PDF generation
    print("\n2. Testing PDF generation...")
    pdf_success = test_pdf_generation()

    # Summary
    print("\n" + "=" * 50)
    print("📋 Test Results:")
    print(f"  Data Parsing: {'✅ PASS' if parsing_success else '❌ FAIL'}")
    print(f"  PDF Generation: {'✅ PASS' if pdf_success else '❌ FAIL'}")

    if parsing_success and pdf_success:
        print("\n🎉 All tests passed! The PDF generator is working correctly.")
        return 0
    else:
        print("\n⚠️  Some tests failed. Please check the errors above.")
        return 1


if __name__ == "__main__":
    sys.exit(main())