package main

import(
	"flag"
	"fmt"
	"strings"
	"time"
)
//custom type of id, id flag
// parse different id's from the command line and then convert it to a string

//basically an interface
//custom type flag
type idsFlag []string
 
func(ids idsFlag) String() string{
	return strings.Join(ids, ",")
}

//this id is passed from command line
// * makes it a pointer receiver
func(ids *idsFlag) Set(id string) error{
	*ids = append(*ids, id)
	return nil
}

//creating struct of type person
type person struct{
	name string
	born time.Time
}

//creating interface
//use Springf when formatting an interface
func (p person) String() string{
	return fmt.Sprintf("my name is: %s, and I am %s", p.name, p.born.String())
}

func(p *person) Set(name string) error{
	p.name = name
	p.born = time.Now()
	return nil
}

func main(){
	var ids idsFlag
	var p person
	//binding variable to a specific flag
	flag.Var(&ids, "id", "the id to be appened to the List")
	flag.Var(&p, "name", "the name of the person")
	flag.Parse()
	fmt.Println(ids)
	fmt.Println(p)
}