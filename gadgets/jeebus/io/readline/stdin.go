//Package readline supports Readline based Gadgets
package readline

import (
	"github.com/golang/glog"
	"github.com/jcw/flow"

	"bufio"
	"os"
)

func init() {
	glog.Info("ReadlineStdIn Init...")
	flow.Registry["ReadlineStdIn"] = func() flow.Circuitry { return new(ReadlineStdIn) }

}

//ReadlineStdIn allows data to enter a flow from the applications stdio
//this means of course that if you wish to use this gadget, stdio must be available to the application
//and stdio should not be hooked by another reader (such as the original HTTPServer, see the LiveReload Gadget
type ReadlineStdIn struct {
	flow.Gadget

	Out flow.Output //This gadget only supports Out (its input is read from stdin)
}

func (g *ReadlineStdIn) Run() {

	exit := false

	sin := bufio.NewScanner(os.Stdin)

	for exit == false {

		var line string
		for sin.Scan() {
			line = sin.Text()
			g.Out.Send(line)
		}

		err := sin.Err()

		if err != nil {
			glog.Errorln("ReadlineStdIn Error:", err)
		}

	}

}
