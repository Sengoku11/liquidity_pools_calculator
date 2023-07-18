package formulas_test

import (
	"math/big"
	"pool_calc/internal/formulas"
	"testing"

	"github.com/stretchr/testify/require"
)

// https://etherscan.io/address/0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852#readContract
// [ getReserves method Response ]
//   _reserve0   uint112 :  15577162790026115051150
//   _reserve1   uint112 :  29638297954875
//	 amountOut   uint256 :  295082

func TestAmountOut(t *testing.T) {
	t.Parallel()

	amountIn, _ := new(big.Int).SetString("155554778123672", 10)
	reserveIn, _ := new(big.Int).SetString("15577162790026115051150", 10)
	reserveOut, _ := new(big.Int).SetString("29638297954875", 10)
	expectedAmountOut, _ := new(big.Int).SetString("295082", 10)
	expectedAmountOut2, _ := new(big.Int).SetString("295230", 10)

	tests := []struct {
		name              string
		amountIn          *big.Int
		reserveIn         *big.Int
		reserveOut        *big.Int
		expectedAmountOut *big.Int
		swapFee           int64
		err               error
	}{
		{
			name:       "nil inputs",
			amountIn:   nil,
			reserveIn:  nil,
			reserveOut: nil,
			swapFee:    997,
			err:        formulas.ErrNilInput,
		}, {
			name:       "zero amountIn",
			amountIn:   big.NewInt(0),
			reserveIn:  reserveIn,
			reserveOut: reserveOut,
			swapFee:    997,
			// expectedAmountOut: nil,
			err: formulas.ErrZeroAmountIn,
		}, {
			name:       "zero reserveIn",
			amountIn:   amountIn,
			reserveIn:  big.NewInt(0),
			reserveOut: reserveOut,
			swapFee:    997,
			// expectedAmountOut: nil,
			err: formulas.ErrZeroReserve,
		}, {
			name:              "valid amountOut",
			amountIn:          amountIn,
			reserveIn:         reserveIn,
			reserveOut:        reserveOut,
			swapFee:           997,
			expectedAmountOut: expectedAmountOut,
		}, {
			name:              "different swapFee",
			amountIn:          amountIn,
			reserveIn:         reserveIn,
			reserveOut:        reserveOut,
			swapFee:           9975,
			expectedAmountOut: expectedAmountOut2,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			amountOut, err := formulas.AmountOut(tt.amountIn, tt.reserveIn, tt.reserveOut, tt.swapFee)
			if err != nil {
				require.ErrorAs(t, tt.err, &err)

				return
			}
			require.EqualValues(t, tt.expectedAmountOut, amountOut)
			require.NoError(t, err)
		})
	}
}
