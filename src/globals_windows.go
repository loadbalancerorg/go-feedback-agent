package main

import (
	"os"
	"path/filepath"
)

var (
	configDir string = filepath.Join(os.Getenv("ProgramData"), "LoadBalancer.org", "LoadBalancer")
)
