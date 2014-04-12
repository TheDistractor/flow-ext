//The logging package provides a basic hierarchical logging fascility.
//It is based upon the premise that a log file uses the pattern YYYY<MM><DD>.ext
//Hierarchical logging is driven by the patterns observed from the previously
//seen log file.Hence a change in pattern MUST be seen to trigger an event.
//A General hierarchy would be:
// logger->path/20140430.txt->add to path/201404.tar
// logger->path/20140501.txt->add to path/201405.tar
// emit path/201404.tar to next hierarchical stage (as a 20140501 (next month) was seen)
//
package logging

import (
	"testing"
	"os"
	"path"
	"fmt"
	"github.com/jcw/flow"
	"math/rand"
	_"time"
)


var InitOK = true


//lets setup some source files
func init() {

	if err:= os.RemoveAll("./log");err!=nil {
		InitOK = false
		return
	}

	if err:= os.MkdirAll("./log/2014",os.ModeDir | os.ModePerm );err != nil {
		InitOK = false
		return
	}
	if err:= os.MkdirAll("./log/2015",os.ModeDir | os.ModePerm );err != nil {
		InitOK = false
		return
	}

	files := []string{"20140331.txt","20140401.txt","20140430.txt","20140501.txt", "2014bad0331.txt", "20140731.txt","20140801.txt", "20150101.txt",  }

	for _,file := range files {
		fd,err := os.Create( path.Join("log",file[:4],file))
		if err == nil {

			size := 100 + rand.Intn(1024-100)

			//fmt.Println("Size:", size)
			bytes := make([]byte, size)
			for i:=0 ; i<size ; i++ {
				bytes[i] = byte(65 + rand.Intn(90-65))
			}

			fd.Write(bytes)

			fd.Close()
		} else	{
			fmt.Println("Oops:",err)
			InitOK = false
		}
	}



}



func TestInit(t *testing.T) {
	if ! InitOK {
		t.Fatal("Initialization failed....other tests may fail")
	}
}


//our mask contains a Daily mask, with .gz, so we expect gzip output files for each input.
func ExampleTGZByDayGZ() {
	g := flow.NewCircuit()
	g.Add("f", "LogArchiverTGZ")
	g.Feed("f.Param", flow.Tag{"-m", "20060102.gz"}	)
	g.Feed("f.Param", flow.Tag{"-v", true}	)
	g.Feed("f.In", "log/2014/20140331.txt")
	g.Feed("f.In", "log/2014/20140401.txt")
	g.Feed("f.In", "log/2014/20140430.txt")
	g.Feed("f.In", "log/2014/20140501.txt")
	g.Run()
	// Output:
	// Lost string: log/2014/20140331.txt.gz
	// Lost string: log/2014/20140401.txt.gz
	// Lost string: log/2014/20140430.txt.gz
	// Lost string: log/2014/20140501.txt.gz

}


// we turn on verbose as we want to see add: messages in addition to a complete tarball message
// this tests the creation of a tarball from a single input file using a Daily mask
// NOTE: we will not see the emission of the last file as tarball because tarball emission rely on the 'next' match
// to cause a tarball cycle, however the base tarball was created (see the last add:)
func ExampleTGZByDayTar() {
	g := flow.NewCircuit()
	g.Add("f", "LogArchiverTGZ")
	g.Feed("f.Param", flow.Tag{"-m", "20060102.tar"}	)
	g.Feed("f.Param", flow.Tag{"-v", true}	)
	g.Feed("f.In", "log/2014/20140331.txt")
	g.Feed("f.In", "log/2014/20140401.txt")
	g.Feed("f.In", "log/2014/20140430.txt")
	g.Feed("f.In", "log/2014/20140501.txt")
	g.Run()
	// Output:
	// Lost string: add:log/2014/20140331.txt to:log/2014/20140331.tar
	// Lost string: add:log/2014/20140401.txt to:log/2014/20140401.tar
	// Lost string: log/2014/20140331.tar
	// Lost string: add:log/2014/20140430.txt to:log/2014/20140430.tar
	// Lost string: log/2014/20140401.tar
	// Lost string: add:log/2014/20140501.txt to:log/2014/20140501.tar
	// Lost string: log/2014/20140430.tar
}


//here we get error as the 1st input file will not match the date parse mask, the next is not present
func ExampleTGZByDayFilenameBadMasks() {
	g := flow.NewCircuit()
	g.Add("f", "LogArchiverTGZ")
	g.Feed("f.Param", flow.Tag{"-m", "20060102.gz"}	)
	g.Feed("f.Param", flow.Tag{"-v", true}	)
	g.Feed("f.In", "log/2014/2014bad0331.txt")
	g.Feed("f.In", "log/2014/2014notexist0331.txt")
	g.Run()
	// Output:
	// Lost string: ignore:log/2014/2014bad0331.txt err:parsing time "2014bad0": month out of range
	// Lost string: err:stat log/2014/2014notexist0331.txt: no such file or directory
}

//we turn verbose on, as we want to see 'add:' messages in addition to tarball emit messages
//this tests the creation of monthly tarballs
//we will not see a 201405 tarball emission until the next input file triggers it by a failed match.
func ExampleTGZByMonthTar() {
	g := flow.NewCircuit()
	g.Add("f", "LogArchiverTGZ")
	g.Feed("f.Param", flow.Tag{"-m", "200601.tar"}	)
	g.Feed("f.Param", flow.Tag{"-v", true}	)
	g.Feed("f.In", "log/2014/20140331.txt")
	g.Feed("f.In", "log/2014/20140401.txt")
	g.Feed("f.In", "log/2014/20140430.txt")
	g.Feed("f.In", "log/2014/20140501.txt")
	g.Run()
	// Output:
	// Lost string: add:log/2014/20140331.txt to:log/2014/201403.tar
	// Lost string: add:log/2014/20140401.txt to:log/2014/201404.tar
	// Lost string: log/2014/201403.tar
	// Lost string: add:log/2014/20140430.txt to:log/2014/201404.tar
	// Lost string: add:log/2014/20140501.txt to:log/2014/201405.tar
	// Lost string: log/2014/201404.tar

}

//here we see tarballs accumulated by year, we see a 2014 tarball emission because the process sees a 2015 file
//and determines the 2014 tarball to be complete
func ExampleTGZByYearTar() {
	g := flow.NewCircuit()
	g.Add("f", "LogArchiverTGZ")
	g.Feed("f.Param",  flow.Tag{"-m", "2006.tar"}	)
	g.Feed("f.Param", flow.Tag{"-v", true}	)
	g.Feed("f.In", "log/2014/20140331.txt")
	g.Feed("f.In", "log/2014/20140401.txt")
	g.Feed("f.In", "log/2014/20140430.txt")
	g.Feed("f.In", "log/2014/20140501.txt")
	g.Feed("f.In", "log/2015/20150101.txt")
	g.Run()
	// Output:
	// Lost string: add:log/2014/20140331.txt to:log/2014/2014.tar
	// Lost string: add:log/2014/20140401.txt to:log/2014/2014.tar
	// Lost string: add:log/2014/20140430.txt to:log/2014/2014.tar
	// Lost string: add:log/2014/20140501.txt to:log/2014/2014.tar
	// Lost string: add:log/2015/20150101.txt to:log/2015/2015.tar
	// Lost string: log/2014/2014.tar


}

func ExampleTGZByYearGZ() {
	g := flow.NewCircuit()
	g.Add("f", "LogArchiverTGZ")
	g.Feed("f.Param",  flow.Tag{"-m", "2006.gz"}	)
	g.Feed("f.Param", flow.Tag{"-v", true}	)
	g.Feed("f.In", "log/2014/2014.tar")
	g.Run()
	// Output:
	// Lost string: log/2014/2014.tar.gz


}


//here we create a monthly tar when the cycle changes from 20140731 to 20140801
//when the tar is emitted, we then create a .gzipped version
func ExampleTGZMonthlyTarGZ() {
	g := flow.NewCircuit()
	g.Add("f1", "LogArchiverTGZ")
	g.Add("f2", "LogArchiverTGZ")
	g.Feed("f1.Param",  flow.Tag{"-m", "200601.tar"}	)
	g.Feed("f1.Param", flow.Tag{"-v", true}	)
	g.Feed("f2.Param",  flow.Tag{"-m", "200601.gz"}	)
	g.Feed("f2.Param", flow.Tag{"-v", true}	)
	g.Connect("f1.Out", "f2.In", 0)

	g.Feed("f1.In", "log/2014/20140731.txt")
	g.Feed("f1.In", "log/2014/20140801.txt")
	g.Run()
	// Output:
	//Lost string: add:log/2014/20140731.txt to:log/2014/201407.tar
	//Lost string: add:log/2014/20140801.txt to:log/2014/201408.tar
	//Lost string: log/2014/201407.tar.gz


}




