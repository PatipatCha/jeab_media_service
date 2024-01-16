package services

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/lab/tests/app/model"
)

func DetectContentTypeFromHeader(fileHeader *multipart.FileHeader, menu string) (string, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", "", err
	}

	// Reset the file offset to the beginning
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", "", err
	}

	menuContainers, err := ReadJSONFile()
	if err != nil {
		fmt.Println(err)

	}

	menuToFind := menu
	menuName, err := FindContainerNameByMenu(menuToFind, menuContainers)
	if err != nil {
		fmt.Println(err)

	}

	// Detect the content type for PDF and PNG
	switch {
	case strings.HasPrefix(http.DetectContentType(buffer), "application/pdf"):
		return "application/pdf", menuName.PDFContainerName, nil
	case strings.HasPrefix(http.DetectContentType(buffer), "image/png"):
		return "image/png", menuName.ImageContainerName, nil
	default:
		return "", "", fmt.Errorf("Unsupported file type")
	}
}

func ReadJSONFile() ([]model.MenuContainer, error) {
	filePath := "app/json/stg-menu.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	// Unmarshal the JSON data into a slice of MenuContainer
	var menuContainers []model.MenuContainer
	err = json.Unmarshal(jsonData, &menuContainers)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return menuContainers, nil
}

func FindContainerNameByMenu(menuToFind string, menuContainers []model.MenuContainer) (model.MenuContainer, error) {
	res := model.MenuContainer{}
	for _, container := range menuContainers {
		if container.Menu == menuToFind {
			res.Menu = container.Menu
			res.ImageContainerName = container.ImageContainerName
			res.PDFContainerName = container.PDFContainerName
			return res, nil
		}
	}
	return res, fmt.Errorf("menu not found: %s", menuToFind)
}
