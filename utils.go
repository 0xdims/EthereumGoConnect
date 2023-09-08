package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func convertType(t abi.Type, value string) (interface{}, error) {
	var returnOutput []byte
	switch t.T {
	case abi.IntTy, abi.UintTy:
		bigint, _ := big.NewInt(0).SetString(value, 10)
		returnOutput = bigint.Bytes()
		return abi.ReadInteger(t, returnOutput)
	case abi.BoolTy:
		boolValue, _ := strconv.ParseBool(value)
		bytes := make([]byte, 1)
		if boolValue {
			bytes[0] = 1
		} else {
			bytes[0] = 0
		}
		returnOutput = bytes
		return readBool(returnOutput)
	case abi.AddressTy:
		addressValue := common.HexToAddress(value)
		returnOutput = addressValue.Bytes()
		return common.BytesToAddress(returnOutput), nil
	case abi.HashTy:
		hashValue := common.HexToHash(value)
		returnOutput = hashValue.Bytes()
		return common.BytesToHash(returnOutput), nil
	default:
		return nil, fmt.Errorf("abi: unknown type %v", t.T)
	}
}

// readBool reads a bool.
func readBool(word []byte) (bool, error) {
	for _, b := range word[:31] {
		if b != 0 {
			return false, fmt.Errorf("bad bool")
		}
	}
	switch word[31] {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("bad bool")
	}
}
