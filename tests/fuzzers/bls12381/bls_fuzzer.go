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
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/bls12381"
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
	for i, precompile := range precompiles {
		bls12381.NoADX = false
		gas1 := precompile.RequiredGas(data)
		out, err := precompile.Run(data)
		if err == nil {
			promote = true
			if len(out) != 128 && len(out) != 256 {
				panic(fmt.Sprintf("precomp %d: Output had strange length: %v %d", i, out, len(out)))
			}
		}
		bls12381.NoADX = true
		gas2 := precompile.RequiredGas(data)
		out2, err2 := precompile.Run(data)
		if err.Error() != err2.Error() {
			panic(fmt.Sprintf("precomp %d: errors not equal %v %v ", i, err, err2))
		}
		if gas1 != gas2 {
			panic(fmt.Sprintf("precomp %d: gas not equal %v %v ", i, gas1, gas2))
		}
		if !bytes.Equal(out, out2) {
			panic(fmt.Sprintf("precomp %d: output not equal %v %v ", i, out, out2))
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
