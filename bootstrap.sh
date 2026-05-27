#!/bin/bash
cd $(dirname $0)/$1
air -c ../.air.toml
