// +build linux darwin

// Interface to serial port devices (only Linux and Mac OSX).
package serial

import (
	"bufio"
	"time"
	"fmt"
	"github.com/chimera/rs232"
	"github.com/jcw/flow"
	serialex "github.com/TheDistractor/flow-ext/gadgets/jeebus/serial/extended"

)

func init() {
	fmt.Println("SerialPortEX registered")
	flow.Registry["SerialPortEx"] = func() flow.Circuitry { return new(serialex.SerialPort) }
}
