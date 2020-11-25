package client

import "fmt"

//helper function to wrap error messages
func wrapError(customMsg string, originarErr error) error{
	return fmt.Errorf("&s :%v, customMsg, originalErr")
}