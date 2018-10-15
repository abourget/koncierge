// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"
	"os"

	"github.com/abourget/koncierge/build"
	"github.com/abourget/koncierge/config"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a previously built image",
	Long: `deploy will re-read the tag (from local git repo or a tag file), and
deploy that to the target environment (see -t, defaults to 'default')`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.WalkConfig()
		if err != nil {
			fmt.Println("error loading configuration:", err)
			os.Exit(101)
		}

		err = conf.Validate()
		if err != nil {
			fmt.Printf("In %q: %s\n\n", conf.FilePath, err.Error())
			os.Exit(101)
		}

		b := build.New(conf)

		target, err := b.TargetWithDefault(cliTarget)
		if err != nil {
			fmt.Println(err)
			os.Exit(102)
		}

		if b.Config.Targets[target].Deployment == nil {
			fmt.Printf("koncierge: target %q: --deploy requested, but no deployment in configuration", target)
			os.Exit(103)
		}

		if err := b.Deploy(target); err != nil {
			fmt.Printf("koncierge: deploy failed: %s\n", err)
			os.Exit(220)
		}

		fmt.Println("koncierge: deploy command terminated successfully")
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
