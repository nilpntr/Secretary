package main

import (
	"github.com/nilpntr/secretary/pkg/action"
	"github.com/nilpntr/secretary/pkg/cli"
	"github.com/spf13/cobra"
)

var settings = cli.New()

func newRootCmd(actionConfig *action.Configuration, args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "secretary",
		Short:        "Secretary is a tool to sync applicable secrets to other namespaces in k8s",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return action.NewRun(actionConfig, settings)
		},
	}
	flags := cmd.PersistentFlags()
	settings.AddFlags(flags)

	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.Parse(args)

	return cmd, nil
}
