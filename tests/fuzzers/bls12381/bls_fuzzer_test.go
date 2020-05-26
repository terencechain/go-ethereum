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
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/bls12381"
	"github.com/influxdata/influxdb/pkg/testing/assert"
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

func testPrec(t *testing.T, input []byte) {
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
	for _, prec := range precompiles {
		prec.RequiredGas(input)
		bls12381.NoADX = false
		bls12381.Fallback = false
		a, _ := prec.Run(input)
		bls12381.NoADX = true
		bls12381.Fallback = false
		b, _ := prec.Run(input)
		bls12381.NoADX = false
		bls12381.Fallback = true
		c, _ := prec.Run(input)
		assert.Equal(t, a, b)
		assert.Equal(t, b, c)
	}
}

func TestG2(t *testing.T) {

	dir, err := ioutil.ReadDir("crashers")
	if err != nil {
		t.Error(err)
	}
	for _, info := range dir {
		name := info.Name()
		var input []byte
		file, err := os.Open("crashers/" + name)
		if err != nil {
			t.Error(err)
		}
		if strings.HasSuffix(name, ".output") || strings.HasSuffix(name, ".quoted") {
			in, err := ioutil.ReadFile("crashers/" + name)
			assert.NoError(t, err)
			input = in
		} else {
			assert.NoError(t, binary.Read(file, binary.LittleEndian, input))
			testPrec(t, input)
			in, err := ioutil.ReadFile("crashers/" + name)
			assert.NoError(t, err)
			testPrec(t, in)
		}

	}
}

func TestString(t *testing.T) {
	input := "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
		"\x00\x00\x00\n\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00!" +
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x00"
	testPrec(t, []byte(input))
	testPrec(t, common.FromHex(input))
}
