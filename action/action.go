package action

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/YangTaeyoung/actions-runner-starter/file"
	"github.com/YangTaeyoung/actions-runner-starter/resolver"
	"github.com/YangTaeyoung/actions-runner-starter/runner"
	"github.com/YangTaeyoung/actions-runner-starter/validator"
	"github.com/samber/lo"
	"github.com/urfave/cli"
	"log"
	"os"
	"path"
)

const DEFAULT = "default"

func Register(ctx cli.Context) {
	var (
		runnerDownloadURL string
		numOfWorkers      int
		err               error
	)
	cfg := runner.Config{}

	prompt := survey.Input{
		Message: "Enter Github Repository URL",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.GithubRepositoryURL, survey.WithValidator(survey.Required)); err != nil {
		log.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Github Token",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.GithubToken, survey.WithValidator(survey.Required)); err != nil {
		log.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Download URL (if you don't know, enter default)",
		Default: resolver.RunnerDownloadURL(),
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &runnerDownloadURL); err != nil {
		log.Println(err.Error())
		return
	}
	if runnerDownloadURL == "" {
		log.Println("unsupported os or architecture")
		return
	}

	prompt = survey.Input{
		Message: "Enter Number of Workers",
		Default: "1",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &numOfWorkers, survey.WithValidator(validator.Positive)); err != nil {
		log.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Group Name",
		Default: DEFAULT,
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.RunnerGroupName); err != nil {
		log.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Name",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.RunnerName, survey.WithValidator(survey.Required)); err != nil {
		log.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Label",
		Default: DEFAULT,
	}

	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.RunnerLabel); err != nil {
		log.Println(err.Error())
		return
	}

	prompt = survey.Input{
		Message: "Enter Runner Work Directory",
		Default: DEFAULT,
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.RunnerWorkDirectory); err != nil {
		log.Println(err.Error())
		return
	}

	// 러너 폴더 생성
	runnersPath := path.Join(os.Getenv("HOME"), "actions-runners")
	if err = os.MkdirAll(runnersPath, os.ModePerm); err != nil {
		log.Println(err.Error())
		return
	}

	// 러너 다운로드
	if err = file.Download(runnersPath+"/actions-runner.tar.gz", runnerDownloadURL); err != nil {
		log.Println(err.Error())
		return
	}

	for i := 0; i < numOfWorkers; i++ {
		// 개별 러너 폴더 생성
		runnerPath := path.Join(runnersPath, fmt.Sprintf("runner-%d", i))
		if err = os.MkdirAll(runnerPath, os.ModePerm); err != nil {
			log.Println(err.Error())
			return
		}

		// 러너 압축 해제
		if err = file.Unzip(runnersPath+"/actions-runner.tar.gz", runnerPath); err != nil {
			log.Println(err.Error())
			return
		}

		// 러너 등록
		r := runner.New(runnerPath, cfg, i)

		if err = r.Register(); err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func Unregister(ctx cli.Context) {
}

func Serve(ctx cli.Context) {

}
