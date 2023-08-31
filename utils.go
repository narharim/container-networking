package main

import (
	"fmt"
	"os"
	"runtime"
)

func CheckRootPrivileges() {
	if os.Getuid() != 0 {
		fmt.Println("Error: user must be root!")
		os.Exit(1)
	}
}

func gracefulExitOnError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		runtime.Goexit()
	}
}
