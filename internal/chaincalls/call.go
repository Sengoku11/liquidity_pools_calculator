package chaincalls

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// A simple contract call. Not suitable if you need to pass From and Value to the call.
func SimpleCall(client *ethclient.Client, inputData string, contractAddress common.Address) (string, error) {
	inputDataBytes, err := hex.DecodeString(inputData)
	if err != nil {
		return "", fmt.Errorf("decode string: %w", err)
	}

	msg := &ethereum.CallMsg{
		To:   &contractAddress, // the destination contract (nil for contract creation)
		Data: inputDataBytes,   // input data, usually an ABI-encoded contract method invocation
	}

	byteResponse, err := client.PendingCallContract(context.Background(), *msg)
	if err != nil {
		return "", fmt.Errorf("call contract: %w", err)
	}

	return hex.EncodeToString(byteResponse), nil
}
