#flow-ext
========

Flow based Extensions, Gadgets, Decoders etc for jcw/flow based apps.


##Decoders
========

###Pressure and Temperature
------------------------

Bmp085 and Bmp085Batt, effectively a decoder for the [Bosch BMP085](http://www.digikey.com/uk/en/ph/bosch/bmp085.html)
and used within the [Jeelabs Pressure Plug](http://jeelabs.net/projects/hardware/wiki/Pressure_Plug)

See [more info about the decoders](https://github.com/TheDistractor/flow-ext/blob/master/decoders/jeelib/bmp085.md)


##Gadgets
=======

###HouseMon focused
-------------------

####LogArchiverTGZ

This Gadget takes the output from the Core Logger gadget as its input. It also takes an input mask to specify its
actions. A typical setup is to set the output mask to monthly (tar) to add each daily log to a monthly tar file.
We then take the output of this gadget, and supply it as the input to another instance of LogArchiverTGZ, but this
time with a mask of monthly (gz), which will take the monthly tar file from the previous gadget instance and turn it
into a tar.gz file. More info to follow, however the package has some documentation already.

####RadioBlippers (Simulation)

This Gadget allows you to simulate a number of radioBlip nodes on specific RF Network groups.

####NodeMap (extended core Gadget)
NodeMap replaces the core NodeMap gadget to support the Band/Frequency parameter. The allows you to use:
        { data: "RFb433g5i2,roomNode,boekenkast JC",  to: "nm.Info" }
        { data: "RFb868g5i2,roomNode,keuken",  to: "nm.Info" }

Important: You must also override the Readings gadget if you want your database to process this extended information
correctly, and is an absolute must if you have the same group (e.g. g5) on both 868 and 433 Mhz, otherwise one bands
data will interleave with another.
(Note: I will be submitting a derivative of this to core shortly)

####OnOffMonitor
OnOffMonitor allows you to manage On/Off events within the context of 'time' and 'duration'. It consumes events you
specify from DataSub and generates one or more additional 'related' events. As an example, you can listen for roomNode
'moved' events and generate another event 20min in the future if the state of the endpoint has not changed.
You can then hook into this event with another appropriate Gadget to handle the new event.


###Jeebus focused
-----------------

####SerialPort (extended core Gadget)
The SerialPortEx Gadget is an extended SerialPort. It contains a .Param pin to allow configuration of such things
as Baud, Databits & StopBits. It otherwise operate directly as per the standard SerialPort gadget and is directly
interchangeable. Its default configuration is as per the core SerialPort Gadget.
This Gadget is appropriate if you wish more control over the serial port, such as if you are using an ATTINY based micro
like the JNu,ATTiny85,ATTiny84 or a GPS or bluetooth module.
There are actually two 'variants' of this Gadget. They both do the "same" thing, but are implemented in different ways
If you load the Gadget via the serial/compat package, it will just replace the standard SerialPort implementation. If
you use the serial/serialex package is will not replace the standard SerialPort, but rather provide a SerialPortEx
Gadget within the Flow Registry. If you use serial/extended you will be able to access the native SerialPortEx Gadget
but it will NOT be added to the Flow Registry as any specific Gadget - in which case registration is up to you.

###MQTTServer
If you choose to use an external MQTT broker like RabbitMQ or Mosquitto, use this Gadget to replace the inbuilt
Gadget. By including this Gadget, the core MQTTServer will be replaced, with a quick check to confirm an external broker
is visible on the chosen url/port. It then steps out of the way allowing your app to talk to the external broker.








