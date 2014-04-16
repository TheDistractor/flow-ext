// OnOffMonitor
// Author: TheDistractor, aka-lightbulb
// A HouseMon Gadget to allow you to react to Boolean (0/1) inputs within a 'time' context
//
// This Gadget came about as I wanted to react to situations when I CAN/CANNOT detect movement over a set
// period - I wanted to make sure lights can be turned off in a rooms where there may be no occupants, just
// like a traditional time-delay switch may operate, or turn heating on when a room seems to have been occupied
// for a while, or switch a relay off after it has been running for a set period of time etc, but with the ability
// to combine this output with other environmental inputs to make a bigger 'decision tree'.
//
// Specifically useful for sensor outputs that emit a 0 or 1 [Off / On] (Inverted values also supported)
// The Gadget takes a Flow.Tag feed that mimic's the Tag/Msg format used with MQTTSub,
// where Tag is the MQTT Topic and Msg is the Value
//
// The use-case is that input using 'sensor' topics will be the primary circuit feed (as per MQTTSub or DataSub),
// but any input conforming to the '<sensor>/<location>/<attribute>/<timestamp> 0|1 should work.
// Infact, <anything>/<filteredname>/<optionals...>/<timestamp> 0!1 should also work
//
// A Typical use would be to listen for sensor 'moved' events (as in RoomNode), or to listen for reed switch
// open/close etc, and then to monitor this state and generate further time-based events
// such as:
//          On T (as soon as an On event is seen)
//          Off T (as soon as an Off event is seen)
//          On-Since T (the attribute has not changed from this state Since T)
//          Off-Since T (ditto)
//          On-For D (the attribute has been at this state for the duration defined by D)
//          Off-For D (ditto)
//
//
//
// emitted messages take the form:
//
//          by/lb/oomon/<location>/<eventName>/<state>[-modifier]   UnixTimeMS | Duration
//
// And would typically be used as input to the MQTTPub component
//
//
// where:
//   <location> relates to the sensors 'Location' as defined by conforming input above.
//   (as stated above, a 'sensor' message typically looks like: sensor/garage/moved/<timestamp> 0|1 )
//
//   <eventName> is a name you (as the user) provide as a parameter. 'onoff' is used by default
//   I use 'motion' for roomnode 'moved' detection, but you can use anything.
//
//   <state> will be either 'On' or 'Off' (an Inverted flag can be used as feed input to invert the usage).
//
//   [-modifier] will be one of:
//     '-Since' - descibed above
//     '-For' - described above
//
// The Value of each emitted message will be the 'reference' time of the event in UnixTime(Millisecond) format or
// standard go 'Duration' syntax for -For messages.
//
// See the example circuits for usage
package statemanagement

import (
	"errors"
	"fmt"
	"github.com/TheDistractor/flow-ext/go-helpers/int64utils"
	"github.com/golang/glog"
	"github.com/jcw/flow"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

//Add a basic OnOffMonitor to the flow registry
func init() {
	flow.Registry["OnOffMonitor"] = func() flow.Circuitry { return NewOnOffMonitor() }
}

//create a new instance of OnOffMonitor with appropriate construction
func NewOnOffMonitor() *OnOffMonitor {
	o := &OnOffMonitor{}
	return o
}

type OnOffMonitor struct {
	flow.Gadget
	Filter    flow.Input //Feed to set the Filters of Messages we act upon
	Param     flow.Input //Feed to setup basic Parameters
	Threshold flow.Input //Feed to setup Duration thresholds for each .Filter input

	In  flow.Input  //Inboud Flow circuit messages
	Out flow.Output //Outbound Flow circuit messages

}

type OnOffMonitorInst struct {
	Watched  map[string]*OnOffState //we manage each of these 'watches' as a separate state.
	baseName string                 //default base namespace for output events
	birth    int64                  //when this instance was first created (ms)

	stateOff     float64
	stateOn      float64
	stateUnknown float64
}

func NewOnOffMonitorInst() (*OnOffMonitorInst) {

	o := &OnOffMonitorInst{}

	o.Watched = make(map[string]*OnOffState)

	o.baseName = "by/ll/oomon/"
	o.birth = UnixMs(time.Now()) //ms

	//pass these down to state upon construction
	o.stateOff = float64(0)
	o.stateOn = float64(1)
	o.stateUnknown = float64(-1)

	return o
}

type OnOffState struct {
	Name       string                        //a moniker for our state
	On         int64                         //last time we saw On
	Off        int64                         //last time we saw Off
	Current    float64                       //the current state of the flag
	Thresholds map[string]*ThresholdDuration //map of Durations used for this State Item
	Remaining  []int64                       //map of remaining durations

	stateOff     float64
	stateOn      float64
	stateUnknown float64
}

//allows us to emit the user input duration string rather than the stringified duration
type ThresholdDuration struct {
	time.Duration
	Text string //the original string representation eg. 90s vs 1m30s
}

const (
	MaskSince = "%s/%s/%s/%s-Since" //the mask we use to generate a 'Since'  basename/<location>/<eventname>/<direction>-Since
	MaskOnOff = "%s/%s/%s/%s"      //the mask we use to generate On/Off basename/<location>/<eventname>/<direction>
	MaskFor   = "%s/%s/%s/%s-For"   //the mask we use to generate For basename/<location>/<eventname>/<direction>-For
)

//create a new State with correct initial construction
func (w *OnOffMonitorInst) NewState(createDate int64, name string) *OnOffState {

	state := &OnOffState{Name: name, On: createDate, Off: createDate, Current: w.stateUnknown, Thresholds: make(map[string]*ThresholdDuration), Remaining: []int64{}}

	//TODO:refactor names
	state.stateOn = w.stateOn
	state.stateOff = w.stateOff
	state.stateUnknown = w.stateUnknown

	//add this state to the overall 'watches'
	w.Watched[name] = state

	return state
}

// Start listening to the circuit inputs
func (w *OnOffMonitor) Run() {
	//float64 used as the common base type in these flow messages (as thats the type of the inboud message value for ON/OFF)

	if glog.V(2) {
		glog.Info("OnOffMonitor.Run")
	}

	wi := NewOnOffMonitorInst()


	checkSince := time.Second * 20 //how often we emit -Since - overriden by .Param

	//TODO:expose this as parameter
	treatUnknownAs := int(-1) //we can treat unknowns (states we have not yet seen) as a specific type.
	_ = treatUnknownAs

	eventName := "onoff" //the default eventName for this Gadget - overridden by .Param

	invert := false //invert the meaning of 0 | 1   (0=On,1=Off) - overridden by .Param

	//our input feed may well provide us with sensor data from lots of <locations>
	//we use Input Feeds to .Filter to create a set of <Locations> we want to use, rest are ignored
	//    { data: "Garage",  to: "oo.Filter" }
	// so...if you don't specify an Input .Filter, you wont get any outputs


	//process .Param inputs
	for param := range w.Param {

		p := param.(flow.Tag)

		if glog.V(2) {
			glog.Info("Param:", p)
		}

		switch p.Tag {
			case "invert":
				invert = p.Msg.(bool)
			case "eventname":
				eventName = p.Msg.(string)
			case "unknown":
				treatUnknownAs = p.Msg.(int)
			case "checkperiod":

				d, err := time.ParseDuration(p.Msg.(string))
				if err == nil {
					checkSince = d
				} else {
					if glog.V(2) {
						glog.Info("Invalid checkperiod", p.Msg.(string))
					}
				}

		}

	}

	//Does this Gadget instance invert meaning of 0 & 1
	if invert {
		wi.stateOff = float64(1)
		wi.stateOn = float64(0)
	}

	//process .Filter inputs
	for filter := range w.Filter {
		_ = wi.NewState(wi.birth, filter.(string))
		//TODO:ideally we should now pull last seen data from DB, instead of Now() so we get a LKV startup
		//but we need an appropriate stable DB API - revisit
	}

	//process .Thresholds  .Tag is <Location>, .Msg is a Duration like 2m30s
	for threshold := range w.Threshold {
		t := threshold.(flow.Tag)
		_ = wi.Watched[t.Tag].AddThreshold(t.Msg.(string))
	}
	if glog.V(3) {
		for wk, wv := range wi.Watched {
			glog.Info("Watched:", wk, wv.Thresholds)
		}
	}

	//we use this timer to provide -Since
	timerSince := time.NewTimer(checkSince)

	//we use this timer to provide -For
	timerFor := time.NewTimer(2)
	timerFor.Stop() //we limp to stop as we only want to start when we have a 'For' to shoot at.


	for {

		select {

		case f := <-timerFor.C:
			//check which states we can cover
			//timerFor.Stop()
			for sk,sv := range wi.Watched {
				fired , _ := sv.ExpiredThresholds(f)

				for _,d := range fired {
					stateDirection := "Off"
					if sv.Current == sv.stateOn {
						stateDirection = "On"
					}

					w.Out.Send(flow.Tag{
						fmt.Sprintf(MaskFor, wi.baseName, sk, eventName, stateDirection), d})

				}

			}
			//we have sent matching events, now reset timer for next fire
			next,err := wi.RecalcNextFor()
			if err == nil {
				timerFor.Reset(next.Sub(time.Now()))
			}

		case t := <-timerSince.C:

			for wk, wv := range wi.Watched {

				currentState := wv.Current

				var stateDirection = "Off"
				if currentState == wv.stateOn {
					stateDirection = "On"
				}

				if currentState == wv.stateOff { //how long off
					if wv.Off <= UnixMs(t.Add(-checkSince)) {

						w.Out.Send(flow.Tag{
							fmt.Sprintf(MaskSince, wi.baseName, wk, eventName, stateDirection), wv.Off})
					}
				} else if currentState == wv.stateOn { //how long On
					if wv.On <= UnixMs(t.Add(-checkSince)) {

						w.Out.Send(flow.Tag{
							fmt.Sprintf(MaskSince, wi.baseName, wk, eventName, stateDirection), wv.On})
					}
				}
			}

			timerSince.Reset(checkSince)

		case m := <-w.In: //process inbound message
			if data, ok := m.(flow.Tag); ok {

				if glog.V(2) {
					glog.Info("Inbound Message:", data, data.Tag, data.Msg)
				}

				//we expect this to be a 'sensor' reading 'sensor/<location>/attribute/timestamp VALUE'
				//but we only need <location> at position 1 and the timestamp at tail (2+).
				parts := strings.Split(data.Tag, "/")

				if len(parts) < 3 {
					if glog.V(2) {
						glog.Info("Not enough parts to be a valid message for this Gadget:", parts)
					}
					continue
				}

				location := parts[1]
				timestr := parts[len(parts)-1:][0] //we only check the tail part

				//when is milliseconds
				when, err := strconv.ParseInt(timestr, 10, 64)
				if err != nil { //the format did not contain timestring, substitute now
					when = time.Now().UnixNano() / 1e6
				}


				if match, ok := wi.Watched[location]; ok {

					stateChange := false
					//change in state?
					if match.Current != data.Msg.(float64) {
						match.Current = data.Msg.(float64)
						stateChange = true

						//TODO:don't refactor - we will integrate match.Unknown
						if match.Current == match.stateOff {
							if glog.V(2) {
								glog.Info("Setting OFF time:" + location)
							}
							match.Off = when
						}
						if match.Current == match.stateOn {
							if glog.V(2) {
								glog.Info("Setting ON time:" + location)
							}
							match.On = when
						}

					}
					_ = stateChange

					currentState := match.Current
					var stateDirection = "Off"
					if currentState == match.stateOn {
						stateDirection = "On"
					}

					//send our On or Off message
					w.Out.Send(flow.Tag{
						fmt.Sprintf(MaskOnOff, wi.baseName, location, eventName, stateDirection), when})


					wi.Watched[location] = match

					//if we have a stateChange we must rebuild remainings from thresholds
					if stateChange {
						wi.Watched[location].ResetRemaining()
						//and because we have new timeslots we must reset the timerFor
						_ = timerFor.Stop()
						next,err := wi.RecalcNextFor()
						var then time.Duration
						if err == nil {
							then = next.Sub(time.Now())
							timerFor.Reset(then)
							if glog.V(2) {
								glog.Info(  fmt.Sprintf("Timers Reset by %s, next event:%s", location, then))
							}
						}

					}

				}

			} //if

		} //select
	}
}


//find the closest remaining timeslot
func (s *OnOffMonitorInst) RecalcNextFor() (time.Time, error) {
	min := int64(math.MaxInt64)

	for _, sv := range s.Watched {
		near, err := sv.NextRemaining()
		if err != nil {
			continue
		}

		if near < min {
			min = near
		}
	}
	if min == math.MaxInt64 {
		return MsUnix(0,min), errors.New("No more events")
	}
	return MsUnix(0,min), nil
}


//attempt to add this input threshold to the state
func (s *OnOffState) AddThreshold(text string) bool {

	d, err := time.ParseDuration(text)
	if err == nil {
		td := &ThresholdDuration{d, text}
		//store using the 'normalized' stringified representation
		s.Thresholds[fmt.Sprintf("%s", d)] = td
		//watched[t.Tag].Thresholds[t.Msg.(string)] = d

	} else {
		if glog.V(2) {
			glog.Info("Invalid Threshold:", text)
		}
		return false
	}
	return true
}

//rebuild the Remaining Stack for the location
func (s *OnOffState) ResetRemaining() {

	remaining := []int64{}

	var when int64 //the reference time (used to create Remaining) so we get back the correct durations
	if s.Current == s.stateOn {
		when = s.On
	} else {
		when = s.Off
	}

	for tk, tv := range s.Thresholds {
		d := MsUnix(0, when)
		_ = tk
		remaining = append(remaining, UnixMs(d.Add(tv.Duration)))
	}

	sort.Sort(int64utils.Int64Array(remaining))
	s.Remaining = remaining

}

//peek the next timeslot
func (s *OnOffState) NextRemaining() (int64, error) {
	if len(s.Remaining) > 0 {
		return s.Remaining[0], nil
	}
	return 0, errors.New("Empty")
}

//remove and return any durations that have 'passed' as a series of durations
//now is  'common' across all OnOffStates
func (s *OnOffState) ExpiredThresholds(now time.Time) ([]string, error) {

	results := []string{} //the durations that have passed

	var tmp []int64 //whats left after removing passed durations

	var when int64 //the reference time (used to create Remaining) so we get back the correct durations
	if s.Current == s.stateOn {
		when = s.On
	} else {
		when = s.Off
	}

	for _, ms := range s.Remaining {
		if ms <= UnixMs(now) { //which slots have passed the input time
			dur := fmt.Sprintf("%v", MsUnix(0, ms).Sub(MsUnix(0, when)))
			td, ok := s.Thresholds[dur]
			if ok {
				results = append(results, td.Text)
			}
		} else {
			tmp = append(tmp, ms)
		}
	}

	s.Remaining = tmp

	return results, nil //
}


//Time to Unix with ms resolution
func UnixMs(t time.Time) int64 {
	return int64(t.UnixNano() / 1e6)
}

//Create Time using Unix seconds and MilliSeconds
func MsUnix(s int64, m int64) time.Time {
	return time.Unix(s, int64(m*1e6))
}
