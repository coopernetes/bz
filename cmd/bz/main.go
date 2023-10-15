// Copyright 2022 bonzai-example Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"

	example "github.com/coopernetes/bz/pkg/example"
	install "github.com/coopernetes/bz/pkg/installer"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

func init() {
	Z.Conf.Init()
	Z.Vars.Init()
}

// tree grown from branch
func main() {
	log.SetFlags(0)

	Z.AllowPanic = true

	Cmd.Run()
}

var Cmd = &Z.Cmd{
	Name:      `bz`,
	Summary:   `coop's bonzai command tree`,
	Copyright: `Copyright 2023 Thomas Cooper`,
	Version:   `v0.0.1`,
	License:   `Apache-2.0`,
	Source:    `github.com/coopernetes/bz`,
	UseConf:   true,
	UseVars:   true,
	Commands: []*Z.Cmd{
		// standard external branch imports (see rwxrob/{help,conf,vars})
		help.Cmd, conf.Cmd, vars.Cmd,
		// my local commands
		example.BarCmd, install.GithubCmd,
	},
	Description: `
		This is a Bonzai command tree '{{ cmd .Name }}'. The original creator of Bonzai is
		Robert S Muhlestein (https://github.com/rwxrob/bonzai). There are a number
		of bash functions and shorthands that I use on a regular basis,
		which I hope that this Bonzai tree will help to replace. This project
		is really interesting because of all the scaffolding that it provides
		when compared to the built-in Go flags package (which is completely
		inadequate) or Cobra (which is too heavy and complex for my needs).
		
		I will also reuse a number of commands that Rob has written and uses himself
		via https://github.com/rwxrob/z since there is no point rewriting something that
		already works.`,
}
