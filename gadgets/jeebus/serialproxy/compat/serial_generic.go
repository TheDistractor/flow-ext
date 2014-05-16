// +build linux darwin

// Interface to serial port devices (Linux Mac OSX AND Windows).
// The serialproxy/compat package replaces the standard jeebus SerialPort with the Extended Serial Port package 'SerialPortEx'
// which itself just has more configuration options
// in addition this version does so via a proxy to support windows
package serial

import (
	"github.com/jcw/flow"
	_"github.com/jcw/jeebus/gadgets/serial" //load existing so we can force replace
	serialex "github.com/TheDistractor/flow-ext/gadgets/jeebus/serialproxy/extended"
)

//Automatically override the standard SerialPort from core with Extended version from this package
func init() {
	flow.Registry["SerialPort"] = func() flow.Circuitry { return new(serialex.SerialPort) }
}



