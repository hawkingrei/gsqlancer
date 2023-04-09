package main

import "flag"

const nmConfig = "config"

var configPath = flag.String(nmConfig, "", "config file path")

func main() {
	flag.Parse()
}
