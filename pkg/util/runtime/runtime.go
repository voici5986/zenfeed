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

package runtime

// Must panics if err is not nil.
// It is useful for handling errors in initialization code where recovery is not possible.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must1 is like Must but returns the value if err is nil.
// It is useful for handling errors in initialization code where recovery is not possible
// and a value needs to be returned.
func Must1[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	
	return v
}

// Must2 is like Must but returns two values if err is nil.
// It is useful for handling errors in initialization code where recovery is not possible
// and two values need to be returned.
func Must2[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}

	return v1, v2
}
