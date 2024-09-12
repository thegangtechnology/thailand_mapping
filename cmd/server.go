package cmd

import (
	"fmt"
	"go-template/config"
	"go-template/internal"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func getServerCmd() cobra.Command {
	return cobra.Command{
		Use:   "server",
		Short: "Starts the server",
		Long:  `Connects to the database and starts the server.`,
		Run: func(cmd *cobra.Command, args []string) {
			const properDelay = 10
			time.Sleep(properDelay * time.Second)

			server := internal.CreateServer()

			port, err := strconv.Atoi(os.Getenv(config.PortAPI))
			if err != nil {
				port = 32423
			}

			server.Logger.Fatal(server.Start(fmt.Sprintf("0.0.0.0:%d", port)))
		},
	}
}
