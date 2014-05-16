// +build linux darwin windows

// Interface to serial port devices (Linux Mac OSX and Windows).
// SerialPortEx is an extended SerialPort with access to port params
// It is NOT auto instantiated to the flow Registry from this package, use serialproxy/serialex to autoadd to Registry
package serialproxy

import (
	"bufio"
	"time"

	"github.com/TheDistractor/goserialproxy" //temp x-platform proxy
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

	initdata := make([]interface{},0) //initialization data sequence (if supplied)

	for param := range w.Param {

		p := param.(flow.Tag)

		switch p.Tag {
		case "baud":
			baud = uint32(p.Msg.(float64))
		case "databits":
			databits = uint8(p.Msg.(float64))
		case "stopbits":
			stopbits = uint8(p.Msg.(float64))
		case "init":
			initdata = append(initdata,p.Msg) //initialization sequence
		}

	}




	if port, ok := <-w.Port; ok {
		opt :=goserialproxy.Options{BitRate: baud, DataBits: databits, StopBits: stopbits}
		//opt := rs232.Options{BitRate: baud, DataBits: databits, StopBits: stopbits}
		dev, err := goserialproxy.NewSerialProxy(port.(string), opt)
		flow.Check(err)

		// try to avoid kernel panics due to that wretched buggy FTDI driver!
		// defer func() {
		// 	time.Sleep(time.Second)
		// 	dev.Close()
		// }()
		// time.Sleep(time.Second)


		//handle initialization data (this loop only happens once)
		go func() {
			for _,data := range initdata {
				switch data.(type) {
					case map[string]interface{}:
						hash := data.(map[string]interface{})

						for k,v := range hash {
							switch k {
								case "delay":
									if d,ok := v.(float64);ok {
										<-time.After(time.Millisecond*time.Duration(int(d)))
									}
							}

						}
					default:
						_,_ = writeHandler(dev, data)
				}
			}
		}()

		// separate process to copy data out to the serial port
		go func() {
			for m := range w.To {
				_,_ = writeHandler(dev, m)
			}
		}()

		scanner := bufio.NewScanner(dev)
		for scanner.Scan() {
			msg:=scanner.Text()
			w.From.Send(msg)
		}
	}
}

//writeHandler is a generic type converter for Serial Input (used by both .Param and .To )
//func writeHandler(dev *rs232.Port,  m interface{} ) (int,error) {
func writeHandler(dev goserialproxy.SerialPort,  m interface{} ) (int,error) {

	switch v := m.(type) {
	case string:
		return dev.Write([]byte(v + "\n"))
	case []byte:
		return dev.Write(v)
	case int:
		if err:= dev.SetDTR(true); err==nil { // pulse DTR to reset
			time.Sleep(time.Duration(v) * time.Millisecond)
			return 0,dev.SetDTR(false)
		}
	case bool:
		return 0,dev.SetRTS(v)
	}

	return 0,nil //TODO: we lost this type - perhaps log?
}
