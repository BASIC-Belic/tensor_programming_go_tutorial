package main 

import (
	"fmt"
	"net/http"
	"os"
	"golang.org/x/net/html" // helps to better parse html 
	"golang.org/x/net/html/atom" // to break down html down into atoms
	"strings"
	log "github.com/llimlib/loglevel" // log level package to help build a log for err handling
)