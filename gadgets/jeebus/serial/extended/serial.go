// +build linux darwin

// Interface to serial port devices (only Linux and Mac OSX).
// SerialPortEx is an extended SerialPort with access to port params
// It is NOT auto instanciated to the flow Registry from this package, use serial/serialex to autoadd to Registry
package serial

import (
	"bufio"
	"time"
	"github.com/chimera/rs232"
	"github.com/jcw/flow"
)


// Line-oriented serial port, opened once the Port input is set.
type SerialPort struct {
	flow.Gadget
	Param flow.Input
	Port flow.Input
	To   flow.Input
	From flow.Output
}

// Start processing text lines to and from the serial interface.
// Send a bool to adjust RTS or an int to pulse DTR for that many milliseconds.
// Registers as "SerialPort".
func (w *SerialPort) Run() {

	baud := uint32(57600)
	databits := uint8(8)
	stopbits := uint8(1)

	for param := range w.Param {

		p := param.(flow.Tag)

		switch p.Tag {
		case "baud":
			baud = uint32(p.Msg.(float64))
		case "databits":
			databits = uint8(p.Msg.(float64))
		case "stopbits":
			stopbits = uint8(p.Msg.(float64))
		}

	}




	if port, ok := <-w.Port; ok {
		opt := rs232.Options{BitRate: baud, DataBits: databits, StopBits: stopbits}
		dev, err := rs232.Open(port.(string), opt)
		flow.Check(err)

		// try to avoid kernel panics due to that wretched buggy FTDI driver!
		// defer func() {
		// 	time.Sleep(time.Second)
		// 	dev.Close()
		// }()
		// time.Sleep(time.Second)

		// separate process to copy data out to the serial port
		go func() {
			for m := range w.To {
				switch v := m.(type) {
				case string:
					dev.Write([]byte(v + "\n"))
				case []byte:
					dev.Write(v)
				case int:
					dev.SetDTR(true) // pulse DTR to reset
					time.Sleep(time.Duration(v) * time.Millisecond)
					dev.SetDTR(false)
				case bool:
					dev.SetRTS(v)
				}
			}
		}()

		scanner := bufio.NewScanner(dev)
		for scanner.Scan() {
			msg:=scanner.Text()
			w.From.Send(msg)
		}
	}
}

