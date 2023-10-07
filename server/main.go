package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("x-access-token")
	router.Use(cors.New(config))

	// Define routes
	router.GET("/testing", test)
	router.POST("/upload",Predict)

	// Run the server on port 8800
	router.Run(":8800")
}

func Predict(c *gin.Context) {
	// Retrieve the uploaded file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// Generate a unique filename for the uploaded file
	extension := filepath.Ext(file.Filename)
	uniqueFilename := "dataset_" + extension
	filePath := filepath.Join("Dataset", uniqueFilename)

	// Save the uploaded file to the "Dataset" folder
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}

	// Validate the CSV file
	// if !isValidCSV(filePath) {
	// 	// Delete the file if it's not valid
	// 	_ = os.Remove(filePath)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV file format"})
	// 	return
	// }

	// Run the Python script to analyze emotions and get the result
	// Run the Python script to analyze emotions and get the result

    cmd := exec.Command("python", "../ml-model/main.py", filePath)  // Updated path
    var out, errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        fmt.Printf("Error: %v\n", err)
        fmt.Printf("Python script error: %s\n", errBuf.String())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze emotions"})
        return
    }


	emotionResult := out.String()

    // Split the emotion result string into individual emotions
    emotions := parseEmotions(emotionResult)

    // Create a JSON response
    response := gin.H{
        "message":   "File uploaded and saved successfully",
        "filename":  uniqueFilename,
        "emotions":  emotions,
    }

    c.JSON(http.StatusOK, response)
}

func parseEmotions(emotionResult string) map[string]float64 {
    emotions := make(map[string]float64)
    lines := strings.Split(emotionResult, "\n")

    // Loop through the lines and extract emotion values
    for _, line := range lines {
        if strings.Contains(line, ":") {
            parts := strings.Split(line, ":")
            emotionName := strings.TrimSpace(parts[0])
            emotionValue, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
            if err == nil {
                emotions[emotionName] = emotionValue
            }
        }
    }

    return emotions
}

// func isValidCSV(filePath string) bool {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return false
// 	}
// 	defer file.Close()

// 	reader := csv.NewReader(file)
// 	header, err := reader.Read()
// 	if err != nil || len(header) != 1 || header[0] != "test" {
// 		return false
// 	}

// 	for {
// 		record, err := reader.Read()
// 		if err == io.EOF {
// 			break
// 		} else if err != nil {
// 			return false
// 		}

// 		if len(record) != 1 || len(record[0]) == 0 {
// 			return false
// 		}
// 	}

// 	return true
// }


func test(c *gin.Context) {
	result := "wooow"
	c.JSON(http.StatusOK, gin.H{"message": result})
}