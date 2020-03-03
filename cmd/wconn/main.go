// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2017 Canonical Ltd.  All rights reserved.

package main

import (
	"fmt"
	"os"

	"github.com/greenbrew/wconn/cli"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
}

func run() error {
	// Parse the command line arguments and execute the command
	parser := flags.NewParser(&cli.Command, flags.HelpFlag)
	_, err := parser.Parse()

	if err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp || e.Type == flags.ErrCommandRequired {
				parser.WriteHelp(os.Stdout)
				return nil
			}
		}
		fmt.Println(err)
	}

	return err
}
