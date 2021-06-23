package miner

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func getGasPriceInt() *big.Int {
	cases := []*big.Int{
		new(big.Int).Lsh(common.Big1, 256),
		new(big.Int).Add(new(big.Int).Lsh(common.Big1, 256), big.NewInt(1)),
		new(big.Int).Sub(new(big.Int).Lsh(common.Big1, 256), big.NewInt(1)),
		nil,
		common.Big0,
		common.Big1,
	}
	// Add a lot of "good" gas prices
	for i := 0; i < 50; i++ {
		cases = append(cases,
			big.NewInt(rand.Int63n(1000000000000000000)), // rand up to 1 ETH
		)
	}
	return cases[rand.Intn(len(cases))]
}

func getToField() *common.Address {
	var (
		ADDR   = "0xb02A2EdA1b317FBd16760128836B0Ac59B560e9D"
		myself = common.HexToAddress(ADDR)
	)
	cases := []*common.Address{
		&myself, // sender itself
		nil,     // contract creation
		{},
	}
	// Precompiles
	for i := 0; i < 1024; i++ {
		addr := common.BigToAddress(big.NewInt(int64(i)))
		cases = append(cases, &addr)
	}
	return cases[rand.Intn(len(cases))]
}

func getChainID() *big.Int {
	/*
		switch i := rand.Intn(100); {
		case i == 99:
			return getGasPriceInt()
		case i < 5:
			return big.NewInt(int64(i))
		case i < 10:
			return params.MainnetChainConfig.ChainID
		default:
			return big.NewInt(1337)
			//return params.CalaverasChainConfig.ChainID
		}*/
	return big.NewInt(1337)
}

func getGas() uint64 {
	switch i := rand.Intn(100); {
	case i == 0:
		return 0
	case i == 1:
		return uint64(^uint(0))
	case i < 10:
		return 21000
	default:
		return uint64(rand.Int63n(10000000))
	}
}

func getSigner() types.Signer {
	chainid := getChainID()
	switch i := rand.Intn(100); {
	case i == 0:
		return types.HomesteadSigner{}
	case i == 1:
		return types.FrontierSigner{}
	case i == 2:
		return types.NewEIP155Signer(chainid)
	case i == 1:
		return types.NewEIP2930Signer(chainid)
	default:
		return types.NewLondonSigner(chainid)
	}
}
