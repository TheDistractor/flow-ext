//package rfdata provides both replacement (enhancements) to the core housemon rfdata package and new functionality
//not found in the existing core package
package rfdata

import (
	"fmt"
	"strconv"
	"strings"

	"errors"
	"github.com/jcw/flow"
	_ "github.com/jcw/housemon/gadgets/rfdata"
)

//We override the NodeMap found in the core housemon/rfdata package
func init() {
	fmt.Println("Installing ext NodeMap")
	flow.Registry["NodeMap"] = func() flow.Circuitry { return &NodeMap{} }
}

// Lookup the group/node information to determine what decoder to use.
// Registers as "NodeMap".
type NodeMap struct {
	flow.Gadget
	Info    flow.Input
	In      flow.Input
	Out     flow.Output
	Missing flow.Output //provides ability to capture data we have not got a nodeMap entry for
}

//Composite key of 3 ints, provides additional <band> support allowing the input nodeMap to contain frequency/band
//which in turn allows us to support inbound packets from different networks such as 433Mhz and 868Mhz.
type NodeMapKey struct {
	band  int
	group int
	node  int
}

func (k *NodeMapKey) String() string {
	return fmt.Sprintf("RFb%dg%di%d", k.band, k.group, k.node)
}

// Start looking up node ID's in the node map.
func (w *NodeMap) Run() {

	defaultBand := 868 //TODO:Change this to input parameter

	nodeMap := map[NodeMapKey]string{}
	locations := map[NodeMapKey]string{}
	for m := range w.Info {
		f := strings.Split(m.(string), ",")

		key := NodeMapKey{}
		if ok, err := key.Unmarshal(f[0]); !ok {
			flow.Check(err)
		}

		//for the case where the default data has not been changed as in:
		// { data: "RFg5i2,roomNode,boekenkast JC",  to: "nm.Info" }
		//this will automatically be incorporated into the defaultBand network.
		if key.band == 0 {
			key.band = defaultBand
		}

		nodeMap[key] = f[1]

		if len(f) > 2 {
			locations[key] = f[2]
		}
	}

	key := NodeMapKey{}
	for m := range w.In {
		w.Out.Send(m)

		if data, ok := m.(map[string]int); ok {

			switch {
			case data["<RF12demo>"] > 0:
				key.group = data["group"]
				key.band = data["band"]
			case data["<node>"] > 0:
				key.node = data["<node>"]
				if loc, ok := locations[key]; ok {
					w.Out.Send(flow.Tag{"<location>", loc})
				} else {
					w.Missing.Send(key)
					//fmt.Printf("Location NOT found:%+v", key)
				}
				if tag, ok := nodeMap[key]; ok {
					w.Out.Send(flow.Tag{"<dispatch>", tag})
				} else {
					//fmt.Printf("NodeMap NOT found:%+v", key)
					w.Missing.Send(key)
					w.Out.Send(flow.Tag{"<dispatch>", ""})
				}
			}
		}
	}
}

//Unmarshal unpacks a nodeMap entry into the key structure
func (k *NodeMapKey) Unmarshal(s string) (bool, error) {

	prefix := "RF"

	ok := strings.HasPrefix(s, prefix)
	if !ok {
		return false, errors.New("RF NodeMapping must begin 'RF'..." + s)
	}

	rdr := strings.NewReader(s[len(prefix):])

	token := ""
	band := ""
	group := ""
	node := ""

	var err error

	for err == nil {
		var c byte
		c, err = rdr.ReadByte()

		switch string(c) {
		case "b", "g", "i":
			token = string(c)
			continue
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if token == "b" {
				band += string(c)
			} else if token == "g" {
				group += string(c)
			} else {
				node += string(c)
			}
			continue
		}

	}

	if band != "" {
		if k.band, err = strconv.Atoi(band); err != nil {
			return false, err
		}
	}

	if group != "" {
		if k.group, err = strconv.Atoi(group); err != nil {
			return false, err
		}
	}

	if node != "" {
		if k.node, err = strconv.Atoi(node); err != nil {
			return false, err
		}
	}

	return true, nil

}
