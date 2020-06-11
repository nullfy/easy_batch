package git

import (
	"fmt"
	"io/ioutil"
	"log"
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
	if len(args) >= 2 {
		last_arg := args[len(args)-1]
		current_path := GetCurrentShellWd()
		direct := true
		if IsDir(last_arg) {
			direct = false
			if runtime.GOOS == "windows" {
				if last_arg[0:1] == "." {
					current_path = current_path + "/" + last_arg
				} else {
					current_path = last_arg
				}
			} else {
				if last_arg[0:1] == "/" {
					current_path = last_arg
				} else {
					current_path = current_path + "/" + last_arg
				}
			}
		} else {
			fmt.Printf("%s\t %s\n",  color.YellowString("[Warning]"), "please input correct path")
			return
		}
		stdouts := GetAllLocalRepo(current_path)
		repos := strings.Split(stdouts, "\n")
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
					git_path := strings.Replace(path, ".git", "", 1)
					git_path = strings.Replace(git_path, "//", "/", 2)
					git_cmd := "cd " + git_path + " && " + cmd
					fmt.Println(cmd, git_path)
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(dir)
	return dir
}

func IsDir(path string) (b bool) {
	fmt.Println("Dir",path)
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func GetAllLocalRepo(path string) string {
	cmd := "find " + path + " -name .git"
	ret := ExecCmd(cmd)
	ret = strings.Replace(ret, "//", "/", 999)
	return ret
}

func ExecCmd(command string) (str string) {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return
	}
	return string(bytes)
}
