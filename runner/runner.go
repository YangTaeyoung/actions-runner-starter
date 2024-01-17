package runner

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/YangTaeyoung/actions-runner-starter/file"
	"github.com/YangTaeyoung/actions-runner-starter/resolver"
	"github.com/YangTaeyoung/actions-runner-starter/validator"
	"github.com/briandowns/spinner"
	"github.com/samber/lo"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const DEFAULT = "default"

var (
	configPath      = path.Join(os.Getenv("HOME"), "actions-runners", ".config.yaml")
	defaultReplacer = strings.NewReplacer(DEFAULT, "")
	runnersPath     = path.Join(os.Getenv("HOME"), "actions-runners")
)

type Config struct {
	GithubRepositoryURL string `yaml:"github_repository_url"`
	GithubToken         string `yaml:"github_token"`
	RunnerGroupName     string `yaml:"runner_group_name"`
	RunnerDownloadURL   string `yaml:"runner_download_url"`
	RunnerName          string `yaml:"runner_name"`
	RunnerLabel         string `yaml:"runner_label"`
	RunnerWorkDirectory string `yaml:"runner_work_directory"`
	NumOfWorkers        int    `yaml:"num_of_workers"`
}

type Runner struct {
	config Config
}

func New() Runner {
	return Runner{}
}

func (r *Runner) Configure() error {
	var err error
	cfg := Config{}

	prompt := survey.Input{
		Message: "Enter Github Repository URL",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.GithubRepositoryURL, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	prompt = survey.Input{
		Message: "Enter Github Token",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.GithubToken, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	prompt = survey.Input{
		Message: "Enter Runner Download URL (if you don't know, enter default)",
		Default: resolver.RunnerDownloadURL(),
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.RunnerDownloadURL); err != nil {
		return err
	}
	if cfg.RunnerDownloadURL == "" {
		return errors.New("unsupported os or architecture")
	}

	prompt = survey.Input{
		Message: "Enter Number of Workers (default: 1)",
		Default: "1",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &r.config.NumOfWorkers, survey.WithValidator(validator.Positive)); err != nil {
		return err
	}

	prompt = survey.Input{
		Message: "Enter Runner Group Name (Enter default)",
		Default: DEFAULT,
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.RunnerGroupName); err != nil {
		return err
	}

	prompt = survey.Input{
		Message: "Enter Runner Name",
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.RunnerName, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	prompt = survey.Input{
		Message: "Enter Runner Label",
		Default: DEFAULT,
	}

	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.RunnerLabel); err != nil {
		return err
	}

	prompt = survey.Input{
		Message: "Enter Runner Work Directory",
		Default: DEFAULT,
	}
	if err = survey.AskOne(lo.ToPtr(prompt), &cfg.RunnerWorkDirectory); err != nil {
		return err
	}

	if err = file.ExportConfig(configPath, &cfg); err != nil {
		return err
	}

	r.config = cfg

	return nil
}

func (r *Runner) Register() error {
	var err error

	// 설정 파일 가져오기
	if err = file.ImportConfig(configPath, &r.config); err != nil {
		return err
	}

	// 러너 폴더 생성
	if err = os.MkdirAll(runnersPath, os.ModePerm); err != nil {
		return err
	}

	// 압축 파일이 없는 경우 다운로드
	if !file.Exists(runnersPath + "/actions-runner.tar.gz") {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = "Downloading runner... "
		s.Start()
		// 러너 다운로드
		if err = file.Download(runnersPath+"/actions-runner.tar.gz", r.config.RunnerDownloadURL); err != nil {
			return err
		}
		s.Stop()

		log.Println("Downloaded runner successfully")
	} else {
		log.Println("Runner already downloaded")
	}

	bar := progressbar.Default(int64(r.config.NumOfWorkers), "Registering runners...")
	for i := 0; i < r.config.NumOfWorkers; i++ {
		var stdin io.WriteCloser

		// 개별 러너 폴더 생성
		runnerPath := path.Join(runnersPath, fmt.Sprintf("runner-%d", i))
		if err = os.MkdirAll(runnerPath, os.ModePerm); err != nil {
			return err
		}

		// 러너 압축 해제
		if err = file.Unzip(runnersPath+"/actions-runner.tar.gz", runnerPath); err != nil {
			return err
		}

		// 권한 설정
		cmd := exec.Command("chmod", "777", runnerPath+"/config.sh")
		if err = cmd.Run(); err != nil {
			return err
		}
		cmd = exec.Command("chmod", "777", runnerPath+"/bin/Runner.Listener")
		if err = cmd.Run(); err != nil {
			return err
		}

		// 명령어와 인자 설정
		cmd = exec.Command(runnerPath+"/config.sh", "--url", r.config.GithubRepositoryURL, "--token", r.config.GithubToken)

		// 표준 입력 파이프 생성
		stdin, err = cmd.StdinPipe()
		if err != nil {
			return err
		}

		// 명령어 실행
		if err = cmd.Start(); err != nil {
			return err
		}

		// 러너 그룹 명
		if _, err = io.WriteString(stdin, defaultReplacer.Replace(r.config.RunnerGroupName)+"\n"); err != nil {
			return err
		}

		// 러너명
		if _, err = io.WriteString(stdin, r.config.RunnerName+"-"+strconv.Itoa(i)+"\n"); err != nil {
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

		if err = stdin.Close(); err != nil {
			return err
		}

		// 실행 완료 대기
		if err = cmd.Wait(); err != nil {
			return err
		}

		bar.Add(1)
	}

	return nil
}

func (r *Runner) Unregister() error {
	var (
		err error
	)

	// 설정 파일 가져오기
	if err = file.ImportConfig(configPath, &r.config); err != nil {
		return err
	}
	if lo.IsEmpty(r.config.GithubToken) {
		return errors.New("github token is empty. please configure again")
	}

	// 폴더 내 액션 러너 디렉토리 리스트 가져오기
	runnersPath := path.Join(os.Getenv("HOME"), "actions-runners")
	runnerDirs, err := file.RunnerDirs(runnersPath)
	if err != nil {
		return err
	}

	// 액션 러너 등록 해제 및 폴더 삭제
	for _, runnerDir := range runnerDirs {
		runnerPath := path.Join(runnersPath, runnerDir)

		// Github 에서 러너 등록 해제
		cmd := exec.Command(runnerPath+"/config.sh", "remove", "--token", r.config.GithubToken)
		if err = cmd.Run(); err != nil {
			return err
		}

		// 러너 폴더 삭제
		if err = os.RemoveAll(runnerPath); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) Serve() error {
	var (
		err error
	)
	// 설정 파일 가져오기
	if err = file.ImportConfig(configPath, &r.config); err != nil {
		return err
	}

	runnerDirs, err := file.RunnerDirs(runnersPath)
	if err != nil {
		return err
	}

	p := progressbar.Default(int64(len(runnerDirs)), "Serving runners...")
	g, _ := errgroup.WithContext(context.Background())
	for _, runnerDir := range runnerDirs {
		runnerPath := path.Join(runnersPath, runnerDir)

		cmd := exec.Command("chmod", "777", runnerPath+"/run-helper.sh")
		if err = cmd.Run(); err != nil {
			return err
		}

		cmd = exec.Command("chmod", "777", runnerPath+"/run.sh")
		if err = cmd.Run(); err != nil {
			return err
		}

		g.Go(func() error {
			cmd = exec.Command(runnerPath + "/run.sh")
			if err = cmd.Run(); err != nil {
				return err
			}

			return nil
		})

		p.Add(1)
	}

	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}
