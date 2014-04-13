//Package logging implements various data logging gadgets that enhance the core housemon logger.
package logging

import (
	"archive/tar"
	_ "archive/zip"
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/golang/glog"
	"github.com/jcw/flow"
)

//Automatically add to flow registry
func init() {
	flow.Registry["LogArchiverTGZ"] = func() flow.Circuitry { return &LogArchiverTGZ{} }
}

//LogArchiverTGZ supports either [tar]ing and/or [gzip]ing files using an input mask to determine output characteristics
type LogArchiverTGZ struct {
	flow.Gadget
	Param  flow.Input
	In     flow.Input  //the filename we see for input
	Out    flow.Output //what we generate
	Info   flow.Output //output some information about whats happening
	Reject flow.Output //TODO:move key fails to here
}

const (
	sumToYear  = 4
	sumToMonth = 6
	sumToDay   = 8
)

//Run is the main flow gadget entry point.
func (w *LogArchiverTGZ) Run() {

	//TODO:implement 'rm'
	rm := false      //remove source after operation
	verbose := false //emit some data on .Info pin

	mask := "20060102" //The default mask unless overridden - just .gz input files
	for t := range w.Param {
		switch m := t.(type) {
		case flow.Tag:
			switch m.Tag {
			case "-m":
				mask = m.Msg.(string)
			case "-v":
				verbose = m.Msg.(bool)
			case "-d":
				rm = m.Msg.(bool)
			}
		}

	}

	maskbase := strings.ToLower(path.Ext(mask))

	mask = mask[:len(path.Base(mask))-len(maskbase)]

	if maskbase == "" {
		maskbase = "tar" //default action to tar
	}

	var prevFile string
	var prevDate time.Time
	var curFile string

	sumTo := len(mask)

	for m := range w.In {
		curFile = m.(string)

		base := path.Base(curFile)
		ext := path.Ext(curFile)
		dir := path.Dir(curFile)

		fi, err := os.Stat(curFile)
		if err != nil {
			w.Out.Send(fmt.Sprintf("err:%s", err))
			continue
		}

		part := base[:sumTo]

		curDate, err := time.Parse(mask, part)
		if err != nil {
			w.Reject.Send(fmt.Sprintf("ignore:%s err:%s", curFile, err))
			continue //we ignore it if not matching input mask
		}

		if prevFile == "" { //one-time init loop
			switch sumTo {
			case sumToDay:
				prevDate = curDate.AddDate(0, 0, -1)
				prevFile = path.Join(dir, prevDate.Format(mask)+ext)
			case sumToMonth:
				prevDate = curDate.AddDate(0, -1, 0)
				prevFile = path.Join(dir, prevDate.Format(mask)+ext)
			case sumToYear:
				prevDate = curDate.AddDate(-1, 0, 0)
				prevFile = path.Join(dir, prevDate.Format(mask)+ext)

			}
		}

		//if its a gzip mask, we gzip and re-emit
		if maskbase == ".gz" {
			gzFile := curFile + ".gz"

			fdin, err := os.Open(curFile)

			if err != nil {
				w.Info.Send(fmt.Sprintf("err:%s", err))
				continue
			}
			bufin := bufio.NewReader(fdin)

			mode := os.O_RDWR | os.O_APPEND | os.O_CREATE
			fdout, err := os.OpenFile(gzFile, mode, os.ModePerm)
			if err != nil {
				w.Info.Send(fmt.Sprintf("err:%s", err))
				continue
			}

			f := gzip.NewWriter(fdout)
			_, err = bufin.WriteTo(f)
			if err != nil {
				w.Info.Send(fmt.Sprintf("err:%s", err))
				continue
			}

			f.Flush()
			f.Close()

			if rm {
				//del source
				if verbose {
					w.Info.Send("rm:" + curFile)
				}
			}

			w.Out.Send(gzFile)

		} else { //we are building up a tar archive

			tarFile := path.Join(dir, curDate.Format(mask)+".tar")
			mode := os.O_RDWR | os.O_CREATE
			fdout, err := os.OpenFile(tarFile, mode, os.ModePerm)
			if err != nil {
				w.Info.Send(fmt.Sprintf("err:%s", err))
				continue
			}

			defer fdout.Close()

			//am I a new Tar file
			//if so, we must 'hack' the standard NewWriter as it does not recognise incremental additions.
			//so we help it by skipping to the <tail> of the tar file
			x, err := fdout.Stat()
			if x.Size() != 0 {
				//pos, err := fdout.Seek(0,os.SEEK_CUR)
				if _, err = fdout.Seek(-2<<9, os.SEEK_END); err != nil {
					w.Info.Send(fmt.Sprintf("err:%s", err))
					continue
				}
				//pos, err = fdout.Seek(0,os.SEEK_CUR)
			}

			cab := tar.NewWriter(fdout)

			fdin, err := os.Open(curFile)
			if err != nil {
				w.Info.Send(fmt.Sprintf("err:%s", err))
				continue
			}
			defer fdin.Close()

			bufin := bufio.NewReader(fdin)

			hdr, err := tar.FileInfoHeader(fi, curFile)
			if err != nil {
				w.Info.Send(fmt.Sprintf("err:%s", err))
				continue
			}

			if err := cab.WriteHeader(hdr); err != nil {
				//log.Fatalln(err)
				//list.Items[source] = err
				fmt.Println("Adding:", curFile, tarFile)
				w.Info.Send(fmt.Sprintf("err:%s", err))
				continue
			}

			_, err = bufin.WriteTo(cab)

			if err != nil {
				w.Info.Send(fmt.Sprintf("err:%s", err))
				continue
			}

			//if _, err := cab.(*tar.Writer).Write( buf.Bytes() ); err != nil {
			//	//list.Items[source] = err
			//	continue
			//}
			if err := cab.Close(); err != nil {
				w.Info.Send(fmt.Sprintf("err:%s", err))
				continue
			}

			if verbose {
				info := fmt.Sprintf("add:%s to:%s", curFile, tarFile)

				if rm {
					//del source
					if verbose {
						w.Info.Send("rm:" + curFile)
					}
				}

				w.Info.Send(info)
			}
			if curDate != prevDate {
				prevdir := path.Dir(prevFile)

				prevTarFile := path.Join(prevdir, prevDate.Format(mask)+".tar")
				if _, err := os.Stat(prevTarFile); err == nil {
					w.Out.Send(prevTarFile)
				} else {
					//fmt.Println("missing:", prevTarFile)
				}

			}
		}

		//before we get next message, we set prevFile
		prevFile = curFile
		prevDate = curDate
	}
}
