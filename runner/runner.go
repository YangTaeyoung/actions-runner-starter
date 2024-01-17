package runner

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

var defaultReplacer = strings.NewReplacer("default", "")

type Config struct {
	GithubRepositoryURL string
	GithubToken         string
	RunnerGroupName     string
	RunnerName          string
	RunnerLabel         string
	RunnerWorkDirectory string
}

type Runner struct {
	runnerPath string
	config     Config
	index      int
}

func New(runnerPath string, config Config, index int) Runner {
	return Runner{
		runnerPath: runnerPath,
		config:     config,
		index:      index,
	}
}

func (r Runner) Register() error {
	// 권한 설정
	cmd := exec.Command("chmod", "777", r.runnerPath+"/config.sh")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("chmod", "777", r.runnerPath+"/bin/Runner.Listener")
	if err := cmd.Run(); err != nil {
		return err
	}

	// 명령어와 인자 설정
	cmd = exec.Command(r.runnerPath+"/config.sh", "--url", r.config.GithubRepositoryURL, "--token", r.config.GithubToken)

	// 표준 입력 파이프 생성
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	// 명령어 실행
	if err = cmd.Start(); err != nil {
		return err
	}

	// 러너 그룹 명
	if _, err = io.WriteString(stdin, defaultReplacer.Replace(r.config.RunnerGroupName)+"\n"); err != nil {
		return err
	}

	// 러너명
	if _, err = io.WriteString(stdin, r.config.RunnerName+"-"+strconv.Itoa(r.index)+"\n"); err != nil {
		fmt.Println("WriteString Error:", err)
		return err
	}

	// 러너 라벨
	if _, err = io.WriteString(stdin, defaultReplacer.Replace(r.config.RunnerLabel)+"\n"); err != nil {
		fmt.Println("WriteString Error:", err)
		return err
	}

	// 폴더 이름
	if _, err = io.WriteString(stdin, defaultReplacer.Replace(r.config.RunnerWorkDirectory)+"\n"); err != nil {
		return err
	}

	// 실행 완료 대기
	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
