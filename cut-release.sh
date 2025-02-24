#!/bin/bash -ex

go tool go-cross-compile --config-file .config/go-cross-compile.json
go tool go-github-release --github-owner dearing --github-repo mywebserver --tag-name v1.0.4
