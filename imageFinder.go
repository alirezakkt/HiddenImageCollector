package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var destDir string

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run script.go <source-directory> [saveDestenition-directory]")
	}

	srcDir := os.Args[1]
	if len(os.Args) > 2 {
		destDir = os.Args[2]
	} else {
		destDir = "collected_images"
	}


	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	err := filepath.Walk(srcDir, visit)
	if err != nil {
		log.Fatal(err)
	}
}

func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("Error accessing path %s: %v\n", path, err)
		return nil
	}

	if info.IsDir() {
		return nil 
	}

	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", path, err)
		return nil
	}
	defer file.Close()

	//Read first 8 bytes
	buffer := make([]byte, 8)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		log.Printf("Error reading file %s: %v\n", path, err)
		return nil
	}

	var fileType string
	switch {
	case isJPEG(buffer[:n]):
		fileType = "jpg"
	case isPNG(buffer[:n]):
		fileType = "png"
	default:
		return nil 
	}


	originalName := filepath.Base(path)
	ext := filepath.Ext(originalName)
	baseWithoutExt := strings.TrimSuffix(originalName, ext)

	var newExt string
	switch fileType {
	case "jpg":
		newExt = ".jpg"
	case "png":
		newExt = ".png"
	}

	desiredName := baseWithoutExt + newExt
	destPath := getUniqueFileName(destDir, desiredName)

	// CopyFile
	if err := copyFile(path, destPath); err != nil {
		log.Printf("Error copying %s to %s: %v\n", path, destPath, err)
	} else {
		fmt.Printf("Copied %s to %s\n", path, destPath)
	}

	return nil
}

func isJPEG(content []byte) bool {
	return len(content) >= 2 && content[0] == 0xFF && content[1] == 0xD8
}

func isPNG(content []byte) bool {
	pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	return len(content) >= 8 && bytes.Equal(content[:8], pngSignature)
}

func getUniqueFileName(dir, baseName string) string {
	name := baseName
	counter := 1

	for {
		destPath := filepath.Join(dir, name)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			return destPath
		}

		ext := filepath.Ext(baseName)
		base := strings.TrimSuffix(baseName, ext)
		name = fmt.Sprintf("%s_%d%s", base, counter, ext)
		counter++
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
