# db-dumper-cli-plugin

Cloud foundry cli plugin to use [db-dumper-service](https://github.com/Orange-OpenSource/db-dumper-service) in more convenient way.

## Available commands

```
NAME:
   db-dumper - Help you to manipulate db-dumper service

USAGE:
   db-dumper [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   create, c	Create a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)
   restore, r	Restore a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)
   delete, d	Delete a instance and all his dumps (dumps can be retrieve during a period)
   list, l	List dumps for an instance
   download, dl	Download a dump to your drive
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --name, -n "db-dumper-service"	Help you to manipulate db-dumper service
   --verbose, --vvv, --vv		Set the flag if you want to see all output from cloudfoundry cli
   --help, -h				show help
   --version, -v			print the version
```

### Create

```
NAME:
   db-dumper create - Create a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)

USAGE:
   db-dumper create [service-name-or-url-of-your-db]
```

### Restore

```
NAME:
   db-dumper restore - Restore a dump from a database service or database uri (e.g: mysql://admin:admin@mybase.com:3306/mysuperdb)

USAGE:
   db-dumper restore [command options] [service-name-or-url-of-your-db]

OPTIONS:
   --recent	Restore from the most recent dump
```

### Delete

```
NAME:
   db-dumper delete - Delete a instance and all his dumps (dumps can be retrieve during a period)

USAGE:
   db-dumper delete [arguments...]
```

### List

```
NAME:
   db-dumper list - List dumps for an instance

USAGE:
   db-dumper list [command options] [service instance](optionnal)

OPTIONS:
   --show-url, -s	If you want to see download url and dashboard url
```

### Download

```
NAME:
   db-dumper download - Download a dump to your drive

USAGE:
   db-dumper download [command options] [service instance](optionnal)

OPTIONS:
   --skip-ssl-validation, -k	Skip the ssl validation (for self-signed certificate mainly)
   --recent			Download from the most recent dump
   --dump-number, -p 		Download from the number showed when using 'db-dumper list'
   --stdout, -o			Show file directly in stdout (service instance is no more optionnal and you need to use flag --dump-number or --recent)
```

### Open

```
NAME:
   db-dumper open - Open dump in your browser

USAGE:
   db-dumper open [command options] [service instance](optionnal)

OPTIONS:
   --recent		Open dump page from the most recent dump
   --dump-number, -p 	Open from the number showed when using 'db-dumper list'
```