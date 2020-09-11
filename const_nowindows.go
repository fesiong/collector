// +build !windows

package collector

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

// NotarySupportedSignals contains the signals we would like to capture:
// - SIGUSR1, indicates a increment of the log level.
// - SIGUSR2, indicates a decrement of the log level.
var NotarySupportedSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
	syscall.SIGUSR1,
	syscall.SIGUSR2,
	syscall.SIGTSTP,
	syscall.SIGUSR1,
	syscall.SIGUSR2,
}

// LogLevelSignalHandle will increase/decrease the logging level via the signal we get.
func LogLevelSignalHandle(sig os.Signal) {
	switch sig {
	case syscall.SIGUSR1:
		if err := AdjustLogLevel(true); err != nil {
			fmt.Printf("Attempt to increase log level failed, will remain at %s level, error: %s\n", logrus.GetLevel(), err)
			return
		}
	case syscall.SIGUSR2:
		if err := AdjustLogLevel(false); err != nil {
			fmt.Printf("Attempt to decrease log level failed, will remain at %s level, error: %s\n", logrus.GetLevel(), err)
			return
		}
	}

	fmt.Println("Successfully setting log level to", logrus.GetLevel())
}

// AdjustLogLevel increases/decreases the log level, return error if the operation is invalid.
func AdjustLogLevel(increment bool) error {
	lvl := logrus.GetLevel()

	// The log level seems not possible, in the foreseeable future,
	// out of range [Panic, Debug]
	if increment {
		if lvl == logrus.DebugLevel {
			return fmt.Errorf("log level can not be set higher than %s", "Debug")
		}
		lvl++
	} else {
		if lvl == logrus.PanicLevel {
			return fmt.Errorf("log level can not be set lower than %s", "Panic")
		}
		lvl--
	}

	logrus.SetLevel(lvl)
	return nil
}
