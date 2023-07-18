package main

import (
	"pool_calc/cmd"

	"github.com/joho/godotenv"
)

const (
	uniswapFee = 997 // 997 equals to 0.03% fee
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	rootCmd := cmd.NewRootCmd(uniswapFee)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
