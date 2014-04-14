//package rfdata provides both replacement (enhancements) to the core housemon rfdata package and new functionality
//not found in the existing core package
package rfdata

import (
	"fmt"
	_"errors"
	"time"
	"github.com/jcw/flow"
	"github.com/jcw/jeebus/gadgets"
	_ "github.com/jcw/housemon/gadgets/rfdata"
)

//We override the PutReadings found in the core housemon/rfdata package
func init() {
	//fmt.Println("Installing ext PutReadings")
	flow.Registry["PutReadings"] = func() flow.Circuitry { return &PutReadings{} }
}

// Save readings in database.
type PutReadings struct {
	flow.Gadget
	In  flow.Input
	Out flow.Output
}

// Convert each loosely structured reading object into a strict map for storage.
func (g *PutReadings) Run() {
	for m := range g.In {
		r := m.(map[string]flow.Message)

		values := r["reading"].(map[string]int)
		asof, ok := r["asof"].(time.Time)
		if !ok {
			asof = time.Now()
		}
		node := r["node"].(map[string]int)
		if node["rssi"] != 0 {
			values["rssi"] = node["rssi"]
		}
		rf12 := r["rf12"].(map[string]int)
		location, _ := r["location"].(string)
		decoder, _ := r["decoder"].(string)

		id := fmt.Sprintf("RF12:%d:%d:%d", rf12["band"], rf12["group"], node["<node>"])
		data := map[string]interface{}{
			"ms":  jeebus.TimeToMs(asof),
			"val": values,
			"loc": location,
			"typ": decoder,
			"id":  id,
		}
		g.Out.Send(flow.Tag{"/reading/" + id, data})
	}
}
