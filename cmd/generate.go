package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	download("https://github.com/igoogolx/lux-geo-data/releases/download/v0.0.4/geoData.tar.gz", filepath.Join("components", "geo", "geoData.tar.gz"))
	download("https://github.com/igoogolx/lux-client/releases/download/v0.2.8/dist-ui.tar.gz", filepath.Join("hub", "routes", "dist.tar.gz"))
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
	defer file.Close()

	// Download the archive
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

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
	defer file.Close()

	// Create a gzip reader to read the compressed data
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

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

		outputPath := filepath.Join(parentDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(outputPath, 0755); err != nil {
				fmt.Errorf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(outputPath)
			if err != nil {
				fmt.Errorf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				fmt.Errorf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		default:
			fmt.Errorf(
				"ExtractTarGz: uknown type: %s in %s",
				header.Typeflag,
				header.Name)
		}
	}

	return nil
}
