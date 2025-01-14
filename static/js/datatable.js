
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

const renderTable = async () => {
    // Fetch data
    const data = await fetchPackageList();

    // Clear existing table rows to avoid duplication
    tbody.innerHTML = "";

    // Populate table rows with data
    data.forEach((item) => {
        const row = document.createElement("tr");
        tbody.appendChild(row);

        const fields = [
            { key: "packageIdCell", value: item.package_id },
            { key: "totalRecordsCell", value: item.total_records },
            { key: "packageNameCell", value: item.package_name },
            { key: "completed", value: item.completed },
            { key: "fullNameCell", value: item.full_name },
            { key: "transferAgentCell", value: item.transfer_agent },
            { key: "typeOfTransfer", value: item.type_of_transfer },
            { key: "createdAtCell", value: item.created },
            { key: "updatedAtCell", value: item.updated },
            "action",
        ];

        const buttons = ["KYC", "Download"];

        fields.forEach((field) => {
            const cell = document.createElement("td");
            cell.style.padding = "10px";
            cell.style.textAlign = "center"

            if (field === "action") {
                cell.style.display = "flex";
                buttons.forEach((buttonText) => {
                    const button = document.createElement("button");
                    button.className = "button pure-button is-primary pure-radio";
                    button.style.margin = "5px";
                    button.textContent = buttonText; // Set button text dynamically
                    button.type = "button";
                    cell.appendChild(button);

                    if (buttonText == "KYC") {
                        button.id = "kyc-submit-button"

                        button.addEventListener("click", async (e) => {
                            try {
                                response = await fetch(`http://localhost/process-kyc/${item.package_id}`)
                                console.log("deu bom!")
                            } catch (err) {
                                Swal.fire({
                                    icon: "error",
                                    title: "Error",
                                    text: "An error occurred while uploading the file.",
                                });
                            }
                        })
                    }
                });
            } else {
                cell.textContent = field.value || "-"; // Set field value or placeholder
            }

            row.appendChild(cell);
        });
    });
};


// Render the table when the script runs
renderTable();