package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladlosev/k8s-apiserver-metrics/pkg/client"
	"github.com/vladlosev/k8s-apiserver-metrics/pkg/server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logLevel string

// NewServerCommand returns a command that will launch a server re-publiching
// the Kubernetes API server's /metrics endpoint.
func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Launch a server to re-publish API server /metrics endpoint",
		RunE:  startServer,
	}

	cmd.PersistentFlags().StringVar(
		&logLevel,
		"log-level",
		"info",
		"Log level. One of: error, warn, info, degug.",
	)

	return cmd
}

func startServer(cmd *cobra.Command, args []string) error {
	parsedLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(parsedLevel)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	metricsClient, err := client.New()
	if err != nil {
		return err
	}
	metricsScraper := server.New(metricsClient)
	go func() {
		<-stopChan
		metricsScraper.Shutdown(context.Background())
	}()
	logrus.WithField("address", metricsScraper.Addr).Info("Launchiing server")
	err = metricsScraper.ListenAndServe()
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}
