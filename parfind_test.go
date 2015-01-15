package main

/*

	parfind
	=======
	A parallel, simplified version of find(1).

	See README.md for usage and licensing information.

*/

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"testing"
)

const (
	DefVersion = false
	DefWorkers = DEFAULT_WORKERS
	DefRoot    = "."
	DefPrint0  = false

	NonExistingDirectory   = "./non_existing_directory"
	TestingDirectory       = "./test_dir"
	TestingSubDirectory    = TestingDirectory + "/subdir"
	TestingSubSubDirectory = TestingSubDirectory + "/subdir"
	TestingFile            = TestingDirectory + "/test_file.txt"
	TestingFileUnicode     = TestingDirectory + "/変.txt"
	TestingFileQuoting     = TestingDirectory + "/\r\n\"'`.txt"
	TestingFileSymlink     = TestingDirectory + "/symlink"
	TestingFileSubDir      = TestingSubDirectory + "/test_file.txt"
	TestingFileSubSubDir   = TestingSubSubDirectory + "/test_file.txt"
)

var TestDirs = []string{
	TestingDirectory,
	TestingSubDirectory,
	TestingSubSubDirectory,
}

var TestFiles = []string{
	TestingFile,
	TestingFileUnicode,
	TestingFileQuoting,
	TestingFileSubDir,
	TestingFileSubSubDir,
}

var TestParams = []struct {
	version bool
	workers int
	root    string
	print0  bool
}{
	{DefVersion, DefWorkers, TestingDirectory, false},
	{DefVersion, DefWorkers, TestingDirectory, true},
	{DefVersion, -1, TestingDirectory, false},
	{DefVersion, -1, TestingDirectory, true},
	{DefVersion, MAX_WORKERS + 1, TestingDirectory, false},
	{DefVersion, MAX_WORKERS + 1, TestingDirectory, true},
	{DefVersion, DefWorkers, NonExistingDirectory, false},
	{DefVersion, DefWorkers, NonExistingDirectory, true},
}

var (
	FmtDefault      = regexp.MustCompile("[dflspCDu] [-0-9]+ [0-9]+ (\"[\x20-\x7E]+\")\r?\n")
	FmtPrint0       = regexp.MustCompile("[dflspCDu]\x00[-0-9]+\x00[0-9]+\x00([^\x00]+)\x00")
	RequiredResults = 0
)

func createTestingDirectory() {
	RequiredResults = 0
	removeTestingDirectory()
	for _, d := range TestDirs {
		if err := os.Mkdir(d, 0600); err == nil {
			RequiredResults++
		}
	}
	for _, f := range TestFiles {
		if _, err := os.Create(f); err == nil {
			RequiredResults++
		}
	}
	if err := os.Symlink(TestingFile, TestingFileSymlink); err == nil {
		RequiredResults++
	}
}

func removeTestingDirectory() {
	os.RemoveAll(TestingDirectory)
	os.RemoveAll(TestingDirectory)
	os.RemoveAll(TestingDirectory)
}

func runParfind(version bool, workers int, root string, print0 bool, t *testing.T) (sout string, serr string) {
	out := new(bytes.Buffer)
	err := new(bytes.Buffer)
	defer func() {
		recover()
		if out.Len() > 0 {
			t.Logf("STDOUT\n%s\n", out)
		}
		if err.Len() > 0 {
			t.Logf("STDERR\n%s\n", err)
		}
		sout = out.String()
		serr = err.String()
	}()
	parfind(version, workers, root, print0, out, err)
	return
}

func TestVersion(t *testing.T) {
	out, err := runParfind(true, DefWorkers, DefRoot, DefPrint0, t)
	if out != VERSION+"\n" || err != "" {
		t.Fail()
	}
}

func TestParfind(t *testing.T) {
	createTestingDirectory()
	for _, params := range TestParams {
		out, err := runParfind(params.version, params.workers, params.root, params.print0, t)
		if params.root == TestingDirectory && err != "" {
			t.Fatalf("Unexpected output on stderr")
		}
		f := FmtDefault
		if params.print0 {
			f = FmtPrint0
		}
		res := f.FindAllStringSubmatch(out, -1)
		if params.root == TestingDirectory && len(res) != RequiredResults {
			t.Fatalf("Expected %d results, got %d", RequiredResults, len(res))
		}
		for _, v := range res {
			f := v[1]
			if params.print0 == false {
				f, _ = strconv.Unquote(f)
			}
			if _, e := os.Stat(f); e != nil {
				t.Errorf("Can't stat result %s", f)
			}
		}
	}
}
