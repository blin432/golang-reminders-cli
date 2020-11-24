package client

import
(
	"fmt"
	"os"
) 

//interface
type BackendHTTPClient interface{

}
//constructor function
//returns type switch
//HTPclient created in http.go
func NewSwitch(uri string) Switch{
	httpClient := NewHTTPClient(uri)
	s := Switch{
		client: httpClient,
		backendAPIURL: uri,
	}
	//returns a function which recives a string and returns an error
	//s."somethings" referencing methods on a specific type on the switch type just created
	//not usual functions because they will have a reciever so they will be methods
	s.commands= map[string]func() func(string) error{
		"create": s.create,
		"edit": s.edit,
		"fetch": s.fetch,
		"delete": s.delete,
		"health": s.health,
	}
	return s
}

//every command is going to be a function that returns another runction
type Switch struct{
	client BackendHTTPClient
	backendAPIURL string
	commands map[string] func() func(string) error 
}

//function to parse command
func(s Switch) Switch() error{
	cmdName := os.Args[1]
	//look at the map of commands and make sure it exits
	cmd, ok := s.commands[cmdName]
	if !ok{
		return fmt.Errorf("invalid command `%s`\n", cmdName)
	}
	return cmd()(cmdName)
}

func(s Switch) Help() {
	var help string
	for name := range s.commands{
		help += name + "\t --help\n"
	}
	fmt.Printf("usage of: %s:\n <command> [<args>]\n%s", os.Args[0], help)
}

//not a usual function , a function that returns another function
func (s Switch) create() func(string) error{
	return func(cmd string) error{
		fmt.Println("create reminder")
		return nil
	}
}

func (s Switch) edit() func(string) error{
	return func(cmd string) error{
		fmt.Println("edit reminder")
		return nil
	}
}

func (s Switch) fetch() func(string) error{
	return func(cmd string) error{
		fmt.Println("fetch reminder")
		return nil
	}
}

func (s Switch) delete() func(string) error{
	return func(cmd string) error{
		fmt.Println("delete reminder")
		return nil
	}
}

func (s Switch) health() func(string) error{
	return func(cmd string) error{
		fmt.Println("health reminder")
		return nil
	}
}
