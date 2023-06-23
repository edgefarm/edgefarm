package extract

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func OpenAndExtractTarGz(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("openAndExtractTarGz: Open() failed: %s", err.Error())
	}
	defer file.Close()

	if err := ExtractTarGz(file); err != nil {
		return fmt.Errorf("openAndExtractTarGz: ExtractTarGz() failed: %s", err.Error())
	}
}

func ExtractTarGz(gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("extractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("extractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return fmt.Errorf("extractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return fmt.Errorf("extractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("extractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		default:
			return fmt.Errorf("extractTarGz: uknown type: %s in %s", header.Typeflag, header.Name)
		}

	}
}
