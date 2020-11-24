package main

//declaring variables which we will be commands we will parse

import(
	"flag"
	"os"
	"fmt"

	"github.com/benjaminlin/reminder-cli/client"
)

var (
	backendURIFlag = flag.String("backend", "http://localhost:8080", "backend API url")
	helpFlag = flag.Bool("help", false, "display a helpful message")
)

func main() {
	flag.Parse()
	s := client.NewSwitch(*backendURIFlag)
	if *helpFlag || len(os.Args) == 1 {
		s.Help()
		return
	}

	err := s.Switch()
		if err != nil {
			fmt.Printf("cmd switch error:%s", err)
			os.Exit(2)
		}
}