package main

import "github.com/chrisfentiman/dbx/cmd"

func main() {
	cmd.ScaffoldFS = ScaffoldFS
	cmd.Execute()
}
