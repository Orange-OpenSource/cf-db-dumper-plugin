package main

import (
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"github.com/daviddengcn/go-colortext"
	"github.com/Orange-OpenSource/db-dumper-cli-plugin/db_dumper"
	"github.com/codegangsta/cli"
	"github.com/cloudfoundry/cli/cf/errors"
)

/*
*	This is the struct implementing the interface defined by the core CLI. It can
*	be found at  "github.com/cloudfoundry/cli/plugin/plugin.go"
*
 */
type BasicPlugin struct{}
var version_major int = 1
var version_minor int = 0
var version_build int = 0
var helpText string = "Help you to manipulate db-dumper service"
var serviceName string
func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command basic-plugin-command
	if args[0] != "db-dumper" {
		return
	}
	app := cli.NewApp()
	app.Name = "db-dumper"
	app.Version = fmt.Sprintf("%d.%d.%d", version_major, version_minor, version_build)
	app.Usage = "Help you to manipulate db-dumper service"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "name, n",
			Value: db_dumper.SERVICE_NAME,
			Usage: helpText,
			Destination: &serviceName,
		},
	}
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}
	app.Commands = []cli.Command{
		{
			Name:      "create",
			Aliases:     []string{"c"},

			Usage:     "Create a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)",
			ArgsUsage: "[service-name-or-url-of-your-db]",
			Action: func(cg *cli.Context) {
				if len(cg.Args()) == 0 {
					checkError(errors.New("you must provide a service name or an url to a database"))
				}
				dbDumperManager := db_dumper.NewDbDumperManager(serviceName, cliConnection)
				err := dbDumperManager.CreateDump(cg.Args().First())
				checkError(err)
			},
		},
		{
			Name:      "restore",
			Aliases:     []string{"r"},
			Usage:     "Restore a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)",
			ArgsUsage: "[service-name-or-url-of-your-db]",
			Action: func(cg *cli.Context) {
				if len(cg.Args()) == 0 {
					checkError(errors.New("you must provide a service name or an url to a target database"))
				}
				dbDumperManager := db_dumper.NewDbDumperManager(serviceName, cliConnection)
				err := dbDumperManager.RestoreDump(cg.Args().First())
				checkError(err)
			},
		},
		{
			Name:      "delete",
			Aliases:     []string{"d"},
			Usage:     "Delete a instance and all his dumps (dumps can be retrieve during a period)",
			Action: func(cg *cli.Context) {
				dbDumperManager := db_dumper.NewDbDumperManager(serviceName, cliConnection)
				err := dbDumperManager.DeleteDump()
				checkError(err)
			},
		},
	}
	app.Run(args)
}

func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "db-dumper",
		Version: plugin.VersionType{
			Major: version_major,
			Minor: version_minor,
			Build: version_build,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "db-dumper",
				HelpText: "helpText",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "db-dumper\n   Run db-dumper help to see usage",
				},
			},
		},
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Print("error: ")
		ct.Foreground(ct.Red, false)
		fmt.Println(fmt.Sprintf("%v", err))
		ct.ResetColor()
		os.Exit(1)

	}
}
func main() {
	plugin.Start(new(BasicPlugin))
}
