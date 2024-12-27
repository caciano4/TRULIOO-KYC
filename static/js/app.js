const fileInput = document.getElementById("fileInput");
const fileName = document.getElementById("fileName");
const uploadButton = document.getElementById("uploadButton");
const result = document.getElementById("result");
const label = document.getElementById("file-name");

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
        Swal.fire({
            icon: "error",
            title: "No File Selected",
            text: "Please select a file before uploading.",
        });
        return;
    }

    const reader = new FileReader();

    reader.onload = async () => {
        const formData = new FormData();
        formData.append("file_name", selectedFile.name);
        formData.append("file", selectedFile);

        try {
            const response = await fetch("http://localhost/kyc-request", {
                method: "POST",
                body: formData,
            });

            if (response.ok) {
                let data = await response.json()
                Swal.fire({
                    icon: "success",
                    title: "Upload Successful",
                    text: data.message,
                });
                result.textContent = data.message;
            } else {
                Swal.fire({
                    icon: "error",
                    title: "Upload Failed",
                    text: "Failed to upload file. Please try again.",
                });
                result.textContent = "Failed to upload file.";
            }
        } catch (error) {
            Swal.fire({
                icon: "error",
                title: "Error",
                text: "An error occurred while uploading the file.",
            });
            result.textContent = "Error while uploading file.";
        }
    };

    reader.readAsDataURL(selectedFile);
});
