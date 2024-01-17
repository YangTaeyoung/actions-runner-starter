package file

import (
	"archive/tar"
	"compress/gzip"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func ImportConfig(filepath string, cfg any) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// 파일을 읽어서 cfg에 저장합니다.
	if err = yaml.Unmarshal(fileBytes, cfg); err != nil {
		return err
	}

	return nil
}

func ExportConfig(configPath string, cfg any) error {
	if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var (
		buf []byte
	)

	// cfg를 yaml 형식으로 변환합니다.
	if buf, err = yaml.Marshal(cfg); err != nil {
		return err
	}

	// buf를 파일에 저장합니다.
	if _, err = file.Write(buf); err != nil {
		return err
	}

	return nil
}

// Unzip 함수는 tar.gz 파일의 경로와 압축 해제할 위치를 받아 해당 위치에 압축을 해제합니다.
// source: 압축 해제할 파일 경로
// destination: 압축 해제할 위치
func Unzip(source string, destination string) error {
	var (
		header  *tar.Header
		outFile *os.File
		err     error
	)

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

		header, err = tarReader.Next()
		if err == io.EOF {
			break // 파일의 끝에 도달
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destination, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err = os.Create(target)
			if err != nil {
				return err
			}
			if _, err = io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}

			if err = outFile.Close(); err != nil {
				return err
			}
		default:
		}
	}

	return nil
}

func RunnerDirs(runnersPath string) ([]string, error) {
	var dirs []string

	entries, err := os.ReadDir(runnersPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), "runner") {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil { // 파일이 존재하는 경우
		return true
	}

	if os.IsNotExist(err) { // 파일이 존재하지 않는 경우
		return false
	}

	// 파일이 존재하지 않거나 다른 오류 발생
	return false
}
