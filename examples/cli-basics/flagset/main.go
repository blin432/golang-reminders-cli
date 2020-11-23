package main

import(
	"fmt"
	"os"
	"flag"
	"log"
)

//flag set example to simplify parsing of flag
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
		//creating new flagset
		greetCmd := flag.NewFlagSet("greet", flag.ExitOnError)
		msgFlag := greetCmd.String("msg", "CLI basics - reminders cli", "the message for greet Command")
		err := greetCmd.Parse(os.Args[2:])
		if err != nil{
			log.Fatal(err.Error())
		}
		//have to use * because of pointer, need to get value
		fmt.Printf("hello and welcome: %s\n", *msgFlag)
	case "help":
		fmt.Println("some help message")
	default: 
		fmt.Printf("unknown command: %s\n", cmd)
	}
}