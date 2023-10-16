// Import necessary modules and styles
"use client";
import React, { useState } from "react";
import Papa, { ParseResult } from "papaparse"; // Import papaparse
import Image from "next/image";
import axios from "axios";
import styles from "./page.module.css";

export default function Home() {
  const [fileMessage, setFileMessage] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      if (file.name.endsWith(".csv")) {
        Papa.parse(file, {
          complete: (result: ParseResult) => {
            if (result.data.length > 1) {
              const firstRow = result.data[0] as string[]; // Assert the type of the first row
              if (firstRow.length === 1 && typeof firstRow[0] === "string") {
                setFileMessage(`File "${file.name}" is a valid CSV file.`);
                setSelectedFile(file);
                // You can process the CSV file here or set it in state for further processing.
              } else {
                setFileMessage("CSV file should have one column with text.");
              }
            } else {
              setFileMessage("CSV file should have at least one row.");
            }
          },
          header: false, // Don't treat the first row as headers
        });
      } else {
        setFileMessage("Please select a valid CSV file.");
      }
    }
  };

  const handleUpload = () => {
    if (selectedFile) {
      const formData = new FormData();
      formData.append("file", selectedFile);

      axios
        .post("http://localhost:8800/upload", formData)
        .then((response) => {
          console.log("Upload success:", response);
          // You can handle the response here.
        })
        .catch((error) => {
          console.error("Upload error:", error);
          // Handle the error here.
        });
    } else {
      setFileMessage("Please select a valid CSV file before uploading.");
    }
  };

  return (
    <main className={styles.main}>
      <div className={styles.title}>
        <h1>Sentiment Analysis Platform for Client Feedback</h1>
      </div>
      <div className={styles.description}>
        <p>
          The <strong>Sentiment Analysis Platform for Client Feedback</strong>{" "}
          is a machine learning project designed to empower businesses with the
          ability to extract valuable insights from their clients' feedback
          data.
        </p>
      </div>
      <div>
        <label htmlFor="file" className={styles.label}>
          CSV <ion-icon name="cloud-upload-outline"></ion-icon>
          <input
            type="file"
            id="file"
            accept=".csv"
            onChange={handleFileUpload}
            className={styles.inputFile}
          />
        </label>
        <p>{fileMessage}</p>
      </div>
      <button className={styles.uploadButton} onClick={handleUpload}>
        Upload
      </button>
    </main>
  );
}
