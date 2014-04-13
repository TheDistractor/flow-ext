// +build linux darwin

// Interface to serial port devices (only Linux and Mac OSX).
package serial

import (
	"bufio"
	"time"
	"github.com/chimera/rs232"
	"github.com/jcw/flow"
	serialex "github.com/TheDistractor/flow-ext/gadgets/jeebus/serial/extended"

)

//Automatic addition to flow registry
func init() {
	flow.Registry["SerialPortEx"] = func() flow.Circuitry { return new(serialex.SerialPort) }
}
