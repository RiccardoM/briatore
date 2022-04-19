package start

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/riccardom/briatore/apis"
	"github.com/riccardom/briatore/types"
)

const (
	flagPort = "port"
)

// GetStartCmd returns the command used to start the APIs
func GetStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the APIs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(os.Stdout)

			// Get the configuration
			cfg, err := types.ReadConfig(cmd)
			if err != nil {
				return err
			}

			// Run the Gin server
			r := gin.Default()
			r.Use(gin.Recovery())

			// Setup CORS
			ginCfg := cors.DefaultConfig()
			ginCfg.AllowAllOrigins = true
			r.Use(cors.New(ginCfg))

			// Register the endpoints
			r.GET("/report", apis.GetReportHandler(cfg))

			port, _ := cmd.Flags().GetUint(flagPort)
			return r.Run(fmt.Sprintf(":%d", port))
		},
	}

	cmd.Flags().Uint(flagPort, 8080, "Port to which the APIs will bind")

	return cmd
}
