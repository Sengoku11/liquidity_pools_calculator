package cmd

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"

	"pool_calc/internal/chaincalls"
	"pool_calc/internal/formulas"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
)

const (
	decBase            = 10 // base for amountIn
	contractLen        = 42 // length of contract address in Ethereum
	amountOutArguments = 4  // [amountIn] [poolAddress] [tokenIn] [tokenOut]
)

var (
	ErrWrongArguments  = errors.New("wrong amount of arguments")
	ErrInvalidContract = errors.New("invalid address input")
	ErrInvalidToken    = errors.New("no such token in the pair")
	ErrInvalidAmountIn = errors.New("invalid amountIn")
)

func NewRootCmd(uniswapFee int64) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "./pool_calc",
		Short: "LP Calculator",
		Long:  "A simple liquidity pool calculator",
	}

	// Добавление подкоманды
	rootCmd.AddCommand(&cobra.Command{
		Use:   "amountOut [amountIn] [poolAddress] [tokenIn] [tokenOut]",
		Short: "Calculates amountOut",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != amountOutArguments {
				return ErrWrongArguments
			}

			if (len(args[1]) + len(args[2]) + len(args[3])) != 3*contractLen {
				return ErrInvalidContract
			}

			amountIn, ok := new(big.Int).SetString(args[0], decBase)
			if !ok {
				return ErrInvalidAmountIn
			}

			if amountIn.Cmp(big.NewInt(1)) == -1 {
				return ErrInvalidAmountIn
			}

			pair := common.HexToAddress(args[1])
			tokenIn := common.HexToAddress(args[2])
			tokenOut := common.HexToAddress(args[3])

			rpcClient, err := rpc.Dial(os.Getenv("GETH_PATH"))
			if err != nil {
				panic(err)
			}

			gethClient := ethclient.NewClient(rpcClient)

			pairInfo, err := chaincalls.PairInfo(gethClient, pair)
			if err != nil {
				return fmt.Errorf("pair info: %w", err)
			}

			if tokenIn != pairInfo[0] && tokenIn != pairInfo[1] {
				return ErrInvalidToken
			}

			if tokenOut != pairInfo[0] && tokenOut != pairInfo[1] {
				return ErrInvalidToken
			}

			reserveIn, reserveOut, err := chaincalls.SortedReserves(
				gethClient,
				pair,
				tokenIn,
				tokenOut,
			)
			if err != nil {
				return fmt.Errorf("sorted reserves: %w", err)
			}

			amountOut, err := formulas.AmountOut(amountIn, reserveIn, reserveOut, uniswapFee)
			if err != nil {
				return fmt.Errorf("amountOut: %w", err)
			}

			log.Printf("Pair: %v\n", pair)
			log.Printf("ReserveIn: %v\n", reserveIn)
			log.Printf("ReserveOut: %v\n", reserveOut)
			log.Printf("AmountIn: %v\n", amountIn)
			log.Printf("AmountOut: %v\n", amountOut)

			return nil
		},
	})

	return rootCmd
}
