package main

import "github.com/christopher-fentiman/dbx/cmd"

func main() {
	cmd.ScaffoldFS = ScaffoldFS
	cmd.Execute()
}
