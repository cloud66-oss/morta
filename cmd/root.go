package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	argTypePid     = "pid"
	argTypeSignal  = "signal"
	argTypeSeconds = "seconds"

	defaultPollingPeriod = 1
)

var (
	// used to store flags
	processPID       int
	shutdownSequence string
	pollingPeriod    int

	rootCmd = &cobra.Command{
		Use:   "shutdown-sequencer",
		Short: "Perform a shutdown sequence against a process",
		Long: `Given a process PID and a sequence of alternating signals and sleep durations,
shutdown-sequencer will perform the sequence against the PID until either the process is dead,
or the sequence has completed.`,
		RunE: run,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func init() {
	rootCmd.Flags().IntVarP(&processPID, "pid", "p", 0, "PID of the process to shut down (required)")
	rootCmd.Flags().StringVarP(&shutdownSequence, "sequence", "s", "", "shutdown sequence to perform against the process (required)")
	rootCmd.Flags().IntVarP(&pollingPeriod, "polling-period", "z", 1, "number of seconds to sleep before checking if the process is still alive")

	rootCmd.MarkFlagRequired("pid")
	rootCmd.MarkFlagRequired("sequence")
}

func run(cmd *cobra.Command, args []string) error {
	log.WithFields(log.Fields{
		"pid":            processPID,
		"sequence":       shutdownSequence,
		"polling-period": pollingPeriod,
	}).Info("Arguments received")

	parsedShutdownSequence, err := parseShutdownSequence()
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"pid":            processPID,
		"sequence":       parsedShutdownSequence,
		"polling-period": pollingPeriod,
	}).Info("Parsed arguments")

	currentArgType := argTypeSignal
	for _, arg := range parsedShutdownSequence {
		processHandled := false
		switch currentArgType {
		case argTypeSignal:
			processHandled, err = handleSignalArg(arg, processPID)
			if err != nil {
				return err
			}

			if processHandled {
				break
			}

			currentArgType = argTypeSeconds
		case argTypeSeconds:
			processHandled, err = handleSecondsArg(arg, processPID)
			if err != nil {
				return err
			}

			if processHandled {
				break
			}

			currentArgType = argTypeSignal
		}
		if processHandled {
			break
		}
	}

	return nil
}

func handleSignalArg(arg, pid int) (bool, error) {
	signal := syscall.Signal(arg)

	log.WithFields(log.Fields{
		"pid":    pid,
		"signal": arg,
	}).Info("Attempting to send signal to process")
	err := syscall.Kill(pid, signal)
	if err == nil {
		log.WithFields(log.Fields{
			"pid":    pid,
			"signal": arg,
		}).Info("Successfully sent sent signal to process")
		return false, nil
	} else {
		if err.Error() == "no such process" {
			log.WithFields(log.Fields{
				"pid":    pid,
				"signal": arg,
			}).Info("Signal sent to process, but process no longer exists")
			return true, nil
		} else {
			log.WithFields(log.Fields{
				"pid":    pid,
				"signal": arg,
			}).Error("Encountered an error while sending signal to process")
			return false, err
		}
	}
}

func handleSecondsArg(arg, pid int) (bool, error) {
	log.WithFields(log.Fields{
		"pid":     pid,
		"seconds": arg,
	}).Info("Waiting for process to exit")

	timeStart := time.Now()
	timeEnd := timeStart.Add(time.Duration(arg) * time.Second)
	for {
		processExists, err := process.PidExists(int32(arg))
		if err != nil {
			log.WithFields(log.Fields{
				"pid":     pid,
				"seconds": arg,
			}).Error("Encountered error when trying to determine if process still exists")
			return false, err
		}

		if !processExists {
			log.WithFields(log.Fields{
				"pid":     pid,
				"seconds": arg,
			}).Info("Process has exited and no longer exists")
			return true, nil
		}

		timeNow := time.Now()
		if timeNow.After(timeEnd) {
			log.WithFields(log.Fields{
				"pid":     pid,
				"seconds": arg,
			}).Info("Process still exists after wait duration of shutdown sequence")
			return false, nil
		}

		log.WithFields(log.Fields{
			"pid":     pid,
			"seconds": arg,
		}).Info("Process is still alive, continuing polling...")
		time.Sleep(time.Duration(pollingPeriod) * time.Second)
	}
}

func parseShutdownSequence() ([]int, error) {
	splitShutdownSequence := strings.Split(shutdownSequence, ":")
	if len(splitShutdownSequence)%2 == 0 {
		return nil, fmt.Errorf("shutdown-sequencer: shutdown sequence must have an odd number of elements")
	}

	result := make([]int, len(splitShutdownSequence))

	currentArgType := argTypeSignal
	for i, arg := range splitShutdownSequence {
		switch currentArgType {
		case argTypeSignal:
			res, err := parseSignalArgument(arg)
			if err != nil {
				return nil, err
			}
			result[i] = res
			currentArgType = argTypeSeconds
		case argTypeSeconds:
			res, err := parseSecondsArgument(arg)
			if err != nil {
				return nil, err
			}
			result[i] = res
			currentArgType = argTypeSignal
		}
	}

	return result, nil
}

// NOTE: list of signals can be found here: https://golang.org/pkg/syscall/#Signal
func parseSignalArgument(arg string) (int, error) {
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
		return 0, fmt.Errorf("shutdown-sequencer: invalid signal %s", arg)
	}

	return int(signal), nil
}

func parseSecondsArgument(arg string) (int, error) {
	seconds, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		return 0, err
	}

	if seconds < 0 {
		return 0, fmt.Errorf("shutdown-sequencer: number of seconds in shutdown sequence (%s) must be non-negative", arg)
	}

	return int(seconds), nil
}
