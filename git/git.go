package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"
)

const (
	GitTopic = "git"
)

func HandlerArgs(args []string) {
	if !CheckIsBash() {
		fmt.Printf("%s %s\n", "[Error]", color.RedString("Please Use Bash Shell"))
		return
	}
	if len(args) >= 2 {
		current_path := GetCurrentShellWd()
		direct := true
		last_arg := os.Args[len(os.Args)-1]
		if runtime.GOOS == "darwin" && IsDir(last_arg) {
			if last_arg[0:1] == "/" {
				current_path = last_arg
			} else {
				current_path = current_path + "/" + last_arg
			}
			direct = false
		}
		repos := GetAllLocalRepo(current_path)
		var wg sync.WaitGroup
		for _, r := range repos {
			if IsDir(r) && len(r) > 0 {
				wg.Add(1)
				go func(lock sync.WaitGroup, path string) {
					defer wg.Done()
					var cmd string
					for n, arg := range args {
						if !direct && n == len(args)-1 {
							continue
						}
						cmd += arg + " "
					}
					git_path := strings.Replace(path, ".git", "", -1)
					git_path = strings.Replace(git_path, "//", "/", -1)
					git_path = strings.Replace(git_path, `\`, `/`, -1)
					git_cmd := "cd " + git_path + " && " + cmd
					out := ExecCmd(git_cmd)
					fmt.Printf("%s %s\n %s\n", "[PATH]", color.YellowString(git_path), out)
				}(wg, r)
			}
		}
		wg.Wait()
	} else {
		fmt.Println("ops! wrong syntax")
	}
}

func GetCurrentShellWd() string {
	pwd, _ := os.Getwd()
	pwd = strings.Replace(pwd, `\`, `/`, -1)
	return pwd
}

func IsDir(path string) (b bool) {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func GetAllLocalRepo(dir string) []string {
	repos := make([]string, 0)
	if runtime.GOOS == "windows" {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".git") {
				repos = append(repos, strings.ReplaceAll(path, ".git", ""))
			} else {
				return nil
			}
			return nil
		})
		if err != nil {
			fmt.Printf("GetAllLocalRepo:%s\n err:%s\n", dir, err)
		}
	} else {
		cmd := "find " + dir + " -name .git"
		ret := ExecCmd(cmd)
		repos = strings.Split(ret, "\n")
	}
	return repos
}

func ExecCmd(command string) (str string) {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		//fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}
	if err := cmd.Start(); err != nil {
		//fmt.Println("Error:The command is err,", err)
		return err.Error()
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		//fmt.Println("ReadAll Stdout:", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		//fmt.Println("wait:", err.Error())
		return
	}
	return string(bytes)
}

func CheckIsBash() bool {
	s := ExecCmd("which bash")
	if strings.Contains(s, "/bash") {
		return true
	}
	return false
}
