package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	download("https://github.com/igoogolx/lux-rules/releases/download/v2.4.1/rules.tar.gz", filepath.Join("internal", "cfg", "distribution", "rule_engine", "rules.tar.gz"))
	download("https://github.com/igoogolx/lux-client/releases/download/v1.13.0-beat.0/dist-ui.tar.gz", filepath.Join("api", "routes", "dist.tar.gz"))
}

func download(url string, outputPath string) {
	err := downloadFile(url, outputPath)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}

	err = unarchiveFile(outputPath)
	if err != nil {
		fmt.Println("Error unarchiving file:", err)
		return
	}

	fmt.Println("File downloaded and unarchived successfully")
}

func downloadFile(url string, filePath string) error {
	// Create the file to which the archive will be downloaded
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("fail to close file: %v\n", filePath)
		}
	}(file)

	// Download the archive
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("fail to close body: %v\n", filePath)
		}
	}(response.Body)

	// Write the archive to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func unarchiveFile(archiveFilePath string) error {
	// Open the archive file
	file, err := os.Open(archiveFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("fail to close file: %v\n", archiveFilePath)
		}
		err = os.Remove(archiveFilePath)
		if err != nil {
			fmt.Printf("fail to remove file: %v\n", archiveFilePath)
		}
	}(file)

	// Create a gzip reader to read the compressed data
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			fmt.Printf("fail to close gzip reader: %v\n", archiveFilePath)
		}
	}(gzipReader)

	// Create a tar reader to read the contents of the archive
	tarReader := tar.NewReader(gzipReader)

	parentDir := filepath.Dir(archiveFilePath)

	// Extract each file in the archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of archive
			break
		}
		if err != nil {
			return err
		}
		if strings.Contains(header.Name, "..") {
			return fmt.Errorf("invalid file name: %s", header.Name)
		}
		outputPath := filepath.Join(parentDir, header.Name)
		err = os.RemoveAll(outputPath)
		if err != nil {
			fmt.Printf("fail to remove path: %s", outputPath)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(outputPath, 0755); err != nil {
				fmt.Printf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(outputPath)
			if err != nil {
				fmt.Printf("ExtractTarGz: Create() failed: %s", err.Error())
				return err

			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				fmt.Printf("ExtractTarGz: Copy() failed: %s", err.Error())
				return err
			}
			err = outFile.Close()
			if err != nil {
				fmt.Printf("fail to close file: %v\n", outputPath)
			}

		default:
			fmt.Printf(
				"ExtractTarGz: uknown type: %v in %s",
				header.Typeflag,
				header.Name)
		}
	}

	return nil
}
