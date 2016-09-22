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
	"strings"
	"bytes"
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
var version_minor int = 4
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
var force bool
var inStdout bool
var dumpNumber string
var verboseMode bool
var commandPluginHelpUsage = `   {{.HelpName}} command{{if .Flags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if .Flags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

const (
	CF_HOME_KEY = "CF_HOME"
	SUFFIX_COMMAND = "-dump"
	APP_NAME = "db-dumper"
)

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command basic-plugin-command
	if !strings.HasSuffix(args[0], SUFFIX_COMMAND) {
		return
	}
	args[0] = strings.TrimSuffix(args[0], SUFFIX_COMMAND)
	app := cli.NewApp()
	if !configFileExist() {
		registerConfig(&model.Config{
			Target: db_dumper.SERVICE_NAME,
		})
	}
	c.config = retrieveConfig()
	c.dbDumperManager = db_dumper.NewDbDumperManager(c.config.Target, cliConnection, &verboseMode)
	app.Name = APP_NAME
	app.Version = fmt.Sprintf("%d.%d.%d", version_major, version_minor, version_build)
	app.Usage = "Help you to manipulate db-dumper service"
	app.Flags = []cli.Flag{}
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}
	app.Commands = generateCommand(c)
	finalArgs := append([]string{APP_NAME}, args...)
	app.Run(finalArgs)
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
		err := c.dbDumperManager.DownloadDumpFromInstanceName(instanceName, recent, inStdout, original, dumpNumber, seeAllDumps, tags)
		checkError(err)
	} else {
		err := c.dbDumperManager.DownloadDump(recent, inStdout, original, dumpNumber, seeAllDumps, tags)
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
	err = c.dbDumperManager.RestoreDump(cg.Args().First(), recent, instanceName, org, space, seeAllDumps, tags, force)
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

	pluginCommands := []plugin.Command{}
	commands := generateCommand(c)
	for _, command := range commands {
		bufferedHelp := new(bytes.Buffer)
		cli.HelpPrinter(bufferedHelp, commandPluginHelpUsage, command)
		pluginCommands = append(pluginCommands, plugin.Command{
			Name:     command.Name + SUFFIX_COMMAND,
			HelpText: command.Usage,
			UsageDetails: plugin.Usage{
				Usage: bufferedHelp.String(),
			},
		})
	}
	return plugin.PluginMetadata{
		Name: APP_NAME,
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
		Commands: pluginCommands,
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
	cfHome := os.Getenv(CF_HOME_KEY)
	if cfHome != "" {
		userDir = cfHome
	}
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
