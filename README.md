# Liquidity Pool Calculator
A Golang application that accepts the address of a Uniswap V2 pool in Ethereum and calculates the resulting 
output amount that the pool would return if you attempted to swap a specified input amount.

## Requirements 

You need to have access to a Go-Ethereum full node in the Ethereum network to make blockchain requests. 
You can use a free node at https://www.infura.io. 

Then put the URL path in `.env` file. See `.env.example`.

## Usage

```
./pool_calc amountOut [amountIn] [poolAddress] [tokenIn] [tokenOut] [flags]
```

You can also call the `./pool_calc amountOut --help` command for a built-in example or simply run `make example`.
