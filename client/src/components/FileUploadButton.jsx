// FileUploadButton.tsx
"use client";
import { useState } from "react";

const FileUploadButton = ({ onFileSelect }) => {
  const [selectedFile, setSelectedFile] = useState(null);

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    setSelectedFile(file);
    onFileSelect(file);
  };

  return (
    <label className="uploadButton">
      Upload .CSV File
      <input type="file" accept=".csv" onChange={handleFileChange} />
    </label>
  );
};

export default FileUploadButton;
