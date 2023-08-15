package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/tsuyoshiwada/go-gitcmd"
	"github.com/urfave/cli/v2"
)

var version string

// CreateApp creates and initializes CLI application
// with description, flags, version, etc.
func CreateApp(actionFunc cli.ActionFunc) *cli.App {
	ttl := color.New(color.FgYellow).SprintFunc()

	cli.AppHelpTemplate = fmt.Sprintf(`
%s
  {{.Name}} [options] <tag query>

    There are the following specification methods for <tag query>.

    1. <old>..<new> - MR contained in <old> tags from <new>.
    2. <name>..     - MR from the <name> to the latest tag.
    3. ..<name>     - MR from the oldest tag to <name>.
    4. <name>       - MR contained in <name>.

%s
  {{range .Flags}}{{.}}
  {{end}}
%s

  $ {{.Name}}

    If <tag query> is not specified, it corresponds to all tags.
    This is the simplest example.

  $ {{.Name}} 1.0.0..2.0.0

    The above is a command to generate CHANGELOG including MR of 1.0.0 to 2.0.0.

  $ {{.Name}} 1.0.0

    The above is a command to generate CHANGELOG including MR of only 1.0.0.

  $ {{.Name}} $(git describe --tags $(git rev-list --tags --max-count=1))

    The above is a command to generate CHANGELOG with the MR included in the latest tag.

  $ {{.Name}} --output CHANGELOG.md

    The above is a command to output to CHANGELOG.md instead of standard output.

  $ {{.Name}} --config custom/dir/config.yml

    The above is a command that uses a configuration file placed other than ".chglog/config.yml".

  $ {{.Name}} --path path/to/my/component --output CHANGELOG.component.md

    Filter commits by specific paths or files in git and output to a component specific changelog.
  $ {{.Name}} --bot 
    Push mr-changelog  Changelog to Feishu Group
  $ {{.Name}} --ai 
    Use ai create CHANGELOG
`,
		ttl("USAGE:"),
		ttl("OPTIONS:"),
		ttl("EXAMPLE:"),
	)

	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		cli.HelpPrinterCustom(colorable.NewColorableStdout(), templ, data, nil)
	}

	app := cli.NewApp()
	app.Name = "mr-changelog "
	app.Usage = "todo usage for mr-changelog "
	app.Version = version

	app.Flags = []cli.Flag{
		// init
		&cli.BoolFlag{
			Name:  "init",
			Usage: "generate the mr-changelog  configuration file in interactive",
		},
		&cli.StringFlag{
			Name:  "app_id",
			Usage: "feishu robot app secret",
		},
		&cli.StringFlag{
			Name:  "app_secret",
			Usage: "feishu robot app secret",
		},
		&cli.StringFlag{
			Name:  "chat_id",
			Usage: "feishu robot send group chat_id,Please use , to separate multiple",
		},
		&cli.StringFlag{
			Name:  "bot_title",
			Usage: "feishu robot send release title",
		},
		// path
		&cli.StringSliceFlag{
			Name:  "path",
			Usage: "Filter commits by path(s). Can use multiple times.",
		},

		// config
		&cli.StringFlag{
			Name:    "config, c",
			Aliases: []string{"c"},
			Usage:   "specifies a different configuration file to pick up",
			Value:   ".chglog/mr_config.yml",
		},

		// template
		&cli.StringFlag{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "specifies a template file to pick up. If not specified, use the one in config",
		},

		// repository url
		&cli.StringFlag{
			Name:  "repository-url",
			Usage: "specifies git repo URL. If not specified, use 'repository_url' in config",
		},

		// repository url
		&cli.StringFlag{
			Name:  "token",
			Usage: "specifies git repo token. If not specified, use 'token' in config",
		},

		// output
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "output path and filename for the changelogs. If not specified, output to stdout",
		},

		&cli.StringFlag{
			Name:  "next-tag",
			Usage: "treat unreleased commits as specified tags (EXPERIMENTAL)",
		},

		// // tag-filter-pattern
		// &cli.StringFlag{
		// 	Name:  "tag-filter-pattern",
		// 	Usage: "Regular expression of tag filter. Is specified, only matched tags will be picked",
		// },

		// sort
		// &cli.StringFlag{
		// 	Name:        "sort",
		// 	Usage:       "Specify how to sort tags; currently supports \"date\" or by \"semver\"",
		// 	DefaultText: "date",
		// },

		&cli.StringFlag{
			Name:  "create-tag",
			Usage: "create tag by CHANGELOG",
		},

		&cli.BoolFlag{
			Name:  "ai",
			Usage: "use ai create CHANGELOG",
		},
		&cli.BoolFlag{
			Name:  "update",
			Usage: "update tag release",
		},
		&cli.StringFlag{
			Name:        "ai-type",
			Usage:       "which ai API to use create CHANGELOG",
			DefaultText: "poe",
			Value:       "poe",
		},

		// bot
		&cli.BoolFlag{
			Name:  "bot",
			Usage: "push mr-changelog  changelog to feishu group",
		},

		// help & version
		cli.HelpFlag,
		cli.VersionFlag,
	}

	app.Action = actionFunc

	return app
}

func AppAction(c *cli.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get working directory", err)
		os.Exit(1)
	}
	// initializer
	if c.Bool("init") {
		initializer := NewInitializer(
			&InitContext{
				WorkingDir: wd,
				Stdout:     colorable.NewColorableStdout(),
				Stderr:     colorable.NewColorableStderr(),
			},
			fs,
			NewQuestioner(
				gitcmd.New(&gitcmd.Config{
					Bin: "git",
				}),
				fs,
			),
			NewConfigBuilder(),
			templateBuilderFactory,
		)

		os.Exit(initializer.Run())
	}
	chglogCLI := NewCLI(
		&CLIContext{
			WorkingDir: wd,
			Stdout:     colorable.NewColorableStdout(),
			Stderr:     colorable.NewColorableStderr(),
			ConfigPath: c.String("config"),
			Template:   c.String("template"),
			OutputPath: c.String("output"),
			Query:      c.Args().First(),
			NextTag:    c.String("next-tag"),
			// TagFilterPattern: c.String("tag-filter-pattern"),
			// Sort:             c.String("sort"),
			AI:            c.Bool("ai"),
			AIType:        c.String("ai-type"),
			PushBot:       c.Bool("bot"),
			RepositoryURL: c.String("repository-url"),
			Token:         c.String("token"),
			AppID:         c.String("app_id"),
			AppSecret:     c.String("app_secret"),
			ChatID:        strings.Split(c.String("chat_id"), ","),
			BotTitle:      c.String("bot_title"),
			Update:        c.Bool("update"),
		},
		fs,
		NewConfigLoader(),
		NewGenerator(),
	)

	os.Exit(chglogCLI.Run())

	return nil
}
