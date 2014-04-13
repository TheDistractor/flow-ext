//Package radioblippers provides a gadget to simulate a network of radioBLIP sketches without needing to deploy
//either radioBLIP nodes or a receiver node.
package radioblippers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	_ "github.com/golang/glog"
	"github.com/jcw/flow"
	"time"
)

//The main flow gadget, handles serial simulation as RFDemo.10 on the [.From] pin.
type RadioBlippers struct {
	flow.Gadget
	Param flow.Input //Feed to setup basic Parameters

	To   flow.Input  //Inboud Flow circuit messages
	From flow.Output //Outbound Flow circuit messages

}

//Simulates a radioBLIP end node
type RadioBlipper struct {
	int
	payload uint32
}

//Next increments the ping counter
func (r *RadioBlipper) Next() {
	r.payload += 1
}

//Automatic incorporation into flow registry
func init() {
	flow.Registry["RadioBlippers"] = func() flow.Circuitry { return new(RadioBlippers) }
}

var radioBands = map[int]interface{}{433: nil, 868: nil, 915: nil}

//Run is the main RadioBlippers gadget entry point.
//This gadget is used to simulate 1 to 30 radioBlip sketches on a specific band/group
//You may incorpoate this Gadget multiple times using different band/group combinations.
//Note: does NOT currently simulate the 'contention' issues that can be experienced on a real RF network.
//Use this to establish numerous 'fake' nodes on a netgroup. Don't forget to add them
//to your node/driver cross reference lookup tables.
func (g *RadioBlippers) Run() {

	band := int(-1)
	group := int(0)

	nodes := make(map[string]*RadioBlipper)

	//read params
	for param := range g.Param {

		p := param.(flow.Tag)

		switch p.Tag {
		case "band":
			band = int(p.Msg.(float64))
		case "group":
			group = int(p.Msg.(float64))
		case "node":
			node := int(p.Msg.(float64))
			if !(node >= 1 && node <= 30) {
				flow.Check(errors.New(fmt.Sprintf("Node %d is out of range 1-30", node)))
				continue
			}
			nodes[fmt.Sprintf("%d", node)] = &RadioBlipper{node, 0}
		}

	}

	if _, ok := radioBands[band]; !ok {
		flow.Check(errors.New(fmt.Sprintf("Band unsupported:%d (433,868,915)", band)))
	}

	if group < 1 || group > 250 {
		flow.Check(errors.New(fmt.Sprintf("Group unsupported:%d (1-250)", group)))
	}

	if len(nodes) == 0 {
		flow.Check(errors.New(fmt.Sprintf("No nodes loaded")))
	}

	<-time.After(time.Millisecond * 500)
	g.From.Send(fmt.Sprintf("[RF12demo.10] _ i31* g%d @ %d MHz", group, band)) //we immitate a collector on node 31

	receiver := time.NewTicker(1 * time.Minute)

	for {
		select {

		case <-receiver.C: //simulate RFDEMO incomming - radioBlips every 1 min
			//we send output messages that simulate the RadioBlip sketch via RF12Demo
			for _, v := range nodes {
				v.Next()
				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.LittleEndian, v.payload)

				bytes := buf.Bytes()

				msg := fmt.Sprintf("OK %d %d %d %d %d", v.int, bytes[0], bytes[1], bytes[2], bytes[3])
				<-time.After(time.Millisecond * time.Duration(v.int*500)) //vary how quickly they come in over the minute
				g.From.Send(msg)

			}

		}

	}

}
