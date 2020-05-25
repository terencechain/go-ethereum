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

	"github.com/ethereum/go-ethereum/crypto/bls12381"
)

func Fuzz(data []byte) int {
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

	cpy := make([]byte, len(data))
	copy(cpy, data)
	for i, precompile := range precompiles {
		bls12381.NoADX = false
		bls12381.Fallback = false
		gas1 := precompile.RequiredGas(cpy)
		out, err := precompile.Run(cpy)
		if err == nil {
			promote = true
			switch i {
			case 0, 1, 2, 6:
				if len(out) != 128 {
					panic(fmt.Sprintf("precomp %d: Output had strange length: %v %d", i, out, len(out)))
				}
			case 3, 4, 5, 7:
				if len(out) != 256 {
					panic(fmt.Sprintf("precomp %d: Output had strange length: %v %d", i, out, len(out)))
				}
			case 8:
				if len(out) != 32 {
					panic(fmt.Sprintf("precomp %d: Output had strange length: %v %d", i, out, len(out)))
				}
			}
		}
		bls12381.NoADX = true
		bls12381.Fallback = false
		gas2 := precompile.RequiredGas(cpy)
		out2, err2 := precompile.Run(cpy)
		if err != nil && err.Error() != err2.Error() {
			panic(fmt.Sprintf("precomp %d: errors not equal %v %v ", i, err, err2))
		}
		if gas1 != gas2 {
			panic(fmt.Sprintf("precomp %d: gas not equal %v %v ", i, gas1, gas2))
		}
		if !bytes.Equal(out, out2) {
			panic(fmt.Sprintf("precomp %d: output not equal %v %v ", i, out, out2))
		}
		bls12381.NoADX = false
		bls12381.Fallback = true
		gas3 := precompile.RequiredGas(cpy)
		out3, err3 := precompile.Run(cpy)
		if err != nil && err.Error() != err2.Error() {
			panic(fmt.Sprintf("precomp %d: fallback: errors not equal %v %v ", i, err, err3))
		}
		if gas1 != gas3 {
			panic(fmt.Sprintf("precomp %d: fallback: gas not equal %v %v ", i, gas1, gas3))
		}
		if !bytes.Equal(out, out3) {
			panic(fmt.Sprintf("precomp %d: fallback: output not equal %v %v ", i, out, out3))
		}
	}
	if !bytes.Equal(data, cpy) {
		panic(fmt.Sprintf("someone modified data: %v %v", data, cpy))
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
