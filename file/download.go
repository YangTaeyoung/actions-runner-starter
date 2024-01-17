package file

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Download(filepath string, url string) error {
	// HTTP GET 요청을 통해 파일을 가져옵니다.
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 지정된 경로에 파일을 생성합니다.
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 응답 본문을 파일에 복사합니다.
	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}

// Unzip 함수는 tar.gz 파일의 경로와 압축 해제할 위치를 받아 해당 위치에 압축을 해제합니다.
// source: 압축 해제할 파일 경로
// destination: 압축 해제할 위치
func Unzip(source string, destination string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipStream, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipStream.Close()

	tarReader := tar.NewReader(gzipStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 파일의 끝에 도달
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destination, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		default:
		}
	}

	return nil
}
