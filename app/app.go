package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	debugMode       = true
	cancelableCtx   context.Context
	canceler        context.CancelFunc
	once            sync.Once
	shutdownHandle  func(context.Context)
	shutdownTimeout time.Duration = time.Second * 60
)

func init() {
	// TODO: debugMode
}

// Context returns a context.Context, that is global unique and cancelable.
func Context() context.Context {
	if cancelableCtx == nil {
		log.Fatal(errors.New("app: not running yet."))
	}
	return cancelableCtx
}

func doLaunch() {
	if cancelableCtx != nil {
		log.Fatal(errors.New("app: already was launched."))
	}

	cancelableCtx, canceler = context.WithCancel(context.Background())
	defer canceler()

	_, hasShow := commandMap["default -v"]
	_, hasShowFull := commandMap["default --version"]
	if !hasShow && !hasShowFull {
		Command("default -v,--version", false, func(v interface{}) {
			if v := v.(bool); !v {
				return
			}
			PrintVersion()
			os.Exit(0)
		}, "Print version information and quit")
		commands.MoveToFront(commands.Back())
	}

	if _, ok := commandMap["default --help"]; !ok {
		Command("default --help", false, func(v interface{}) {
			if v := v.(bool); !v {
				return
			}
			PrintUsage("default")
			os.Exit(0)
		}, func(*CommandDirective) string { return "Display full usage information and quit" })
		commands.MoveToFront(commands.Back())
	}

	c := make(chan os.Signal)
	go func() {
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		var sigcnt int
		for {
			// s := <-c
			// switch s {
			// case syscall.SIGINT:
			// 	sigcnt++
			// case syscall.SIGTERM:
			// 	sigcnt = 3
			// }
			<-c
			sigcnt++
			if sigcnt > 1 {
				exit(1, "app: exit now.")
			}
			fmt.Println("press ctrl+c again to force exit")
			Shutdown()
		}
	}()

	// flag.Parse()
	args := parse(os.Args[1:])
	elem := commands.Front()
	for elem != nil {
		cmd := elem.Value.(*CommandDirective)
		var (
			v   string
			cok bool
		)
		for _, key := range cmd.keys {
			if iv, ok := args[key]; ok {
				if cok {
					log.Fatalf("app: Parameter conflict, only one of %v can be passed in.", cmd.keys)
				}
				cok = true
				v = iv
			}
		}
		if v == "=" {
			if cmd.dataType == "bool" {
				v = "true"
			} else {
				v = ""
			}
		}
		if v, err := cmd.parser(v); err != nil {
			log.Fatalf("app: failed to resolve the Parameter,%v=%s", cmd.keys, v)
		} else {
			cmd.handle(v)
		}

		elem = elem.Next()
	}

	if defaultCommand != nil {
		defaultCommand()
	}
	wg.Wait()
	exit(0, "exit now.")

}

// Launch this applications.
func Launch(appName, appVersion, gitCommit, builds string) {
	envApplication.Set(appName)
	envVersion.Set(appVersion)
	envGoVersion.Set(strings.Trim(runtime.Version(), "go"))
	envGitCommit.Set(gitCommit)
	envBuilds.Set(builds)

	if commands.Len() == 0 && defaultCommand == nil {
		log.Fatal(errors.New("app: no commands for Launch, and must be declared with app.Command."))
	}

	once.Do(doLaunch)
}

func timeoutWithExit() context.Context {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(shutdownTimeout))
	go func() {
		defer cancel()
		<-ctx.Done()
		log.Println("app: forced exit due to shutting down timeout.")
		os.Exit(2)
	}()
	return ctx
}

func exit(status int, msg string) {
	if shutdownHandle != nil {
		shutdownHandle(timeoutWithExit())
	}
	log.Println(msg)
	os.Exit(status)
}

func Shutdown() context.Context {
	defer canceler()
	return timeoutWithExit()
}

func SetShutdownTimeout(d time.Duration) {
	shutdownTimeout = d
}

func HandleShutdown(h func(context.Context)) {
	shutdownHandle = h
}
