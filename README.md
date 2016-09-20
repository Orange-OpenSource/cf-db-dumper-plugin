# db-dumper-cli-plugin

Cloud foundry cli plugin to use [db-dumper-service](https://github.com/Orange-OpenSource/db-dumper-service) in more convenient way.

## Available commands

```
NAME:
   db-dumper - Help you to manipulate db-dumper service

USAGE:
   db-dumper [global options] command [command options] [arguments...]

VERSION:
   1.3.0

COMMANDS:
   target, t   	Target a db-dumper service
   create, c   	Create a dump from a database service (e.g.: mydb) or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)
   restore, r  	Restore a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)
   delete, d   	Delete a instance and all its dumps (dumps can be retrieve during a period)
   list, l     	List dumps for an instance
   download, dl	Download a dump to your drive
   open, o     	Open dump in your browser
   help, h     	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, --vvv, --vv      	Set the flag if you want to see all output from cloudfoundry cli
   --help, -h  			show help
   --version, -v       		print the version
```

### Create

```
NAME:
   db-dumper create - Create a dump from a database service (e.g.: mydb) or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)

USAGE:
   db-dumper create [command options] [service-name-or-url-of-your-db]

OPTIONS:
   --plan      	Choose the plan to use. If not set it will ask you to choose one from a list
   --tags      	Pass a list of tags to create a dump with these tags (e.g.: --tags=mytag1,mytag2...)
```

### Restore

```
NAME:
   db-dumper restore - Restore a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)

USAGE:
   db-dumper restore [command options] [service-name-or-url-of-your-db]

OPTIONS:
   --recent    		Restore from the most recent dump
   --see-all-dumps     	Restore by selecting one of dumps available for the targetted database (made by all db-dumper service instance linked to this database)
   --tags      		Restore by selecting one of dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)
   --org, -o   		Org which contain the service name passed
   --space, -s 		Space which contain the service name passed
   --source-instance   	The db dumper service instance where dumps should be retrieved (this can be the service instance you passed in create e.g.: mydb)
```

### Delete

```
NAME:
   db-dumper delete - Delete a instance and all its dumps (dumps can be retrieve during a period)

USAGE:
   db-dumper delete [command options] [service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)

OPTIONS:
   --force, -f 	Force deletion without confirmation
```

### List

```
NAME:
   db-dumper list - List dumps for an instance

USAGE:
   db-dumper list [command options] [service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)

OPTIONS:
   --show-url, -s      	If you want to see download url and dashboard url
   --see-all-dumps     	See all dumps for the database (made by all db-dumper service instance linked to this database)
   --tags      		See dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)
```

### Download

```
NAME:
   db-dumper download - Download a dump to your drive

USAGE:
   db-dumper download [command options] [service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)

OPTIONS:
   --skip-ssl-validation, -k   	Skip the ssl validation (for self-signed certificate mainly)
   --recent    			Download from the most recent dump
   --original  			Download the original file (e.g.: download directly sql file instead of compressed file)
   --dump-number, -p   		Download from the number showed when using 'db-dumper list'
   --stdout, -o			Show file directly in stdout (service instance is no more optionnal and you need to use flag --dump-number or --recent)
   --see-all-dumps     		Download dumps by selecting one of dumps available for the targetted database (made by all db-dumper service instance linked to this database)
   --tags      			Download dumps by selecting one of dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)
```

### Open

```
NAME:
   db-dumper open - Open dump in your browser

USAGE:
   db-dumper open [command options] [service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)

OPTIONS:
   --recent    		Open dump page from the most recent dump
   --dump-number, -p   	Open from the number showed when using 'db-dumper list'
   --see-all-dumps     	Show dumps by selecting one of dumps available for the targetted database (made by all db-dumper service instance linked to this database)
   --tags      		Show dumps by selecting one of dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)
```

### Target

```
NAME:
   db-dumper target - Target a db-dumper service

USAGE:
   db-dumper target [db-dumper-service]

DESCRIPTION:
   Pass the db-dumper service name, default is db-dumper-service (e.g.: db-dumper-service-dev)
```