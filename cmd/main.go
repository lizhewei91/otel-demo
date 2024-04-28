/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lizw91/otel-demo/pkg/client"
	"github.com/lizw91/otel-demo/pkg/observability"
	"github.com/lizw91/otel-demo/pkg/server"
	"github.com/lizw91/otel-demo/pkg/worker"
)

var (
	// client
	serverAddr string
	interval   int
)

func main() {
	cmd := &cobra.Command{
		Use:   "oteldemo",
		Short: "otel demo",
	}
	clientcmd := &cobra.Command{
		Use:   "client",
		Short: "client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return client.Run(serverAddr, interval)
		},
	}
	clientcmd.Flags().StringVarP(&serverAddr, "server-addr", "s", "http://localhost:8080", "")
	clientcmd.Flags().IntVarP(&interval, "interval", "i", 10, "")

	servercmd := &cobra.Command{
		Use:   "server",
		Short: "api server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := observability.Init("server"); err != nil {
				return err
			}
			if err := server.Init(); err != nil {
				return err
			}
			return server.Run()
		},
	}

	servercmd.Flags().StringVarP(&server.Opts.WorkerAddr, "worker-addr", "w", "http://localhost:8081", "")
	servercmd.Flags().StringVarP(&server.Opts.MysqlAddr, "mysql-addr", "", "localhost:3306", "")
	servercmd.Flags().StringVarP(&server.Opts.MysqlUserName, "mysql-user-name", "", "root", "")
	servercmd.Flags().StringVarP(&server.Opts.MysqlPassword, "mysql-password", "", "Lizw220@1208", "")
	servercmd.Flags().StringVarP(&server.Opts.MysqlDBName, "mysql-db-name", "", "oteldemo", "")

	workercmd := &cobra.Command{
		Use:   "worker",
		Short: "worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := observability.Init("worker"); err != nil {
				return err
			}
			if err := worker.Init(); err != nil {
				return err
			}
			return worker.Run()
		},
	}

	cmd.AddCommand(clientcmd, servercmd, workercmd)
	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
