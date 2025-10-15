// KYC Clients Management
window.clientsData = window.clientsData || [];

// Load KYC clients when page is shown
function initKYCClients() {
    if (document.getElementById('checkKYC').classList.contains('active')) {
        loadKYCClients();
    }
}

// Load clients data from API
async function loadKYCClients() {
    const loadingElement = document.getElementById('loadingClients');
    const contentElement = document.getElementById('clientsTableContent');
    const statsElement = document.getElementById('clientsStats');

    // Show loading
    loadingElement.style.display = 'block';
    contentElement.style.display = 'none';
    statsElement.style.display = 'none';

    try {
        // Fetch clients data
        const response = await fetch('/api/kyc-clients');

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        window.clientsData = data.clients || [];

        // Render table
        renderClientsTable(window.clientsData);
        updateClientsStats(window.clientsData);

        // Show content
        loadingElement.style.display = 'none';
        contentElement.style.display = 'block';
        statsElement.style.display = 'block';

    } catch (error) {
        console.error('Error loading clients:', error);

        // Show error message
        loadingElement.innerHTML = `
            <div style="text-align: center; padding: 3rem; color: var(--error-color);">
                <div style="font-size: 3rem; margin-bottom: 1rem;">❌</div>
                <p>Erro ao carregar dados dos clientes</p>
                <p style="font-size: 0.9rem; margin-top: 0.5rem;">${error.message}</p>
                <button onclick="loadKYCClients()" class="btn-primary" style="margin-top: 1rem;">
                    🔄 Tentar Novamente
                </button>
            </div>
        `;
    }
}

// Render clients table
function renderClientsTable(clients) {
    const tbody = document.getElementById('clientsTableBody');
    const noClientsMessage = document.getElementById('noClientsMessage');

    if (!clients || clients.length === 0) {
        tbody.innerHTML = '';
        noClientsMessage.style.display = 'block';
        return;
    }

    noClientsMessage.style.display = 'none';

    tbody.innerHTML = clients.map(client => {
        const statusBadge = getKYCStatusBadge(client);
        const fullName = `${client.first_name || ''} ${client.middle_name || ''} ${client.last_name || ''}`.trim();
        const formattedDate = client.updated_at ? new Date(client.updated_at).toLocaleDateString('pt-BR') : 'N/A';
        const matchDetails = getMatchDetails(client);
        const reportsColumn = getReportsColumn(client);

        const allReportsColumn = getAMReportsColumn(client);

        return `
            <tr style="border-bottom: 1px solid var(--gray-200);">
                <td style="padding: 1rem; color: var(--gray-700);">#${client.id}</td>
                <td style="padding: 1rem; color: var(--gray-900); font-weight: 500;">${fullName}</td>
                <td style="padding: 1rem; color: var(--gray-700); font-family: monospace;">${client.client_reference_id || 'N/A'}</td>
                <td style="padding: 1rem; color: var(--gray-700);">${client.email || 'N/A'}</td>
                <td style="padding: 1rem;">${statusBadge}</td>
                <td style="padding: 0.5rem; font-size: 0.9rem;">${matchDetails}</td>
                <td style="padding: 0.5rem; font-size: 0.9rem;">${reportsColumn}</td>
                <td style="padding: 0.5rem; font-size: 0.9rem;">${allReportsColumn}</td>
                <td style="padding: 1rem; color: var(--gray-700);">${formattedDate}</td>
                <td style="padding: 1rem; text-align: center;">
                    <div style="display: flex; gap: 0.5rem; justify-content: center; align-items: center;">
                        <button onclick="downloadClientPDF(${client.id}, this)"
                                class="btn-secondary"
                                style="padding: 0.5rem 1rem; font-size: 0.8rem; display: flex; align-items: center; gap: 0.3rem;"
                                title="Baixar Relatório PDF">
                            📄 PDF
                        </button>
                        ${(client.am_reports_count > 0 || client.wl_reports_count > 0 || client.pep_reports_count > 0) ? `
                        <button onclick="downloadAllReports(${client.id})"
                                class="btn-secondary"
                                style="padding: 0.5rem 1rem; font-size: 0.8rem; display: flex; align-items: center; gap: 0.3rem; background: #10b981; color: white;"
                                title="Baixar Reports">
                            📥 Reports
                        </button>
                        ` : ''}
                        <button onclick="viewClientDetails(${client.id})"
                                class="btn-outline"
                                style="padding: 0.5rem 1rem; font-size: 0.8rem; display: flex; align-items: center; gap: 0.3rem;"
                                title="Ver Detalhes">
                            👁️ Ver
                        </button>
                        ${client.complete_kyc ? `
                        <button onclick="processClientKYC(${client.id})"
                                class="btn-primary"
                                style="padding: 0.5rem 1rem; font-size: 0.8rem; display: flex; align-items: center; gap: 0.3rem;"
                                title="Reprocessar KYC">
                            🔄 KYC
                        </button>
                        ` : `
                        <button onclick="processClientKYC(${client.id})"
                                class="btn-primary"
                                style="padding: 0.5rem 1rem; font-size: 0.8rem; display: flex; align-items: center; gap: 0.3rem;"
                                title="Processar KYC">
                            ▶️ KYC
                        </button>
                        `}
                    </div>
                </td>
            </tr>
        `;
    }).join('');
}

// Get all reports column for client
function getAMReportsColumn(client) {
    // Use persisted data if available
    if (client.am_reports_count !== undefined || client.wl_reports_count !== undefined || client.pep_reports_count !== undefined) {
        const amCount = client.am_reports_count || 0;
        const wlCount = client.wl_reports_count || 0;
        const pepCount = client.pep_reports_count || 0;
        const totalCount = amCount + wlCount + pepCount;
        
        if (totalCount > 0) {
            return `
                <div style="display: flex; flex-direction: column; gap: 0.25rem;">
                    <div style="font-size: 0.75rem; color: var(--gray-600);">
                        ${amCount > 0 ? `AM: ${amCount}` : ''}
                        ${wlCount > 0 ? `${amCount > 0 ? ' | ' : ''}WL: ${wlCount}` : ''}
                        ${pepCount > 0 ? `${(amCount > 0 || wlCount > 0) ? ' | ' : ''}PEP: ${pepCount}` : ''}
                    </div>
                    <button onclick="downloadAllReports(${client.id})"
                            style="background: none; border: 1px solid var(--gray-300); border-radius: 4px; padding: 0.25rem 0.5rem; font-size: 0.7rem; cursor: pointer; color: var(--gray-600);"
                            title="Baixar Todos Reports">
                        📥 Download
                    </button>
                </div>
            `;
        }
        return '<span style="color: var(--gray-400); font-style: italic;">Sem reports</span>';
    }

    // Fallback to parsing response data
    if (!client.response) {
        return '<span style="color: var(--gray-400); font-style: italic;">Sem reports</span>';
    }

    try {
        const responseData = JSON.parse(client.response);
        const counts = extractReportsCounts(responseData);
        const totalCount = counts.am + counts.wl + counts.pep;
        
        if (totalCount > 0) {
            return `
                <div style="display: flex; flex-direction: column; gap: 0.25rem;">
                    <div style="font-size: 0.75rem; color: var(--gray-600);">
                        ${counts.am > 0 ? `AM: ${counts.am}` : ''}
                        ${counts.wl > 0 ? `${counts.am > 0 ? ' | ' : ''}WL: ${counts.wl}` : ''}
                        ${counts.pep > 0 ? `${(counts.am > 0 || counts.wl > 0) ? ' | ' : ''}PEP: ${counts.pep}` : ''}
                    </div>
                    <button onclick="downloadAllReports(${client.id})"
                            style="background: none; border: 1px solid var(--gray-300); border-radius: 4px; padding: 0.25rem 0.5rem; font-size: 0.7rem; cursor: pointer; color: var(--gray-600);"
                            title="Baixar Todos Reports">
                        📥 Download
                    </button>
                </div>
            `;
        }
        
        return '<span style="color: var(--gray-400); font-style: italic;">Sem reports</span>';
        
    } catch (error) {
        console.error('Error parsing reports data:', error);
        return '<span style="color: var(--gray-400); font-style: italic;">Erro nos dados</span>';
    }
}

// Extract reports counts from response data
function extractReportsCounts(responseData) {
    let amCount = 0, wlCount = 0, pepCount = 0;
    
    if (responseData.flowData) {
        const flowId = Object.keys(responseData.flowData)[0];
        const flowData = responseData.flowData[flowId];
        
        if (flowData && flowData.serviceData) {
            flowData.serviceData.forEach(service => {
                if (service.fullServiceDetails && service.fullServiceDetails.Record && service.fullServiceDetails.Record.DatasourceResults) {
                    service.fullServiceDetails.Record.DatasourceResults.forEach(datasource => {
                        if (datasource.DatasourceFields) {
                            datasource.DatasourceFields.forEach(field => {
                                if (field.FieldName === 'WatchlistHitDetails' && field.Data) {
                                    if (field.Data.AM_results) amCount += field.Data.AM_results.length;
                                    if (field.Data.WL_results) wlCount += field.Data.WL_results.length;
                                    if (field.Data.PEP_results) pepCount += field.Data.PEP_results.length;
                                }
                            });
                        }
                    });
                }
            });
        }
    }
    
    return { am: amCount, wl: wlCount, pep: pepCount };
}

// Download all reports as text file
function downloadAllReports(clientId) {
    const client = window.clientsData.find(c => c.id === clientId);
    if (!client) return;

    try {
        const fullName = `${client.first_name || ''} ${client.middle_name || ''} ${client.last_name || ''}`.trim();
        let reportsText = `ALL REPORTS - ${fullName} (ID: ${client.id})\n`;
        reportsText += `Data: ${new Date().toLocaleString('pt-BR')}\n`;
        reportsText += `${'='.repeat(80)}\n\n`;

        let totalReports = 0;

        // Use persisted data if available
        if (client.reports_data) {
            try {
                const reportsData = JSON.parse(client.reports_data);
                
                ['AM_results', 'WL_results', 'PEP_results'].forEach(reportType => {
                    if (reportsData[reportType] && reportsData[reportType].length > 0) {
                        reportsText += `${reportType.replace('_results', '').toUpperCase()} REPORTS\n`;
                        reportsText += `${'='.repeat(40)}\n\n`;
                        
                        reportsData[reportType].forEach((report, index) => {
                            totalReports++;
                            reportsText += `${reportType.replace('_results', '').toUpperCase()} REPORT ${index + 1}\n`;
                            reportsText += `${'-'.repeat(30)}\n`;
                            reportsText += `Score: ${report.score || 'N/A'}\n`;
                            reportsText += `Subject Matched: ${report.subjectMatched || 'N/A'}\n`;
                            reportsText += `URL: ${report.URL || 'N/A'}\n\n`;
                            reportsText += `Text:\n${report.text || 'N/A'}\n\n`;
                            reportsText += `${'-'.repeat(60)}\n\n`;
                        });
                    }
                });
            } catch (e) {
                console.error('Error parsing persisted reports data:', e);
            }
        }

        // Fallback to parsing response data
        if (totalReports === 0 && client.response) {
            const responseData = JSON.parse(client.response);
            
            if (responseData.flowData) {
                const flowId = Object.keys(responseData.flowData)[0];
                const flowData = responseData.flowData[flowId];
                
                if (flowData && flowData.serviceData) {
                    flowData.serviceData.forEach(service => {
                        if (service.fullServiceDetails && service.fullServiceDetails.Record && service.fullServiceDetails.Record.DatasourceResults) {
                            service.fullServiceDetails.Record.DatasourceResults.forEach(datasource => {
                                if (datasource.DatasourceFields) {
                                    datasource.DatasourceFields.forEach(field => {
                                        if (field.FieldName === 'WatchlistHitDetails' && field.Data) {
                                            ['AM_results', 'WL_results', 'PEP_results'].forEach(reportType => {
                                                if (field.Data[reportType] && field.Data[reportType].length > 0) {
                                                    reportsText += `${reportType.replace('_results', '').toUpperCase()} REPORTS\n`;
                                                    reportsText += `${'='.repeat(40)}\n\n`;
                                                    
                                                    field.Data[reportType].forEach((report, index) => {
                                                        totalReports++;
                                                        reportsText += `${reportType.replace('_results', '').toUpperCase()} REPORT ${index + 1}\n`;
                                                        reportsText += `${'-'.repeat(30)}\n`;
                                                        reportsText += `Score: ${report.score || 'N/A'}\n`;
                                                        reportsText += `Subject Matched: ${report.subjectMatched || 'N/A'}\n`;
                                                        reportsText += `URL: ${report.URL || 'N/A'}\n\n`;
                                                        reportsText += `Text:\n${report.text || 'N/A'}\n\n`;
                                                        reportsText += `${'-'.repeat(60)}\n\n`;
                                                    });
                                                }
                                            });
                                        }
                                    });
                                }
                            });
                        }
                    });
                }
            }
        }

        if (totalReports === 0) {
            reportsText += 'Nenhum report encontrado.\n';
        }

        // Create and download file
        const blob = new Blob([reportsText], { type: 'text/plain;charset=utf-8' });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `all_reports_${client.id}_${fullName.replace(/\s+/g, '_')}.txt`;
        document.body.appendChild(a);
        a.click();
        
        // Clean up
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);

        // Show success message
        if (typeof Swal !== 'undefined') {
            Swal.fire({
                icon: 'success',
                title: 'Reports Baixados!',
                text: `${totalReports} report(s) foram salvos em arquivo de texto.`,
                timer: 3000,
                showConfirmButton: false
            });
        }

    } catch (error) {
        console.error('Error downloading reports:', error);
        if (typeof Swal !== 'undefined') {
            Swal.fire({
                icon: 'error',
                title: 'Erro',
                text: 'Não foi possível baixar os reports.',
            });
        }
    }
}

// Download AM reports as text file
function downloadAMReports(clientId) {
    const client = window.clientsData.find(c => c.id === clientId);
    if (!client || !client.response) return;

    try {
        const responseData = JSON.parse(client.response);
        const fullName = `${client.first_name || ''} ${client.middle_name || ''} ${client.last_name || ''}`.trim();
        let amReportsText = `AM REPORTS - ${fullName} (ID: ${client.id})\n`;
        amReportsText += `Data: ${new Date().toLocaleString('pt-BR')}\n`;
        amReportsText += `${'='.repeat(80)}\n\n`;

        let reportCount = 0;

        if (responseData.flowData) {
            const flowId = Object.keys(responseData.flowData)[0];
            const flowData = responseData.flowData[flowId];
            
            if (flowData && flowData.serviceData) {
                flowData.serviceData.forEach((service, serviceIndex) => {
                    if (service.fullServiceDetails && service.fullServiceDetails.Record && service.fullServiceDetails.Record.DatasourceResults) {
                        service.fullServiceDetails.Record.DatasourceResults.forEach((datasource, dsIndex) => {
                            if (datasource.DatasourceFields) {
                                datasource.DatasourceFields.forEach(field => {
                                    if (field.FieldName === 'WatchlistHitDetails' && field.Data && field.Data.AM_results) {
                                        field.Data.AM_results.forEach((amResult, amIndex) => {
                                            reportCount++;
                                            amReportsText += `REPORT ${reportCount}\n`;
                                            amReportsText += `${'-'.repeat(40)}\n`;
                                            amReportsText += `Score: ${amResult.score || 'N/A'}\n`;
                                            amReportsText += `Subject Matched: ${amResult.subjectMatched || 'N/A'}\n`;
                                            amReportsText += `URL: ${amResult.URL || 'N/A'}\n\n`;
                                            amReportsText += `Text:\n${amResult.text || 'N/A'}\n\n`;
                                            amReportsText += `${'='.repeat(80)}\n\n`;
                                        });
                                    }
                                });
                            }
                        });
                    }
                });
            }
        }

        if (reportCount === 0) {
            amReportsText += 'Nenhum AM report encontrado.\n';
        }

        // Create and download file
        const blob = new Blob([amReportsText], { type: 'text/plain;charset=utf-8' });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `am_reports_${client.id}_${fullName.replace(/\s+/g, '_')}.txt`;
        document.body.appendChild(a);
        a.click();
        
        // Clean up
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);

        // Show success message
        if (typeof Swal !== 'undefined') {
            Swal.fire({
                icon: 'success',
                title: 'AM Reports Baixados!',
                text: `${reportCount} report(s) foram salvos em arquivo de texto.`,
                timer: 3000,
                showConfirmButton: false
            });
        }

    } catch (error) {
        console.error('Error downloading AM reports:', error);
        if (typeof Swal !== 'undefined') {
            Swal.fire({
                icon: 'error',
                title: 'Erro',
                text: 'Não foi possível baixar os AM reports.',
            });
        }
    }
}

// Get reports column for client
function getReportsColumn(client) {
    if (!client.response) {
        return '<span style="color: var(--gray-400); font-style: italic;">Sem relatórios</span>';
    }

    try {
        const responseData = JSON.parse(client.response);
        
        if (responseData.flowData) {
            const flowId = Object.keys(responseData.flowData)[0];
            const flowData = responseData.flowData[flowId];
            
            if (flowData && flowData.serviceData && flowData.serviceData.length > 0) {
                return `
                    <div style="display: flex; flex-direction: column; gap: 0.25rem;">
                        <span style="color: var(--gray-700); font-size: 0.8rem;">${flowData.serviceData.length} relatório(s)</span>
                        <button onclick="showReports(${client.id})"
                                style="background: none; border: 1px solid var(--gray-300); border-radius: 4px; padding: 0.25rem 0.5rem; font-size: 0.7rem; cursor: pointer; color: var(--gray-600);"
                                title="Ver relatórios">
                            📊 Relatórios
                        </button>
                    </div>
                `;
            }
        }
        
        return '<span style="color: var(--gray-400); font-style: italic;">Processando...</span>';
        
    } catch (error) {
        console.error('Error parsing reports data:', error);
        return '<span style="color: var(--gray-400); font-style: italic;">Erro nos dados</span>';
    }
}

// Show reports for a client
function showReports(clientId) {
    const client = window.clientsData.find(c => c.id === clientId);
    if (!client || !client.response) return;

    try {
        const responseData = JSON.parse(client.response);
        const fullName = `${client.first_name || ''} ${client.middle_name || ''} ${client.last_name || ''}`.trim();

        let reportsHtml = `<div style="text-align: left; max-width: 600px; margin: 0 auto;">`;

        if (responseData.flowData) {
            const flowId = Object.keys(responseData.flowData)[0];
            const flowData = responseData.flowData[flowId];

            if (flowData && flowData.serviceData) {
                flowData.serviceData.forEach((service, index) => {
                    reportsHtml += `
                        <div style="margin-bottom: 1rem; padding: 1rem; border: 1px solid var(--gray-200); border-radius: 8px; background: var(--gray-50);">
                            <h4 style="margin: 0 0 0.5rem 0; color: var(--gray-800);">Relatório ${index + 1}: ${service.nodeTitle || 'Verificação'}</h4>
                            <p style="margin: 0 0 0.5rem 0; font-size: 0.9rem;"><strong>Status:</strong> <span style="color: ${service.serviceStatus === 'COMPLETED' ? '#10b981' : '#ef4444'}">${service.serviceStatus}</span></p>
                            <p style="margin: 0; font-size: 0.9rem;"><strong>Match:</strong> <span style="color: ${service.match ? '#10b981' : '#ef4444'}">${service.match ? 'Sim' : 'Não'}</span></p>
                        </div>
                    `;
                });
            }
        }

        reportsHtml += `</div>`;

        if (typeof Swal !== 'undefined') {
            Swal.fire({
                title: `Relatórios - ${fullName}`,
                html: reportsHtml,
                showCloseButton: true,
                confirmButtonText: 'Fechar',
                width: '700px'
            });
        } else {
            console.log('Reports for client:', client);
        }

    } catch (error) {
        console.error('Error showing reports:', error);
        if (typeof Swal !== 'undefined') {
            Swal.fire({
                icon: 'error',
                title: 'Erro',
                text: 'Não foi possível carregar os relatórios.',
            });
        }
    }
}

// Get match details from Trulioo response
function getMatchDetails(client) {
    if (!client.response) {
        return '<span style="color: var(--gray-400); font-style: italic;">Dados não disponíveis</span>';
    }

    try {
        const responseData = JSON.parse(client.response);

        // Extract match data from flowData service results
        if (responseData.flowData) {
            const flowId = Object.keys(responseData.flowData)[0];
            const flowData = responseData.flowData[flowId];

            if (flowData && flowData.serviceData) {
                let totalMatches = 0;
                let totalFields = 0;
                let missingFields = 0;
                let matchedFields = [];
                let missingFieldsList = [];

                // Process each service data
                flowData.serviceData.forEach(service => {
                    if (service.fullServiceDetails && service.fullServiceDetails.Record && service.fullServiceDetails.Record.DatasourceResults) {
                        service.fullServiceDetails.Record.DatasourceResults.forEach(datasource => {
                            if (datasource.DatasourceFields) {
                                datasource.DatasourceFields.forEach(field => {
                                    totalFields++;
                                    if (field.Status === 'match') {
                                        totalMatches++;
                                        matchedFields.push(field.FieldName);
                                    } else if (field.Status === 'missing') {
                                        missingFields++;
                                        missingFieldsList.push(field.FieldName);
                                    }
                                });
                            }
                        });
                    }
                });

                if (totalFields > 0) {
                    const matchPercentage = Math.round((totalMatches / totalFields) * 100);
                    const matchColor = matchPercentage >= 80 ? '#10b981' : matchPercentage >= 60 ? '#f59e0b' : '#ef4444';

                    return `
                        <div style="display: flex; flex-direction: column; gap: 0.25rem;">
                            <div style="display: flex; align-items: center; gap: 0.5rem;">
                                <span style="color: ${matchColor}; font-weight: 600;">${totalMatches}/${totalFields}</span>
                                <span style="color: var(--gray-500); font-size: 0.8rem;">(${matchPercentage}%)</span>
                            </div>
                            ${missingFields > 0 ? `
                            <div style="color: #ef4444; font-size: 0.75rem;">
                                ${missingFields} campo(s) ausente(s)
                            </div>
                            ` : ''}
                            <button onclick="showMatchDetails(${client.id})"
                                    style="background: none; border: 1px solid var(--gray-300); border-radius: 4px; padding: 0.25rem 0.5rem; font-size: 0.7rem; cursor: pointer; color: var(--gray-600);"
                                    title="Ver detalhes completos">
                                🔍 Detalhes
                            </button>
                        </div>
                    `;
                }
            }
        }

        // Fallback - check basic status
        if (responseData.status === 'ACCEPTED') {
            return `
                <div style="display: flex; flex-direction: column; gap: 0.25rem;">
                    <span style="color: #10b981; font-weight: 500;">✅ Verificado</span>
                    <button onclick="showMatchDetails(${client.id})"
                            style="background: none; border: 1px solid var(--gray-300); border-radius: 4px; padding: 0.25rem 0.5rem; font-size: 0.7rem; cursor: pointer; color: var(--gray-600);"
                            title="Ver detalhes completos">
                        🔍 Detalhes
                    </button>
                </div>
            `;
        }

        return '<span style="color: var(--gray-400);">Processando...</span>';

    } catch (error) {
        console.error('Error parsing match details:', error);
        return '<span style="color: var(--gray-400); font-style: italic;">Erro ao analisar dados</span>';
    }
}

// Show detailed match information in a modal
function showMatchDetails(clientId) {
    const client = window.clientsData.find(c => c.id === clientId);
    if (!client || !client.response) return;

    try {
        const responseData = JSON.parse(client.response);
        const fullName = `${client.first_name || ''} ${client.middle_name || ''} ${client.last_name || ''}`.trim();

        let detailsHtml = `<div style="text-align: left; max-width: 600px; margin: 0 auto;">`;

        // Show basic status information
        detailsHtml += `
            <div style="margin-bottom: 1.5rem; padding: 1rem; border: 1px solid var(--gray-200); border-radius: 8px; background: var(--gray-50);">
                <h4 style="margin: 0 0 1rem 0; color: var(--gray-800);">Status Geral</h4>
                <p style="margin: 0 0 0.5rem 0; font-size: 0.9rem;"><strong>Status:</strong> <span style="color: ${responseData.status === 'ACCEPTED' ? '#10b981' : '#ef4444'}">${responseData.status}</span></p>
                <p style="margin: 0 0 0.5rem 0; font-size: 0.9rem;"><strong>Tipo:</strong> ${responseData.profileType}</p>
                <p style="margin: 0; font-size: 0.9rem;"><strong>Data:</strong> ${new Date(responseData.created * 1000).toLocaleString('pt-BR')}</p>
            </div>
        `;

        if (responseData.flowData) {
            const flowId = Object.keys(responseData.flowData)[0];
            const flowData = responseData.flowData[flowId];

            // Show field data information
            if (flowData && flowData.fieldData) {
                detailsHtml += `
                    <div style="margin-bottom: 1.5rem; padding: 1rem; border: 1px solid var(--gray-200); border-radius: 8px; background: var(--gray-50);">
                        <h4 style="margin: 0 0 1rem 0; color: var(--gray-800);">Dados Processados</h4>
                        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.5rem; font-size: 0.85rem;">
                `;

                Object.values(flowData.fieldData).forEach(field => {
                    if (field.name && field.value && field.value[0]) {
                        detailsHtml += `
                            <div style="display: flex; justify-content: space-between; padding: 0.25rem; border-bottom: 1px solid var(--gray-200);">
                                <span style="color: var(--gray-700);">${field.name}</span>
                                <span style="color: var(--gray-900); font-weight: 500;">${field.value[0]}</span>
                            </div>
                        `;
                    }
                });

                detailsHtml += `</div></div>`;
            }

            if (flowData && flowData.serviceData) {
                flowData.serviceData.forEach((service, serviceIndex) => {
                    detailsHtml += `
                        <div style="margin-bottom: 1.5rem; padding: 1rem; border: 1px solid var(--gray-200); border-radius: 8px; background: var(--gray-50);">
                            <h4 style="margin: 0 0 1rem 0; color: var(--gray-800);">${service.nodeTitle || `Serviço ${serviceIndex + 1}`}</h4>
                            <p style="margin: 0 0 0.5rem 0; font-size: 0.9rem;"><strong>Status:</strong> <span style="color: ${service.serviceStatus === 'COMPLETED' ? '#10b981' : '#ef4444'}">${service.serviceStatus}</span></p>
                            <p style="margin: 0 0 1rem 0; font-size: 0.9rem;"><strong>Match:</strong> <span style="color: ${service.match ? '#10b981' : '#ef4444'}">${service.match ? 'Sim' : 'Não'}</span></p>
                    `;

                    if (service.fullServiceDetails && service.fullServiceDetails.Record && service.fullServiceDetails.Record.DatasourceResults) {
                        service.fullServiceDetails.Record.DatasourceResults.forEach(datasource => {
                            if (datasource.DatasourceFields && datasource.DatasourceFields.length > 0) {
                                detailsHtml += `<h5 style="margin: 1rem 0 0.5rem 0; color: var(--gray-700);">${datasource.DatasourceName}:</h5>`;
                                detailsHtml += `<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.5rem; font-size: 0.85rem;">`;

                                datasource.DatasourceFields.forEach(field => {
                                    const statusColor = field.Status === 'match' ? '#10b981' : field.Status === 'missing' ? '#ef4444' : '#f59e0b';
                                    const statusIcon = field.Status === 'match' ? '✅' : field.Status === 'missing' ? '❌' : '⚠️';

                                    detailsHtml += `
                                        <div style="display: flex; justify-content: space-between; padding: 0.25rem; border-bottom: 1px solid var(--gray-200);">
                                            <span style="color: var(--gray-700);">${field.FieldName}</span>
                                            <span style="color: ${statusColor};">${statusIcon} ${field.Status}</span>
                                        </div>
                                    `;
                                });

                                detailsHtml += `</div>`;
                            }
                        });
                    }

                    detailsHtml += `</div>`;
                });
            }
        }

        detailsHtml += `</div>`;

        if (typeof Swal !== 'undefined') {
            Swal.fire({
                title: `Detalhes do Match - ${fullName}`,
                html: detailsHtml,
                showCloseButton: true,
                confirmButtonText: 'Fechar',
                width: '800px'
            });
        } else {
            alert('Detalhes do match disponíveis no console');
            console.log('Match details for client:', client);
        }

    } catch (error) {
        console.error('Error showing match details:', error);
        if (typeof Swal !== 'undefined') {
            Swal.fire({
                icon: 'error',
                title: 'Erro',
                text: 'Não foi possível carregar os detalhes do match.',
            });
        }
    }
}

// Get status badge HTML
function getKYCStatusBadge(client) {
    let status = 'pending';
    let statusText = 'Pendente';

    if (client.complete_kyc) {
        status = 'completed';
        statusText = 'Processado';
    }

    // Try to extract status from response if available
    if (client.response) {
        try {
            const responseData = JSON.parse(client.response);
            if (responseData.status === 'ACCEPTED') {
                status = 'accepted';
                statusText = 'Aceito';
            }
        } catch (e) {
            // Ignore JSON parse errors
        }
    }

    const statusConfig = {
        'pending': { color: '#f59e0b', bg: '#fef3c7', text: 'Pendente' },
        'processing': { color: '#3b82f6', bg: '#dbeafe', text: 'Processando' },
        'completed': { color: '#10b981', bg: '#d1fae5', text: 'Processado' },
        'accepted': { color: '#10b981', bg: '#d1fae5', text: 'Aceito' },
        'rejected': { color: '#ef4444', bg: '#fee2e2', text: 'Rejeitado' },
        'error': { color: '#ef4444', bg: '#fee2e2', text: 'Erro' }
    };

    const config = statusConfig[status] || statusConfig['pending'];

    return `
        <span style="
            display: inline-flex;
            align-items: center;
            padding: 0.25rem 0.75rem;
            font-size: 0.75rem;
            font-weight: 500;
            border-radius: 9999px;
            color: ${config.color};
            background-color: ${config.bg};
        ">
            ${config.text}
        </span>
    `;
}

// Update clients statistics
function updateClientsStats(clients) {
    const total = clients.length;
    const completed = clients.filter(c => c.complete_kyc).length;
    const pending = total - completed;

    document.getElementById('totalClients').textContent = total;
    document.getElementById('completedClients').textContent = completed;
    document.getElementById('pendingClients').textContent = pending;
}

// Download PDF report for a client
async function downloadClientPDF(clientId, buttonElement) {
    const button = buttonElement;
    const originalText = button.innerHTML;

    try {
        // Show loading state
        button.innerHTML = '⏳ Gerando...';
        button.disabled = true;

        // Make request to PDF endpoint
        const response = await fetch(`/kyc-pdf/${clientId}`);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        // Get the PDF blob
        const blob = await response.blob();

        // Create download link
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `kyc_report_${clientId}.pdf`;
        document.body.appendChild(a);
        a.click();

        // Clean up
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);

        // Show success message
        if (typeof Swal !== 'undefined') {
            Swal.fire({
                icon: 'success',
                title: 'PDF Gerado!',
                text: 'O relatório PDF foi baixado com sucesso.',
                timer: 3000,
                showConfirmButton: false
            });
        }

    } catch (error) {
        console.error('Error downloading PDF:', error);

        if (typeof Swal !== 'undefined') {
            Swal.fire({
                icon: 'error',
                title: 'Erro ao Gerar PDF',
                text: 'Não foi possível gerar o relatório PDF. Tente novamente.',
                confirmButtonText: 'OK'
            });
        } else {
            alert('Erro ao gerar PDF: ' + error.message);
        }
    } finally {
        // Restore button state
        button.innerHTML = originalText;
        button.disabled = false;
    }
}

// View client details
function viewClientDetails(clientId) {
    const client = window.clientsData.find(c => c.id === clientId);
    if (!client) return;

    const fullName = `${client.first_name || ''} ${client.middle_name || ''} ${client.last_name || ''}`.trim();

    if (typeof Swal !== 'undefined') {
        Swal.fire({
            title: `Cliente: ${fullName}`,
            html: `
                <div style="text-align: left; max-width: 400px; margin: 0 auto;">
                    <p><strong>ID:</strong> #${client.id}</p>
                    <p><strong>Referência:</strong> ${client.client_reference_id || 'N/A'}</p>
                    <p><strong>Email:</strong> ${client.email || 'N/A'}</p>
                    <p><strong>Telefone:</strong> ${client.personal_phone_number || 'N/A'}</p>
                    <p><strong>Data de Nascimento:</strong> ${client.date_of_birth_day ? new Date(client.date_of_birth_day).toLocaleDateString('pt-BR') : 'N/A'}</p>
                    <p><strong>Endereço:</strong> ${client.street_address || 'N/A'}</p>
                    <p><strong>Cidade:</strong> ${client.city || 'N/A'}</p>
                    <p><strong>Estado:</strong> ${client.letter_state || 'N/A'}</p>
                    <p><strong>CEP:</strong> ${client.postal || 'N/A'}</p>
                    <p><strong>País:</strong> ${client.letter_country || 'N/A'}</p>
                    <p><strong>KYC Processado:</strong> ${client.complete_kyc ? 'Sim' : 'Não'}</p>
                </div>
            `,
            showCloseButton: true,
            confirmButtonText: 'Fechar'
        });
    } else {
        alert(`Detalhes do Cliente: ${fullName}\nID: #${client.id}\nEmail: ${client.email || 'N/A'}`);
    }
}

// Process KYC for a client
function processClientKYC(clientId) {
    if (typeof Swal !== 'undefined') {
        Swal.fire({
            title: 'Processar KYC',
            text: 'Deseja iniciar o processamento KYC para este cliente?',
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Sim, Processar',
            cancelButtonText: 'Cancelar'
        }).then((result) => {
            if (result.isConfirmed) {
                // TODO: Implement KYC processing
                Swal.fire({
                    icon: 'info',
                    title: 'Funcionalidade em Desenvolvimento',
                    text: 'O processamento KYC individual será implementado em breve.',
                    timer: 3000,
                    showConfirmButton: false
                });
            }
        });
    } else {
        if (confirm('Deseja iniciar o processamento KYC para este cliente?')) {
            alert('Funcionalidade em desenvolvimento');
        }
    }
}

// Filter clients based on search and status
function filterClients() {
    const searchTerm = document.getElementById('clientSearch').value.toLowerCase();
    const statusFilter = document.getElementById('statusFilter').value;

    let filteredClients = window.clientsData;

    // Apply search filter
    if (searchTerm) {
        filteredClients = filteredClients.filter(client => {
            const fullName = `${client.first_name || ''} ${client.middle_name || ''} ${client.last_name || ''}`.toLowerCase();
            const clientRef = (client.client_reference_id || '').toLowerCase();
            const email = (client.email || '').toLowerCase();

            return fullName.includes(searchTerm) ||
                   clientRef.includes(searchTerm) ||
                   email.includes(searchTerm) ||
                   client.id.toString().includes(searchTerm);
        });
    }

    // Apply status filter
    if (statusFilter) {
        filteredClients = filteredClients.filter(client => {
            if (statusFilter === 'pending') return !client.complete_kyc;
            if (statusFilter === 'completed') return client.complete_kyc;
            // TODO: Add more status filtering based on actual response data
            return true;
        });
    }

    renderClientsTable(filteredClients);
}

// Add event listeners for filters
document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('clientSearch');
    const statusSelect = document.getElementById('statusFilter');

    if (searchInput) {
        searchInput.addEventListener('input', filterClients);
    }

    if (statusSelect) {
        statusSelect.addEventListener('change', filterClients);
    }
});

// Auto-load clients when the page becomes active
const observer = new MutationObserver(function(mutations) {
    mutations.forEach(function(mutation) {
        if (mutation.type === 'attributes' && mutation.attributeName === 'class') {
            const target = mutation.target;
            if (target.id === 'checkKYC' && target.classList.contains('active')) {
                setTimeout(initKYCClients, 100);
            }
        }
    });
});

// Start observing
document.addEventListener('DOMContentLoaded', function() {
    const checkKYCPage = document.getElementById('checkKYC');
    if (checkKYCPage) {
        observer.observe(checkKYCPage, { attributes: true });
    }
});