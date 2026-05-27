#!/bin/bash
cd $(dirname $0)/$1
goimports -w .
air -c ../.air.toml
