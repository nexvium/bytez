/*
	MIT License

	Copyright (c) 2019 Javier Alvarado
*/

package bytez

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	var negative = []struct {
		in string
	}{
		{""},
		{"mb"},
		{"2.5"},
		{"2.mb"},
		{"2.9mb"},
		{"2\tmb"},
		{"2  mb"},
	}

	for _, test := range negative {
		_, err := Parse(test.in)
		require.Error(t, err)
	}

	var positive = []struct {
		in  string
		out uint64
	}{
		{"4321", uint64(4321)},
		{"4kb", 4 * Kilobyte},
		{"4.0k", 4 * Kilobyte},
		{"4.5k", 4*Kilobyte + Kilobyte/2},
		{"4 GiB", 4 * Gibibyte},
		{"4.0 GiB", 4 * Gibibyte},
		{"4.5 GiB", 4*Gibibyte + Gibibyte/2},
	}

	for _, test := range positive {
		out, err := Parse(test.in)
		require.NoError(t, err)
		require.Equal(t, test.out, out)
	}
}

func TestSize_AsString(t *testing.T) {
	var tests = []struct {
		in  uint64
		out string
	}{
		{1, "1"},
		{500, "500"},
		{1000, "1kb"},
		{1500, "1.5kb"},
		{5000, "5kb"},
		{5500, "5.5kb"},
		{5000000000, "5gb"},
		{5500000000, "5.5gb"},
	}

	for _, test := range tests {
		out := AsString(test.in)
		require.Equal(t, test.out, out)
	}
}
