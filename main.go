package main

import (
	// "fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
)

const argTypePid = "pid"
const argTypeSignal = "signal"
const argTypeSeconds = "seconds"

const exitCodeSuccess = 0
const exitCodeInvalidUsage = 10
const exitCodeInvalidArgument = 11

func main() {
	parsedArguments := parseArguments(os.Args[1:])

	pid := parsedArguments[0]
	currentArgType := argTypeSignal
	for _, arg := range parsedArguments[1:] {
		switch currentArgType {
		case argTypeSignal:
			processHandled, err := handleSignalArg(arg, pid)
			if err != nil {
				// TODO: error
			}

			if processHandled {
				break
			}

			currentArgType = argTypeSeconds
		case argTypeSeconds:
			processHandled, err := handleSecondsArg(arg, pid)
			if err != nil {
				// TODO: error
			}

			if processHandled {
				break
			}

			currentArgType = argTypeSignal
		}
	}

	os.Exit(exitCodeSuccess)
}

func handleSignalArg(arg, pid int) (bool, error) {
	signal := syscall.Signal(arg)
	err := syscall.Kill(pid, signal)
	if err == nil {
		return false, nil
	} else {
		if err.Error() == "no such process" {
			return true, nil
		} else {
			return false, err
		}
	}
}

func handleSecondsArg(arg, pid int) (bool, error) {
	pollingPeriod := 1 * time.Second // TODO: make variable
	timeStart := time.Now()
	timeEnd := timeStart.Add(time.Duration(arg))
	for {
		processExists, err := process.PidExists(int32(arg))
		if err != nil {
			return false, err
		}

		if !processExists {
			return true, nil
		}

		timeNow := time.Now()
		if timeNow.After(timeEnd) {
			return false, nil
		}

		time.Sleep(pollingPeriod)
	}
}

func printUsage() {
	// TODO: print usage to stderr
	// fmt.Fprintf(os.Stderr, "number of foo: %d", nFoo)
}

func parseArguments(args []string) []int {
	if len(args) < 2 {
		printUsage()
		os.Exit(exitCodeInvalidUsage)
	}

	result := make([]int, len(args))
	currentArgType := argTypePid
	for i, arg := range args {
		switch currentArgType {
		case argTypePid:
			result[i] = parsePidArgument(arg)
			currentArgType = argTypeSignal
		case argTypeSignal:
			result[i] = parseSignalArgument(arg)
			currentArgType = argTypeSeconds
		case argTypeSeconds:
			result[i] = parseSecondsArgument(arg)
			currentArgType = argTypeSignal
		}
	}

	if currentArgType != argTypeSignal {
		printUsage()
		os.Exit(exitCodeInvalidUsage)
	}

	return result
}

func parsePidArgument(arg string) int {
	pid, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		// TODO: error
		os.Exit(exitCodeInvalidArgument)
	}

	if pid < 0 {
		// TODO: error
		os.Exit(exitCodeInvalidArgument)
	}

	return int(pid)
}

// NOTE: list of signals can be found here: https://golang.org/pkg/syscall/#Signal
func parseSignalArgument(arg string) int {
	var signal syscall.Signal

	switch arg {
	case "abrt":
		signal = syscall.SIGABRT
	case "alrm":
		signal = syscall.SIGALRM
	case "bus":
		signal = syscall.SIGBUS
	case "chld":
		signal = syscall.SIGCHLD
	// case "cld":
	// 	signal = syscall.SIGCLD
	case "cont":
		signal = syscall.SIGCONT
	case "fpe":
		signal = syscall.SIGFPE
	case "hup":
		signal = syscall.SIGHUP
	case "ill":
		signal = syscall.SIGILL
	case "int":
		signal = syscall.SIGINT
	case "io":
		signal = syscall.SIGIO
	case "iot":
		signal = syscall.SIGIOT
	case "kill":
		signal = syscall.SIGKILL
	case "pipe":
		signal = syscall.SIGPIPE
	// case "poll":
	// 	signal = syscall.SIGPOLL
	case "prof":
		signal = syscall.SIGPROF
	// case "pwr":
	// 	signal = syscall.SIGPWR
	case "quit":
		signal = syscall.SIGQUIT
	case "segv":
		signal = syscall.SIGSEGV
	// case "stkflt":
	// 	signal = syscall.SIGSTKFLT
	case "stop":
		signal = syscall.SIGSTOP
	case "sys":
		signal = syscall.SIGSYS
	case "term":
		signal = syscall.SIGTERM
	case "trap":
		signal = syscall.SIGTRAP
	case "tstp":
		signal = syscall.SIGTSTP
	case "ttin":
		signal = syscall.SIGTTIN
	case "ttou":
		signal = syscall.SIGTTOU
	// case "unused":
	// 	signal = syscall.SIGUNUSED
	case "urg":
		signal = syscall.SIGURG
	case "usr1":
		signal = syscall.SIGUSR1
	case "usr2":
		signal = syscall.SIGUSR2
	case "vtalrm":
		signal = syscall.SIGVTALRM
	case "winch":
		signal = syscall.SIGWINCH
	case "xcpu":
		signal = syscall.SIGXCPU
	case "xfsz":
		signal = syscall.SIGXFSZ
	default:
		// TODO: error
		os.Exit(exitCodeInvalidArgument)
	}

	return int(signal)
}

func parseSecondsArgument(arg string) int {
	seconds, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		// TODO: error
		os.Exit(exitCodeInvalidArgument)
	}

	if seconds < 0 {
		// TODO: error
		os.Exit(exitCodeInvalidArgument)
	}

	return int(seconds)
}
