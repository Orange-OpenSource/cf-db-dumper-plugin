package main

import (
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"github.com/daviddengcn/go-colortext"
	"github.com/orange-cloudfoundry/db-dumper-cli-plugin/db_dumper"
	"github.com/codegangsta/cli"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/mitchellh/go-homedir"
	"github.com/orange-cloudfoundry/db-dumper-cli-plugin/db_dumper/model"
	"encoding/json"
	"path"
	"io/ioutil"
)

/*
*	This is the struct implementing the interface defined by the core CLI. It can
*	be found at  "github.com/cloudfoundry/cli/plugin/plugin.go"
*
 */
type BasicPlugin struct {
	config          *model.Config
	dbDumperManager *db_dumper.DbDumperManager
}

var version_major int = 1
var version_minor int = 3
var version_build int = 0
var helpText string = "Help you to manipulate db-dumper service"
var configFile string = ".db-dumper"
var sourceInstance string
var tags string
var seeAllDumps bool
var plan string
var org string
var space string
var showUrl bool
var recent bool
var original bool
var skipInsecure bool
var force bool
var inStdout bool
var dumpNumber string
var verboseMode bool

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command basic-plugin-command
	if args[0] != "db-dumper" {
		return
	}

	app := cli.NewApp()
	if !configFileExist() {
		registerConfig(&model.Config{
			Target: db_dumper.SERVICE_NAME,
		})
	}
	c.config = retrieveConfig()
	c.dbDumperManager = db_dumper.NewDbDumperManager(c.config.Target, cliConnection, verboseMode)
	app.Name = "db-dumper"
	app.Version = fmt.Sprintf("%d.%d.%d", version_major, version_minor, version_build)
	app.Usage = "Help you to manipulate db-dumper service"
	app.Flags = []cli.Flag{
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
			Name:      "target",
			Aliases:     []string{"t"},

			Usage:     "Target a db-dumper service",
			ArgsUsage: "[db-dumper-service]",
			Description: "Pass the db-dumper service name, default is db-dumper-service (e.g.: db-dumper-service-dev)",
			Action: c.targetCommand,
		},
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
				cli.StringFlag{
					Name: "tags",
					Usage: "Pass a list of tags to create a dump with these tags (e.g.: --tags=mytag1,mytag2...)",
					Destination: &tags,
				},
			},
			Action: c.createCommand,
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
				cli.BoolFlag{
					Name: "see-all-dumps",
					Usage: "Restore by selecting one of dumps available for the targetted database (made by all db-dumper service instance linked to this database)",
					Destination: &seeAllDumps,
				},
				cli.StringFlag{
					Name: "tags",
					Usage: "Restore by selecting one of dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)",
					Destination: &tags,
				},
				cli.StringFlag{
					Name: "org, o",
					Usage: "Org which contain the service name passed",
					Destination: &org,
				},
				cli.StringFlag{
					Name: "space, s",
					Usage: "Space which contain the service name passed",
					Destination: &space,
				},
				cli.StringFlag{
					Name: "source-instance",
					Value: "",
					Usage: "The db dumper service instance where dumps should be retrieved (this can be the service instance you passed in create e.g.: mydb)",
					Destination: &sourceInstance,
				},
			},
			Usage:     "Restore a dump to a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)",
			ArgsUsage: "[service-name-or-url-of-your-db]",
			Action: c.restoreCommand,
		},
		{
			Name:      "delete",
			ArgsUsage: "[service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)",
			Aliases:     []string{"d"},
			Usage:     "Delete a instance and all its dumps (dumps can be retrieve during a period)",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "force, f",
					Usage: "Force deletion without confirmation",
					Destination: &force,
				},
			},
			Action: c.deleteCommand,
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
				cli.BoolFlag{
					Name: "see-all-dumps",
					Usage: "See all dumps for the database (made by all db-dumper service instance linked to this database)",
					Destination: &seeAllDumps,
				},
				cli.StringFlag{
					Name: "tags",
					Usage: "See dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)",
					Destination: &tags,
				},
			},
			Action: c.listCommand,
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
				cli.BoolFlag{
					Name: "see-all-dumps",
					Usage: "Download dumps by selecting one of dumps available for the targetted database (made by all db-dumper service instance linked to this database)",
					Destination: &seeAllDumps,
				},
				cli.StringFlag{
					Name: "tags",
					Usage: "Download dumps by selecting one of dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)",
					Destination: &tags,
				},
			},
			Action: c.downloadCommand,
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
				cli.BoolFlag{
					Name: "see-all-dumps",
					Usage: "Show dumps by selecting one of dumps available for the targetted database (made by all db-dumper service instance linked to this database)",
					Destination: &seeAllDumps,
				},
				cli.StringFlag{
					Name: "tags",
					Usage: "Show dumps by selecting one of dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)",
					Destination: &tags,
				},
			},
			Action: c.showCommand,
		},
	}
	app.Run(args)
}
func (c *BasicPlugin) showCommand(cg *cli.Context) {
	if len(cg.Args()) > 0 {
		instanceName := c.getInstanceName(cg.Args().First())
		err := c.dbDumperManager.ShowDumpFromInstanceName(instanceName, recent, dumpNumber, seeAllDumps, tags)
		checkError(err)

	} else {
		err := c.dbDumperManager.ShowDump(recent, dumpNumber, seeAllDumps, tags)
		checkError(err)
	}

}
func (c *BasicPlugin) downloadCommand(cg *cli.Context) {
	if len(cg.Args()) > 0 {
		instanceName := c.getInstanceName(cg.Args().First())
		err := c.dbDumperManager.DownloadDumpFromInstanceName(instanceName, skipInsecure, recent, inStdout, original, dumpNumber, seeAllDumps, tags)
		checkError(err)
	} else {
		err := c.dbDumperManager.DownloadDump(skipInsecure, recent, inStdout, original, dumpNumber, seeAllDumps, tags)
		checkError(err)
	}

}
func (c *BasicPlugin) listCommand(cg *cli.Context) {
	if len(cg.Args()) > 0 {
		instanceName := c.getInstanceName(cg.Args().First())
		err := c.dbDumperManager.ListFromInstanceName(instanceName, showUrl, seeAllDumps, tags)
		checkError(err)
	} else {
		err := c.dbDumperManager.List(showUrl, seeAllDumps, tags)
		checkError(err)
	}

}
func (c *BasicPlugin) getInstanceName(instanceNameArg string) string {
	prefix, err := c.dbDumperManager.GetNamePrefix()
	checkError(err)
	suffix, _ := c.dbDumperManager.GetNameSuffix()
	instanceName := prefix + instanceNameArg + suffix
	err = c.dbDumperManager.CheckIsDbDumperInstance(instanceName)
	if err != nil {
		instanceName = instanceNameArg
		err = c.dbDumperManager.CheckIsDbDumperInstance(instanceName)
		checkError(err)
	}
	return instanceName
}
func (c *BasicPlugin) deleteCommand(cg *cli.Context) {
	if len(cg.Args()) > 0 {
		instanceName := c.getInstanceName(cg.Args().First())
		err := c.dbDumperManager.DeleteDump(instanceName, force)
		checkError(err)
	} else {
		err := c.dbDumperManager.DeleteDump("", force)
		checkError(err)
	}

}
func (c *BasicPlugin) restoreCommand(cg *cli.Context) {
	if len(cg.Args()) == 0 {
		checkError(errors.New("you must provide a service name or an url to a target database"))
	}
	instanceName := ""
	prefix, err := c.dbDumperManager.GetNamePrefix()
	checkError(err)
	suffix, _ := c.dbDumperManager.GetNameSuffix()
	if sourceInstance != "" {
		instanceName = prefix + sourceInstance + suffix
		err = c.dbDumperManager.CheckIsDbDumperInstance(instanceName)
	}
	if err != nil && instanceName != "" {
		instanceName = sourceInstance
		err = c.dbDumperManager.CheckIsDbDumperInstance(instanceName)
		checkError(err)
	}
	err = c.dbDumperManager.RestoreDump(cg.Args().First(), recent, instanceName, org, space, seeAllDumps, tags)
	checkError(err)
}
func (c *BasicPlugin) createCommand(cg *cli.Context) {
	if len(cg.Args()) == 0 {
		checkError(errors.New("you must provide a service name or an url to a database"))
	}
	err := c.dbDumperManager.CreateDump(cg.Args().First(), plan, tags)
	checkError(err)
}
func (c *BasicPlugin) targetCommand(cg *cli.Context) {
	if len(cg.Args()) == 0 {
		checkError(errors.New("you must provide a db-dumper-service service name"))
	}
	c.config.Target = cg.Args().First()
	registerConfig(c.config)
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
func registerConfig(config *model.Config) {
	data, err := json.Marshal(config)
	checkError(err)
	ioutil.WriteFile(getConfigFileLocation(), data, 0644)
}
func configFileExist() bool {
	if _, err := os.Stat(getConfigFileLocation()); os.IsNotExist(err) {
		return false
	}
	return true
}
func retrieveConfig() *model.Config {
	data, err := ioutil.ReadFile(getConfigFileLocation())
	checkError(err)
	config := &model.Config{}
	json.Unmarshal(data, config)
	checkError(err)
	return config
}
func getConfigFileLocation() string {
	userDir, err := homedir.Dir()
	checkError(err)
	return path.Join(userDir, configFile)
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
