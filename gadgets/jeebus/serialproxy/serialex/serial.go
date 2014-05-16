// +build linux darwin windows

// Interface to serial port devices (Linux Mac OSX AND Windows).
package serialproxy

import (
	"github.com/jcw/flow"
	serialex "github.com/TheDistractor/flow-ext/gadgets/jeebus/serialproxy/extended"

)

//Automatic addition to flow registry
func init() {
	flow.Registry["SerialPortEx"] = func() flow.Circuitry { return new(serialex.SerialPort) }
}
