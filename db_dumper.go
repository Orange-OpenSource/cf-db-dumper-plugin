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
var version_minor int = 1
var version_build int = 0
var helpText string = "Help you to manipulate db-dumper service"
var serviceName string
var sourceInstance string
var plan string
var showUrl bool
var recent bool
var original bool
var skipInsecure bool
var inStdout bool
var dumpNumber string
var verboseMode bool

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
		cli.BoolFlag{
			Name: "verbose, vvv, vv",
			Usage: "Set the flag if you want to see all output from cloudfoundry cli",
			Destination: &verboseMode,
		},
	}
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}
	app.Commands = []cli.Command{
		{
			Name:      "create",
			Aliases:     []string{"c"},

			Usage:     "Create a dump from a database service (e.g.: mydb) or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)",
			ArgsUsage: "[service-name-or-url-of-your-db]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "plan",
					Usage: "Choose the plan to use. If not set it will ask you to choose one from a list",
					Destination: &plan,
				},
			},
			Action: func(cg *cli.Context) {
				if len(cg.Args()) == 0 {
					checkError(errors.New("you must provide a service name or an url to a database"))
				}
				dbDumperManager := db_dumper.NewDbDumperManager(serviceName, cliConnection, verboseMode)
				err := dbDumperManager.CreateDump(cg.Args().First(), plan)
				checkError(err)
			},
		},
		{
			Name:      "restore",
			Aliases:     []string{"r"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "recent",
					Usage: "Restore from the most recent dump",
					Destination: &recent,
				},
				cli.StringFlag{
					Name: "source-instance",
					Usage: "The db dumper service instance where dumps should be retrieved (this can be the service instance you passed in create e.g.: mydb)",
					Destination: &sourceInstance,
				},
			},
			Usage:     "Restore a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)",
			ArgsUsage: "[service-name-or-url-of-your-db]",
			Action: func(cg *cli.Context) {
				if len(cg.Args()) == 0 {
					checkError(errors.New("you must provide a service name or an url to a target database"))
				}
				dbDumperManager := db_dumper.NewDbDumperManager(serviceName, cliConnection, verboseMode)
				prefix, err := dbDumperManager.GetNamePrefix()
				checkError(err)
				suffix, _ := dbDumperManager.GetNameSuffix()
				err = dbDumperManager.RestoreDump(cg.Args().First(), recent, prefix + sourceInstance + suffix)
				if err != nil {
					warning("Trying on a non generated service", err)
					err = dbDumperManager.RestoreDump(cg.Args().First(), recent, sourceInstance)
					checkError(err)
				}
			},
		},
		{
			Name:      "delete",
			ArgsUsage: "[service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)",
			Aliases:     []string{"d"},
			Usage:     "Delete a instance and all its dumps (dumps can be retrieve during a period)",
			Action: func(cg *cli.Context) {
				dbDumperManager := db_dumper.NewDbDumperManager(serviceName, cliConnection, verboseMode)
				if len(cg.Args()) > 0 {
					prefix, err := dbDumperManager.GetNamePrefix()
					checkError(err)
					suffix, _ := dbDumperManager.GetNameSuffix()
					err = dbDumperManager.DeleteDump(prefix + cg.Args().First() + suffix)
					if err != nil {
						warning("Trying on a non generated service", err)
						err = dbDumperManager.DeleteDump(cg.Args().First())
						checkError(err)
					}
				} else {
					err := dbDumperManager.DeleteDump("")
					checkError(err)
				}

			},
		},
		{
			Name:      "list",
			Aliases:     []string{"l"},
			Usage:     "List dumps for an instance",
			ArgsUsage: "[service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "show-url, s",
					Usage: "If you want to see download url and dashboard url",
					Destination: &showUrl,
				},
			},
			Action: func(cg *cli.Context) {
				dbDumperManager := db_dumper.NewDbDumperManager(serviceName, cliConnection, verboseMode)

				if len(cg.Args()) > 0 {
					prefix, err := dbDumperManager.GetNamePrefix()
					checkError(err)
					suffix, _ := dbDumperManager.GetNameSuffix()
					err = dbDumperManager.ListFromInstanceName(prefix + cg.Args().First() + suffix, showUrl)
					if err != nil {
						warning("Trying on a non generated service", err)
						err = dbDumperManager.ListFromInstanceName(cg.Args().First(), showUrl)
						checkError(err)
					}
				} else {
					err := dbDumperManager.List(showUrl)
					checkError(err)
				}

			},
		},
		{
			Name:      "download",
			Aliases:     []string{"dl"},
			Usage:     "Download a dump to your drive",
			ArgsUsage: "[service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "skip-ssl-validation, k",
					Usage: "Skip the ssl validation (for self-signed certificate mainly)",
					Destination: &skipInsecure,
				},
				cli.BoolFlag{
					Name: "recent",
					Usage: "Download from the most recent dump",
					Destination: &recent,
				},
				cli.BoolFlag{
					Name: "original",
					Usage: "Download the original file (e.g.: download directly sql file instead of compressed file)",
					Destination: &original,
				},
				cli.StringFlag{
					Name: "dump-number, p",
					Usage: "Download from the number showed when using 'db-dumper list'",
					Value: "",
					Destination: &dumpNumber,
				},
				cli.BoolFlag{
					Name: "stdout, o",
					Usage: "Show file directly in stdout (service instance is no more optionnal and you need to use flag --dump-number or --recent)",
					Destination: &inStdout,
				},
			},
			Action: func(cg *cli.Context) {
				dbDumperManager := db_dumper.NewDbDumperManager(serviceName, cliConnection, verboseMode)
				if len(cg.Args()) > 0 {
					prefix, err := dbDumperManager.GetNamePrefix()
					checkError(err)
					suffix, _ := dbDumperManager.GetNameSuffix()
					err = dbDumperManager.DownloadDumpFromInstanceName(prefix + cg.Args().First() + suffix, skipInsecure, recent, inStdout, original, dumpNumber)
					if err != nil {
						warning("Trying on a non generated service", err)
						err = dbDumperManager.DownloadDumpFromInstanceName(cg.Args().First(), skipInsecure, recent, inStdout, original, dumpNumber)
						checkError(err)
					}
				} else {
					err := dbDumperManager.DownloadDump(skipInsecure, recent, inStdout, original, dumpNumber)
					checkError(err)
				}

			},
		},
		{
			Name:      "open",
			Aliases:  []string{"o"},
			Usage:     "Open dump in your browser",
			ArgsUsage: "[service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "recent",
					Usage: "Open dump page from the most recent dump",
					Destination: &recent,
				},
				cli.StringFlag{
					Name: "dump-number, p",
					Usage: "Open from the number showed when using 'db-dumper list'",
					Value: "",
					Destination: &dumpNumber,
				},
			},
			Action: func(cg *cli.Context) {
				dbDumperManager := db_dumper.NewDbDumperManager(serviceName, cliConnection, verboseMode)
				if len(cg.Args()) > 0 {
					prefix, err := dbDumperManager.GetNamePrefix()
					checkError(err)
					suffix, _ := dbDumperManager.GetNameSuffix()
					err = dbDumperManager.ShowDumpFromInstanceName(prefix + cg.Args().First() + suffix, recent, dumpNumber)
					if err != nil {
						warning("Trying on a non generated service", err)
						err = dbDumperManager.ShowDumpFromInstanceName(cg.Args().First(), recent, dumpNumber)
						checkError(err)
					}
				} else {
					err := dbDumperManager.ShowDump(recent, dumpNumber)
					checkError(err)
				}

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
				HelpText: helpText,

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "db-dumper\n   Run db-dumper help to see usage",
				},
			},
		},
	}
}
func warning(message string, err error) {
	fmt.Print("Warning: ")
	ct.Foreground(ct.Yellow, false)
	fmt.Println(fmt.Sprintf("%v", err))
	fmt.Println(message)
	ct.ResetColor()
}
func checkError(err error) {
	if err == nil {
		return
	}
	fmt.Print("Error: ")
	ct.Foreground(ct.Red, false)
	fmt.Println(fmt.Sprintf("%v", err))
	ct.ResetColor()
	if verboseMode == false {
		fmt.Println("Tip: use --verbose to see what's going wrong with cf core command")
	}
	os.Exit(1)
}
func main() {
	plugin.Start(new(BasicPlugin))
}
