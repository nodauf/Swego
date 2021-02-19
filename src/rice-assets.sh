#!/bin/bash
go get github.com/GeertJohan/go.rice github.com/GeertJohan/go.rice/rice
rice embed-go -i ./src/controllers
