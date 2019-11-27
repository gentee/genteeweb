#!/usr/local/bin/gentee

// This script builds and runs GenteeWeb server
// It uses Gentee programming language - https://github.com/gentee/gentee

run {
    str env = $ go env
    arr.arr.str ret &= FindRegExp(env, `GOPATH="?([^"|\n|\r]*)`)
	if ret? : $GOPATH = ret[0][1]

    $ go install 
    $ ${GOPATH}/bin/genteeweb examples/genteeweb.yaml
}