document.addEventListener('DOMContentLoaded', () => {
    // Get references to the drag-and-drop area and output element
    const dropArea = document.querySelector('div.border-dashed');
    const output = document.getElementById('output');
    const fileInfo = document.getElementById('fileInfo'); // New element for file info

    // Prevent default drag behaviors
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropArea.addEventListener(eventName, preventDefaults, false);    
        document.body.addEventListener(eventName, preventDefaults, false);    
    });

    // Highlight the drop area when dragging a file over it
    ['dragenter', 'dragover'].forEach(eventName => {
        dropArea.addEventListener(eventName, highlight, false);
    });

    // Remove highlight when leaving the drop area
    ['dragleave', 'drop'].forEach(eventName => {
        dropArea.addEventListener(eventName, unhighlight, false);
    });

    // Handle dropped files
    dropArea.addEventListener('drop', handleDrop, false);

    // Prevent default behavior (Prevent file from being opened)
    function preventDefaults(e) {
        console.log("Preventing default behavior for event:", e.type);
        e.preventDefault();
        e.stopPropagation();
    }

    // Highlight the drop area
    function highlight() {
        console.log("Highlighting drop area");
        dropArea.classList.add('bg-ukg-light-blue');
    }

    // Remove highlight
    function unhighlight() {
        console.log("Removing highlight from drop area");
        dropArea.classList.remove('bg-ukg-light-blue');
    }

    // Handle the dropped files
    function handleDrop(e) {
        console.log("Handling drop event");
        const dt = e.dataTransfer;
        const files = dt.files;

        console.log("Files dropped:", files);

        // Process the first file
        if (files.length > 0) {
            const file = files[0];
            console.log("Processing file:", file.name, "Type:", file.type, "Size:", file.size);

            // Display file information
            fileInfo.textContent = `Filename: ${file.name}, Type: ${file.type}, Size: ${file.size} bytes`;

            // Check for accepted file types
            if (file.type === "text/plain" || file.name.endsWith('.md')) {
                const reader = new FileReader();
                reader.onload = function(event) {
                    console.log("File read successfully");
                    output.textContent = event.target.result; // Display file content
                };
                reader.onerror = function(event) {
                    console.error("Error reading file:", event);
                };
                reader.readAsText(file);
            } else {
                alert("Please drop a valid text file (.txt) or Markdown file (.md).");
                console.warn("Invalid file type dropped:", file.type);
            }
        } else {
            console.warn("No files were dropped.");
        }
    }
});