
const tbody = document.querySelector("#data-table tbody");

// Improved function to fetch package list with error handling
const fetchPackageList = async () => {
    try {
        const response = await fetch("http://localhost/kyc-package-list", { method: "GET" });

        if (response.ok) {
            const { data } = await response.json();
            return data;
        } else {
            Swal.fire({
                title: "Error",
                icon: "error",
                text: "An error occurred while fetching data",
            });
            return []; // Return an empty array in case of failure
        }
    } catch (error) {
        Swal.fire({
            title: "Error",
            icon: "error",
            text: `An unexpected error occurred: ${error.message}`,
        });
        return []; // Return an empty array on error
    }
};

// Function to format date
const formatDate = (dateString) => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('pt-BR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};

// Function to get status badge
const getStatusBadge = (completed, totalRecords) => {
    if (completed === 0) {
        return '<span class="status-badge status-pending">Pendente</span>';
    } else if (completed === totalRecords) {
        return '<span class="status-badge status-completed">Completo</span>';
    } else {
        return '<span class="status-badge status-processing">Processando</span>';
    }
};

// Function to create action buttons
const createActionButtons = (item) => {
    const actionsContainer = document.createElement('div');
    actionsContainer.className = 'actions-container';

    // Process KYC Button
    const processBtn = document.createElement('button');
    processBtn.className = 'action-btn btn-process';
    processBtn.setAttribute('data-tooltip', 'Processar verificação KYC');
    processBtn.innerHTML = `
        <svg style="width: 14px; height: 14px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        KYC
    `;

    processBtn.addEventListener('click', async (e) => {
        e.preventDefault();

        // Show confirmation dialog
        const result = await Swal.fire({
            title: 'Processar KYC',
            text: `Deseja processar o pacote ${item.package_name}?`,
            icon: 'question',
            showCancelButton: true,
            confirmButtonColor: 'var(--primary-color)',
            cancelButtonColor: 'var(--gray-400)',
            confirmButtonText: 'Sim, processar',
            cancelButtonText: 'Cancelar'
        });

        if (result.isConfirmed) {
            // Add loading state
            processBtn.classList.add('loading');
            processBtn.disabled = true;
            processBtn.innerHTML = 'Processando...';

            try {
                const response = await fetch(`${window.location.origin}/process-kyc/${item.package_id}`, {
                    method: 'GET'
                });

                if (response.ok) {
                    Swal.fire({
                        icon: 'success',
                        title: 'Sucesso!',
                        text: 'Processamento KYC iniciado com sucesso.',
                        confirmButtonColor: 'var(--success-color)'
                    });

                    // Refresh table
                    renderTable();
                } else {
                    throw new Error('Falha na requisição');
                }
            } catch (error) {
                console.error('KYC processing error:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro no Processamento',
                    text: 'Ocorreu um erro ao processar o KYC. Tente novamente.',
                    confirmButtonColor: 'var(--danger-color)'
                });
            } finally {
                // Remove loading state
                processBtn.classList.remove('loading');
                processBtn.disabled = false;
                processBtn.innerHTML = `
                    <svg style="width: 14px; height: 14px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                    KYC
                `;
            }
        }
    });

    // Download Button
    const downloadBtn = document.createElement('button');
    downloadBtn.className = 'action-btn btn-download';
    downloadBtn.setAttribute('data-tooltip', 'Baixar relatório');
    downloadBtn.innerHTML = `
        <svg style="width: 14px; height: 14px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
        </svg>
        Download
    `;

    downloadBtn.addEventListener('click', async (e) => {
        e.preventDefault();

        // Add loading state
        downloadBtn.classList.add('loading');
        downloadBtn.disabled = true;
        downloadBtn.innerHTML = 'Baixando...';

        try {
            // Simulate download (replace with actual download logic)
            await new Promise(resolve => setTimeout(resolve, 1500));

            Swal.fire({
                icon: 'success',
                title: 'Download Iniciado',
                text: 'O arquivo será baixado em breve.',
                confirmButtonColor: 'var(--success-color)'
            });
        } catch (error) {
            Swal.fire({
                icon: 'error',
                title: 'Erro no Download',
                text: 'Não foi possível baixar o arquivo.',
                confirmButtonColor: 'var(--danger-color)'
            });
        } finally {
            // Remove loading state
            downloadBtn.classList.remove('loading');
            downloadBtn.disabled = false;
            downloadBtn.innerHTML = `
                <svg style="width: 14px; height: 14px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                Download
            `;
        }
    });

    // View Details Button
    const viewBtn = document.createElement('button');
    viewBtn.className = 'action-btn btn-view';
    viewBtn.setAttribute('data-tooltip', 'Ver detalhes');
    viewBtn.innerHTML = `
        <svg style="width: 14px; height: 14px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
        </svg>
        Ver
    `;

    viewBtn.addEventListener('click', (e) => {
        e.preventDefault();

        Swal.fire({
            title: 'Detalhes do Pacote',
            html: `
                <div style="text-align: left; padding: 1rem;">
                    <p><strong>ID:</strong> ${item.package_id}</p>
                    <p><strong>Nome:</strong> ${item.package_name}</p>
                    <p><strong>Registros:</strong> ${item.completed}/${item.total_records}</p>
                    <p><strong>Agente:</strong> ${item.transfer_agent}</p>
                    <p><strong>Tipo:</strong> ${item.type_of_transfer}</p>
                    <p><strong>Criado:</strong> ${formatDate(item.created)}</p>
                    <p><strong>Atualizado:</strong> ${formatDate(item.updated)}</p>
                </div>
            `,
            confirmButtonColor: 'var(--primary-color)',
            confirmButtonText: 'Fechar'
        });
    });

    actionsContainer.appendChild(processBtn);
    actionsContainer.appendChild(downloadBtn);
    actionsContainer.appendChild(viewBtn);

    return actionsContainer;
};

const renderTable = async () => {
    // Show loading state
    tbody.innerHTML = `
        <tr>
            <td colspan="10" style="text-align: center; padding: 2rem;">
                <div class="loading" style="margin: 0 auto 1rem;"></div>
                <p style="color: var(--gray-600);">Carregando dados...</p>
            </td>
        </tr>
    `;

    // Fetch data
    const data = await fetchPackageList();

    // Clear existing table rows
    tbody.innerHTML = "";

    if (data.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="10" style="text-align: center; padding: 2rem; color: var(--gray-600);">
                    <svg style="width: 48px; height: 48px; margin: 0 auto 1rem; opacity: 0.5;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2 2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
                    </svg>
                    <p>Nenhum pacote encontrado</p>
                    <p style="font-size: 0.875rem; opacity: 0.7;">Faça upload de um arquivo para começar</p>
                </td>
            </tr>
        `;
        return;
    }

    // Populate table rows with data
    data.forEach((item, index) => {
        const row = document.createElement("tr");
        row.style.animationDelay = `${index * 0.05}s`;

        // Add alternating row colors
        if (index % 2 === 0) {
            row.style.backgroundColor = 'var(--white)';
        }

        const fields = [
            { key: "packageId", value: item.package_id, type: "text" },
            { key: "totalRecords", value: item.total_records, type: "number" },
            { key: "completed", value: getStatusBadge(item.completed, item.total_records), type: "html" },
            { key: "packageName", value: item.package_name, type: "text" },
            { key: "fullName", value: item.full_name, type: "text" },
            { key: "transferAgent", value: item.transfer_agent, type: "text" },
            { key: "typeOfTransfer", value: item.type_of_transfer, type: "text" },
            { key: "created", value: formatDate(item.created), type: "text" },
            { key: "updated", value: formatDate(item.updated), type: "text" },
            { key: "actions", value: createActionButtons(item), type: "element" }
        ];

        fields.forEach((field) => {
            const cell = document.createElement("td");

            if (field.type === "html") {
                cell.innerHTML = field.value;
            } else if (field.type === "element") {
                cell.appendChild(field.value);
            } else if (field.type === "number") {
                cell.textContent = field.value || "0";
                cell.style.fontWeight = "600";
                cell.style.color = "var(--primary-color)";
            } else {
                cell.textContent = field.value || "-";
            }

            if (field.key === "packageName") {
                cell.style.fontWeight = "600";
                cell.style.color = "var(--gray-800)";
            }

            row.appendChild(cell);
        });

        tbody.appendChild(row);
    });
};


// Render the table when the script runs
renderTable();