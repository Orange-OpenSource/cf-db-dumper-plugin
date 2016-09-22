# db-dumper-cli-plugin

Cloud foundry cli plugin to use [db-dumper-service](https://github.com/Orange-OpenSource/db-dumper-service) in more convenient way.

## Available commands

```
target-dump   	Target a db-dumper service
create-dump   	Create a dump from a database service (e.g.: mydb) or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)
restore-dump  	Restore a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)
delete-dump   	Delete a instance and all its dumps (dumps can be retrieve during a period)
list-dump       List dumps for an instance
download-dump  	Download a dump to your drive
open-dump       Open dump in your browser
```

### Create

```
NAME:
   create-dump - Create a dump from a database service (e.g.: mydb) or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)

USAGE:
   create-dump [command options] [service-name-or-url-of-your-db]

OPTIONS:
   --plan      	Choose the plan to use. If not set it will ask you to choose one from a list
   --tags      	Pass a list of tags to create a dump with these tags (e.g.: --tags=mytag1,mytag2...)
   --verbose, --vvv, --vv      	Set the flag if you want to see all output from cloudfoundry cli
```

### Restore

```
NAME:
   restore-dump - Restore a dump to a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)

USAGE:
   restore-dump [command options] [service-name-or-url-of-your-db]

OPTIONS:
   --recent    		Restore from the most recent dump
   --see-all-dumps     	Restore by selecting one of dumps available for the targetted database (made by all db-dumper service instance linked to this database)
   --tags      		Restore by selecting one of dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)
   --org, -o   		Org which contain the service name passed
   --space, -s 		Space which contain the service name passed
   --source-instance   	The db dumper service instance where dumps should be retrieved (this can be the service instance you passed in create e.g.: mydb)
   --force, -f 	Force restore without confirmation
   --verbose, --vvv, --vv      	Set the flag if you want to see all output from cloudfoundry cli
```

### Delete

```
NAME:
   delete-dump - Delete a instance and all its dumps (dumps can be retrieve during a period)

USAGE:
   delete-dump [command options] [service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)

OPTIONS:
   --force, -f 	Force deletion without confirmation
   --verbose, --vvv, --vv      	Set the flag if you want to see all output from cloudfoundry cli
```

### List

```
NAME:
   list-dump - List dumps for an instance

USAGE:
   list-dump [command options] [service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)

OPTIONS:
   --show-url, -s      	If you want to see download url and dashboard url
   --see-all-dumps     	See all dumps for the database (made by all db-dumper service instance linked to this database)
   --tags      		See dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)
   --verbose, --vvv, --vv      	Set the flag if you want to see all output from cloudfoundry cli
```

### Download

```
NAME:
   download-dump - Download a dump to your drive

USAGE:
   download-dump [command options] [service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)

OPTIONS:
   --recent    			Download from the most recent dump
   --original  			Download the original file (e.g.: download directly sql file instead of compressed file)
   --dump-number, -p   		Download from the number showed when using 'db-dumper list'
   --stdout, -o			Show file directly in stdout (service instance is no more optionnal and you need to use flag --dump-number or --recent)
   --see-all-dumps     		Download dumps by selecting one of dumps available for the targetted database (made by all db-dumper service instance linked to this database)
   --tags      			Download dumps by selecting one of dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)
   --verbose, --vvv, --vv      	Set the flag if you want to see all output from cloudfoundry cli
```

### Open

```
NAME:
   open-dump - Open dump in your browser

USAGE:
   open-dump [command options] [service instance](*optional*, this can be the service instance you passed in create e.g.: mydb)

OPTIONS:
   --recent    		Open dump page from the most recent dump
   --dump-number, -p   	Open from the number showed when using 'db-dumper list'
   --see-all-dumps     	Show dumps by selecting one of dumps available for the targetted database (made by all db-dumper service instance linked to this database)
   --tags      		Show dumps by selecting one of dumps marked with tag(s) (e.g.: --tags=mytag1,mytag2...)
   --verbose, --vvv, --vv      	Set the flag if you want to see all output from cloudfoundry cli
```

### Target

```
NAME:
   target-dump - Target a db-dumper service

USAGE:
   target-dump [db-dumper-service]

DESCRIPTION:
   Target a db-dumper service, default is db-dumper-service (e.g.: db-dumper-service-dev)

OPTIONS:
    --verbose, --vvv, --vv      	Set the flag if you want to see all output from cloudfoundry cli
```