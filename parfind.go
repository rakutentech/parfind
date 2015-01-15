package main

/*

  parfind
  =======
  A parallel, simplified version of find(1).

  See README.md for usage and licensing information.

*/

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	VERSION         = "parfind-0.2"
	MAX_WORKERS     = 128
	DEFAULT_WORKERS = 16
)

// CLI parameters
var root = flag.String("root", ".", "The directory to start scanning from")
var workers = flag.Int("workers", DEFAULT_WORKERS, "How many workers to use")
var version = flag.Bool("version", false, "Show version information")
var print0 = flag.Bool("print0", false, "Use NUL as field and record separator (for use with xargs -0)")

type FileDesc struct {
	path  string
	mode  os.FileMode
	size  int64
	mtime time.Time
}

type Parfind struct {
	// workers -> output channel
	result_chan chan FileDesc

	// semaphore channel to limit number of running workers
	worker_chan chan struct{}

	// unit of work counter (how many UOWs are currently in flight)
	uows int64

	// output -> main channel to signal completion of the output loop
	output_complete_chan chan struct{}

	// where to write the output
	out io.Writer

	print0 bool
}

// recursively walk the contents of a directory (dir)
func (p *Parfind) find(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}
	for _, f := range files {
		fp, err := filepath.Abs(path.Join(dir, f.Name()))
		if err != nil {
			log.Println(err)
			continue
		}
		if f.IsDir() && f.Mode()&os.ModeSymlink == 0 {
			p.enqueue_find(fp)
		}
		p.result_chan <- getFileDesc(fp, f)
	}
}

// enqueue a directory to be scanned
// NOTE: this function must be called synchronously to correctly account for the
// number of inflight UOWs
func (p *Parfind) enqueue_find(dir string) {
	// synchronously signal that we have a new unit of work to process
	atomic.AddInt64(&p.uows, int64(1))
	go func() {
		// wait until there are less than *workers workers running
		<-p.worker_chan
		p.find(dir)
		p.worker_chan <- struct{}{} // signal this worker has ended
		if atomic.AddInt64(&p.uows, int64(-1)) == 0 {
			close(p.result_chan) // signal that we have finished the find job
		}
	}()
}

func modeToType(mode os.FileMode) byte {
	switch {
	case mode.IsDir():
		return 'd'
	case mode.IsRegular():
		return 'f'
	case mode&os.ModeSymlink != 0:
		return 'l'
	case mode&os.ModeSocket != 0:
		return 's'
	case mode&os.ModeNamedPipe != 0:
		return 'p'
	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice != 0:
		return 'C'
	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice == 0:
		return 'D'
	default:
		return 'u'
	}
}

func getFileDesc(fp string, fi os.FileInfo) FileDesc {
	return FileDesc{
		path:  fp,
		mode:  fi.Mode(),
		size:  fi.Size(),
		mtime: fi.ModTime(),
	}
}

func (p *Parfind) output() {
	F := "%c %d %d %+q\n"
	if p.print0 {
		F = "%c\x00%d\x00%d\x00%s\x00"
	}
	for file := range p.result_chan {
		fmt.Fprintf(p.out, F, modeToType(file.mode), file.mtime.Unix(), file.size, file.path)
	}
	close(p.output_complete_chan)
}

func parfind(version bool, workers int, root string, print0 bool, o io.Writer, e io.Writer) {
	log.SetOutput(e)

	if version {
		fmt.Fprintln(o, VERSION)
		return
	}

	if workers > MAX_WORKERS {
		workers = MAX_WORKERS
	} else if workers < 1 {
		workers = DEFAULT_WORKERS
	}

	p := &Parfind{
		result_chan:          make(chan FileDesc, workers),
		worker_chan:          make(chan struct{}, workers),
		output_complete_chan: make(chan struct{}),
		out:                  o,
		print0:               print0,
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < workers; i++ {
		p.worker_chan <- struct{}{}
	}

	path, err := filepath.Abs(root)
	if err != nil {
		log.Println(err)
		return
	}

	fi, err := os.Stat(path)
	if err != nil {
		log.Println(err)
		return
	}

	p.result_chan <- getFileDesc(path, fi)
	p.enqueue_find(path)
	go p.output()

	<-p.output_complete_chan // wait for all workers to terminate
}

func main() {
	flag.Parse()
	parfind(*version, *workers, *root, *print0, os.Stdout, os.Stderr)
}
