package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/YangTaeyoung/actions-runner-starter/file"
	"github.com/YangTaeyoung/actions-runner-starter/resolver"
	"github.com/YangTaeyoung/actions-runner-starter/validator"
	"github.com/samber/lo"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

var DEFAULT = "default"

func main() {
	var (
		githubRepositoryURL string
		githubToken         string
		runnerDownloadURL   string
		numOfWorkers        int
		runnerGroupName     string
		runnerName          string
		runnerLabel         string
		runnerWorkDirectory string
		err                 error
	)

	prompt := survey.Input{
		Message: "Enter Github Repository URL",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &githubRepositoryURL, survey.WithValidator(survey.Required)); err != nil {
		fmt.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Github Token",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &githubToken, survey.WithValidator(survey.Required)); err != nil {
		fmt.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Download URL (if you don't know, enter default)",
		Default: DEFAULT,
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &runnerDownloadURL); err != nil {
		fmt.Println(err.Error())
		return
	}
	if runnerDownloadURL == DEFAULT {
		if runnerDownloadURL = resolver.RunnerDownloadURL(); runnerDownloadURL == "" {
			fmt.Println("Not supported OS or Architecture")
			return
		}
	}

	prompt = survey.Input{
		Message: "Enter Number of Workers",
		Default: "1",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &numOfWorkers, survey.WithValidator(validator.Positive)); err != nil {
		fmt.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Group Name",
		Default: DEFAULT,
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &runnerGroupName); err != nil {
		fmt.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Name",
		Default: DEFAULT,
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &runnerName); err != nil {
		fmt.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Label",
		Default: DEFAULT,
	}

	if err = survey.AskOne(lo.ToPtr(prompt), &runnerLabel); err != nil {
		fmt.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Work Directory",
		Default: DEFAULT,
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &runnerWorkDirectory); err != nil {
		fmt.Println(err.Error())
		return
	}

	runnersPath := path.Join(os.Getenv("HOME"), "actions-runners")
	if err = os.MkdirAll(runnersPath, os.ModePerm); err != nil {
		fmt.Println(err.Error())
		return
	}

	runnerExt := "tar.gz"
	if strings.Contains(runnerDownloadURL, ".zip") {
		runnerExt = "zip"
	}

	// 러너 다운로드
	if err = file.Download(runnersPath, "actions-runner", runnerExt, runnerDownloadURL); err != nil {
		fmt.Println(err.Error())
		return
	}

	for i := 0; i < numOfWorkers; i++ {
		// 개별 러너 폴더 생성
		runnerPath := path.Join(runnersPath, fmt.Sprintf("runner-%d", i))
		if err = os.MkdirAll(runnerPath, os.ModePerm); err != nil {
			fmt.Println(err.Error())
			return
		}

		// 러너 압축 해제

	}

	// 명령어와 인자 설정
	home := os.Getenv("HOME")
	cmd := exec.Command(home+"/actions-runner/config.sh", "--url", githubRepositoryURL, "--token", githubToken)

	// 표준 입력 파이프 생성
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("StdinPipe Error:", err)
		return
	}

	// 명령어 실행
	err = cmd.Start()
	if err != nil {
		fmt.Println("Start Error:", err)
		return
	}

	// 러너 그룹 명
	if _, err = io.WriteString(stdin, "\n"); err != nil {
		fmt.Println("WriteString Error:", err)
		return
	}

	// 러너명
	if _, err = io.WriteString(stdin, "hello-world\n"); err != nil {
		fmt.Println("WriteString Error:", err)
		return
	}

	// 러너 라벨
	if _, err = io.WriteString(stdin, "test-runner\n"); err != nil {
		fmt.Println("WriteString Error:", err)
		return
	}

	// 폴더 이름
	if _, err = io.WriteString(stdin, "\n"); err != nil {
		fmt.Println("WriteString Error:", err)
		return
	}

	stdin.Close()

	// 실행 완료 대기
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Wait Error:", err)
		return
	}

	fmt.Println("Command executed successfully")
}
