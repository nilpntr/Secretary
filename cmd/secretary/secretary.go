package main

import (
	"github.com/nilpntr/secretary/pkg/action"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func main() {
	actionConfig := new(action.Configuration)

	cmd, err := newRootCmd(actionConfig, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(func() {
		if err := actionConfig.Init(); err != nil {
			log.Fatal(err)
		}
	})

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}