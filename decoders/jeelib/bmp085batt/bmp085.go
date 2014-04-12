// Decoder for the "BMP085demoBatt.ino" sketch as defined: http://github.com/TheDistractor/ino/tree/master/jeelib/examples/Ports/bmp085demobatt/bmp085demobatt.ino
//
//This decoder is slightly different than the base BMP085 decoder (also in this package) in that it allows for the node
//to send back a Battery low flag (the value of which is calculated in your sketch using perhaps the internal bandgap method, as per example above)
//
// Registers as "Node-Bmp085Batt".
package decoders

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/glog"
	"github.com/jcw/flow"
)

func init() {
	flow.Registry["Node-Bmp085Batt"] = func() flow.Circuitry { return &Bmp085Batt{} }

}

type Bmp085Batt struct {
	flow.Gadget
	In  flow.Input
	Out flow.Output
}

// Note:see jeelib for why 'Press' is 32bit
type Bmp085BattData struct {
	Node  uint8
	Temp  uint16
	Press uint32
	Lobat uint8
}

// Start decoding Bmp085Batt packets.
func (w *Bmp085Batt) Run() {
	if glog.V(2) {
		glog.Infoln("BMP085Batt starts")
	}

	for m := range w.In {

		if v, ok := m.([]byte); ok && len(v) >= 8 {
			buf := bytes.NewReader(v)
			var data Bmp085BattData
			_ = binary.Read(buf, binary.LittleEndian, &data)

			m = map[string]int{
				"<reading>": 1,
				"temp":      int(data.Temp),
				"pressure":  int(data.Press),
				"lobat":     int(data.Lobat),
			}
		}

		w.Out.Send(m)
	}
}
