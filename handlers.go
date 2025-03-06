package main

import (
	"errors"
	"fmt"
)

func HandlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("no username given")
	}
	
	err := s.cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	
	fmt.Printf("Username has been set to %v.\n", cmd.Arguments[0])
	return nil
}