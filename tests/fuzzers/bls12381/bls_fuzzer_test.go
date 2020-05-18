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
