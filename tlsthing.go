package main

import (
	"os"
	"strconv"
	dh "fatgo/docopt_helpers"
)

// Default values may be overridden by build system
var build_mode = "Debug"
var app_name = "tlsthing"
var app_version = "v0.0.0"
var internal_repo_path = "/var/tlsthing/repo"

var default_refresh = 60

func main() {
	loglevel := LL_INFO
	if build_mode == "Debug" {
		loglevel = LL_DEBUG
	}

	view := &View{loglevel: loglevel}
	view.begin()
	defer view.end()

	s_refresh := strconv.Itoa(default_refresh)

	options := map[string]dh.HelpOption{}
	options["--path"] = dh.HelpOption{
		"<config>",
		"Config file path",
		""}
	options["--repo"] = dh.HelpOption{
		"<url>",
		"Config repo URL",
		""}
	options["--refresh"] = dh.HelpOption{
		"<seconds>",
		"Config refresh in seconds",
		s_refresh}

	usages := []string{}
	usages = append(usages, "--path [--repo] [--refresh]")
	uses, opts := dh.MakeUses(usages, options)

	args, err := dh.NoExitParser.ParseArgs(
		dh.BuildUsageString(uses, opts),
		nil, // this nil becomes os.Args[1:]
		(app_name + " " + app_version))

	if err == nil {
		var config Args
		err = args.Bind(&config)
		if err == nil {
			app := Control{
				args: config,
				alive: dh.IsSet(args, "--path"),
				view: view}
			app.main()
		}
	}

	if err != nil {
		view.log(
			LL_ERROR,
			err.Error())
		os.Exit(1)
	}
}
