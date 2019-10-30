package simdjson

import (
	"reflect"
	"testing"
	"math/bits"
)

func flatten_bits_incremental(base *[INDEX_SIZE]uint32, base_index *int, mask uint64, carried *int) {

	shifts := 0
	for {
		zeros := bits.TrailingZeros64(mask)
		if zeros == 64 {
			*carried += 64 - shifts
			return
		}
		zeros++
		(*base)[*base_index] = uint32(zeros + *carried)
		*base_index += 1
		mask = mask >> zeros
		shifts += zeros
		*carried = 0
	}
}

func TestFlattenBitsIncremental(t *testing.T) {

	testCases := []struct {
		masks    []uint64
		expected []uint32
	}{
		// Single mask
		{[]uint64{0x0},[]uint32{}},
		{[]uint64{0x11},[]uint32{0x1, 0x4}},
		{[]uint64{0x100100100100},[]uint32{0x9, 0x14-0x8, 0x20-0x14, 0x2c-0x20}},
		{[]uint64{0x100100100300},[]uint32{0x9, 0x9-0x8, 0x14-0x9, 0x20-0x14, 0x2c-0x20}},
		{[]uint64{0x8101010101010101},[]uint32{0x1, 0x8, 0x10-0x8, 0x18-0x10, 0x20-0x18, 0x28-0x20, 0x30-0x28, 0x38-0x30, 0x3f-0x38}},
		{[]uint64{0xf000000000000000},[]uint32{0x3d, 0x1, 0x1, 0x1}},
		{[]uint64{0xffffffffffffffff},[]uint32{
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
		}},
		//
		// Multiple masks
		{[]uint64{0x1, 0x1},[]uint32{0x1, 0x40}},
		{[]uint64{0x1, 0x8000000000000000},[]uint32{0x1, 0x7f}},
		{[]uint64{0x1, 0x0, 0x8000000000000000},[]uint32{0x1, 0xbf}},
		{[]uint64{0x1, 0x0, 0x0, 0x8000000000000000},[]uint32{0x1, 0xff}},
		{[]uint64{0x100100100100100, 0x100100100100100},[]uint32{0x9, 0xc, 0xc, 0xc, 0xc, 0x10, 0xc, 0xc, 0xc, 0xc}},
		{[]uint64{0xffffffffffffffff, 0xffffffffffffffff},[]uint32{
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
			0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
		}},
	}

	for i, tc := range testCases {

		index := indexChan{}
		index.indexes = &[INDEX_SIZE]uint32{}
		carried := 0

		for _, mask := range tc.masks {
			flatten_bits_incremental(index.indexes, &index.length, mask, &carried)
		}

		compare := make([]uint32, 0, 1024)
		for idx := 0; idx < index.length; idx++ {
			compare = append(compare, index.indexes[idx])
		}

		if !reflect.DeepEqual(compare, tc.expected) {
			t.Errorf("TestFlattenBitsIncremental(%d): got: %v want: %v", i, compare, tc.expected)
		}
	}
}

func TestFlattenBits(t *testing.T) {

	testCases := []struct {
		bits     uint64
		expected []uint32
	}{
		{0x11,[]uint32{0x0, 0x4}},
		{0x100100100100,[]uint32{0x8, 0x14, 0x20, 0x2c}},
		{0x8101010101010101,[]uint32{0x0, 0x8, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38, 0x3f}},
		{0xf000000000000000,[]uint32{0x3c, 0x3d, 0x3e, 0x3f}},
		{0xffffffffffffffff,[]uint32{
			0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf,
			0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
			0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
		}},
	}

	for i, tc := range testCases {

		index := indexChan{}
		index.indexes = &[INDEX_SIZE]uint32{}

		flatten_bits(index.indexes, &index.length, uint64(64), tc.bits)

		if index.length != len(tc.expected) {
			t.Errorf("TestFlattenBitsIncremental(%d): got: %d want: %d", i, index.length, len(tc.expected))
		}

		compare := make([]uint32, 0, 1024)
		for idx := 0; idx < index.length; idx++ {
			compare = append(compare, index.indexes[idx])
		}

		if !reflect.DeepEqual(compare, tc.expected) {
			t.Errorf("TestFlattenBits(%d): got: %v want: %v", i, compare, tc.expected)
		}
	}
}
