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

package vector

import (
	"math"
)

func Quantize(vec []float32) (quantized []int8, min, scale float32) {
	// Find the minimum and maximum values.
	min, max := float32(math.MaxFloat32), float32(-math.MaxFloat32)
	for _, v := range vec {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// Calculate the quantization scale.
	scale = float32(255) / (max - min)

	// Quantize the data.
	quantized = make([]int8, len(vec))
	for i, v := range vec {
		quantized[i] = int8(math.Round(float64((v-min)*scale - 128)))
	}

	return quantized, min, scale
}

func Dequantize(quantized []int8, min, scale float32) []float32 {
	vec := make([]float32, len(quantized))
	for i, v := range quantized {
		vec[i] = (float32(v)+128)/scale + min
	}

	return vec
}
