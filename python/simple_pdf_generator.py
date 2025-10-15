#!/usr/bin/env python3
"""
Simple KYC PDF Report Generator using ReportLab
Generates PDF reports from Trulioo JSON request/response data
"""

import json
import os
import sys
from datetime import datetime
from pathlib import Path
from typing import Dict, Any, List

from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.units import inch
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle
from reportlab.lib import colors
from reportlab.lib.enums import TA_LEFT, TA_CENTER, TA_RIGHT


class SimpleKYCPDFGenerator:
    def __init__(self):
        """Initialize simple PDF generator"""
        self.styles = getSampleStyleSheet()
        self.setup_custom_styles()

    def setup_custom_styles(self):
        """Setup custom paragraph styles"""
        self.styles.add(ParagraphStyle(
            name='CustomTitle',
            parent=self.styles['Heading1'],
            fontSize=24,
            spaceAfter=30,
            alignment=TA_CENTER
        ))

        self.styles.add(ParagraphStyle(
            name='SectionHeader',
            parent=self.styles['Heading2'],
            fontSize=16,
            spaceAfter=12,
            spaceBefore=20,
            textColor=colors.black,
            borderPadding=10
        ))

        self.styles.add(ParagraphStyle(
            name='StatusAccepted',
            parent=self.styles['Normal'],
            fontSize=14,
            textColor=colors.green,
            alignment=TA_CENTER,
            fontName='Helvetica-Bold'
        ))

        self.styles.add(ParagraphStyle(
            name='StatusRejected',
            parent=self.styles['Normal'],
            fontSize=14,
            textColor=colors.red,
            alignment=TA_CENTER,
            fontName='Helvetica-Bold'
        ))

    def load_json_file(self, file_path: str) -> Dict[str, Any]:
        """Load and parse JSON file"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                # Remove comments (lines starting with //)
                content = []
                for line in f:
                    stripped = line.strip()
                    if not stripped.startswith('//'):
                        content.append(line)
                return json.loads(''.join(content))
        except Exception as e:
            raise ValueError(f"Error loading JSON file {file_path}: {e}")

    def parse_trulioo_response(self, response_data: Dict[str, Any]) -> Dict[str, Any]:
        """Parse Trulioo response data and extract relevant information"""
        try:
            flow_data = response_data.get("flowData", {})
            main_flow = list(flow_data.values())[0] if flow_data else {}
            field_data = main_flow.get("fieldData", {})
            service_data = main_flow.get("serviceData", [])

            # Extract basic information
            full_name_parts = []
            first_name = ""
            middle_name = ""
            last_name = ""

            for field_id, field_info in field_data.items():
                role = field_info.get("role", "")
                value = field_info.get("value", [])
                field_value = value[0] if value else ""

                if role == "first_name":
                    first_name = field_value
                    full_name_parts.append(field_value)
                elif role == "last_name":
                    last_name = field_value
                    full_name_parts.append(field_value)
                elif field_info.get("normalizedName") == "MiddleName":
                    middle_name = field_value
                    full_name_parts.insert(-1, field_value)  # Insert before last name

            full_name = " ".join(filter(None, full_name_parts))

            # Extract address components
            address_parts = []
            for field_id, field_info in field_data.items():
                role = field_info.get("role", "")
                value = field_info.get("value", [])
                field_value = value[0] if value else ""

                if role in ["address_1", "address_city", "address_state", "address_zip", "address_country"]:
                    if field_value:
                        address_parts.append(field_value)

            full_address = ", ".join(address_parts)

            # Extract other fields
            dob = ""
            national_id = ""
            client_ref = ""

            for field_id, field_info in field_data.items():
                role = field_info.get("role", "")
                value = field_info.get("value", [])
                field_value = value[0] if value else ""

                if role == "dob":
                    dob = field_value
                elif role == "social_service_number":
                    national_id = field_value
                elif role == "external_customer_id":
                    client_ref = field_value

            # Extract watchlist results
            wl_hits = 0
            am_hits = 0
            pep_hits = 0

            for service in service_data:
                if service.get("nodeType") == "trulioo_person_wl":
                    watchlist_results = service.get("watchlistResults", {})
                    advanced_watchlist = watchlist_results.get("Advanced Watchlist", {})
                    hit_details = advanced_watchlist.get("watchlistHitDetails", {})

                    wl_hits = hit_details.get("wlHitsNumber", 0)
                    am_hits = hit_details.get("amHitsNumber", 0)
                    pep_hits = hit_details.get("pepHitsNumber", 0)
                    break

            # Extract processing steps
            processing_steps = []
            step_number = 1

            for service in service_data:
                step = [
                    f"Step {step_number}",
                    service.get("nodeTitle", "Unknown Service"),
                    service.get("serviceStatus", "Unknown"),
                    "Yes" if service.get("match", False) else "No"
                ]
                processing_steps.append(step)
                step_number += 1

            return {
                "ApplicationID": response_data.get("id", "N/A"),
                "ClientReferenceID": client_ref or "N/A",
                "FullName": full_name or "N/A",
                "FullAddress": full_address or "N/A",
                "DateOfBirth": dob or "N/A",
                "NationalID": national_id or "N/A",
                "Status": response_data.get("status", "Unknown"),
                "WatchlistHits": str(wl_hits),
                "AdverseMediaHits": str(am_hits),
                "PEPHits": str(pep_hits),
                "ProcessingSteps": processing_steps,
                "ProcessedDate": datetime.fromtimestamp(response_data.get("lastModified", 0)).strftime("%Y-%m-%d %H:%M:%S") if response_data.get("lastModified") else "",
                "GeneratedDate": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                "ClientName": full_name or "Unknown Client"
            }

        except Exception as e:
            raise ValueError(f"Error parsing Trulioo response: {e}")

    def create_pdf_elements(self, data: Dict[str, Any]) -> List:
        """Create PDF elements from parsed data"""
        elements = []

        # Title
        title = Paragraph("KYC Report", self.styles['CustomTitle'])
        elements.append(title)
        elements.append(Spacer(1, 12))

        # Generated date
        date_para = Paragraph(f"Generated: {data['GeneratedDate']}", self.styles['Normal'])
        elements.append(date_para)
        elements.append(Spacer(1, 20))

        # Application Details Section
        elements.append(Paragraph("Application Details", self.styles['SectionHeader']))

        app_details_data = [
            ['Field', 'Value'],
            ['Application ID', data['ApplicationID']],
            ['Client Reference', data['ClientReferenceID']],
            ['Full Name', data['FullName']],
            ['Address', data['FullAddress']],
            ['Date of Birth', data['DateOfBirth']],
            ['National ID / SSN', data['NationalID']],
        ]

        app_table = Table(app_details_data, colWidths=[2*inch, 4*inch])
        app_table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, 0), colors.grey),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
            ('ALIGN', (0, 0), (-1, -1), 'LEFT'),
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('FONTSIZE', (0, 0), (-1, 0), 12),
            ('BOTTOMPADDING', (0, 0), (-1, 0), 12),
            ('BACKGROUND', (0, 1), (-1, -1), colors.beige),
            ('GRID', (0, 0), (-1, -1), 1, colors.black)
        ]))

        elements.append(app_table)
        elements.append(Spacer(1, 20))

        # Screening Results Section
        elements.append(Paragraph("Screening Results", self.styles['SectionHeader']))

        screening_data = [
            ['Type', 'Hits'],
            ['Watch List', data['WatchlistHits']],
            ['Adverse Media', data['AdverseMediaHits']],
            ['PEP (Politically Exposed Persons)', data['PEPHits']]
        ]

        screening_table = Table(screening_data, colWidths=[3*inch, 2*inch])
        screening_table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, 0), colors.grey),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
            ('ALIGN', (0, 0), (-1, -1), 'LEFT'),
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('FONTSIZE', (0, 0), (-1, 0), 12),
            ('BOTTOMPADDING', (0, 0), (-1, 0), 12),
            ('BACKGROUND', (0, 1), (-1, -1), colors.beige),
            ('GRID', (0, 0), (-1, -1), 1, colors.black)
        ]))

        elements.append(screening_table)
        elements.append(Spacer(1, 20))

        # Processing Steps Section
        if data['ProcessingSteps']:
            elements.append(Paragraph("Processing Steps", self.styles['SectionHeader']))

            steps_data = [['Step', 'Service', 'Status', 'Match']] + data['ProcessingSteps']

            steps_table = Table(steps_data, colWidths=[1*inch, 2.5*inch, 1.5*inch, 1*inch])
            steps_table.setStyle(TableStyle([
                ('BACKGROUND', (0, 0), (-1, 0), colors.grey),
                ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
                ('ALIGN', (0, 0), (-1, -1), 'LEFT'),
                ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
                ('FONTSIZE', (0, 0), (-1, 0), 12),
                ('BOTTOMPADDING', (0, 0), (-1, 0), 12),
                ('BACKGROUND', (0, 1), (-1, -1), colors.beige),
                ('GRID', (0, 0), (-1, -1), 1, colors.black)
            ]))

            elements.append(steps_table)
            elements.append(Spacer(1, 20))

        # Overall Status
        elements.append(Paragraph("Overall Status", self.styles['SectionHeader']))

        status_style = 'StatusAccepted' if data['Status'] == 'ACCEPTED' else 'StatusRejected'
        status_para = Paragraph(f"<b>{data['Status']}</b>", self.styles[status_style])
        elements.append(status_para)
        elements.append(Spacer(1, 20))

        # Footer information
        if data['ProcessedDate']:
            elements.append(Paragraph(f"Processed Date: {data['ProcessedDate']}", self.styles['Normal']))

        # Watermark
        elements.append(Spacer(1, 30))
        watermark = Paragraph("Generated by KORE KYC System", self.styles['Normal'])
        elements.append(watermark)

        return elements

    def generate_pdf(self, request_file: str, response_file: str, output_file: str = None) -> str:
        """Generate PDF report from request and response JSON files"""

        # Load JSON files
        print(f"Loading request file: {request_file}")
        request_data = self.load_json_file(request_file)

        print(f"Loading response file: {response_file}")
        response_data = self.load_json_file(response_file)

        # Parse data
        print("Parsing Trulioo response data...")
        template_data = self.parse_trulioo_response(response_data)

        # Generate output filename if not provided
        if output_file is None:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            client_name = template_data.get("ClientName", "Unknown").replace(" ", "_")
            output_file = f"kyc_report_{client_name}_{timestamp}.pdf"

        # Ensure output directory exists
        output_path = Path(output_file)
        output_path.parent.mkdir(parents=True, exist_ok=True)

        # Create PDF
        print(f"Generating PDF: {output_file}")
        doc = SimpleDocTemplate(
            str(output_file),
            pagesize=A4,
            rightMargin=72,
            leftMargin=72,
            topMargin=72,
            bottomMargin=18
        )

        # Create PDF elements
        elements = self.create_pdf_elements(template_data)

        # Build PDF
        doc.build(elements)

        print(f"PDF generated successfully: {output_file}")
        return output_file


def main():
    """Main function for command line usage"""
    if len(sys.argv) < 3:
        print("Usage: python simple_pdf_generator.py <request.json> <response.json> [output.pdf]")
        print("Example: python simple_pdf_generator.py request.json response.json report.pdf")
        sys.exit(1)

    request_file = sys.argv[1]
    response_file = sys.argv[2]
    output_file = sys.argv[3] if len(sys.argv) > 3 else None

    try:
        generator = SimpleKYCPDFGenerator()
        output_path = generator.generate_pdf(request_file, response_file, output_file)
        print(f"\n✅ Success! PDF generated: {output_path}")

    except Exception as e:
        print(f"\n❌ Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()