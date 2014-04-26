# flow-ext
---------

Flow based Extensions, Gadgets, Decoders etc for jcw/flow based apps.


* **Update 2014/04/23:** Added RangeMap Gadget.
* **Update 2014/04/21:** Updated SerialPortEx Gadget.
* **Update 2014/04/26:** MQTTServerEx packages reorganised to make them more flexible like SerialPortEx.


## Decoders
----------

### Pressure and Temperature
------------------------

Bmp085 and Bmp085Batt, effectively a decoder for the [Bosch BMP085](http://www.digikey.com/uk/en/ph/bosch/bmp085.html)
and used within the [Jeelabs Pressure Plug](http://jeelabs.net/projects/hardware/wiki/Pressure_Plug)

See [more info about the decoders](https://github.com/TheDistractor/flow-ext/blob/master/decoders/jeelib/bmp085.md)


## Gadgets
---------

### Generic
----------

#### ReadlineStdIn
Sometimes you may want to pipe data into Jeebus/Housemon from the commandline. Perhaps you have a python script that
outputs data to its stdout. Use this Gadget to capture line orientated input to your apps stdin.
To use, simply include this line in your imports

```go
	_ "github.com/TheDistractor/flow-ext/gadgets/jeebus/io/readline"  //ReadlineStdIn
```

You now have a ReadlineStdIn Gadget whos' **.Out** pin represents the line orientated data read from stdin. You you can
hook this pin up to an Input pin of another gadget.


#### RangeMap
It is often a requirement to map a value from one range constraint to another, for instance mapping a 10bit ADC (0-1023)
to an 8bit PWM (0-255). The RangeMap Gadget helps you achieve just that. My main use for this is to actually map
a 8bit (light intensity) LDR value from roomnode's to an on/off (0/1).

Usage:

Add to imports:

```
	_ "github.com/TheDistractor/flow-ext/gadgets/generic/conversions/rangemap"  //RangeMapper
```

Provide the input range values as flow.Tag's such as:

```json
    { tag: "fromlow", data: 0, to: "map.Param" }
    { tag: "fromhi", data: 255, to: "map.Param" }
    { tag: "tolow", data: 0, to: "map.Param" }
    { tag: "tohi", data: 1023, to: "map.Param" }
```
(This example will map a typical PWM value back to an 10bit ADC range).

Then wire up the **.In** and **.Out** pins within your circuit.
The **.In** pin could be the value from a sensor feed (via the MQTTSub Gadget).
The **.Out** pin could be to the **.To** pin of a Serial Gadget connected to an arduino/jeenode controlling another device.

If your arduino was driving a relay, you could use the light intensity to switch it on/off using something like:

```json
    { tag: "fromlow", data: 255, to: "map.Param" }
    { tag: "fromhi", data: 0, to: "map.Param" }
    { tag: "tolow", data: 0, to: "map.Param" }
    { tag: "tohi", data: 1, to: "map.Param" }
```
i.e If the light intensity is more than 50% (128) the output would be 1 and less that 50%, it would be 0. You could
reverse the login by switching the 'fromlow' and 'fromhi' values.



### Flow focused
---------------

#### DeadOutPin

This utility Gadget can be used within a circuit to sit on an input PIN of another Gadget that would normally expect
Input flows, but in cases where your circuit does not want to supply any. This will stop your Gadget from spinning
on nil input values. I would expect this Gadget to be un-necessary in future revisions of 'flow', but for now it can
help Gadget development in specific cases. I use it in gadgets that can have multiple inputs, that in some circuits only
'some' of the inputs are connected.

Add to imports:

```
	_ "github.com/TheDistractor/flow-ext/gadgets/flow/deadoutpin"  //DeadOutPin
```

Incorporate into a circuit:

```json
    { from: "<DeadOutPin>.Out", to:"<TargetGadget>.To"} #stop it complaining
```


### HouseMon focused
-------------------

#### LogArchiverTGZ

This Gadget takes the output from the Core Logger gadget as its input. It also takes an input mask to specify its
actions. A typical setup is to set the output mask to monthly (tar) to add each daily log to a monthly tar file.
We then take the output of this gadget, and supply it as the input to another instance of LogArchiverTGZ, but this
time with a mask of monthly (gz), which will take the monthly tar file from the previous gadget instance and turn it
into a tar.gz file. More info to follow, however the package has some documentation already.

#### RadioBlippers (Simulation)

This Gadget allows you to simulate a number of radioBlip nodes on specific RF Network groups.

#### NodeMap (extended core Gadget)
NodeMap replaces the core NodeMap gadget to support the Band/Frequency parameter. This will allow you to use:

        { data: "RFb433g5i2,roomNode,boekenkast JC",  to: "nm.Info" }
        { data: "RFb868g5i2,roomNode,keuken",  to: "nm.Info" }

**Important**: You must also override the Readings gadget if you want your database to process this extended information
correctly, and is an absolute must if you have the same group (e.g. g5) on both 868 and 433 Mhz, otherwise one bands
data will interleave with another.

( **Note**: I will be submitting a derivative of this to core shortly)

#### PutReadings (extended core Gadget)
PutReadings replaces the core PutReadings gadget to support the Band/Frequency parameter. This will allow you to use:

        { data: "RFb433g5i2,roomNode,boekenkast JC",  to: "nm.Info" }
        { data: "RFb868g5i2,roomNode,keuken",  to: "nm.Info" }
**Important**: You will need the revised NodeMap (see above).

**Note**:I have also published [convert-rf-readings](https://github.com/TheDistractor/convert-rf-readings) which allows
you to convert between the two formats. convert-rf-readings has basic documentation.

#### OnOffMonitor
OnOffMonitor allows you to manage On/Off events within the context of 'time' and 'duration'. It consumes events you
specify from DataSub and generates one or more additional 'related' events. As an example, you can listen for roomNode
'moved' events and generate another event 20min in the future if the state of the endpoint has not changed.
You can then hook into this event with another appropriate Gadget to handle the new event.


### Jeebus focused
-----------------

#### SerialPortEx (extended core Gadget)
The SerialPortEx Gadget is an extended SerialPort. It contains a .Param pin to allow configuration of such things
as Baud, Databits & StopBits. It otherwise operates directly as per the standard SerialPort gadget and is directly
interchangeable. Its default configuration is as per the core SerialPort Gadget.
This Gadget is appropriate if you wish more control over the serial port, such as if you are using an ATTINY based micro
like the JNu,ATTiny85,ATTiny84 or a GPS or bluetooth module.

There are actually two 'variants' of this Gadget. They both do the "same" thing, but are implemented in different ways:

* If you load the Gadget via the serial/compat package, it will just replace the standard SerialPort implementation.

```
	_ "github.com/TheDistractor/flow-ext/gadgets/jeebus/serial/compat"  //override default SerialPort with SerialPortEx

```
* If you use the serial/serialex package it will not replace the standard SerialPort, but rather provide a SerialPortEx
Gadget within the Flow Registry.

```
	_ "github.com/TheDistractor/flow-ext/gadgets/jeebus/serial/serialex"  //provide SerialPortEx
```

* If you use serial/extended you will be able to access the native SerialPortEx Gadget, but it will NOT be added to the
Flow Registry as a Gadget - in which case registration is up to you to add it to the Registry/Circuit yourself.

```
	_ "github.com/TheDistractor/flow-ext/gadgets/jeebus/serial/extended"  //SerialPortEx but NOT added to Registry
```

**Update 2014/04/20** - The .Param pin now specifically supports an 'init' parameter, that you can use to send *initial*
data to the serial port in a one off manner (perhaps used to set things up)

```json
    { tag:"init", data: "v", to: "sp.Param" }
```

The 'init' parameter also supports the concept of a delay (in ms), that can be used to separate time sensitive init sequences:

```json
    { tag:"init", data: {delay:20}, to: "sp.Param" }
```

These 'init' sequences are replayed in the order they are received.

( **Note**: I will be submitting a derivative of this to core shortly)

#### MQTTServerEx
If you choose to use an external MQTT broker like RabbitMQ or Mosquitto, use this Gadget to replace the inbuilt
Gadget. it simply provides you with a quick check to confirm an external broker is visible on the chosen url/port.
It then steps out of the way allowing your app to talk to the external broker.

This Gadget can be used in three different ways like the SerialPortEx above.

* Including package as mqtt/compat  will replace the core MQTTServer with MQTTServerEx, effectively using a remote broker.

```
	_ "github.com/TheDistractor/flow-ext/gadgets/network/mqtt/compat"  //override the default server with remote broker
```

* Including package as mqtt/mqttex  will include MQTTServerEx in the Registry and will not replace MQTTServer.

```
	_ "github.com/TheDistractor/flow-ext/gadgets/network/mqtt/mqttex"  //override the default server with remote broker
```

* Lastly, including package as mqtt/extended will provide the Gadget MQTTServerEx but NOT directly add it to flow
Registry. You must add it yourself using whatever alias you prefer.


#### HTTPServer
This version of HTTPServer supports HTTP(S):// and WS(S)://. It can be loaded to override the existing HTTPServer
implementation within the core packages, simply import as:

```go
	_ "github.com/TheDistractor/flow-ext/gadgets/network/http"  //HTTPServer with https and ws-protocol selection
```

An additional "Feed" input of 'Param' is supplied where you can include the
following parameters as flow.Tag{} entries:

'certfile'

'certkey'

```json

   feeds: [
     { tag:"certfile", data: "/path/to/cert.pem", to: "http.Param" }
     { tag:"certkey", data: "/path/to/cert.key", to: "http.Param" }
     #...more feeds
   ]

```

If you want to use this in Housemon, its HTTPServer is defined in main.go (and overrides the one in json circuit) so you
should add the appropriate parameters there:

```go
main.go:

	c.Feed("http.Param", flow.Tag{"certfile", "/path/to/cert.pem"})
	c.Feed("http.Param", flow.Tag{"certkey", "/path/to/cert.key"})
```

*Note*: When you add a valid certificate/key your server will switch to HTTPS:// and your websocket support will also
switch to WSS:// using whatever PORT you defined.
(Your browser/client may warn you if your server certificate is untrusted, you should use the appropriate commands
for your os/client to enable this trust)

**Important**: This revised gadget has been submitted to the core jeebus/housemon repo's , until its features are incorporated
you should also check and add the following line to 'jeebus.coffee' within jeebus to enable wss:// support.
[see here, jeebus.coffee lines 49:50](https://github.com/TheDistractor/jeebus/commit/7cd3c80eb2fe158ae597c4daa02203ef3471f28e#diff-f4e44c99773d98dee8fb4e934fad59e5R5)

```go
-      ws = new WebSocket "ws://#{location.host}/ws", [appTag]
 +      wsProto = (if "https:" is document.location.protocol then "wss://" else "ws://")
 +      ws = new WebSocket "#{wsProto}#{location.host}/ws", [appTag]
```










