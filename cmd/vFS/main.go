package main

import (
	"github.com/KScaesar/IsCoolLab2024/pkg/adapters/cli"
)

func main() {
	command := cli.NewRootCommand()
	command.Execute()
}
