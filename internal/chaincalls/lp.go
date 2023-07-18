package chaincalls

import (
	"errors"
	"math/big"
	"pool_calc/internal/helpers"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	reservesOutputLen = 192 // length of reserves output string
	hexBase           = 16  // numbers stored in HEX
)

var (
	ErrInvalidReserveOutput = errors.New("reserves output len != 192")
)

// Returns token0, token1 and factory of the given pair address.
func PairInfo(client *ethclient.Client, pair common.Address) ([]common.Address, error) {
	// "0dfe1681": "token0()"
	// "d21220a7": "token1()"
	// "c45a0155": "factory()"
	token0, err := SimpleCall(client, "0dfe1681", pair)
	if err != nil {
		return nil, err
	}

	token1, err := SimpleCall(client, "d21220a7", pair)
	if err != nil {
		return nil, err
	}

	factory, err := SimpleCall(client, "c45a0155", pair)
	if err != nil {
		return nil, err
	}

	pairInfo := []common.Address{
		common.HexToAddress(token0),
		common.HexToAddress(token1),
		common.HexToAddress(factory),
	}

	return pairInfo, nil
}

// Returns the reserves for tokenIn and tokenOut in this order,
// instead of returning reserve0 and reserve1.
func SortedReserves(client *ethclient.Client, pair, tokenIn, tokenOut common.Address) (*big.Int, *big.Int, error) {
	reserve0, reserve1, err := Reserves(client, pair)
	if err != nil {
		return nil, nil, err
	}

	tk0sorted, _ := helpers.SortAddresses(tokenIn, tokenOut)
	if tk0sorted == tokenIn {
		return reserve0, reserve1, nil
	}

	return reserve1, reserve0, nil
}

// Calls pair (liquid pool) contract and returns its reserves.
func Reserves(client *ethclient.Client, pair common.Address) (*big.Int, *big.Int, error) {
	var (
		reserve0 = new(big.Int)
		reserve1 = new(big.Int)
	)

	reserves, err := SimpleCall(client, "0902f1ac", pair)
	if err != nil {
		return nil, nil, err
	}

	if len(reserves) != reservesOutputLen {
		return reserve0, reserve1, ErrInvalidReserveOutput
	}

	// Reserves output example (in one line):
	// 00000000000000000000000000000000000000000000001796b195b1ffe1701e
	// 00000000000000000000000000000000000000000000000ce5b145f8fa5dc02c
	// 0000000000000000000000000000000000000000000000000000000061d71c29
	reserve0.SetString(reserves[0:64], hexBase)
	reserve1.SetString(reserves[64:128], hexBase)

	return reserve0, reserve1, nil
}
