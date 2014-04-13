// Numerous Gadgets and Decoders for the Flow based Jeebus/Housemon
//
// This is a convenience package to wrap all the Gadgets/Decoders available in the Library, you can however
// reference many as separate entities.
//
package flowext

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/jcw/flow"
	_ "github.com/jcw/flow/gadgets"

	_ "github.com/jcw/jeebus/gadgets/network"
	_ "github.com/TheDistractor/flow-ext/gadgets/network/extmqtt"

	_ "github.com/TheDistractor/flow-ext/decoders/jeelib/bmp085"
	_ "github.com/TheDistractor/flow-ext/decoders/jeelib/bmp085batt"

	_ "github.com/TheDistractor/flow-ext/gadgets/housemon/logging"
	_ "github.com/TheDistractor/flow-ext/gadgets/jeebus/serial/compat"

	_ "github.com/TheDistractor/flow-ext/gadgets/housemon/statemanagement"


)

var Version = "0.9.0"

func init() {

}


