package formulas

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"pool_calc/internal/helpers"
)

var (
	ErrNilInput     = errors.New("input is nil")
	ErrZeroAmountIn = errors.New("insufficient input amount, contract will revert")
	ErrZeroReserve  = errors.New("insufficient liquidity, contract will revert")
)

// Calculates the amountOut based on the given amountIn and
// non-zero reserves in a liquid pool.
func AmountOut(amountIn, reserveIn, reserveOut *big.Int, swapFee int64) (*big.Int, error) {
	switch {
	case amountIn == nil:
		return nil, fmt.Errorf("amountIn is nil: %w", ErrNilInput)
	case reserveIn == nil:
		return nil, fmt.Errorf("reserveIn is nil: %w", ErrNilInput)
	case reserveOut == nil:
		return nil, fmt.Errorf("reserveOut is nil: %w", ErrNilInput)
	}

	zeroInt := big.NewInt(0)
	if amountIn.Cmp(zeroInt) == 0 {
		return nil, ErrZeroAmountIn
	}

	if reserveIn.Cmp(zeroInt) == 0 {
		return nil, fmt.Errorf("reserveIn is 0: %w", ErrZeroReserve)
	}

	if reserveOut.Cmp(zeroInt) == 0 {
		return nil, fmt.Errorf("reserveOut is 0: %w", ErrZeroReserve)
	}

	// 	uint amountInWithFee = amountIn.mul(997);
	// 	uint numerator = amountInWithFee.mul(reserveOut);
	// 	uint denominator = reserveIn.mul(1000).add(amountInWithFee);
	// 	amountOut = numerator / denominator;

	digits := helpers.Digits(int(swapFee))
	feeBase := math.Pow10(digits)

	amountInWithFee := new(big.Int)
	amountInWithFee.Mul(amountIn, big.NewInt(swapFee))

	numerator := new(big.Int)
	numerator.Mul(amountInWithFee, reserveOut)

	denominator := new(big.Int)
	denominator.Mul(reserveIn, big.NewInt(int64(feeBase))).Add(denominator, amountInWithFee)

	amountOut := numerator.Div(numerator, denominator)

	return amountOut, nil
}
