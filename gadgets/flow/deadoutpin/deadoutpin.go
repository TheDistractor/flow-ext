//provide basic two way rf gateway services
package gateway


import (
	"time"

	"github.com/golang/glog"
	"github.com/jcw/flow"
)


//register with flow registry
func init() {
	glog.V(2).Info("DeadOutPin loaded...")
	flow.Registry["DeadOutPin"] = func() flow.Circuitry { return new(DeadOutPin) }
}

type DeadOutPin struct {
	flow.Gadget

	Out flow.Output

}

func (g *DeadOutPin) Run() {

	for {
		<-time.After(time.Hour * 1)
	}
}
