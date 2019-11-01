// Copyright 2014 The go-ethereum Authors
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

package common

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	checker "gopkg.in/check.v1"
)

type BytesSuite struct{}

var _ = checker.Suite(&BytesSuite{})

func (s *BytesSuite) TestCopyBytes(c *checker.C) {
	data1 := []byte{1, 2, 3, 4}
	exp1 := []byte{1, 2, 3, 4}
	res1 := CopyBytes(data1)
	c.Assert(res1, checker.DeepEquals, exp1)
}

func (s *BytesSuite) TestLeftPadBytes(c *checker.C) {
	val1 := []byte{1, 2, 3, 4}
	exp1 := []byte{0, 0, 0, 0, 1, 2, 3, 4}

	res1 := LeftPadBytes(val1, 8)
	res2 := LeftPadBytes(val1, 2)

	c.Assert(res1, checker.DeepEquals, exp1)
	c.Assert(res2, checker.DeepEquals, val1)
}

func (s *BytesSuite) TestRightPadBytes(c *checker.C) {
	val := []byte{1, 2, 3, 4}
	exp := []byte{1, 2, 3, 4, 0, 0, 0, 0}

	resstd := RightPadBytes(val, 8)
	resshrt := RightPadBytes(val, 2)

	c.Assert(resstd, checker.DeepEquals, exp)
	c.Assert(resshrt, checker.DeepEquals, val)
}

func TestFromHex(t *testing.T) {
	input := "0x01"
	expected := []byte{1}
	result := FromHex(input)
	if !bytes.Equal(expected, result) {
		t.Errorf("Expected %x got %x", expected, result)
	}
}

func TestIsHex(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"", true},
		{"0", false},
		{"00", true},
		{"a9e67e", true},
		{"A9E67E", true},
		{"0xa9e67e", false},
		{"a9e67e001", false},
		{"0xHELLO_MY_NAME_IS_STEVEN_@#$^&*", false},
	}
	for _, test := range tests {
		if ok := isHex(test.input); ok != test.ok {
			t.Errorf("isHex(%q) = %v, want %v", test.input, ok, test.ok)
		}
	}
}

func TestFromHexOddLength(t *testing.T) {
	input := "0x1"
	expected := []byte{1}
	result := FromHex(input)
	if !bytes.Equal(expected, result) {
		t.Errorf("Expected %x got %x", expected, result)
	}
}

func TestNoPrefixShortHexOddLength(t *testing.T) {
	input := "1"
	expected := []byte{1}
	result := FromHex(input)
	if !bytes.Equal(expected, result) {
		t.Errorf("Expected %x got %x", expected, result)
	}
}

func TestToHex(t *testing.T) {
	testBytes := []byte{
		105, 24, 76, 225, 150, 125, 28, 144, 68, 17, 185, 70, 162, 62, 105, 42,
		16, 46, 238, 27, 148, 229, 81, 36, 136, 115, 27, 151, 68, 77, 195, 216,
	}
	strHex := "0x69184ce1967d1c904411b946a23e692a102eee1b94e5512488731b97444dc3d8"
	assert.Equal(t, strHex, ToHex(testBytes))

	testBytes = []byte{}
	strHex = "0x0"
	assert.Equal(t, strHex, ToHex(testBytes))
}

func TestCopyBytes(t *testing.T) {
	originBytes := []byte{
		105, 24, 76, 225, 150, 125, 28, 144, 68, 17, 185, 70, 162, 62, 105, 42,
		16, 46, 238, 27, 148, 229, 81, 36, 136, 115, 27, 151, 68, 77, 195, 216,
	}
	copiedBytes := CopyBytes(originBytes)
	assert.Equal(t, originBytes, copiedBytes)

	originBytes = nil
	copiedBytes = CopyBytes(originBytes)
	assert.Equal(t, originBytes, copiedBytes)

}

func TestHex2BytesFixed(t *testing.T) {
	str := "0x69184ce1967d1c904411b946a23e692a102eee1b94e5512488731b97444dc3d8"
	flen := len(str)
	fmt.Println(Hex2BytesFixed(str, flen))
}
