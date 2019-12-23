// Copyright © 2019 Percona, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Percona-Lab/percona-dbaas-cli/dbaas-cli/cmd/mysql"
	"github.com/Percona-Lab/percona-dbaas-cli/dbaas-cli/cmd/mongo"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "percona-dbaas",
	Short: "The simplest DBaaS tool in the world",
	Long: `    Hello, it is the simplest DBaaS tool in the world,
	please use commands below to manage your DBaaS.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := detectFormat(cmd)
		if err != nil {
			log.Error("detect format:", err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("demo", false, "demo mode (no spinners)")
	rootCmd.PersistentFlags().MarkHidden("demo")
	rootCmd.PersistentFlags().String("output", "", `Answers format. Can be "json" or "text". "text" is set by default`)
	rootCmd.AddCommand(mysql.PXCCmd)
	rootCmd.AddCommand(mongo.MongoCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func detectFormat(cmd *cobra.Command) error {
	format, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	switch format {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			DisableTimestamp: true,
		})
	default:
		log.SetFormatter(&cliTextFormatter{log.TextFormatter{}})
	}
	return nil
}

type cliTextFormatter struct {
	log.TextFormatter
}

func (f *cliTextFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	if entry.Level == log.ErrorLevel {
		b.WriteString("[Error] ")
	}
	if entry.Message != "" {
		b.WriteString(entry.Message)
	}

	if len(entry.Data) == 0 {
		b.WriteString("\n")
		return b.Bytes(), nil
	}

	for _, v := range entry.Data {
		fmt.Fprint(b, "\n", v)
	}
	b.WriteString("\n")
	return b.Bytes(), nil
}
