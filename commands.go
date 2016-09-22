package main

import (
	"github.com/codegangsta/cli"
)

func generateCommand(c *BasicPlugin) []cli.Command {
	commands := []cli.Command{
		{
			Name:      "target",
			Aliases:     []string{"t"},

			Usage:     "Target a db-dumper service",
			ArgsUsage: "[db-dumper-service]",
			Description: "Target a db-dumper service, default is db-dumper-service (e.g.: db-dumper-service-dev)",
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
				cli.BoolFlag{
					Name: "force, f",
					Usage: "Force restore without confirmation",
					Destination: &force,
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

	for index, command := range commands {
		commands[index].Flags = append(command.Flags, cli.BoolFlag{
			Name: "verbose, vvv, vv",
			Usage: "Set the flag if you want to see all output from cloudfoundry cli",
			Destination: &verboseMode,
		})
	}
	return commands
}