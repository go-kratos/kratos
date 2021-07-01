package change

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/base"
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "changelog",
	Short: "Get a kratos change log",
	Long:  "Get a kratos release or commits info. Example: kratos changelog dev or kratos changelog {version}",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	version := "latest"
	if len(args) > 0 {
		version = args[0]
	}
	if version == "dev" {
		info, err := base.GetCommitsInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		for _, value := range info {
			fmt.Printf(
				"\033[36mcommit: %s \033[0m\nAuthor: %s\nDate: %s\nUrl: %s\n\n%s\n\n",
				value.Sha, value.Author.Login,
				value.Commit.Author.Date,
				value.HtmlUrl,
				value.Commit.Message,
			)
		}
	} else {
		info, err := base.GetReleaseInfo(version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		reg := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z|<[\S\s]+?>`)
		body := reg.ReplaceAll([]byte(info.Body), []byte(""))
		if string(body) == "" {
			body = []byte("no release info")
		}
		splitters := "--------------------------------------------"
		fmt.Printf(
			"Author: %s\nDate: %s\nUrl: %s\n\n%s\n\n%s\n\n%s\n",
			info.Author.Login,
			info.PublishedAt,
			info.HtmlUrl,
			splitters,
			body,
			splitters,
		)
	}
}
