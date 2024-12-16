const fileInput = document.getElementById("fileInput");
const fileName = document.getElementById("fileName");
const uploadButton = document.getElementById("uploadButton");
const result = document.getElementById("result");
const label = document.getElementById("file-name")


let selectedFile = null;

// Exibe o nome do arquivo ao selecionar
fileInput.addEventListener("change", (event) => {
    selectedFile = event.target.files[0];
    if (selectedFile) {
        label.textContent = selectedFile.name;
    } else {
        label.textContent = "No file selected";
    }
});

// Envia o arquivo ao clicar no botão
uploadButton.addEventListener("click", async () => {
    if (!selectedFile) {
        result.textContent = "Please select a file before uploading.";
        return;
    }


    const reader = new FileReader();

    reader.onload = async () => {

        const formData = new FormData();
        formData.append("file_name", selectedFile.name);
        formData.append("file", selectedFile);
        // const base64File = reader.result.split(",")[1]; // Take only part base64
        // formData.append("base_64", base64File); If you wanna send base64 stead a file

        try {
            const response = await fetch("http://localhost/kyc-request", {
                method: "POST",
                body: formData,
            });

            if (response.ok) {
                result.textContent = "File uploaded successfully!";
            } else {
                result.textContent = "Failed to upload file.";
            }
        } catch (error) {
            result.textContent = "Error while uploading file.";
        }
    };

    reader.readAsDataURL(selectedFile); // Converte para Base64
});