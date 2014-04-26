package mqtt

import (
	"net"

	"github.com/golang/glog"
	"github.com/jcw/flow"
	mqttex "github.com/TheDistractor/flow-ext/gadgets/network/mqtt/extended"

)

//Add additional MQTTServerEx functionality
func init() {
	if glog.V(2) {
		glog.Infoln("Loading MQTTServerEx as MQTTServerEx into Registry...")
	}
	flow.Registry["MQTTServerEx"] = func() flow.Circuitry { return &mqttex.RemoteMQTTServer{} }
}
