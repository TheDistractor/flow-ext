package network

import (
	_"bufio"
	_ "crypto/tls"
	_"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/golang/glog"
	"github.com/jcw/flow"
	_ "github.com/jcw/flow/gadgets/pipe" //http needs Pipe for WebSocket-default, next step to split load gadgets!!
	_ "github.com/jcw/jeebus/gadgets/network" //pull in Jeebus so we can replace HTTPServer and still get WSLiveReload,default

)

func init() {
	flow.Registry["HTTPServer"] = func() flow.Circuitry { return new(HTTPServer) }
	fmt.Println("Extended HTTP(s) Server loaded")
}

var wsClients = map[string]*websocket.Conn{}

// HTTPServer is a .Feed( which sets up an HTTP server.
type HTTPServer struct {
	flow.Gadget
	Handlers flow.Input
	Param    flow.Input
	Port     flow.Input
	Out      flow.Output
}

type flowHandler struct {
	h http.Handler
	s *HTTPServer
}

func (fh *flowHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fh.s.Out.Send(req.URL)
	fh.h.ServeHTTP(w, req)
}


//information describing an http endpoint
type HttpEndpointInfo struct {
	uri url.URL //
	pem string  //path to pem file
	key string  //path to key file

}

func NewHttpEndpointInfo(addr, pem, key string) (*HttpEndpointInfo, error) {

	info := &HttpEndpointInfo{}

	//we create a uri to hold structured data
	host, port, _ := net.SplitHostPort(addr)
	if host == "" {
		host = "localhost"
	}

	info.uri = url.URL{Scheme: "http", Host: net.JoinHostPort(host, port)}

	info.pem = pem
	info.key = key

	if (pem != "") && (key != "") {
		info.uri.Scheme = "https"
	}

	return info, nil
}

// Set up the handlers, then start the server and start processing requests.
func (w *HTTPServer) Run() {
	mux := http.NewServeMux() // don't use default to allow multiple instances

	port := getInputOrConfig(w.Port, "HTTP_PORT") //TODO:This is dependant upon mqtt func, needs moving - lightbulb

	pem := ""
	key := ""

	for param := range w.Param {

		switch param.(type) {

		case flow.Tag:
			switch param.(flow.Tag).Tag {
			case "certfile":
				f := param.(flow.Tag).Msg.(string)
				if _, err := os.Stat(f); err == nil {
					glog.Infoln("Using Certfile:", f)
					pem = f
				}
			case "certkey":
				f := param.(flow.Tag).Msg.(string)
				if _, err := os.Stat(f); err == nil {
					glog.Infoln("Using Keyfile:", f)
					key = f
				}
			}
		}
	}

	info, _ := NewHttpEndpointInfo(port, pem, key)

	for m := range w.Handlers {
		tag := m.(flow.Tag)
		switch v := tag.Msg.(type) {
		case string:
			h := createHandler(tag.Tag, v, info)
			mux.Handle(tag.Tag, &flowHandler{h, w})
		case http.Handler:
			mux.Handle(tag.Tag, &flowHandler{v, w})
		}
	}

	go func() {
		// will stay running until an error is returned or the app ends
		defer flow.DontPanic()
		var err error
		if info.uri.Scheme == "https" {
			err = http.ListenAndServeTLS(info.uri.Host, info.pem, info.key, mux)
		} else {
			err = http.ListenAndServe(info.uri.Host, mux)
		}
		glog.Fatal(err)
		glog.Infoln("http started on", info.uri.Host)
	}()
	// TODO: this is a hack to make sure the server is ready
	// better would be to interlock the goroutine with the listener being ready
	time.Sleep(50 * time.Millisecond)
}

func createHandler(tag, s string, info *HttpEndpointInfo) http.Handler {
	// TODO: hook gadget in as HTTP handler
	// if _, ok := flow.Registry[s]; ok {
	// 	return http.Handler(reqHandler)
	// }
	if s == "<websocket>" {
		var wsConfig *websocket.Config
		var err error
		//TODO: use wss:// and TlsConfig if wanting secure websockets outside https
		wsproto := "ws://"
		if info.uri.Scheme == "https" {
			wsproto = "wss://"
		}
		if wsConfig, err = websocket.NewConfig(wsproto+info.uri.Host+tag, info.uri.String()); err != nil {
			glog.Fatal(err)
		}

		hsfunc := func(ws *websocket.Config, req *http.Request) error {

			tag := ""
			for _, v := range ws.Protocol { //check for first supported WebSocket- (circuit) protocol
				if flow.Registry["WebSocket-"+v] != nil {
					tag = v
					break
				}
			}
			ws.Protocol = []string{tag} //let client know we picked one

			return nil //errors.New("Protocol Unsupported")
		}
		wsHandshaker := websocket.Server{Handler: wsHandler,
			Config:    *wsConfig,
			Handshake: hsfunc,
		}
		return wsHandshaker
	}

	if !strings.ContainsAny(s, "./") {
		glog.Fatalln("cannot create handler for:", s)
	}
	h := http.FileServer(http.Dir(s))
	if s != "/" {
		h = http.StripPrefix(tag, h)
	}
	if tag != "/" {
		return h
	}
	// special-cased to return main page unless the URL has an extension
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if path.Ext(r.URL.Path) == "" {
			r.URL.Path = "/"
		}
		h.ServeHTTP(w, r)
	})
}

//wsHandler now used ws.Config as protocol handshake now supported
func wsHandler(ws *websocket.Conn) {
	defer flow.DontPanic()
	defer ws.Close()

	hdr := ws.Request().Header

	// keep track of connected clients for reload broadcasting
	id := hdr.Get("Sec-Websocket-Key")
	wsClients[id] = ws
	defer delete(wsClients, id)

	// the protocol name is used as tag to locate the proper circuit
	//lightbulb: We use the protocol provided by ws, rather than header, as that contains server accepted value
	tag := ws.Config().Protocol[0]

	fmt.Println("WS Protocol Selected:", tag)

	if tag == "" { //no specific protocol, lets opt for 'default' which just echoes (or return with no circuit!)
		tag = "default"
	}

	g := flow.NewCircuit()
	g.AddCircuitry("head", &wsHead{ws: ws})
	g.Add("ws", "WebSocket-"+tag) //the client has negotiated this support
	g.AddCircuitry("tail", &wsTail{ws: ws})
	g.Connect("head.Out", "ws.In", 0)
	g.Connect("ws.Out", "tail.In", 0)
	g.Run()
}

type wsHead struct {
	flow.Gadget
	Out flow.Output

	ws *websocket.Conn
}

func (w *wsHead) Run() {
	for {
		var msg interface{}
		err := websocket.JSON.Receive(w.ws, &msg)
		if err == io.EOF {
			break
		}
		flow.Check(err)
		if s, ok := msg.(string); ok {
			id := w.ws.Request().Header.Get("Sec-Websocket-Key")
			fmt.Println("msg <"+id[:4]+">:", s)
		} else {
			w.Out.Send(msg)
		}
	}
}

type wsTail struct {
	flow.Gadget
	In flow.Input

	ws *websocket.Conn
}

func (w *wsTail) Run() {
	for m := range w.In {
		err := websocket.JSON.Send(w.ws, m)
		flow.Check(err)
	}
}

func getInputOrConfig(vin flow.Input, vname string) string {
	// if a port is given, use it, else use the default from the configuration
	value := flow.Config[vname]
	if m := <-vin; m != nil {
		value = m.(string)
	}
	if value == "" {
		glog.Errorln("no value given for:", vname)
	}
	return value
}
