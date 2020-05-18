// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bls

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func Fuzz(data []byte) int {
	data = common.FromHex(string(data))
	promote := false
	precompiles := []precompile{
		new(bls12381G1Add),
		new(bls12381G1Mul),
		new(bls12381G1MultiExp),
		new(bls12381G2Add),
		new(bls12381G2Mul),
		new(bls12381G2MultiExp),
		new(bls12381MapG1),
		new(bls12381MapG2),
		new(bls12381Pairing),
	}
	for _, precompile := range precompiles {
		precompile.RequiredGas(data)
		out, err := precompile.Run(data)
		if err == nil {
			promote = true
			if len(out) != 128 && len(out) != 256 {
				panic(fmt.Sprintf("Output had strange length: %v %d", out, len(out)))
			}
		}
	}
	if promote {
		return 1
	}
	return 0
}

type precompile interface {
	RequiredGas(input []byte) uint64
	Run(input []byte) ([]byte, error)
}
