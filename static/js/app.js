// Modern File Upload with Drag & Drop
document.addEventListener('DOMContentLoaded', function() {
    const fileInput = document.getElementById("fileInput");
    const uploadButton = document.getElementById("uploadButton");
    const result = document.getElementById("result");
    const fileNameDisplay = document.getElementById("file-name");
    const dropZone = document.getElementById("dropZone");

    let selectedFile = null;

    // File validation
    function validateFile(file) {
        const validTypes = ['application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'application/vnd.ms-excel'];
        const maxSize = 10 * 1024 * 1024; // 10MB

        if (!validTypes.includes(file.type)) {
            throw new Error('Formato de arquivo inválido. Apenas arquivos Excel (.xlsx, .xls) são aceitos.');
        }

        if (file.size > maxSize) {
            throw new Error('Arquivo muito grande. Tamanho máximo: 10MB.');
        }

        return true;
    }

    // Update UI when file is selected
    function handleFileSelect(file) {
        try {
            validateFile(file);
            selectedFile = file;

            // Update file name display
            fileNameDisplay.textContent = `📄 ${file.name} (${formatFileSize(file.size)})`;
            fileNameDisplay.style.color = 'var(--success-color)';

            // Enable upload button
            uploadButton.disabled = false;
            uploadButton.textContent = 'Enviar Arquivo';

            // Update drop zone appearance
            dropZone.style.borderColor = 'var(--success-color)';
            dropZone.style.background = 'rgba(16, 185, 129, 0.05)';

        } catch (error) {
            showError(error.message);
            resetFileSelection();
        }
    }

    // Reset file selection
    function resetFileSelection() {
        selectedFile = null;
        fileInput.value = '';
        fileNameDisplay.textContent = '';
        uploadButton.disabled = true;
        uploadButton.textContent = 'Selecione um arquivo';
        dropZone.style.borderColor = 'var(--gray-300)';
        dropZone.style.background = '';
    }

    // Format file size
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    // Show success message
    function showSuccess(message) {
        result.innerHTML = `<div class="success">${message}</div>`;
    }

    // Show error message
    function showError(message) {
        result.innerHTML = `<div class="error">${message}</div>`;
    }

    // File input change event
    fileInput.addEventListener("change", (event) => {
        const file = event.target.files[0];
        if (file) {
            handleFileSelect(file);
        } else {
            resetFileSelection();
        }
    });

    // Drag and drop events
    dropZone.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropZone.classList.add('dragover');
    });

    dropZone.addEventListener('dragleave', (e) => {
        e.preventDefault();
        dropZone.classList.remove('dragover');
    });

    dropZone.addEventListener('drop', (e) => {
        e.preventDefault();
        dropZone.classList.remove('dragover');

        const files = e.dataTransfer.files;
        if (files.length > 0) {
            const file = files[0];
            handleFileSelect(file);
        }
    });

    // Upload button click event
    uploadButton.addEventListener("click", async () => {
        if (!selectedFile) {
            Swal.fire({
                icon: "warning",
                title: "Nenhum arquivo selecionado",
                text: "Por favor, selecione um arquivo antes de enviar.",
                confirmButtonColor: 'var(--primary-color)'
            });
            return;
        }

        // Show loading state
        uploadButton.disabled = true;
        uploadButton.innerHTML = '<div class="loading"></div> Enviando...';
        result.innerHTML = '';

        const formData = new FormData();
        formData.append("file_name", selectedFile.name);
        formData.append("file", selectedFile);

        try {
            const response = await fetch(`${window.location.origin}/kyc-request`, {
                method: "POST",
                body: formData,
            });

            if (response.ok) {
                const data = await response.json();

                Swal.fire({
                    icon: "success",
                    title: "Upload Realizado com Sucesso!",
                    text: data.message,
                    confirmButtonColor: 'var(--success-color)'
                });

                showSuccess(`✅ ${data.message}`);
                resetFileSelection();

            } else {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Falha no upload do arquivo');
            }
        } catch (error) {
            console.error('Upload error:', error);

            Swal.fire({
                icon: "error",
                title: "Erro no Upload",
                text: error.message || "Ocorreu um erro ao enviar o arquivo. Tente novamente.",
                confirmButtonColor: 'var(--danger-color)'
            });

            showError(`❌ ${error.message || 'Erro ao enviar arquivo'}`);

            // Reset button state
            uploadButton.disabled = false;
            uploadButton.innerHTML = `
                <svg style="width: 16px; height: 16px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 12l2 2 4-4"></path>
                </svg>
                Enviar Arquivo
            `;
        }
    });

    // Initial state
    resetFileSelection();
});
