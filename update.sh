#!/usr/bin/env bash

go mod edit -replace=gopkg.in/russross/blackfriday.v2@v2.0.1=github.com/russross/blackfriday/v2@v2.0.1

go get github.com/vannnnish/yeego
go get github.com/vannnnish/easyweb
