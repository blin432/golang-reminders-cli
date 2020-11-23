package main

import(
	"fmt"
	"os"
	"strings"
)


//example before flagset to show how the flag is parsed
func main(){
	//checking if there are two arguements
	if len(os.Args) < 2 {
		fmt.Println("no command provided")
		os.Exit(2)
	}
	//command switch
	cmd := os.Args[1]
	switch cmd {
	case "greet":
		//default msg
		msg := "reminders cli - cli basics"
		if len(os.Args) > 2{
			//parsing flag message
			f := strings.Split(os.Args[2], "=")
			if len(f) == 2 && f[0] == "--msg" {
			msg = f[1]
			}
		}
		fmt.Printf("hello and welcome: %s\n", msg)
	case "help":
		fmt.Println("some help message")
	default: 
		fmt.Printf("unknown command: %s\n", cmd)
	}
}