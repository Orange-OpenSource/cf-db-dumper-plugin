#!/usr/bin/env bash
CURRENTDIR=`pwd`
PLUGIN_PATH="$CURRENTDIR/out/db-dumper"

$CURRENTDIR/bin/build
cf uninstall-plugin db-dumper
cf install-plugin "$PLUGIN_PATH" -f