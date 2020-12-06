/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package new

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// CmdNew represents the init command
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a service template template",
	Long:  `Create a service project using the repository template. Example: kratos new helloworld`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p := &Project{Name: args[0]}
	if err := p.Generate(ctx, wd); err != nil {
		log.Fatal(err)
	}
}
