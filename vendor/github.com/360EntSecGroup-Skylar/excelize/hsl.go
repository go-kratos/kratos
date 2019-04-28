/*
Copyright (c) 2012 Rodrigo Moraes. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

	 * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
	 * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
	 * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package excelize

import (
	"image/color"
	"math"
)

// HSLModel converts any color.Color to a HSL color.
var HSLModel = color.ModelFunc(hslModel)

// HSL represents a cylindrical coordinate of points in an RGB color model.
//
// Values are in the range 0 to 1.
type HSL struct {
	H, S, L float64
}

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
// for the HSL.
func (c HSL) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := HSLToRGB(c.H, c.S, c.L)
	return uint32(r) * 0x101, uint32(g) * 0x101, uint32(b) * 0x101, 0xffff
}

// hslModel converts a color.Color to HSL.
func hslModel(c color.Color) color.Color {
	if _, ok := c.(HSL); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	h, s, l := RGBToHSL(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	return HSL{h, s, l}
}

// RGBToHSL converts an RGB triple to a HSL triple.
func RGBToHSL(r, g, b uint8) (h, s, l float64) {
	fR := float64(r) / 255
	fG := float64(g) / 255
	fB := float64(b) / 255
	max := math.Max(math.Max(fR, fG), fB)
	min := math.Min(math.Min(fR, fG), fB)
	l = (max + min) / 2
	if max == min {
		// Achromatic.
		h, s = 0, 0
	} else {
		// Chromatic.
		d := max - min
		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}
		switch max {
		case fR:
			h = (fG - fB) / d
			if fG < fB {
				h += 6
			}
		case fG:
			h = (fB-fR)/d + 2
		case fB:
			h = (fR-fG)/d + 4
		}
		h /= 6
	}
	return
}

// HSLToRGB converts an HSL triple to a RGB triple.
func HSLToRGB(h, s, l float64) (r, g, b uint8) {
	var fR, fG, fB float64
	if s == 0 {
		fR, fG, fB = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - s*l
		}
		p := 2*l - q
		fR = hueToRGB(p, q, h+1.0/3)
		fG = hueToRGB(p, q, h)
		fB = hueToRGB(p, q, h-1.0/3)
	}
	r = uint8((fR * 255) + 0.5)
	g = uint8((fG * 255) + 0.5)
	b = uint8((fB * 255) + 0.5)
	return
}

// hueToRGB is a helper function for HSLToRGB.
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}
