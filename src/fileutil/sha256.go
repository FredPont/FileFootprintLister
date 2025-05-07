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
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func calcSHA256(file string) string {

	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Println(err)
	}

	sha256 := fmt.Sprintf("%x", h.Sum(nil))
	return sha256
}
