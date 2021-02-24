package main

import (
	"flag"
	"fmt"

	"github.com/chatto-extensions/weatherapi/internal/ext"
	"github.com/chatto-extensions/weatherapi/internal/version"
	"github.com/jaimeteb/chatto/extension"
)

func main() {
	vers := flag.Bool("version", false, "Display version.")
	flag.Parse()

	if *vers {
		fmt.Println(version.Build())
		return
	}

	extension.ServeREST(ext.RegisteredFuncs)
}
