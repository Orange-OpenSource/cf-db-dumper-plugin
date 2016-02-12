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
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --name, -n "db-dumper-service"	set the service name of your db-dumper-service
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
   db-dumper restore [service-name-or-url-of-your-db]
```

### Delete

```
NAME:
   db-dumper delete - Delete a instance and all his dumps (dumps can be retrieve during a period)

USAGE:
   db-dumper delete [arguments...]
```