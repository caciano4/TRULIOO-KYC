#!/usr/bin/env python3
"""
KYC PDF Report Generator
Generates PDF reports from Trulioo JSON request/response data using HTML templates
"""

import json
import os
import sys
from datetime import datetime
from pathlib import Path
from typing import Dict, Any, List

from jinja2 import Template
from weasyprint import HTML, CSS
from weasyprint.text.fonts import FontConfiguration


class KYCPDFGenerator:
    def __init__(self, template_path: str = None):
        """Initialize PDF generator with template path"""
        if template_path is None:
            # Default template path relative to this script
            self.template_path = Path(__file__).parent.parent / "templates" / "kyc_report.html"
        else:
            self.template_path = Path(template_path)

        if not self.template_path.exists():
            raise FileNotFoundError(f"Template not found: {self.template_path}")

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
            address_1 = ""
            city = ""
            state = ""
            postal = ""
            country = ""

            for field_id, field_info in field_data.items():
                role = field_info.get("role", "")
                value = field_info.get("value", [])
                field_value = value[0] if value else ""

                if role == "address_1":
                    address_1 = field_value
                    address_parts.append(field_value)
                elif role == "address_city":
                    city = field_value
                    address_parts.append(field_value)
                elif role == "address_state":
                    state = field_value
                    address_parts.append(field_value)
                elif role == "address_zip":
                    postal = field_value
                    address_parts.append(field_value)
                elif role == "address_country":
                    country = field_value
                    address_parts.append(field_value)

            full_address = ", ".join(filter(None, address_parts))

            # Extract other fields
            dob = ""
            national_id = ""
            client_ref = ""
            phone = ""
            email = ""

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
                step = {
                    "Step": f"Step {step_number}",
                    "Service": service.get("nodeTitle", "Unknown Service"),
                    "Status": service.get("serviceStatus", "Unknown"),
                    "Match": "Yes" if service.get("match", False) else "No"
                }
                processing_steps.append(step)
                step_number += 1

            # Determine match classes for styling
            def get_match_class(has_match):
                return "match-yes" if has_match else "match-no"

            return {
                "ApplicationID": response_data.get("id", "N/A"),
                "ClientReferenceID": client_ref or "N/A",
                "FullName": full_name or "N/A",
                "FullAddress": full_address or "N/A",
                "DateOfBirth": dob or "N/A",
                "NationalID": national_id or "N/A",
                "Phone": phone or "N/A",
                "Email": email or "N/A",
                "Status": response_data.get("status", "Unknown"),
                "WatchlistHits": str(wl_hits),
                "AdverseMediaHits": str(am_hits),
                "PEPHits": str(pep_hits),
                "ProcessingSteps": processing_steps,
                "TransactionID": response_data.get("id", ""),
                "ProcessedDate": datetime.fromtimestamp(response_data.get("lastModified", 0)).strftime("%Y-%m-%d %H:%M:%S") if response_data.get("lastModified") else "",
                "GeneratedDate": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                "ClientName": full_name or "Unknown Client",

                # Match classes for styling
                "NameMatchClass": get_match_class(bool(full_name)),
                "AddressMatchClass": get_match_class(bool(full_address)),
                "DOBMatchClass": get_match_class(bool(dob)),
                "IDMatchClass": get_match_class(bool(national_id)),

                # Match status text
                "NameMatch": "Yes" if full_name else "No",
                "AddressMatch": "Yes" if full_address else "No",
                "DOBMatch": "Yes" if dob else "No",
                "IDMatch": "Yes" if national_id else "No"
            }

        except Exception as e:
            raise ValueError(f"Error parsing Trulioo response: {e}")

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

        # Load template
        print(f"Loading template: {self.template_path}")
        with open(self.template_path, 'r', encoding='utf-8') as f:
            template_content = f.read()

        # Render template
        print("Rendering HTML template...")
        template = Template(template_content)
        html_content = template.render(**template_data)

        # Generate output filename if not provided
        if output_file is None:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            client_name = template_data.get("ClientName", "Unknown").replace(" ", "_")
            output_file = f"kyc_report_{client_name}_{timestamp}.pdf"

        # Ensure output directory exists
        output_path = Path(output_file)
        output_path.parent.mkdir(parents=True, exist_ok=True)

        # Generate PDF
        print(f"Generating PDF: {output_file}")
        font_config = FontConfiguration()

        # Create CSS for better PDF rendering
        css_content = """
        @page {
            size: A4;
            margin: 1cm;
        }
        body {
            font-family: Arial, sans-serif;
            line-height: 1.4;
        }
        """

        html_doc = HTML(string=html_content)
        css_doc = CSS(string=css_content, font_config=font_config)

        html_doc.write_pdf(output_file, stylesheets=[css_doc], font_config=font_config)

        print(f"PDF generated successfully: {output_file}")
        return output_file


def main():
    """Main function for command line usage"""
    if len(sys.argv) < 3:
        print("Usage: python pdf_generator.py <request.json> <response.json> [output.pdf]")
        print("Example: python pdf_generator.py request.json response.json report.pdf")
        sys.exit(1)

    request_file = sys.argv[1]
    response_file = sys.argv[2]
    output_file = sys.argv[3] if len(sys.argv) > 3 else None

    try:
        generator = KYCPDFGenerator()
        output_path = generator.generate_pdf(request_file, response_file, output_file)
        print(f"\n✅ Success! PDF generated: {output_path}")

    except Exception as e:
        print(f"\n❌ Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()