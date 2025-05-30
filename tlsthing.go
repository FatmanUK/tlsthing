package main

import (
	"os"
	"github.com/docopt/docopt-go"
	"fatgo/docopt_helpers"
)

// Default values may be overridden by build system
var build_mode = "Debug"
var app_name = "tlsthing"
var app_version = "v0.0.0"

var noExitParser = &docopt.Parser{
	HelpHandler:   docopt.PrintHelpOnly,
	OptionsFirst:  false,
	SkipHelpFlags: false,
}

func parse_args() error {
	loglevel := LL_INFO
	if build_mode == "Debug" {
		loglevel = LL_DEBUG
	}

	view := &View{loglevel: loglevel}
	view.begin()
	defer view.end()

	uses := []string{}
	uses = append(uses, "--repo=<url> [--refresh=<seconds>]")

	opts := make(map[string]string)
	opts["--repo=<url>"] = "Config repo URL or file path"
	opts["--refresh=<seconds>"] = "Config refresh in seconds [default: 60]"

	args, err := noExitParser.ParseArgs(
		docopt_helpers.BuildUsageString(uses, opts),
		nil, // this nil becomes os.Args[1:]
		(app_name + " " + app_version))
	if err == nil {
		alive := false
		if args["--repo"] != "" {
			alive = true
		}
		var config Args
		err = args.Bind(&config)
		if err == nil {
			app := Control{
				args: config,
				alive: alive,
				view: view}
			app.main()
		} else {
			view.log(LL_ERROR, err.Error())
		}
	} else {
		view.log(LL_ERROR, err.Error())
	}
	return err
}

func main() {
	if parse_args() != nil {
		os.Exit(1)
	}
}
