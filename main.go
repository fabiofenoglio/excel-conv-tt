package main

import (
	"os"

	"github.com/fabiofenoglio/excelconv/cmd"
	"github.com/fabiofenoglio/excelconv/services"
)

func main() {
	l := services.GetLogger()

	if len(os.Args) == 2 {
		err := services.Run(os.Args[1])
		if err == nil {
			os.Exit(0)
		} else {
			l.Error(err.Error())
			os.Exit(1)
		}
	} else {

		cmd.Execute()
	}
}
