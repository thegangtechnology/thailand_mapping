package cmd

import (
	"go-template/database/duck"

	"github.com/spf13/cobra"
)

func getDuckCommand() cobra.Command {
	return cobra.Command{
		Use:   "duck",
		Short: "Use goose to migrate data",
		Long:  `Connects to the database and migrate tables.`,
		Run: func(cmd *cobra.Command, args []string) {
			duck.UseGoose(args)
		},
	}
}
