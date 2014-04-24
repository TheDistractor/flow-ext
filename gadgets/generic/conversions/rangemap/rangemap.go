//Package conversions provides Gadgets that offer some form of data conversion process
package conversions

import (
	"github.com/golang/glog"
	"github.com/jcw/flow"
	"math"
)

func init() {
	flow.Registry["RangeMap"] = func() flow.Circuitry { return new(RangeMapper) }
	glog.V(2).Infoln("RangeMapper loaded")
}


// RangeMapper helps to convert an input value to an equivalent output value within a set scale
type RangeMapper struct {
	flow.Gadget
	Param    flow.Input
	In       flow.Input
	Out      flow.Output
}


//Gadget loop
func (g *RangeMapper) Run() {

	//need upper and lower bounds input
	//gi := &RangeMapData{fromLow:0,fromHi:1023,toLow:0,toHi:255}  //typical ADC to PWM
	gi := NewRangeMap()
	for param := range g.Param {

		switch param.(type) {

		case flow.Tag:
			switch param.(flow.Tag).Tag {
			case "fromlow":
				gi.SetFromLow( int64(param.(flow.Tag).Msg.(float64) ) )
			case "fromhi":
				gi.SetFromHi( int64(param.(flow.Tag).Msg.(float64) ) )
			case "tolow":
				gi.SetToLow( int64(param.(flow.Tag).Msg.(float64) ) )
			case "tohi":
				gi.SetToHi( int64(param.(flow.Tag).Msg.(float64) ) )
			}
		}
	}


	if !gi.Valid() {
		glog.Fatal()
	}


	for m := range g.In {

		switch m.(type) {
		case flow.Tag: //the input is a flow.Tag
			switch m.(flow.Tag).Msg.(type) {
			case float64:
				if f,ok := m.(flow.Tag).Msg.(float64); ok {
					if v,ok := gi.Map(int64(f));ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:float64(v)})
					}
				}
			case float32:
				if f,ok := m.(flow.Tag).Msg.(float32); ok {
					if v,ok := gi.Map(int64(f));ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:float32(v)})
					}
				}
			case int64:
				if f,ok := m.(flow.Tag).Msg.(int64); ok {
					if v,ok := gi.Map(f); ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:int64(v)})
					}
				}
			case int32:
				if f,ok := m.(flow.Tag).Msg.(int32); ok {
					if v,ok := gi.Map(int64(f));ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:int32(v)})
					}
				}
			case uint32:
				if f,ok := m.(flow.Tag).Msg.(uint32); ok {
					if v,ok := gi.Map(int64(f));ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:uint32(v)})
					}
				}
			case int16:
				if f,ok := m.(flow.Tag).Msg.(int16); ok {
					if v,ok := gi.Map(int64(f));ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:int16(v)})
					}
				}
			case uint16:
				if f,ok := m.(flow.Tag).Msg.(uint16); ok {
					if v,ok := gi.Map(int64(f));ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:uint16(v)})
					}
				}
			case int8:
				if f,ok := m.(flow.Tag).Msg.(int8); ok {
					if v,ok := gi.Map(int64(f));ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:int8(v)})
					}
				}
			case uint8:
				if f,ok := m.(flow.Tag).Msg.(uint8); ok {
					if v,ok := gi.Map(int64(f));ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:uint8(v)})
					}
				}
			case int:
				if f,ok := m.(flow.Tag).Msg.(int); ok {
					if v,ok := gi.Map(int64(f));ok {
						g.Out.Send(flow.Tag{Tag:m.(flow.Tag).Tag , Msg:int(v)})
					}
				}


			}
		default: //the input is a normal flow.Message
			switch m.(type) {
			case float64:
				f:= m.(float64)
				if v,ok := gi.Map(int64(f));ok {
					g.Out.Send( float64(v) )
				}
			case float32:
				f:= m.(float32)
				if v,ok := gi.Map(int64(f));ok {
					g.Out.Send( float32(v) )
				}
			case  int64:
				f:= m.(int64)
				if v,ok := gi.Map(int64(f));ok {
					g.Out.Send( int64(v) )
				}
			case  int32:
				f:= m.(int32)
				if v,ok := gi.Map(int64(f));ok {
					g.Out.Send( int32(v) )
				}
			case  int16:
				f:= m.(int16)
				if v,ok := gi.Map(int64(f));ok {
					g.Out.Send( int16(v) )
				}
			case  int8:
				f:= m.(int8)
				if v,ok := gi.Map(int64(f));ok {
					g.Out.Send( int8(v) )
				}
			case  int:
				f:= m.(int)
				if v,ok := gi.Map(int64(f));ok {
					g.Out.Send( int(v) )
				}

			}
		}

	}



}


//TODO: move this out to a generic package and call back in.
type RangeMap interface {
	Map(int64) int64
	SetToLow(int64) bool
	SetToHi(int64) bool
	SetFromLow(int64) bool
	SetFromHi(int64) bool
}

type RangeMapData struct {
	fromLow int64
	fromHi int64
	toLow int64
	toHi int64
	bit uint8
}

//NewRangeMap creates a RangeMap that must be provided a From/To range later otherwise its not valid
func NewRangeMap() (*RangeMapData) {
	return &RangeMapData{}
}

//NewRangeMapFrom creates a RangeMap with the From/To values provided
func NewRangeMapFrom(fromLo, fromHi, toLo, toHi int64) (*RangeMapData) {
	return &RangeMapData{fromLow: fromLo, fromHi: fromHi, toLow:toLo, toHi: toHi, bit:15}

}

func (r *RangeMapData) Valid() (bool) {
	return r.bit == 15
}

func (r *RangeMapData) SetFromLow(v int64) (bool) {
	r.fromLow = v
	r.setBit(0)
	return true
}
func (r *RangeMapData) SetFromHi(v int64) (bool) {
	r.fromHi = v
	r.setBit(1)
	return true
}
func (r *RangeMapData) SetToLow(v int64) (bool) {
	r.toLow = v
	r.setBit(2)
	return true
}
func (r *RangeMapData) SetToHi(v int64) (bool) {
	r.toHi = v
	r.setBit(3)
	return true
}

func (r *RangeMapData) setBit(pos uint) () {
	r.bit |= 1 << pos
}


//MustMap will panic if the Map function returns false, otherwise it silently returns Map value.
func (r *RangeMapData) MustMap(v int64) (int64) {
	if v,ok := r.Map(v);ok {
		return v
	}
	panic("RangeMap parameters not correct, check inputs")

}

//Map will transform the input (between fromLow, fromHi) to a value within the range (toLow, toHi)
//using this simple formula: = (x - in_min) * (out_max - out_min) / (in_max - in_min) + out_min;
//input is converted to float64 before calculations and the result is rounded back down
//to fit an int64 output, where 1.51 = 2 and 1.49=1
//output is capped to stay within output range
//this is a generalised function for *MY* wide use-cases, there are many ways to get specific results faster and more efficiently
func (r *RangeMapData) Map(v int64) (int64,bool) {

	if ! r.Valid() {
		return 0,false
	}

	f := float64(v-r.fromLow)/float64(r.fromHi-r.fromLow) * float64(r.toHi-r.toLow) + float64(r.toLow)

	fr := RoundPrec(f,0)
	if fr > float64(r.toHi) {
		fr = float64(r.toHi)
	} else if fr < float64(r.toLow) {
		fr = float64(r.toLow)
	}


	return  int64( fr ), true
}

//generic Precision rounding of float64 (go-lang-nuts discussion)
func RoundPrec(x float64, prec int) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}

	sign := 1.0
	if x < 0 {
		sign = -1
		x *= -1
	}

	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)

	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow * sign
}

