package app

import (
	"expvar"
	"fmt"
	"strings"
)

var (
	envApplication = expvar.NewString("Application")
	envVersion     = expvar.NewString("Version")
	envGoVersion   = expvar.NewString("GoVersion")
	envGitCommit   = expvar.NewString("GitCommit")
	envBuilds      = expvar.NewString("Builds")
)

func Env(key string) {

}

func PrintVersion() {
	fmt.Println("Application \t:", strings.Trim(envApplication.String(), `"`))
	fmt.Println("Version \t:", strings.Trim(envVersion.String(), `"`))
	fmt.Println("Go Version \t:", strings.Trim(envGoVersion.String(), `"`))
	fmt.Println("Git Commit \t:", strings.Trim(envGitCommit.String(), `"`))
	fmt.Println("Builds	\t:", strings.Trim(envBuilds.String(), `"`))
}
