// +build windows

package main

import "os"

// NotarySupportedSignals does not contain any signals, because SIGUSR1/2 are not supported on windows
var NotarySupportedSignals = []os.Signal{}

// LogLevelSignalHandle will do nothing, because we aren't currently supporting signal handling in windows
func LogLevelSignalHandle(sig os.Signal) {
}