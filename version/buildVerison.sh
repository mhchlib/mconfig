#!/bin/sh

dist=version/version.go

rm -f $dist


echo "package version\n"            >> $dist
echo "const version = \"$1\"\n"     >> $dist
echo "func GetVersion() string {"   >> $dist
echo "	return version"             >> $dist
echo "}"                            >> $dist

echo "echo project version to version.go success"
