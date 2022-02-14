// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/sdk/gitea"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"strings"
)

// CmdPulls represents to login a gitea server.
var CmdPulls = cli.Command{
	Name:        "pulls",
	Usage:       "Operate with pulls of the repository",
	Description: `Operate with pulls of the repository`,
	Action:      runPulls,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "login, l",
			Usage: "Indicate one login, optional when inside a gitea repository",
		},
		cli.StringFlag{
			Name:  "repo, r",
			Usage: "Indicate one repository, optional when inside a gitea repository",
		},
		cli.BoolFlag{
			Name:  "matchLogin, ml",
			Usage: "Results will be filtered to match the current login value",
		},
	},
}

func runPulls(ctx *cli.Context) error {
	login, owner, repo := initCommand(ctx)

	prs, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		Page:  0,
		State: string(gitea.StateOpen),
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(prs) == 0 {
		fmt.Println("No pull requests left")
		return nil
	}

	matchLogin := ctx.Bool("matchLogin")
	for _, pr := range prs {
		if pr == nil || (matchLogin && login.Name != pr.Poster.Email) {
			continue
		}
		name := pr.Poster.FullName
		if len(name) == 0 {
			name = pr.Poster.UserName
		}

		var jiraTicket string
		idx := strings.Index(pr.Body, "[PRD-")
		if idx > 0 {
			partial := pr.Body[idx+1:]
			toIdx := strings.Index(partial, "]")
			jiraTicket = partial[:toIdx]
		}
		fmt.Printf("#%d\t%s\t%s\t%s\t%s\n", pr.Index, name, pr.Updated.Format("2006-01-02 15:04:05"), pr.Title, jiraTicket)
	}

	return nil
}
