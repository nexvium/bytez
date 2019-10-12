/*
	MIT License

	Copyright (c) 2019 Javier Alvarado
*/

package bytez

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsInt(t *testing.T) {
	var negative = []struct {
		in string
	}{
		{""},
		{"mb"},
		{"2."},
		{"2.5"},
		{"2.mb"},
		{"2.9mb"},
		{"2\tmb"},
		{"2  mb"},
	}

	for _, test := range negative {
		_, err := AsInt(test.in)
		if testing.Verbose() {
			fmt.Printf("\"%v\" ==> %v\n", test.in, err)
		}
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
		out, err := AsInt(test.in)
		if testing.Verbose() {
			fmt.Printf("\"%v\" --> %v\n", test.in, out)
		}
		require.NoError(t, err)
		require.Equal(t, test.out, out)
	}
}

func TestAsStr(t *testing.T) {
	var tests = []struct {
		in  uint64
		out string
	}{
		{0, "0"},
		{1, "1"},
		{500, "500"},
		//
		{1000, "1kb"},
		{1500, "1.5kb"},
		{3000000, "3mb"},
		{3500000, "3.5mb"},
		{5000000000, "5gb"},
		{5500000000, "5.5gb"},
		//
		{1024, "1KiB"},
		{1536, "1.5KiB"},
		{3145728, "3MiB"},
		{3670016, "3.5MiB"},
		{5368709120, "5GiB"},
		{5905580032, "5.5GiB"},
		//
		{314159265359, "314159265359"},
	}

	for _, test := range tests {
		out := AsStr(test.in)
		if testing.Verbose() {
			fmt.Printf("%v --> %v\n", test.in, out)
		}
		require.Equal(t, test.out, out)
	}
}

func TestMarshal(t *testing.T) {
	type conf struct {
		CacheSize Size `json:"cache_size"`
	}

	var cfg conf
	var err error

	cfgStr := `{"cache_size": "50mb"}`
	err = json.Unmarshal([]byte(cfgStr), &cfg)
	if testing.Verbose() {
		fmt.Printf("%+v  -->  %+v\n", cfgStr, cfg)
	}
	require.NoError(t, err)
	require.Equal(t, 50*Megabyte, uint64(cfg.CacheSize))

	cfg = conf{CacheSize: Size(100*Mebibyte + Mebibyte/2)}
	bytes, err := json.Marshal(cfg)
	if testing.Verbose() {
		fmt.Printf("%+v  -->  %+v\n", cfg, string(bytes))
	}
	require.NoError(t, err)
	require.Equal(t, `{"cache_size":"100.5MiB"}`, string(bytes))
}
