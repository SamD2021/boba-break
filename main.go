/*
Copyright © 2024 Samuel Dasilva
*/
package main

import (
	"github.com/SamD2021/boba-break/cmd"
	"os"
)

func main() {
	cmd.Execute()
	_, err := os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
	if err != nil {
		return
	}
}
