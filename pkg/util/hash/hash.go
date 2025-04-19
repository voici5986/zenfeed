// Copyright (C) 2025 wangyusong
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package hash

import "hash/fnv"

func Sum64(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	
	return h.Sum64()
}

func Sum64s(ss []string) uint64 {
	h := fnv.New64a()
	for _, s := range ss {
		h.Write([]byte(s))
		h.Write([]byte{0})
	}

	return h.Sum64()
}
