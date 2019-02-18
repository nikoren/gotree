package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)


var root string
func handleIfErr(msg string, err error) {
	if err != nil {
		log.Fatalf("ERROR: %s %s", msg, err)
	}
}

type Order int

func(o Order)String() string{
	switch o {
	case First:
		return "First"
	case Middle:
		return "Middle"
	case Last:
		return "Last"
	default:
		return "Unknown"
	}
}

const (
	First Order = iota
	Middle
	Last
	Unknown
)

type Format string

const (
	firstOrMidInCurrentFolderFormatter Format = "├───"
	lastInCurrentFolderFormatter       Format = "└───"
	lastInParentFolderFormatter        Format = "	"
	firstOrMidInParrentFolderFormatter Format = "│	"
	empty                              Format = ""
)

func childOrderToFormat(o Order) Format {
		switch o {
		case First:
			return firstOrMidInCurrentFolderFormatter
		case  Middle:
			return firstOrMidInCurrentFolderFormatter
		case Last:
			return lastInCurrentFolderFormatter
		default:
			return empty
	}
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	root = path
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	handleIfErr("Failure to walk the tree: ", err)
}

func dirTree(out io.Writer, path string, printFilesFlag bool) error {
	err := filepath.Walk(path, decoratedWF)
	handleIfErr("Couldn't walk the path", err)
	return nil
}

// empty WalkFunction to print the values as is for debug purposes
func debugWF(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	return nil
}

func decoratedWF(path string, info os.FileInfo, err error) error{
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	if strings.HasPrefix(path,"."){
		return nil
	}
	fmt.Printf("%s", decoratePath(path))
	return nil
}

func decoratePath(up string) string {
	o:= getOrder(up)
	c:= countParents(up)
	if c == 0{
		return zeroLevelFormat(up, o)

	}else{
		return parentFormat(up) + zeroLevelFormat(up, o)
	}
}

func getOrder(p string) Order {
	parent := filepath.Dir(p)
	fs, err := ioutil.ReadDir(parent)
	handleIfErr(fmt.Sprintf("Couldn't list dir %s \n", p), err)

	for i, f := range fs {
		if f.Name() == filepath.Base(p) {
			switch i {
			case 0:
				if len(fs) > 1 {
					return First
				}
				return Last
			case len(fs) - 1:
				return Last
			default:
				return Middle
			}
		}
	}
	return Unknown
}

func parentFormat(p string) string {
	//"testdata/static/css/body.css	"
	//" 	│	 │	  └───body.css (28b)"
	if p == root{
		return ""
	}

	leaf := filepath.Dir(p)
	parent := filepath.Dir(leaf)

	fs, err := ioutil.ReadDir(parent)
	handleIfErr(fmt.Sprintf("Couldn't readDir %s", parent), err)
	var resultFormatter Format
	lok:
	for i, f := range fs {
		if f.Name() == filepath.Base(leaf) {
			switch i {
			case 0:
				if len(fs) > 1 {
					resultFormatter = firstOrMidInParrentFolderFormatter
					break lok
				}
				resultFormatter =  lastInParentFolderFormatter
				break lok

			case len(fs) - 1:
				resultFormatter = lastInParentFolderFormatter
				break lok

			default:
				resultFormatter = firstOrMidInParrentFolderFormatter
				break lok
			}
		}
	}

	return fmt.Sprintf("%s%s", parentFormat(leaf), resultFormatter)
}

func zeroLevelFormat(up string, o Order) string {
	formatter :=  childOrderToFormat(o)
	return fmt.Sprintf("%s%s\n",formatter, filepath.Base(up))
}

func countParents(path string) int {
	if path == "."{
		return 0
	}
	numOfParents := strings.Count(path, string(os.PathSeparator))
	return numOfParents
}

