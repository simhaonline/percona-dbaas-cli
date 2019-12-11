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

package mysql

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	dbaas "github.com/Percona-Lab/percona-dbaas-cli/dbaas-lib"
)

// delCmd represents the list command
var delCmd = &cobra.Command{
	Use:   "delete-db <mysql-cluster-name>",
	Short: "Delete MySQL cluster",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you have to specify resource name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !*forced {
			var yn string
			preservText := "YOUR DATA WILL BE SAVED\n"
			if !*preserve {
				preservText = "ALL YOUR DATA WILL BE LOST. USE '--preserve-data' FLAG FOR SAVE IT\n"
			}
			fmt.Printf("ARE YOU SURE YOU WANT TO DELETE THE DATABASE '%s'? Yes/No\n"+preservText, args[0])
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				yn = strings.TrimSpace(scanner.Text())
				break
			}
			if yn != "yes" && yn != "Yes" && yn != "YES" && yn != "Y" && yn != "y" {
				return
			}
		}
		sp := spinner.New(spinner.CharSets[14], 250*time.Millisecond)
		sp.Color("green", "bold")
		sp.Prefix = "Deleting cluster........."
		sp.FinalMSG = ""
		sp.Start()
		i := dbaas.Instance{
			Name:     args[0],
			Engine:   *delEngine,
			Provider: *delProvider,
		}
		deletePVC := false
		if !*preserve {
			deletePVC = true
		}
		dataStorage, err := dbaas.DeleteDB(i, deletePVC)
		if err != nil {
			log.Error("delete db: ", err)
			return
		}
		sp.Stop()
		log.Println("Deleting done")
		if *preserve {
			log.Println("Your data store in " + dataStorage)
		}
	},
}

var envDlt *string
var delOptions *string
var delProvider *string
var delEngine *string
var forced *bool
var preserve *bool

func init() {
	forced = delCmd.Flags().BoolP("yes", "y", false, "Unswer yes for questions")
	delProvider = delCmd.Flags().String("provider", "", "Provider")
	delEngine = delCmd.Flags().String("engine", "", "Engine")
	preserve = delCmd.Flags().Bool("preserve-data", false, "Do not delete data")

	PXCCmd.AddCommand(delCmd)
}
