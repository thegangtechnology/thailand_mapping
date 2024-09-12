package cmd

import (
	"go-template/config"
	"go-template/database"
	"go-template/utils/fancylogger"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func getRoodCmd() cobra.Command {
	return cobra.Command{
		Use:   "GO-web-server",
		Short: "The root command for GO-web-server",
		Long:  `The root command for GO-web-server`,
	}
}

func loadEnv() {
	env := os.Getenv(config.AppEnvironment)
	if env == "" {
		env = "development"
	}

	gcp := env == config.Production

	fancylogger.SetupLogger(nil, gcp)

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Infof("error loading %s file", ".env."+env)
	}

	database.Connect(database.ConnectParams{})
}

func Execute() error {
	cobra.OnInitialize(loadEnv)

	rootCmd := getRoodCmd()
	serverCmd := getServerCmd()
	migrateCmd := getDuckCommand()

	rootCmd.AddCommand(&serverCmd)
	rootCmd.AddCommand(&migrateCmd)

	return rootCmd.Execute()
}
