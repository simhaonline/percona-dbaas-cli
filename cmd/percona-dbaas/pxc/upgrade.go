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

package pxc

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Percona-Lab/percona-dbaas-cli/dbaas"
	"github.com/Percona-Lab/percona-dbaas-cli/dbaas/pxc"
)

// upgradeCmd represents the edit command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade-db <pxc-cluster-name> <to-version>",
	Short: "Upgrade MySQL cluster",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You have to specify pxc-cluster-name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		switch *upgradeAnswerFormat {
		case "json":
			log.Formatter = new(logrus.JSONFormatter)
		}
		dbservice, err := dbaas.New(*envUpgrd)
		if err != nil {
			log.Errorln("new dbservice:", err)
			return
		}

		app := pxc.New(name, defaultVersion, "")

		sp := spinner.New(spinner.CharSets[14], 250*time.Millisecond)
		sp.Color("green", "bold")
		demo, err := cmd.Flags().GetBool("demo")
		if demo && err == nil {
			sp.UpdateCharSet([]string{""})
		}
		sp.Prefix = "Looking for the cluster..."
		sp.FinalMSG = ""
		sp.Start()
		defer sp.Stop()

		ext, err := dbservice.IsObjExists("pxc", name)
		if err != nil {
			log.Errorln("check if cluster exists:", err)
			return
		}

		if !ext {
			sp.Stop()
			log.Errorln("unable to find cluster pxc/" + name)
			list, err := dbservice.List("pxc")
			if err != nil {
				log.Errorln("list pxc clusters:", err)
				return
			}

			log.Println("avaliable clusters:", list)
			return
		}

		created := make(chan string)
		msg := make(chan dbaas.OutuputMsg)
		cerr := make(chan error)

		oparg := ""
		if len(args) > 1 {
			oparg = args[1]
		}
		appsImg, err := app.Images(oparg, cmd.Flags())
		if err != nil {
			log.Errorln("setup images for upgrade:", err)
			return
		}

		go dbservice.Upgrade("pxc", app, appsImg, created, msg, cerr)
		sp.Lock()
		sp.Prefix = "Upgrading cluster..."
		sp.Unlock()
		for {
			select {
			case <-created:
				okmsg, _ := dbservice.ListName("pxc", name)
				sp.FinalMSG = ""
				sp.Stop()
				log.Println("upgrade cluster done.", okmsg)
				return
			case omsg := <-msg:
				switch omsg.(type) {
				case dbaas.OutuputMsgDebug:
					// fmt.Printf("\n[debug] %s\n", omsg)
				case dbaas.OutuputMsgError:
					sp.Stop()
					log.Errorln("perator log error:", omsg.String())
					sp.Start()
				}
			case err := <-cerr:
				log.Errorln("upgrade pxc:", err)
				sp.HideCursor = true
				return
			}
		}
	},
}

var envUpgrd *string
var upgradeAnswerFormat *string

func init() {
	upgradeCmd.Flags().String("database-image", "", "Custom image to upgrade pxc to")
	upgradeCmd.Flags().String("proxysql-image", "", "Custom image to upgrade proxySQL to")
	upgradeCmd.Flags().String("backup-image", "", "Custom image to upgrade backup to")
	envUpgrd = upgradeCmd.Flags().String("environment", "", "Target kubernetes cluster")

	upgradeAnswerFormat = upgradeCmd.Flags().String("output", "", "Answers format")

	PXCCmd.AddCommand(upgradeCmd)
}
