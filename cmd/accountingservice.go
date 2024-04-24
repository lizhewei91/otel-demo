package cmd

import (
	"context"
	"fmt"
	"github.com/lizw91/otel-demo/pkg/observability"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/lizw91/otel-demo/internal/router"
)

var accountServiceConf = &accountServiceConfig{}

type accountServiceConfig struct {
	serverAddr          string
	serverPort          string
	nextServiceEndpoint string
}

// accountingserviceCmd represents the accountingservice command
var accountingserviceCmd = &cobra.Command{
	Use:   "accountingservice",
	Short: "accounting service",
	RunE:  runAccountingService,
}

const appName = "accounting-service"

func init() {
	rootCmd.AddCommand(accountingserviceCmd)

	accountingserviceCmd.Flags().StringVarP(&accountServiceConf.serverAddr, "server-addr", "s", "0.0.0.0", "Help message for toggle")
	accountingserviceCmd.Flags().StringVarP(&accountServiceConf.serverPort, "server-port", "p", "8080", "Help message for toggle")
	accountingserviceCmd.Flags().StringVarP(&accountServiceConf.nextServiceEndpoint, "next-server-endpoint", "n", "8080", "Help message for toggle")
}

func runAccountingService(cmd *cobra.Command, args []string) error {
	if _, err := observability.InitProvider(appName); err != nil {
		return err
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", accountServiceConf.serverAddr, accountServiceConf.serverPort),
		Handler:      router.Routers(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//start http server
	go func() {
		server.ListenAndServe()
	}()

	log.Printf("server has been started successfully!")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGKILL)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	log.Println("server exiting")
	return nil
}
