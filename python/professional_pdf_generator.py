#!/usr/bin/env python3
"""
Professional KYC PDF Report Generator
Generates PDF reports that match the exact structure and styling of the HTML template
"""

import json
import os
import sys
from datetime import datetime
from pathlib import Path
from typing import Dict, Any, List

from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.units import inch, mm
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle, PageBreak
from reportlab.lib import colors
from reportlab.lib.enums import TA_LEFT, TA_CENTER, TA_RIGHT
from reportlab.pdfgen import canvas
from reportlab.platypus.frames import Frame
from reportlab.platypus.doctemplate import PageTemplate, BaseDocTemplate


class ProfessionalKYCPDFGenerator:
    def __init__(self):
        """Initialize professional PDF generator with exact template styling"""
        self.styles = getSampleStyleSheet()
        self.setup_professional_styles()

    def setup_professional_styles(self):
        """Setup styles that match the HTML template exactly"""

        # Header Title Style
        self.styles.add(ParagraphStyle(
            name='HeaderTitle',
            fontName='Helvetica-Bold',
            fontSize=28,
            textColor=colors.black,
            spaceAfter=15,  # Increased from 5 to 15
            alignment=TA_CENTER  # Changed from TA_LEFT to TA_CENTER
        ))

        # Header Date Style
        self.styles.add(ParagraphStyle(
            name='HeaderDate',
            fontName='Helvetica',
            fontSize=14,
            textColor=colors.Color(0.4, 0.4, 0.4),  # #666
            spaceAfter=30,  # Increased from 20 to 30
            alignment=TA_CENTER  # Changed from TA_LEFT to TA_CENTER
        ))

        # Section Title Style (matches .section-title from template)
        self.styles.add(ParagraphStyle(
            name='SectionTitle',
            fontName='Helvetica-Bold',
            fontSize=16,  # Reduced from 18 to 16
            textColor=colors.black,
            spaceAfter=12,  # Reduced from 20 to 12
            spaceBefore=12,  # Reduced from 20 to 12
            borderWidth=0,
            borderPadding=8,
            alignment=TA_CENTER  # Changed from TA_LEFT to TA_CENTER
        ))

        # Medium Number Style for screening results (reduced from 32px)
        self.styles.add(ParagraphStyle(
            name='ScreeningValue',
            fontName='Helvetica-Bold',
            fontSize=20,  # Reduced from 32 to 20
            textColor=colors.black,
            alignment=TA_CENTER,
            spaceAfter=5
        ))

        # Screening Label Style
        self.styles.add(ParagraphStyle(
            name='ScreeningLabel',
            fontName='Helvetica-Bold',
            fontSize=14,
            textColor=colors.Color(0.2, 0.2, 0.2),  # #333
            alignment=TA_CENTER,
            spaceBefore=0
        ))

        # Status Badge Styles
        self.styles.add(ParagraphStyle(
            name='StatusAccepted',
            fontName='Helvetica-Bold',
            fontSize=14,
            textColor=colors.white,
            alignment=TA_CENTER,
            borderWidth=0,
            borderPadding=8
        ))

        self.styles.add(ParagraphStyle(
            name='StatusRejected',
            fontName='Helvetica-Bold',
            fontSize=14,
            textColor=colors.white,
            alignment=TA_CENTER,
            borderWidth=0,
            borderPadding=8
        ))

        # Normal text
        self.styles.add(ParagraphStyle(
            name='KYCBodyText',
            fontName='Helvetica',
            fontSize=13,
            textColor=colors.Color(0.2, 0.2, 0.2),  # #333
            alignment=TA_LEFT
        ))

        # Watermark
        self.styles.add(ParagraphStyle(
            name='Watermark',
            fontName='Helvetica',
            fontSize=12,
            textColor=colors.Color(0.6, 0.6, 0.6),  # #999
            alignment=TA_CENTER,
            spaceBefore=20
        ))

    def load_json_file(self, file_path: str) -> Dict[str, Any]:
        """Load and parse JSON file"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = []
                for line in f:
                    stripped = line.strip()
                    if not stripped.startswith('//'):
                        content.append(line)
                return json.loads(''.join(content))
        except Exception as e:
            raise ValueError(f"Error loading JSON file {file_path}: {e}")

    def extract_comments_from_request(self, file_path: str) -> Dict[str, str]:
        """Extract phone and email from comments in request JSON"""
        comments_data = {}
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                for line in f:
                    stripped = line.strip()
                    if stripped.startswith('//'):
                        if 'Phone:' in stripped:
                            # Extract: // Phone: 3252421782
                            phone = stripped.split('Phone:')[1].strip()
                            comments_data['phone'] = phone
                        elif 'Email:' in stripped:
                            # Extract: // Email: missinglimbranch@mac.com
                            email = stripped.split('Email:')[1].strip()
                            comments_data['email'] = email
        except Exception as e:
            print(f"Warning: Could not extract comments from {file_path}: {e}")

        return comments_data

    def parse_trulioo_response(self, response_data: Dict[str, Any], comments_data: Dict[str, str] = None) -> Dict[str, Any]:
        """Parse Trulioo response data and extract relevant information"""
        try:
            flow_data = response_data.get("flowData", {})
            main_flow = list(flow_data.values())[0] if flow_data else {}
            field_data = main_flow.get("fieldData", {})
            service_data = main_flow.get("serviceData", [])

            # Extract personal information
            personal_info = {}

            for field_id, field_info in field_data.items():
                role = field_info.get("role", "")
                value = field_info.get("value", [])
                field_value = value[0] if value else ""
                normalized_name = field_info.get("normalizedName", "")

                if role == "first_name":
                    personal_info['first_name'] = field_value
                elif role == "last_name":
                    personal_info['last_name'] = field_value
                elif normalized_name == "MiddleName":
                    personal_info['middle_name'] = field_value
                elif role == "dob":
                    personal_info['dob'] = field_value
                elif role == "address_1":
                    personal_info['address'] = field_value
                elif role == "address_city":
                    personal_info['city'] = field_value
                elif role == "address_state":
                    personal_info['state'] = field_value
                elif role == "address_zip":
                    personal_info['postal'] = field_value
                elif role == "address_country":
                    personal_info['country'] = field_value
                elif role == "social_service_number":
                    personal_info['ssn'] = field_value
                elif role == "external_customer_id":
                    personal_info['client_ref'] = field_value

            # Build full name
            name_parts = []
            if personal_info.get('first_name'):
                name_parts.append(personal_info['first_name'])
            if personal_info.get('middle_name'):
                name_parts.append(personal_info['middle_name'])
            if personal_info.get('last_name'):
                name_parts.append(personal_info['last_name'])
            full_name = " ".join(name_parts)

            # Build full address
            address_parts = []
            for key in ['address', 'city', 'state', 'postal', 'country']:
                if personal_info.get(key):
                    address_parts.append(personal_info[key])
            full_address = ", ".join(address_parts)

            # Extract screening results
            screening_results = {'wl': 0, 'am': 0, 'pep': 0}

            for service in service_data:
                if service.get("nodeType") == "trulioo_person_wl":
                    watchlist_results = service.get("watchlistResults", {})
                    advanced_watchlist = watchlist_results.get("Advanced Watchlist", {})
                    hit_details = advanced_watchlist.get("watchlistHitDetails", {})

                    screening_results['wl'] = hit_details.get("wlHitsNumber", 0)
                    screening_results['am'] = hit_details.get("amHitsNumber", 0)
                    screening_results['pep'] = hit_details.get("pepHitsNumber", 0)
                    break

            # Extract processing steps
            processing_steps = []
            for i, service in enumerate(service_data, 1):
                step = {
                    'step': f"Step {i}",
                    'service': service.get("nodeTitle", "Unknown Service"),
                    'status': service.get("serviceStatus", "Unknown"),
                    'match': "Yes" if service.get("match", False) else "No"
                }
                processing_steps.append(step)

            # Determine match results for each field
            def get_match_status(value):
                return "YES" if value and value != "N/A" else "NO"

            # Get phone and email from comments data or default to N/A
            comments_data = comments_data or {}
            phone_value = comments_data.get('phone', 'N/A')
            email_value = comments_data.get('email', 'N/A')

            return {
                'application_id': response_data.get("id", "N/A"),
                'client_reference': personal_info.get('client_ref', 'N/A'),
                'full_name': full_name or "N/A",
                'full_address': full_address or "N/A",
                'date_of_birth': personal_info.get('dob', 'N/A'),
                'national_id': personal_info.get('ssn', 'N/A'),
                'phone': phone_value,
                'email': email_value,
                'status': response_data.get("status", "Unknown"),
                'watchlist_hits': screening_results['wl'],
                'adverse_media_hits': screening_results['am'],
                'pep_hits': screening_results['pep'],
                'processing_steps': processing_steps,
                'processed_date': datetime.fromtimestamp(response_data.get("lastModified", 0)).strftime("%Y-%m-%d %H:%M:%S") if response_data.get("lastModified") else "",
                'generated_date': datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                'client_name': full_name or "Unknown Client",

                # Match status for styling
                'name_match': get_match_status(full_name),
                'address_match': get_match_status(full_address),
                'dob_match': get_match_status(personal_info.get('dob')),
                'id_match': get_match_status(personal_info.get('ssn')),
                'phone_match': "NP",  # Not Processed
                'email_match': "NP"   # Not Processed
            }

        except Exception as e:
            raise ValueError(f"Error parsing Trulioo response: {e}")

    def create_header_section(self, data: Dict[str, Any]) -> List:
        """Create the professional header section"""
        elements = []

        # Main title
        title = Paragraph("KYC Report", self.styles['HeaderTitle'])
        elements.append(title)
        elements.append(Spacer(1, 10))  # Add space between title and date

        # Date
        date_para = Paragraph(data['generated_date'], self.styles['HeaderDate'])
        elements.append(date_para)

        return elements

    def create_application_details_section(self, data: Dict[str, Any]) -> List:
        """Create the Application Details section with match indicators"""
        elements = []

        # Section title with underline
        title = Paragraph("Application Details", self.styles['SectionTitle'])
        elements.append(title)

        # Add underline manually using a table
        underline_table = Table([['']], colWidths=[6*inch])
        underline_table.setStyle(TableStyle([
            ('LINEBELOW', (0, 0), (-1, -1), 2, colors.black),
            ('TOPPADDING', (0, 0), (-1, -1), 0),
            ('BOTTOMPADDING', (0, 0), (-1, -1), 8),
        ]))
        elements.append(underline_table)

        # Application details table with text wrapping for long addresses
        table_data = [
            ['Field', 'Value', 'Match'],
            ['Application ID', data['application_id'], '-'],
            ['Client Reference', data['client_reference'], '-'],
            ['Full Name', data['full_name'], data['name_match']],
            ['Address', Paragraph(data['full_address'], self.styles['KYCBodyText']), data['address_match']],
            ['Date of Birth', data['date_of_birth'], data['dob_match']],
            ['National ID / SSN', data['national_id'], data['id_match']],
            ['Phone', data['phone'], data['phone_match']],
            ['Email', data['email'], data['email_match']]
        ]

        table = Table(table_data, colWidths=[2*inch, 3.2*inch, 0.8*inch])
        table.setStyle(TableStyle([
            # Header row
            ('BACKGROUND', (0, 0), (-1, 0), colors.Color(0.97, 0.97, 0.97)),  # #f8f8f8
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.Color(0.2, 0.2, 0.2)),      # #333
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('FONTSIZE', (0, 0), (-1, 0), 11),
            ('ALIGN', (0, 0), (-1, 0), 'LEFT'),
            ('ALIGN', (2, 0), (2, -1), 'CENTER'),  # Match column centered

            # Data rows
            ('FONTNAME', (0, 1), (-1, -1), 'Helvetica'),
            ('FONTSIZE', (0, 1), (-1, -1), 11),
            ('TEXTCOLOR', (0, 1), (-1, -1), colors.Color(0.2, 0.2, 0.2)),

            # Alternating row colors
            ('BACKGROUND', (0, 1), (-1, 1), colors.white),
            ('BACKGROUND', (0, 2), (-1, 2), colors.Color(0.98, 0.98, 0.98)),  # #fafafa
            ('BACKGROUND', (0, 3), (-1, 3), colors.white),
            ('BACKGROUND', (0, 4), (-1, 4), colors.Color(0.98, 0.98, 0.98)),
            ('BACKGROUND', (0, 5), (-1, 5), colors.white),
            ('BACKGROUND', (0, 6), (-1, 6), colors.Color(0.98, 0.98, 0.98)),
            ('BACKGROUND', (0, 7), (-1, 7), colors.white),
            ('BACKGROUND', (0, 8), (-1, 8), colors.Color(0.98, 0.98, 0.98)),

            # Borders
            ('GRID', (0, 0), (-1, -1), 1, colors.Color(0.87, 0.87, 0.87)),  # #ddd
            ('VALIGN', (0, 0), (-1, -1), 'TOP'),  # Changed from MIDDLE to TOP for better text wrapping
            ('LEFTPADDING', (0, 0), (-1, -1), 5),   # Reduced padding
            ('RIGHTPADDING', (0, 0), (-1, -1), 5),  # Reduced padding
            ('TOPPADDING', (0, 0), (-1, -1), 5),    # Reduced padding
            ('BOTTOMPADDING', (0, 0), (-1, -1), 5), # Reduced padding

            # Match column styling
            ('FONTNAME', (2, 1), (2, -1), 'Helvetica-Bold'),
        ]))

        elements.append(table)

        # Add legend for NP (Not Processed)
        legend_style = ParagraphStyle(
            name='Legend',
            fontName='Helvetica',
            fontSize=11,
            textColor=colors.Color(0.5, 0.5, 0.5),  # Gray color
            alignment=TA_LEFT,
            spaceBefore=10,
            spaceAfter=20
        )
        legend = Paragraph("<b>*</b> NP = Not Processed (field not sent to Trulioo for verification)", legend_style)
        elements.append(legend)

        return elements

    def create_screening_section(self, data: Dict[str, Any]) -> List:
        """Create the Screening Details section with large numbers"""
        elements = []

        # Section title
        title = Paragraph("Screening Details", self.styles['SectionTitle'])
        elements.append(title)

        # Add underline
        underline_table = Table([['']], colWidths=[6*inch])
        underline_table.setStyle(TableStyle([
            ('LINEBELOW', (0, 0), (-1, -1), 2, colors.black),
            ('TOPPADDING', (0, 0), (-1, -1), 0),
            ('BOTTOMPADDING', (0, 0), (-1, -1), 20),
        ]))
        elements.append(underline_table)

        # Create the three-column layout for screening results
        screening_data = []

        # Row 1: Labels
        screening_data.append([
            Paragraph("Watch List", self.styles['ScreeningLabel']),
            Paragraph("Adverse Media", self.styles['ScreeningLabel']),
            Paragraph("PEP", self.styles['ScreeningLabel'])
        ])

        # Row 2: Values (large numbers)
        screening_data.append([
            Paragraph(str(data['watchlist_hits']), self.styles['ScreeningValue']),
            Paragraph(str(data['adverse_media_hits']), self.styles['ScreeningValue']),
            Paragraph(str(data['pep_hits']), self.styles['ScreeningValue'])
        ])

        screening_table = Table(screening_data, colWidths=[2*inch, 2*inch, 2*inch])
        screening_table.setStyle(TableStyle([
            ('ALIGN', (0, 0), (-1, -1), 'CENTER'),
            ('VALIGN', (0, 0), (-1, -1), 'MIDDLE'),
            ('TOPPADDING', (0, 0), (-1, 0), 10),
            ('BOTTOMPADDING', (0, 0), (-1, 0), 10),
            ('TOPPADDING', (0, 1), (-1, 1), 0),
            ('BOTTOMPADDING', (0, 1), (-1, 1), 20),
        ]))

        elements.append(screening_table)

        return elements

    def create_processing_steps_section(self, data: Dict[str, Any]) -> List:
        """Create the Processing Steps section"""
        elements = []

        if not data['processing_steps']:
            return elements

        # Section title
        title = Paragraph("Processing Steps", self.styles['SectionTitle'])
        elements.append(title)

        # Add underline
        underline_table = Table([['']], colWidths=[6*inch])
        underline_table.setStyle(TableStyle([
            ('LINEBELOW', (0, 0), (-1, -1), 2, colors.black),
            ('TOPPADDING', (0, 0), (-1, -1), 0),
            ('BOTTOMPADDING', (0, 0), (-1, -1), 15),
        ]))
        elements.append(underline_table)

        # Processing steps table (remove Status column, keep only Step, Service, Match)
        table_data = [['Step', 'Service', 'Match']]

        for step in data['processing_steps']:
            table_data.append([
                step['step'],
                Paragraph(step['service'], self.styles['KYCBodyText']),  # Wrap service name
                step['match']
            ])

        table = Table(table_data, colWidths=[1*inch, 4*inch, 1*inch])
        table.setStyle(TableStyle([
            # Header
            ('BACKGROUND', (0, 0), (-1, 0), colors.Color(0.94, 0.94, 0.94)),  # #f0f0f0
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('FONTSIZE', (0, 0), (-1, 0), 12),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.black),
            ('ALIGN', (0, 0), (-1, 0), 'LEFT'),

            # Data rows
            ('FONTNAME', (0, 1), (-1, -1), 'Helvetica'),
            ('FONTSIZE', (0, 1), (-1, -1), 12),
            ('TEXTCOLOR', (0, 1), (-1, -1), colors.black),

            # Borders
            ('GRID', (0, 0), (-1, -1), 1, colors.Color(0.8, 0.8, 0.8)),  # #ccc
            ('VALIGN', (0, 0), (-1, -1), 'TOP'),  # Changed to TOP for better text alignment
            ('LEFTPADDING', (0, 0), (-1, -1), 5),   # Reduced padding
            ('RIGHTPADDING', (0, 0), (-1, -1), 5),  # Reduced padding
            ('TOPPADDING', (0, 0), (-1, -1), 5),    # Reduced padding
            ('BOTTOMPADDING', (0, 0), (-1, -1), 5), # Reduced padding
        ]))

        elements.append(table)
        elements.append(Spacer(1, 10))

        return elements

    def create_status_section(self, data: Dict[str, Any]) -> List:
        """Create the Overall Status section"""
        elements = []

        # Status label (create a centered style for this)
        status_label_style = ParagraphStyle(
            name='StatusLabel',
            fontName='Helvetica-Bold',
            fontSize=16,
            textColor=colors.Color(0.2, 0.2, 0.2),
            alignment=TA_CENTER,
            spaceAfter=10
        )
        status_label = Paragraph("OVERALL STATUS", status_label_style)
        elements.append(status_label)
        elements.append(Spacer(1, 10))

        # Status badge with background color
        status = data['status']
        if status == "ACCEPTED":
            bg_color = colors.Color(0.3, 0.686, 0.314)  # #4caf50 green
        else:
            bg_color = colors.Color(0.827, 0.184, 0.184)  # #d32f2f red

        # Create status badge using a table for background color
        status_data = [[status]]
        status_table = Table(status_data, colWidths=[2*inch])
        status_table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, -1), bg_color),
            ('TEXTCOLOR', (0, 0), (-1, -1), colors.white),
            ('FONTNAME', (0, 0), (-1, -1), 'Helvetica-Bold'),
            ('FONTSIZE', (0, 0), (-1, -1), 14),
            ('ALIGN', (0, 0), (-1, -1), 'CENTER'),
            ('VALIGN', (0, 0), (-1, -1), 'MIDDLE'),
            ('TOPPADDING', (0, 0), (-1, -1), 8),
            ('BOTTOMPADDING', (0, 0), (-1, -1), 8),
            ('LEFTPADDING', (0, 0), (-1, -1), 20),
            ('RIGHTPADDING', (0, 0), (-1, -1), 20),
        ]))

        elements.append(status_table)
        elements.append(Spacer(1, 15))

        return elements

    def create_footer_section(self, data: Dict[str, Any]) -> List:
        """Create the footer section"""
        elements = []

        # Transaction ID and processed date if available
        if data.get('processed_date'):
            processed_para = Paragraph(f"Processed Date: {data['processed_date']}", self.styles['KYCBodyText'])
            elements.append(processed_para)
            elements.append(Spacer(1, 10))

        # Watermark
        watermark = Paragraph(f"Generated by KORE KYC System | {data['generated_date']}", self.styles['Watermark'])
        elements.append(watermark)

        return elements

    def generate_pdf(self, request_file: str, response_file: str, output_file: str = None) -> str:
        """Generate professional PDF report"""

        # Load and parse data
        print(f"Loading request file: {request_file}")
        request_data = self.load_json_file(request_file)

        print(f"Loading response file: {response_file}")
        response_data = self.load_json_file(response_file)

        print("Extracting comments from request file...")
        comments_data = self.extract_comments_from_request(request_file)

        print("Parsing Trulioo response data...")
        data = self.parse_trulioo_response(response_data, comments_data)

        # Generate output filename
        if output_file is None:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            client_name = data.get("client_name", "Unknown").replace(" ", "_")
            output_file = f"professional_kyc_report_{client_name}_{timestamp}.pdf"

        # Ensure output directory exists
        output_path = Path(output_file)
        output_path.parent.mkdir(parents=True, exist_ok=True)

        # Create PDF document
        print(f"Generating professional PDF: {output_file}")
        doc = SimpleDocTemplate(
            str(output_file),
            pagesize=A4,
            rightMargin=25,  # Reduced from 40 to 25
            leftMargin=25,   # Reduced from 40 to 25
            topMargin=20,    # Reduced from 30 to 20
            bottomMargin=20  # Reduced from 30 to 20
        )

        # Build all sections
        elements = []

        # Header
        elements.extend(self.create_header_section(data))
        elements.append(Spacer(1, 15))

        # Application Details
        elements.extend(self.create_application_details_section(data))

        # Screening Details
        elements.extend(self.create_screening_section(data))

        # Processing Steps
        elements.extend(self.create_processing_steps_section(data))

        # Overall Status
        elements.extend(self.create_status_section(data))

        # Footer
        elements.extend(self.create_footer_section(data))

        # Build PDF
        doc.build(elements)

        print(f"Professional PDF generated successfully: {output_file}")
        return output_file


def main():
    """Main function for command line usage"""
    if len(sys.argv) < 3:
        print("Usage: python professional_pdf_generator.py <request.json> <response.json> [output.pdf]")
        print("Example: python professional_pdf_generator.py request.json response.json professional_report.pdf")
        sys.exit(1)

    request_file = sys.argv[1]
    response_file = sys.argv[2]
    output_file = sys.argv[3] if len(sys.argv) > 3 else None

    try:
        generator = ProfessionalKYCPDFGenerator()
        output_path = generator.generate_pdf(request_file, response_file, output_file)
        print(f"\n✅ Success! Professional PDF generated: {output_path}")

    except Exception as e:
        print(f"\n❌ Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()