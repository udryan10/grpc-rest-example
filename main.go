package main

import (
	"fmt"
	"os"

	"github.com/udryan10/grpc-rest-example/cmd"
)

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

}
