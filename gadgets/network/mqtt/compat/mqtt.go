package mqttex

import (

	"github.com/golang/glog"
	"github.com/jcw/flow"
	_ "github.com/jcw/jeebus/gadgets/network"
	mqttex "github.com/TheDistractor/flow-ext/gadgets/network/mqtt/extended"

)

//Automatically override the core MQTTServer
func init() {
	if glog.V(2) {
		glog.Infoln("Loading MQTTServerEx as MQTTServer into Registry...")
	}
	flow.Registry["MQTTServer"] = func() flow.Circuitry { return &mqttex.RemoteMQTTServer{} }
}
