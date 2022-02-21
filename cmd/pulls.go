// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli"
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
		cli.StringFlag{
			Name:  "match, m",
			Usage: "Results will be filtered according to the supplied name value",
		},
	},
}

func runPulls(ctx *cli.Context) error {
	login, owner, repo := initCommand(ctx)
	matchName := ctx.String("match")
	if ctx.Bool("matchLogin") {
		matchName = login.Name
	}

	for i := 1; true; i++ {
		prs, res, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
			ListOptions: gitea.ListOptions{Page: i},
			State:       gitea.StateOpen,
		})
		if err != nil || res.StatusCode >= 400 {
			log.Fatalf("invalid response received: %v", err)
		}
		lenPrs := len(prs)
		if lenPrs == 0 {
			return nil
		}

		for _, pr := range prs {
			name := pr.Poster.Email
			if pr == nil || (matchName != "" && matchName != pr.Poster.Email) {
				continue
			}

			var jiraTicket string
			idx := strings.Index(pr.Body, "[PRD-")
			if idx > 0 {
				partial := pr.Body[idx+1:]
				toIdx := strings.Index(partial, "]")
				if toIdx > 0 {
					jiraTicket = partial[:toIdx]
				}
			}
			fmt.Printf("#%d\t%s\t%s\t%s\t%s\n", pr.Index, name, pr.Updated.Format("2006-01-02 15:04:05"), pr.Title, jiraTicket)
		}
	}

	return nil
}
