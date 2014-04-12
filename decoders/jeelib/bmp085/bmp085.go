// Decoder for the "BMP085demo.ino" sketch as in: http://github.com/jcw/jeelib/tree/master/examples/Ports/bmp085demo/bmp085demo.ino
// Decoder Registers as "Node-Bmp085".
//
// See also BMP085Batt for a decoder with lobat warning
package decoders

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/glog"
	"github.com/jcw/flow"
)

func init() {
	flow.Registry["Node-Bmp085"] = func() flow.Circuitry { return &Bmp085{} }

}

type Bmp085 struct {
	flow.Gadget
	In  flow.Input
	Out flow.Output
}

// Note:see jeelib for why 'Press' is 32bit
type Bmp085Data struct {
	Node  uint8
	Temp  uint16
	Press uint32
}

// Start decoding Bmp085 packets.
func (w *Bmp085) Run() {
	if glog.V(2) {
		glog.Infoln("BMP085 starts")
	}

	for m := range w.In {

		if v, ok := m.([]byte); ok && len(v) >= 8 {
			buf := bytes.NewReader(v)
			var data Bmp085Data
			_ = binary.Read(buf, binary.LittleEndian, &data)

			m = map[string]int{
				"<reading>": 1,
				"temp":      int(data.Temp),
				"pressure":  int(data.Press),
			}
		}

		w.Out.Send(m)
	}
}
