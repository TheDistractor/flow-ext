// +build linux darwin

// Interface to serial port devices (only Linux and Mac OSX).
// The serial/compat package replaces the standard jeebus SerialPort with the Extended Serial Port package 'SerialPortEx'
// which itself just has more configuration options
package serial

import (
	"github.com/jcw/flow"
	_"github.com/jcw/jeebus/gadgets/serial"
	serialex "github.com/TheDistractor/flow-ext/gadgets/jeebus/serial/extended"
)

//Automatically override the standard SerialPort from core with Extended version from this package
func init() {
	flow.Registry["SerialPort"] = func() flow.Circuitry { return new(serialex.SerialPort) }
}



