// Copyright 2016 The go-ethereum Authors
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
package geth

import (
	"encoding/json"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// MakeInterfacesFromJSON creates a new Interface slice that contains the
// prototype of types defined in the input json.
func MakeInterfacesFromJSON(data []byte) (Interfaces, error) {
	var field struct {
		Inputs []abi.Argument
	}
	if err := json.Unmarshal(data, &field); err != nil {
		return Interfaces{}, err
	}
	var out []interface{}
	for _, in := range field.Inputs {
		out = append(out, Interface{reflect.New(in.Type.GetType())})
	}
	return Interfaces{out}, nil
}
