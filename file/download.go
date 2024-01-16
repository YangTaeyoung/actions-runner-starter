package file

import (
	"io"
	"net/http"
	"os"
)

func Download(filepath string, fileName string, fileExt string, url string) error {
	// HTTP GET 요청을 통해 파일을 가져옵니다.
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 지정된 경로에 파일을 생성합니다.
	out, err := os.Create(filepath + "/" + fileName + "." + fileExt)
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

func Unzip(filepath string, fileName string, fileExt string) error {
	// TODO
	return nil
}
