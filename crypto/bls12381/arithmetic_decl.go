// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser eral Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser eral Public License for more details.
//
// You should have received a copy of the GNU Lesser eral Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bls12381

var isX86CharacteristicSet bool = false
var NoADX = false
var Fallback = true

func Init() {
	if NoADX {
		mul = mulNoADX
		mulAssign = mulAssignNoADX
	} else if Fallback {
		mul = Fallbackmul
		mulAssign = FallbackmulAssign
		add = Fallbackadd
		addAssign = FallbackaddAssign
		ladd = Fallbackladd
		laddAssign = FallbackladdAssign
		double = Fallbackdouble
		doubleAssign = FallbackdoubleAssign
		ldouble = Fallbackldouble
		sub = Fallbacksub
		subAssign = FallbacksubAssign
		lsubAssign = FallbacklsubAssign
		_neg = Fallbackneg
	}
}

// Use ADX backend for default
var mul func(c, a, b *fe) = mulADX
var mulAssign func(a, b *fe) = mulAssignADX
var add func(a, b, c *fe) = addX
var addAssign func(a, b *fe) = addAssignX
var ladd func(c, a, b *fe) = laddX
var laddAssign func(c, a *fe) = laddAssignX
var double func(c, a *fe) = doubleX
var doubleAssign func(a *fe) = doubleAssignX
var ldouble func(c, a *fe) = ldoubleX
var sub func(a, b, c *fe) = subX
var subAssign func(a, b *fe) = subAssignX
var lsubAssign func(c, a *fe) = lsubAssignX

var _neg func(c, a *fe) = _negX

func square(c, a *fe) {
	mul(c, a, a)
}

func neg(c, a *fe) {
	if a.isZero() {
		c.set(a)
	} else {
		_neg(c, a)
	}
}

//go:noescape
func addX(c, a, b *fe)

//go:noescape
func addAssignX(a, b *fe)

//go:noescape
func laddX(c, a, b *fe)

//go:noescape
func laddAssignX(a, b *fe)

//go:noescape
func doubleX(c, a *fe)

//go:noescape
func doubleAssignX(a *fe)

//go:noescape
func ldoubleX(c, a *fe)

//go:noescape
func subX(c, a, b *fe)

//go:noescape
func subAssignX(a, b *fe)

//go:noescape
func lsubAssignX(a, b *fe)

//go:noescape
func _negX(c, a *fe)

//go:noescape
func mulNoADX(c, a, b *fe)

//go:noescape
func mulAssignNoADX(a, b *fe)

//go:noescape
func mulADX(c, a, b *fe)

//go:noescape
func mulAssignADX(a, b *fe)
