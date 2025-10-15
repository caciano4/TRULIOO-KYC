#!/usr/bin/env python3
"""
Script para converter dados de usuários em formato JSON
Uso: Lê os dados do arquivo relatori_kyc e converte para JSON
"""

import json
import re
from datetime import datetime
from pathlib import Path

def parse_user_data(data_string):
    """
    Converte uma string de dados do usuário para o formato JSON desejado

    Formato esperado dos dados (separados por tab ou espaços):
    TA_Responsible | Type | Email | UserID | FirstName | MiddleName | LastName | DOB | Phone | Address | City | Postal | State | Country | SSN
    """

    lines = data_string.strip().split('\n')
    users = []

    for line in lines:
        if not line.strip():
            continue

        # Divide os campos (ajuste conforme necessário)
        fields = line.split('\t')

        # Se não houver tabs, tenta dividir por espaços múltiplos
        if len(fields) < 10:
            fields = re.split(r'\s{2,}', line)

        # Se ainda não funcionou, tenta dividir por espaços simples (caso os dados estejam grudados)
        if len(fields) < 10:
            # Para casos onde os dados estão grudados, usa regex mais específico
            # Ajuste este padrão conforme o formato exato dos seus dados
            fields = re.split(r'(?<=\w)(?=[A-Z][a-z])|(?<=\d)(?=[A-Z])', line)

        # Extrai os campos com validação
        ta_responsible = fields[0] if len(fields) > 0 else ""
        transfer_type = fields[1] if len(fields) > 1 else ""
        email = fields[2] if len(fields) > 2 else ""
        user_id = fields[3] if len(fields) > 3 else ""
        first_name = fields[4] if len(fields) > 4 else ""
        middle_name = fields[5] if len(fields) > 5 else ""
        last_name = fields[6] if len(fields) > 6 else ""
        dob = fields[7] if len(fields) > 7 else ""
        phone = fields[8] if len(fields) > 8 else ""
        address = fields[9] if len(fields) > 9 else ""
        city = fields[10] if len(fields) > 10 else ""
        postal = fields[11] if len(fields) > 11 else ""
        state = fields[12] if len(fields) > 12 else ""
        country = fields[13] if len(fields) > 13 else ""
        ssn = fields[14] if len(fields) > 14 else ""

        # Formata a data se necessário (de DD-MM-YYYY para YYYY-MM-DD)
        if dob and '-' in dob:
            parts = dob.split('-')
            if len(parts) == 3 and len(parts[0]) == 2:  # formato DD-MM-YYYY
                dob = f"{parts[2]}-{parts[1]}-{parts[0]}"

        # Cria o client reference ID
        today = datetime.now().strftime('%Y-%m-%d')
        if user_id and user_id.isdigit():
            client_ref = f"manual-{today}-{user_id}"
        else:
            # Se não houver user_id, usa o nome
            name_slug = f"{first_name}-{last_name}".lower().replace(' ', '-')
            client_ref = f"manual-{today}-{name_slug}"

        # Cria o objeto JSON
        user_obj = {
            "67228aef1e5e2108d84020a2": client_ref,
            "6716b75a1287d277472c8d82": first_name,
            "6716b75a1287d277472c8d83": last_name,
            "6716b75a1287d277472c8d84": dob,
            "6716b75a1287d277472c8d86": address,
            "6716b75a1287d277472c8d8c": city,
            "6716b75a1287d277472c8d88": state,
            "6716b75a1287d277472c8d87": postal,
            "6716b75a1287d277472c8d81": country,
        }

        # Adiciona middle name se existir
        if middle_name and middle_name.strip():
            user_obj["674dcb7ce686813288f9045f"] = middle_name

        # Adiciona SSN se existir e não for "none"
        if ssn and ssn.strip() and ssn.lower() != "none":
            user_obj["6744facf99661447b4b58ff7"] = ssn

        # Adiciona comentário no início
        comment = f"""// TA: {ta_responsible}
    // Type: {transfer_type}
    // Email: {email}
    // UserId: {user_id if user_id else '(not provided)'}
    // Name: {first_name} {middle_name + ' ' if middle_name else ''}{last_name}
    // DOB: {dob}
    // Phone: {phone}
    // Address: {address}, {city}, {state}, {postal}, {country}
    // SSN: {ssn if ssn else 'none'}"""

        users.append({
            'comment': comment,
            'data': user_obj
        })

    return users


def format_json_output(users):
    """Formata a saída JSON com comentários"""
    if len(users) == 1:
        # Para um único usuário, usa o formato atual (sem array)
        user = users[0]
        output = "{\n"
        output += "    " + user['comment'].replace('\n', '\n    ') + "\n"

        # Adiciona os campos do JSON
        items = list(user['data'].items())
        for j, (key, value) in enumerate(items):
            comma = "," if j < len(items) - 1 else ""
            output += f'    "{key}": "{value}"{comma}\n'

        output += "}"
    else:
        # Para múltiplos usuários, usa array
        output = "[\n"

        for i, user in enumerate(users):
            output += "  {\n"
            output += "    " + user['comment'].replace('\n', '\n    ') + "\n"

            # Adiciona os campos do JSON
            items = list(user['data'].items())
            for j, (key, value) in enumerate(items):
                comma = "," if j < len(items) - 1 else ""
                output += f'    "{key}": "{value}"{comma}\n'

            if i < len(users) - 1:
                output += "  },\n"
            else:
                output += "  }\n"

        output += "]"

    return output


def process_relatori_kyc_file(file_path="relatori_kyc"):
    """
    Processa o arquivo relatori_kyc e gera os JSONs de request
    """
    try:
        # Tenta ler o arquivo
        with open(file_path, 'r', encoding='utf-8') as f:
            data = f.read()

        print(f"✅ Arquivo {file_path} lido com sucesso")
        print(f"📄 Conteúdo ({len(data)} caracteres):")
        print("=" * 50)
        print(data[:200] + "..." if len(data) > 200 else data)
        print("=" * 50)

        # Processa os dados
        users = parse_user_data(data)

        if not users:
            print("❌ Nenhum usuário encontrado no arquivo")
            return

        print(f"✅ {len(users)} usuário(s) processado(s)")

        # Gera um arquivo para cada usuário ou um arquivo único
        if len(users) == 1:
            # Um único arquivo request.json
            json_output = format_json_output(users)

            with open('request.json', 'w', encoding='utf-8') as f:
                f.write(json_output)

            print(f"✅ Arquivo request.json gerado para: {users[0]['data']['6716b75a1287d277472c8d82']} {users[0]['data']['6716b75a1287d277472c8d83']}")

        else:
            # Múltiplos arquivos
            for i, user in enumerate(users):
                filename = f"request_{i+1}.json"
                json_output = format_json_output([user])

                with open(filename, 'w', encoding='utf-8') as f:
                    f.write(json_output)

                print(f"✅ Arquivo {filename} gerado para: {user['data']['6716b75a1287d277472c8d82']} {user['data']['6716b75a1287d277472c8d83']}")

        # Também salva todos em um arquivo combined
        if len(users) > 1:
            combined_output = format_json_output(users)
            with open('request_combined.json', 'w', encoding='utf-8') as f:
                f.write(combined_output)
            print(f"✅ Arquivo request_combined.json gerado com todos os {len(users)} usuários")

    except FileNotFoundError:
        print(f"❌ Arquivo {file_path} não encontrado")
        print(f"💡 Crie o arquivo {file_path} com os dados dos usuários")

        # Cria um arquivo de exemplo
        example_data = """Wei	Merge	jlz9396@gmail.com	526481	Joseph		Zinsmeyer	1972-02-15	8309315340	8715 W HIGHWAY 71 #3206	Austin	78735	TX	US	452152252
Marie-Lou	Transfer	ryan.m.nelson28@gmail.com	645687	Ryan	Michael	Nelson	1993-06-28	1 303-562-7464	12964 W 64th Drive	Arvada	80004	CO	US	523096930"""

        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(example_data)

        print(f"✅ Arquivo de exemplo {file_path} criado")
        print("📝 Edite o arquivo com seus dados e execute novamente")

    except Exception as e:
        print(f"❌ Erro ao processar arquivo: {e}")


if __name__ == "__main__":
    print("🔄 Processando arquivo relatori_kyc...")
    process_relatori_kyc_file()