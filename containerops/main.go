package main

import (
	"flag"
	"fmt"
	"git.moqi.ai/fingerprint-system/fingerprint-tools/containerops/lib"
	"os"
)

var (
	h bool
	stopall string
	removeall string
	ops string
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: cops -H storage0002 -P 2376
Options:
`)
	flag.PrintDefaults()
	//os.Exit(0)
}

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&lib.Host, "H", "", "docker host")
	flag.StringVar(&lib.Port, "P", "", "docker rest api port")
	flag.StringVar(&stopall,"s","","stop all containers")
	flag.StringVar(&removeall,"r","","remove all containers")
	flag.StringVar(&ops,"ops","","stop or remove a container")
	flag.StringVar(&lib.Name,"name","","container name")
	flag.Usage = usage
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}

	if removeall == "all" {
		lib.RemoveAllContainers()
		return
	}

	if stopall == "all" {
		lib.StopALLContainers()
		return
	}

	if lib.Name != "" {
		if ops == "stop" {
			lib.StopContainer()
		} else if ops == "remove" {
			lib.RemoveContainer()
		}
	} else {
		lib.ListContainer()
	}
}

