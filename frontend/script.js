// Función para cargar un archivo de script
function loadFile() {
    const fileInput = document.getElementById('fileInput');
    fileInput.click(); // Simula el clic en el input de archivo

    fileInput.addEventListener('change', (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                const content = e.target.result;
                document.getElementById('commandInput').value = content;
            };
            reader.readAsText(file);
        }
    });
}

// Función para ejecutar comandos
async function executeCommand() {
    const command = document.getElementById('commandInput').value;
    if (!command.trim()) {
        alert('Por favor, ingresa un comando o carga un archivo de script.');
        return;
    }

    const outputElement = document.getElementById('output');
    const loadingElement = document.getElementById('loading');

    // Mostrar animación de carga
    loadingElement.style.display = 'block';
    outputElement.textContent = '';

    try {
        const response = await fetch('http://localhost:8080/execute', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ command }),
        });

        if (!response.ok) {
            throw new Error('Error en la respuesta del servidor');
        }

        const result = await response.json();
        outputElement.textContent = result.output;
    } catch (error) {
        outputElement.textContent = 'Error: ' + error.message;
    } finally {
        // Ocultar animación de carga
        loadingElement.style.display = 'none';
    }
}

// Función para limpiar el contenido
function clearAll() {
    document.getElementById('commandInput').value = '';
    document.getElementById('output').textContent = '';
}