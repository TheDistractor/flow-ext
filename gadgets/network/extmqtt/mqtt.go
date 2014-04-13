//Package extmqtt implements a stub that can replace the inbuilt Jeebus MQTTServer to allow a remote (out of process)
//MQTT Broker like RabbitMQ or Mosquitto to be used. This Gadget does NOT implement this remote broker, rather
//it checks for the remote host/port is listening and passes its input (Port) parameter through and steps out the way.
//
//This can be useful if you want features that are not within the inbuilt MQTT Server, or you have additional external
//processes that want to share the MQTT broker that your jeebus app is using.
//
//Usage: (just add to your imports in main.go) - it will override the default core MQTTServer.
//
// 	_ "github.com/TheDistractor/flow-ext/gadgets/network/extmqtt"  //override the default server
package extmqtt

import (
	"net"

	"github.com/golang/glog"
	"github.com/jcw/flow"
	_ "github.com/jcw/jeebus/gadgets/network"
)

//Automatically override the core MQTTServer
func init() {
	if glog.V(2) {
		glog.Infoln("Remote Broker is attempting overriding inbuild MQTT...")
	}
	flow.Registry["MQTTServer"] = func() flow.Circuitry { return &RemoteMQTTServer{} }
}

//Our MQTTServer is really just a gated check for a remote MQTT Broker.
//because we pull in the original jeebus package, we make sure its init() runs before this one
//this then allows us to overwrite the MQTTServer before its used
type RemoteMQTTServer struct {
	flow.Gadget
	Port    flow.Input
	PortOut flow.Output
}

// Start the MQTT server.
func (w *RemoteMQTTServer) Run() {
	if glog.V(2) {
		glog.Infoln("RemoteMQTTBroker Run begins...")
	}

	port := getInputOrConfigwithDefault(w.Port, "MQTT_PORT", ":1883")

	//TODO: Perhaps add a 'real' server check by making an MQTT Client connection.
	conn, err := net.Dial("tcp", port)
	if err != nil { //This gives a more specific error log, scoped to file
		glog.Errorln("Error connecting to MQTT Server:", err)
	}
	flow.Check(err) //And then let flow panic.

	defer conn.Close()

	Done := make(chan bool)
	w.PortOut.Send(port)
	<-Done //we dont really need this!
}

//use the config/envvar setting unless overridden by flow param
func getInputOrConfigwithDefault(vin flow.Input, vname string, vdefault string) string {
	// if a port is given, use it, else use the default from the configuration. else use the default input
	var value string
	var ok bool
	if value, ok = flow.Config[vname]; !ok {
		value = vdefault
	}

	if m := <-vin; m != nil {
		value = m.(string)
	}
	if value == "" {
		glog.Errorln("no value given for:", vname)
	}
	return value
}
