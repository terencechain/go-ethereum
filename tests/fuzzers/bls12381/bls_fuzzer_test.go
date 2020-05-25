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
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/bls12381"
)

func TestGenerateCorpus(t *testing.T) {
	dir, err := ioutil.ReadDir("csv")
	if err != nil {
		t.Error(err)
	}
	for _, info := range dir {
		name := info.Name()
		file, err := os.Open("csv/" + name)
		if err != nil {
			t.Error(err)
		}
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Error(err)
		}
		for i, recs := range records {
			for j, rec := range recs {
				filename := fmt.Sprintf("corpus/%v_%d_%d", name, i, j)
				err := ioutil.WriteFile(filename, []byte(rec), 0644)
				if err != nil {
					t.Error(err)
				}
			}
		}
	}
}

func TestZKCryptoVectors(t *testing.T) {

	type dat struct {
		filename string
		size     int
	}
	dats := []dat{
		{"g1_uncompressed_valid_test_vectors.dat", 96},
		{"g1_compressed_valid_test_vectors.dat", 48},
		{"g2_uncompressed_valid_test_vectors.dat", 192},
		{"g2_compressed_valid_test_vectors.dat", 96},
	}
	for _, dat := range dats {
		data, err := ioutil.ReadFile(dat.filename)
		if err != nil {
			t.Error(err)
		}
		for i := 0; i < 1000; i++ {
			vector := data[i*dat.size : (i+1)*dat.size]
			filename := fmt.Sprintf("corpus/%v_%d", dat.filename, i)
			err := ioutil.WriteFile(filename, []byte(common.ToHex(vector)), 0644)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestG2(t *testing.T) {
	input, err := ioutil.ReadFile("crashers/6bfad3af42250a2f5c551148439c85d82dbc5f46")
	if err != nil {
		t.Error(err)
	}

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
	//fmt.Println(prec.RequiredGas(common.FromHex(input)))
	for _, prec := range precompiles {
		bls12381.NoADX = false
		bls12381.Fallback = false
		prec.Run(input)
		bls12381.NoADX = true
		bls12381.Fallback = false
		prec.Run(input)
		bls12381.NoADX = false
		bls12381.Fallback = true
		prec.Run(input)
	}
}
