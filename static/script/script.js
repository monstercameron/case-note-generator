document.addEventListener("DOMContentLoaded", () => {
  console.log("DOM loaded");
  const elements = {
    dropArea: document.querySelector(".border-dashed"),
    output: document.getElementById("output"),
    fileInput: document.querySelector('input[type="file"]'),
    datePicker: document.getElementById("datePicker"),
    generateJiraCommentButton: document.getElementById("generateJiraComment"),
    copyButton: document.getElementById("copyButton"),
  };

  const missingElements = Object.entries(elements)
    .filter(([name, element]) => !element)
    .map(([name]) => name);

  if (missingElements.length > 0) {
    console.error("Missing elements:", missingElements.join(", "));
    return;
  }

  console.log("All elements found");

  elements.datePicker.value = new Date().toISOString().split("T")[0];

  ["dragenter", "dragover", "dragleave", "drop"].forEach((eventName) => {
    elements.dropArea.addEventListener(eventName, preventDefaults, false);
    document.body.addEventListener(eventName, preventDefaults, false);
  });

  ["dragenter", "dragover"].forEach((eventName) => {
    elements.dropArea.addEventListener(
      eventName,
      () => {
        elements.dropArea.classList.add("bg-ukg-light-blue");
        console.log(`${eventName} event triggered`);
      },
      false
    );
  });

  ["dragleave", "drop"].forEach((eventName) => {
    elements.dropArea.addEventListener(
      eventName,
      () => {
        elements.dropArea.classList.remove("bg-ukg-light-blue");
        console.log(`${eventName} event triggered`);
      },
      false
    );
  });

  elements.dropArea.addEventListener("drop", handleDrop, false);
  elements.fileInput.addEventListener("change", handleFileSelect, false);
  elements.generateJiraCommentButton.addEventListener(
    "click",
    generateJiraComment
  );
  elements.copyButton.addEventListener("click", copyToClipboard);

  function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
  }

  function handleDrop(e) {
    console.log("File dropped");
    const files = e.dataTransfer.files;
    handleFiles(files);
  }

  function handleFileSelect(e) {
    console.log("File selected");
    const files = e.target.files;
    handleFiles(files);
  }

  function handleFiles(files) {
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
  }

  async function generateJiraComment() {
    console.log("Generating Jira comment");
    setLoadingState(true);
    try {
      const files = elements.fileInput.files;
      if (files.length === 0) {
        throw new Error("No files selected");
      }
  
      const fileContents = await Promise.all(
        Array.from(files)
          .filter((file) => file.type === "text/plain" || file.name.endsWith(".md"))
          .map(readFile)
      );
  
      console.log("Files read:", fileContents.length);
  
      // Get the selected date and format it as mm/dd/yyyy
      const selectedDate = new Date(elements.datePicker.value + 'T00:00:00');
      const formattedDate = `${(selectedDate.getMonth() + 1).toString().padStart(2, '0')}/${selectedDate.getDate().toString().padStart(2, '0')}/${selectedDate.getFullYear()}`;
  
      console.log("Selected date:", elements.datePicker.value);
      console.log("Formatted date:", formattedDate);
  
      const response = await fetch("/generate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          prompt: fileContents.join("\n\n"),
          date: formattedDate
        }),
      });
  
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      const jiraComment = await response.text();
      const { completion } = JSON.parse(jiraComment);
      elements.output.textContent = completion;
      console.log("Jira comment generated");
    } catch (error) {
      console.error("Error:", error.message);
      elements.output.textContent = `Error: ${error.message}`;
    } finally {
      setLoadingState(false);
    }
  }

  function copyToClipboard() {
    console.log("Copying to clipboard");
    navigator.clipboard
      .writeText(elements.output.textContent)
      .then(() => {
        console.log("Text copied successfully");
        alert("Text copied to clipboard");
      })
      .catch((err) => console.error("Failed to copy text:", err));
  }

  function readFile(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (event) => {
        console.log(`File read: ${file.name}`);
        resolve(event.target.result);
      };
      reader.onerror = (error) => reject(error);
      reader.readAsText(file);
    });
  }

  const generateButton = document.getElementById('generateJiraComment');

  function setLoadingState(isLoading) {
    if (isLoading) {
      generateButton.classList.add('loading');
      generateButton.disabled = true;
    } else {
      generateButton.classList.remove('loading');
      generateButton.disabled = false;
    }
  }

  generateButton.addEventListener('click', generateJiraComment);
});