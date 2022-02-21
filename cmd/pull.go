// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/sdk/gitea"
	"fmt"
	"log"
	"strconv"

	"github.com/urfave/cli"
)

// CmdPull represents to login a gitea server.
var CmdPull = cli.Command{
	Name:        "pull",
	Usage:       "Operate with pull of the repository",
	Description: `Operate with pull of the repository and returns the pull request's body message`,
	Action:      runPull,
	ArgsUsage:   "[<pull request index>]",
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
			Name:  "merge, mg",
			Usage: "Merge a certain pull request",
		},
	},
}

func runPull(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Present() {
		return runPullDetail(ctx, args.First())
	}
	return runPulls(ctx)
}

func runPullDetail(ctx *cli.Context, index string) error {
	login, owner, repo := initCommand(ctx)

	idx, err := strconv.ParseInt(index, 10, 64)
	if err != nil {
		return err
	}

	pr, _, err := login.Client().GetPullRequest(owner, repo, idx)
	if err != nil {
		log.Fatal(err)
	}

	if pr == nil {
		fmt.Println("Pull request not found")
		return nil
	}

	if ctx.Bool("merge") {
		_, _, err := login.Client().MergePullRequest(owner, repo, idx, gitea.MergePullRequestOption{Style: gitea.MergeStyleRebase})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		name := pr.Poster.FullName
		if len(name) == 0 {
			name = pr.Poster.UserName
		}
		fmt.Printf("#%d\t%s\t%s\t%s\t%s\n", pr.Index, name, pr.Updated.Format("2006-01-02 15:04:05"), pr.Title, pr.Body)
	}

	return nil
}
