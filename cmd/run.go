/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"telecart/internal"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := listen(); err != nil {
			log.Println("err", err)
			return errors.Wrap(err, "listen")
		}

		log.Println("finish...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func listen() error {
	sqlDNS := viper.GetString("DNS")
	topic := viper.GetString("MSQTT_TOPIC")
	if sqlDNS == "" || topic == "" {
		return errors.New("pleas specify DNS and MSQTT_TOPIC")
	}

	log.Println(sqlDNS, topic)
	svc, err := internal.NewService(sqlDNS, topic)
	if err != nil {
		return errors.Wrap(err, "new service")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-sigs
		log.Println(sig)
		cancel()
	}()

	if err := svc.Wait(ctx); err != nil {
		return err
	}

	return nil
}
