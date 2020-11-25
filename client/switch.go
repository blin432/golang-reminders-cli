package client

import
(
	"fmt"
	"os"
	"time"
	"flag"
	"strings"
) 


//custom type for ids (needs to parse)
type idsFlag []string

func (list idsFlag) String() string{
	return strings.Join(list, ",")
}

func (list *idsFlag) Set(v string) error{
	*list = append(*list, v)
	return nil
}

//interface
type BackendHTTPClient interface{
	Create(title, message string, duration time.Duration) ([]byte, error)
	Edit(id string, title, message string, duration time.Duration) ([]byte, error)
	Fetch(ids []string) ([]byte, error)
	Delete(ids []string) error
	Healthy(host string) bool
}

//constructor function
//returns type switch
//HTPclient created in http.go
func NewSwitch(uri string) Switch{
	httpClient := NewHTTPClient(uri)
	s := Switch{client: httpClient, backendAPIURL: uri}
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
		//creating flagset
		createCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		t, m, d := s.reminderFlags(createCmd)
		
		if err := s.checkArgs(3); err != nil {
			return err
		}
		if err := s.parseCmd(createCmd); err != nil {
			return err
		}

		res, err := s.client.Create(*t, *m, *d)
		if err != nil{
			return wrapError("could not create reminder", err)
		}
		fmt.Printf("reminder created successfully:\n%s", string(res))
		fmt.Println("create reminder")
		return nil
	}
}

func (s Switch) edit() func(string) error{
	return func(cmd string) error{

		//for ids
		//custom type uses Var
		ids := idsFlag{}
		//creating flagset
		editCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		editCmd.Var(&ids, "id", "the ID (int) of the reminder to edit")
		t, m, d := s.reminderFlags(editCmd)
		
		if err := s.checkArgs(2); err != nil {
			return err
		}
		if err := s.parseCmd(editCmd); err != nil {
			return err
		}
		lastID := ids[len(ids)-1]
		res, err := s.client.Edit(lastID, *t, *m, *d)
		if err != nil{
			return wrapError("could not edit reminder", err)
		}
		fmt.Printf("reminder edit successfully:\n%s", string(res))
		fmt.Println("create edit")
		return nil
	}
}

func (s Switch) fetch() func(string) error{
	return func(cmd string) error{
		//for ids
		//custom type uses Var
		ids := idsFlag{}
		//creating flagset
		fetchCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		fetchCmd.Var(&ids, "id", "list of reminder IDs (int) to fetch")
		
		if err := s.checkArgs(1); err != nil {
			return err
		}
		if err := s.parseCmd(fetchCmd); err != nil {
			return err
		}

		res, err := s.client.Fetch(ids)
		if err != nil{
			return wrapError("could not fetch reminder", err)
		}
		fmt.Printf("reminder fetch successfully:\n%s", string(res))
		fmt.Println("create fetch")
		return nil
	}
}

func (s Switch) delete() func(string) error{
	return func(cmd string) error{
			//for ids
		//custom type uses Var
		ids := idsFlag{}
		//creating flagset
		deleteCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		deleteCmd.Var(&ids, "id", "list of reminder IDs (int) to delete")
		
		if err := s.checkArgs(1); err != nil {
			return err
		}
		if err := s.parseCmd(deleteCmd); err != nil {
			return err
		}

		err := s.client.Delete(ids)
		if err != nil{
			return wrapError("could not delete reminder", err)
		}
		fmt.Printf("reminder delete successfully:\n%v\n", ids)
		fmt.Println("create delete")
		return nil
	}
}

func (s Switch) health() func(string) error{
	return func(cmd string) error{
		var host string
		healthCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		healthCmd.StringVar(&host, "host", s.backendAPIURL, "host to ping for health")
		if err := s.parseCmd(healthCmd); err != nil{
			return err
		}
		if !s.client.Healthy(host){
			fmt.Printf("host %s is down\n", host)
		}else{
			fmt.Printf("host %s is up and running \n", host)
		}
		return nil
	}
}

//creating helper so we don't have to copy logic
func (s Switch) reminderFlags(f *flag.FlagSet) (*string, *string, *time.Duration) {
	t, m, d := "", "", time.Duration(0)
	f.StringVar(&t, "title", "", "Reminder title")
	f.StringVar(&t, "t", "", "Reminder title")
	f.StringVar(&m, "message", "", "Reminder message")
	f.StringVar(&m, "m", "", "Reminder message")
	f.DurationVar(&d, "duration", 0, "Reminder time")
	f.DurationVar(&d, "d", 0, "Reminder time")
	return &t, &m, &d
}

//helper to parse commands
func (s Switch) parseCmd(cmd *flag.FlagSet) error {
	err := cmd.Parse(os.Args[2:])
	if err != nil{
		return wrapError("could not parse '"+cmd.Name()+"' command flags", err)
	}
	return nil
}


//helper to check min of argument an arg has
func (s Switch) checkArgs(minArgs int) error {
	if len(os.Args) == 3 && os.Args[2] == "--help" {
		return nil
	}
	if len(os.Args)-2 <minArgs {
		fmt.Printf("incorrect use of %s\n%s --help\n",
		os.Args[1], os.Args[0], os.Args[1],)
		return fmt.Errorf(
			"%s expects at least %d args, %d provided", os.Args[1], minArgs, len(os.Args)-2,
		)
	}
	return nil
}
