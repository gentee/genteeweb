#!/usr/local/bin/gentee

// This script builds and runs GenteeWeb server
// It uses Gentee programming language - https://github.com/gentee/gentee

run {
    str env = $ go env
    $GOPATH = RegExp(env, `GOPATH="?([^"|\n|\r]*)`)

    $ go install 
    $ ${GOPATH}/bin/genteeweb examples/genteeweb.yaml
}