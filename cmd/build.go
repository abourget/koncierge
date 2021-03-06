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

var doPush bool
var doDeploy bool

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a Koncierge project",
	Long: `Example usage:

    koncierge build
    koncierge build --push
    koncierge build --push --deploy
    koncierge build --deploy   ; if your build script pushes automatically
    koncierge -t prod build
`,
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

		if doDeploy && b.Config.Targets[target].Deployment == nil {
			fmt.Printf("koncierge: target %q: --deploy requested, but no deployment in configuration", target)
			os.Exit(103)
		}

		if err := b.Build(target); err != nil {
			fmt.Printf("koncierge: build failed: %s\n", err)
			os.Exit(200)
		}

		if doPush {
			if err := b.Push(target); err != nil {
				fmt.Printf("koncierge: push failed: %s\n", err)
				os.Exit(210)
			}

		}
		if doDeploy {
			if err := b.Deploy(target); err != nil {
				fmt.Printf("koncierge: deploy failed: %s\n", err)
				os.Exit(220)
			}
		}

		fmt.Println("koncierge: build command terminated successfully")
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVarP(&doPush, "push", "p", false, "Push after a successful build")
	buildCmd.Flags().BoolVarP(&doDeploy, "deploy", "d", false, "Deploy after a successful push. Requires --push")

}
