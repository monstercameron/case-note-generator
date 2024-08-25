/**
 * Initializes the application when the DOM is fully loaded.
 */
document.addEventListener("DOMContentLoaded", initializeApp);

/**
 * Main function to initialize the application.
 */
function initializeApp() {
  console.log("DOM loaded");
  const elements = getElements();
  
  if (!validateElements(elements)) return;
  
  console.log("All elements found");
  initializeDatePicker(elements.datePicker);
  setupFileUpload(elements);
  addEventListeners(elements);
  fetchSystemPrompts(elements.promptSelect);
}

/**
 * Retrieves all necessary DOM elements.
 * @returns {Object} An object containing references to DOM elements.
 */
const getElements = () => ({
  dropArea: document.querySelector(".border-dashed"),
  output: document.getElementById("output"),
  fileInput: document.querySelector('input[type="file"]'),
  datePicker: document.getElementById("datePicker"),
  generateJiraCommentButton: document.getElementById("generateJiraComment"),
  copyButton: document.getElementById("copyButton"),
  generateJiraSummaryButton: document.getElementById("generateJiraSummary"),
  promptSelect: document.getElementById("promptSelect"),
  promptEditor: document.getElementById("promptEditor"),
  updatePrompt: document.getElementById("updatePrompt"),
  refreshLogButton: document.getElementById("refreshLogButton"),
  logContent: document.getElementById("logContent"),
  loadPromptButton: document.getElementById("loadPromptButton"),
});

/**
 * Validates that all required elements are present in the DOM.
 * @param {Object} elements - The object containing DOM element references.
 * @returns {boolean} True if all elements are present, false otherwise.
 */
const validateElements = (elements) => {
  const missingElements = Object.entries(elements)
    .filter(([, element]) => !element)
    .map(([name]) => name);

  if (missingElements.length > 0) {
    console.error("Missing elements:", missingElements.join(", "));
    return false;
  }
  return true;
};

/**
 * Initializes the date picker with the current date.
 * @param {HTMLInputElement} datePicker - The date picker input element.
 */
const initializeDatePicker = (datePicker) => {
  datePicker.value = new Date().toISOString().split("T")[0];
};

/**
 * Sets up drag and drop functionality for the specified area.
 * @param {HTMLElement} dropArea - The element to set up as a drop area.
 */
const setupFileUpload = (elements) => {
  const dropZone = document.getElementById('dropZone');

  dropZone.addEventListener('dragenter', (e) => {
    e.preventDefault();
    dropZone.classList.add('file-hover-animation');
  });

  dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
  });

  dropZone.addEventListener('dragleave', () => {
    dropZone.classList.remove('file-hover-animation');
  });

  dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropZone.classList.remove('file-hover-animation');
    
    // Trigger the drop animation
    dropZone.classList.add('file-drop-animation');
    setTimeout(() => {
      dropZone.classList.remove('file-drop-animation');
    }, 500);

    handleFiles(e.dataTransfer.files, elements);
  });

  elements.fileInput.addEventListener('change', (e) => {
    handleFiles(e.target.files, elements);
  });
};

/**
 * Adds event listeners to various elements.
 * @param {Object} elements - The object containing DOM element references.
 */
const addEventListeners = (elements) => {
  elements.dropArea.addEventListener("drop", (e) => handleDrop(e, elements), false);
  elements.fileInput.addEventListener("change", (e) => handleFileSelect(e, elements), false);
  elements.generateJiraCommentButton.addEventListener("click", () => generateJiraComment(elements));
  elements.copyButton.addEventListener("click", () => copyToClipboard(elements));
  elements.generateJiraSummaryButton.addEventListener("click", () => generateJiraSummary(elements));
  elements.promptSelect.addEventListener("change", () => loadSelectedPrompt(elements));
  elements.updatePrompt.addEventListener("click", () => confirmAndUpdatePrompt(elements));
  elements.refreshLogButton.addEventListener("click", () => refreshLog(elements));
  elements.loadPromptButton.addEventListener("click", () => loadSelectedPrompt(elements));
};

/**
 * Confirms with the user before updating the prompt.
 * @param {Object} elements - The object containing DOM element references.
 */
const confirmAndUpdatePrompt = (elements) => {
  const isConfirmed = confirm("Are you sure you want to update the system prompt? This action cannot be undone.");
  if (isConfirmed) {
    updatePrompt(elements);
  }
};

/**
 * Updates the selected prompt with the content from the prompt editor.
 * @param {Object} elements - The object containing DOM element references.
 */
const updatePrompt = async (elements) => {
  const promptId = elements.promptSelect.value;
  const promptContent = elements.promptEditor.value;
  const promptFile = elements.promptSelect.options[elements.promptSelect.selectedIndex].text;
  
  if (!promptId) {
    alert("Please select a prompt to update.");
    return;
  }

  try {
    setLoadingState(true, elements.updatePrompt);
    const response = await fetch(`/systemprompt`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ 
        prompt: promptContent,
        filename: promptFile
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    alert("Prompt updated successfully");
  } catch (error) {
    console.error("Error updating prompt:", error);
    alert("Error updating prompt. Please try again.");
  } finally {
    setLoadingState(false, elements.updatePrompt);
  }
};

/**
 * Handles the file drop event.
 * @param {DragEvent} e - The drag event object.
 * @param {Object} elements - The object containing DOM element references.
 */
const handleDrop = (e, elements) => {
  console.log("File dropped");
  const files = e.dataTransfer.files;
  handleFiles(files, elements);
};

/**
 * Handles the file selection event.
 * @param {Event} e - The change event object.
 * @param {Object} elements - The object containing DOM element references.
 */
const handleFileSelect = (e, elements) => {
  console.log("File selected");
  const files = e.target.files;
  handleFiles(files, elements);
};

/**
 * Processes the selected files.
 * @param {FileList} files - The list of selected files.
 * @param {Object} elements - The object containing DOM element references.
 */
const handleFiles = (files, elements) => {
  if (files.length === 0) {
    console.log("No files selected");
    return;
  }

  const fileList = Array.from(files)
    .filter((file) => file.type === "text/plain" || file.name.endsWith(".md"))
    .map((file) => `${file.name} (${file.type}, ${file.size} bytes)`)
    .join(", ");

  if (fileList) {
    elements.dropArea.innerHTML = `Selected files: ${fileList}`;
    console.log("Valid files:", fileList);
  } else {
    elements.dropArea.innerHTML =
      "Please select valid text (.txt) or Markdown (.md) files.";
    console.log("No valid files selected");
  }

  elements.fileInput.files = files;
};

/**
 * Generates a Jira comment based on the selected files and date.
 * @param {Object} elements - The object containing DOM element references.
 */
const generateJiraComment = async (elements) => {
  console.log("Generating Jira comment");
  setLoadingState(true, elements.generateJiraCommentButton);
  try {
    const files = elements.fileInput.files;
    if (files.length === 0) throw new Error("No files selected");

    const fileContents = await Promise.all(Array.from(files)
      .filter(file => file.type === "text/plain" || file.name.endsWith(".md"))
      .map(readFile));

    console.log("Files read:", fileContents.length);

    const formattedDate = formatDate(elements.datePicker.value);
    console.log("Selected date:", elements.datePicker.value);
    console.log("Formatted date:", formattedDate);

    const response = await fetch("/generate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ prompt: fileContents.join("\n\n"), date: formattedDate }),
    });

    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

    const { completion } = await response.json();
    elements.output.textContent = completion;
    console.log("Jira comment generated");
  } catch (error) {
    console.error("Error:", error.message);
    elements.output.textContent = `Error: ${error.message}`;
  } finally {
    setLoadingState(false, elements.generateJiraCommentButton);
  }
};

/**
 * Generates a Jira summary based on the selected files.
 * @param {Object} elements - The object containing DOM element references.
 */
const generateJiraSummary = async (elements) => {
  console.log("Generating Jira summary");
  setLoadingState(true, elements.generateJiraSummaryButton);
  try {
    const files = elements.fileInput.files;
    if (files.length === 0) throw new Error("No files selected");

    const fileContents = await Promise.all(Array.from(files)
      .filter(file => file.type === "text/plain" || file.name.endsWith(".md"))
      .map(readFile));

    console.log("Files read:", fileContents.length);

    const response = await fetch("/summary", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ prompt: fileContents.join("\n\n") }),
    });

    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

    const result = await response.json();
    elements.output.textContent = result.summary;
    console.log("Jira summary generated");
  } catch (error) {
    console.error("Error:", error.message);
    elements.output.textContent = `Error: ${error.message}`;
  } finally {
    setLoadingState(false, elements.generateJiraSummaryButton);
  }
};

/**
 * Copies the output text to the clipboard.
 * @param {Object} elements - The object containing DOM element references.
 */
const copyToClipboard = (elements) => {
  console.log("Copying to clipboard");
  navigator.clipboard
    .writeText(elements.output.textContent)
    .then(() => {
      console.log("Text copied successfully");
      alert("Text copied to clipboard");
    })
    .catch((err) => console.error("Failed to copy text:", err));
};

/**
 * Reads the contents of a file.
 * @param {File} file - The file to read.
 * @returns {Promise<string>} A promise that resolves with the file contents.
 */
const readFile = (file) => new Promise((resolve, reject) => {
  const reader = new FileReader();
  reader.onload = (event) => {
    console.log(`File read: ${file.name}`);
    resolve(event.target.result);
  };
  reader.onerror = (error) => reject(error);
  reader.readAsText(file);
});

/**
 * Sets the loading state of a button.
 * @param {boolean} isLoading - Whether the button should be in a loading state.
 * @param {HTMLButtonElement} button - The button element to update.
 */
const setLoadingState = (isLoading, button) => {
  if (isLoading) {
    button.classList.add("loading");
    button.disabled = true;
  } else {
    button.classList.remove("loading");
    button.disabled = false;
  }
};

/**
 * Formats a date string into MM/DD/YYYY format.
 * @param {string} dateString - The date string to format.
 * @returns {string} The formatted date string.
 */
const formatDate = (dateString) => {
  const date = new Date(dateString + "T00:00:00");
  return `${(date.getMonth() + 1).toString().padStart(2, "0")}/${date.getDate().toString().padStart(2, "0")}/${date.getFullYear()}`;
};

/**
 * Fetches system prompts and populates the prompt select element.
 * @param {HTMLSelectElement} promptSelect - The select element for prompts.
 */
const fetchSystemPrompts = async (promptSelect) => {
  try {
    const response = await fetch("/systemprompt");
    const prompts = await response.json();
    prompts.forEach((prompt, index) => {
      const option = document.createElement("option");
      option.value = prompt;
      option.textContent = prompt;
      promptSelect.appendChild(option);
    });
  } catch (error) {
    console.error("Error fetching system prompts:", error);
  }
};

/**
 * Loads the selected prompt into the prompt editor.
 * @param {Object} elements - The object containing DOM element references.
 */
const loadSelectedPrompt = async (elements) => {
  const promptId = elements.promptSelect.value;
  if (!promptId) {
    console.error("No prompt selected");
    return;
  }

  try {
    setLoadingState(true, elements.loadPromptButton);
    const response = await fetch(
      `/systemprompt?file=${encodeURIComponent(promptId)}`
    );
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    if (data.file && data.content) {
      elements.promptEditor.value = data.content;
      elements.promptSelect.value = data.file;
      console.log("Prompt loaded successfully");
    } else {
      console.error("Unexpected response format");
    }
  } catch (error) {
    console.error("Error loading selected prompt:", error);
    alert("Error loading prompt. Please try again.");
  } finally {
    setLoadingState(false, elements.loadPromptButton);
  }
};

/**
 * Refreshes the log content and scrolls to the bottom.
 * @param {Object} elements - The object containing DOM element references.
 */
const refreshLog = async (elements) => {
  try {
    const response = await fetch("/logs");
    if (!response.ok) throw new Error("Failed to fetch logs");
    const logs = await response.text();
    elements.logContent.textContent = logs;
    
    // Scroll to the bottom of the log content
    elements.logContent.scrollTop = elements.logContent.scrollHeight;
  } catch (error) {
    console.error("Error fetching logs:", error);
    elements.logContent.textContent = "Error fetching logs. Please try again.";
  }
};

// Collapsible sections functionality
const collapsibles = document.getElementsByClassName("collapsible");

Array.from(collapsibles).forEach(collapsible => {
    collapsible.addEventListener("click", () => {
        console.log("clicked");
        collapsible.classList.toggle("active");
        const content = collapsible.nextElementSibling;
        content.style.display = content.style.display === "block" ? "none" : "block";
    });
});