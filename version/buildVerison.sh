#!/bin/sh

dist=version/version.go

rm -f $dist

echo "package version\n"            >> $dist
echo "const version = \"$1\"\n"     >> $dist

echo "echo project version to version.go success"
