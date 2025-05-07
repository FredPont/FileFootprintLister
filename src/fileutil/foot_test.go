/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2024
*/

package fileutil

import (
	"fmt"
	"testing"
)

func TestCalcMD5(t *testing.T) {
	tests := []struct {
		file string
		want string
	}{
		{"testfile1.txt", "826e8142e6baabe8af779f5f490cf5f5"},
		//{"testfile1.txt", "826e8142e6baabe8af779f5f490cf5f4"},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := calcMD5(tc.file)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestCalSHA256(t *testing.T) {
	tests := []struct {
		file string
		want string
	}{
		{"testfile1.txt", "c147efcfc2d7ea666a9e4f5187b115c90903f0fc896a56df9a6ef5d8f3fc9f31"},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := calcSHA256(tc.file)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
