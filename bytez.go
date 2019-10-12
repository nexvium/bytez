/*
	MIT License

	Copyright (c) 2019 Javier Alvarado
*/

// bytez package encapsulates functionality for working with large byte sizes in a human-friendly
// way. (The 'z' in the name is not an attempt to be cute but to create a package name that is
// short yet unique.)
//
// Unfortunately, the history is computing is one of inconsistency and ambiguity when it comes to
// specifying large byte sizes. For example, does 1KB mean 1000 bytes or 1024? It has been used to
// mean both. The unambiguous units KiB to indicate binary (powers of two) units were introduced
// in 1998 but are not universally used.
//
// This package handles the ambiguity by using the letter case of the first letter of the units to
// determine the base: lowercase indicates base 10 and uppercase indicates base 2. (The "smaller"
// letter indicates the smaller units.) Thus, the following units all represent 4000 bytes: 4k,
// 4kb, 4kB. And the following units all represent 4096 bytes: 4K, 4KB, 4Kb, 4KiB
//
// When converting from numbers to strings, this package uses the two-letter lowercase units
// (e.g. "mb") for powers of 10 and the three-letter mixed case (e.g. "MiB") for powers of 2.
package bytez

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// Type Size can be used to automatically marshal and unmarshal byte size specifications to and
// from text when parsing or outputting JSON, YAML, etc.
type Size uint64

// Decimal (SI) constants that fit in 64 bits
const (
	Kilobyte uint64 = 1000
	Megabyte        = Kilobyte * 1000
	Gigabyte        = Megabyte * 1000
	Terabyte        = Gigabyte * 1000
	Petabyte        = Terabyte * 1000
	Exabyte         = Petabyte * 1000
)

// Binary (ISO/IEC) constants that fit in 64 bits
const (
	Kibibyte uint64 = 1024
	Mebibyte        = Kibibyte * 1024
	Gibibyte        = Mebibyte * 1024
	Tebibyte        = Gibibyte * 1024
	Pebibyte        = Tebibyte * 1024
	Exbibyte        = Pebibyte * 1024
)

var unitMap = map[string]uint64{
	"k": Kilobyte, "kb": Kilobyte, "kB": Kilobyte,
	"m": Megabyte, "mb": Megabyte, "mB": Megabyte,
	"g": Gigabyte, "gb": Gigabyte, "gB": Gigabyte,
	"t": Terabyte, "tb": Terabyte, "tB": Terabyte,
	"p": Petabyte, "pb": Petabyte, "pB": Petabyte,
	"e": Exabyte, "eb": Exabyte, "eB": Exabyte,

	"K": Kibibyte, "KB": Kibibyte, "Kb": Kibibyte, "KiB": Kibibyte,
	"M": Mebibyte, "MB": Mebibyte, "Mb": Mebibyte, "MiB": Mebibyte,
	"G": Gibibyte, "GB": Gibibyte, "Gb": Gibibyte, "GiB": Gibibyte,
	"T": Tebibyte, "TB": Tebibyte, "Tb": Tebibyte, "TiB": Tebibyte,
	"P": Pebibyte, "PB": Pebibyte, "Pb": Pebibyte, "PiB": Pebibyte,
	"E": Exbibyte, "EB": Exbibyte, "Eb": Exbibyte, "EiB": Exbibyte,
}

var unitsBase2 = []string{"", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
var unitsBase10 = []string{"", "kb", "mb", "gb", "tb", "pb", "eb"}
var valuesBase2 = []uint64{1, Kibibyte, Mebibyte, Gibibyte, Tebibyte, Pebibyte, Exbibyte}
var valuesBase10 = []uint64{1, Kilobyte, Megabyte, Gigabyte, Terabyte, Petabyte, Exabyte}

func (sz *Size) UnmarshalText(bytes []byte) error {
	val, err := AsInt(string(bytes))
	if err != nil {
		return err
	}

	*sz = Size(val)
	return nil
}

func (sz Size) MarshalText() ([]byte, error) {
	return []byte(AsStr(uint64(sz))), nil
}

// AsStr accepts a number of bytes, like 4194304, and returns the byte size as a string,
// like "4MiB". The function tries to return a value that uses one of the supported units but it
// is not guaranteed to do so.
func AsStr(size uint64) string {
	if size < 1000 {
		return strconv.FormatUint(size, 10)
	}

	var base uint64
	var values []uint64
	var units []string

	if size%500 == 0 {
		base = 1000
		values = valuesBase10
		units = unitsBase10
	} else if size%512 == 0 {
		base = 1024
		values = valuesBase2
		units = unitsBase2
	} else {
		return strconv.FormatUint(size, 10)
	}

	var idx int
	for sz := size; sz >= base; sz /= base {
		idx++
	}

	str := strconv.FormatUint(size/values[idx], 10)
	if size%values[idx] != 0 {
		str += ".5"
	}
	str += units[idx]

	return str
}

// AsInt accepts a byte size, like "4MiB", and returns the exact number of bytes, like 4194304.
// The leading number should be a whole number, but as a special case the fractions ".0" and ".5"
// are allowed, like "1.5mb" to indicate 1,500,000 bytes. A single space is allowed between
// the number and the units.
func AsInt(str string) (uint64, error) {
	var num uint64
	var idx int

	str = strings.Trim(str, " \t\r\n")
	for idx = 0; idx < len(str); idx++ {
		if str[idx] < '0' || str[idx] > '9' {
			break
		} else {
			num = num*10 + uint64(str[idx]-'0')
		}
	}

	if idx == 0 {
		return 0, errors.New("no number in string")
	}

	// If the number has no units label, it is an exact number of bytes.
	if idx == len(str) {
		return num, nil
	}

	// Special case: allow ".5" to specify half units like 2.5GiB, and ".0" for parity.
	var addHalf uint64
	if str[idx] == '.' {
		if idx < len(str)-1 && str[idx:idx+2] == ".5" {
			addHalf = 1
			idx += 2
		} else if idx < len(str)-1 && str[idx:idx+2] == ".0" {
			idx += 2
		} else {
			return 0, errors.New("invalid fractional part")
		}
	}

	// A single space, not a tab or two spaces, is allowed.
	if idx < len(str) && str[idx] == ' ' {
		idx++
	}

	if str[idx:] == "" {
		return 0, errors.New("missing units")
	} else if !unicode.IsLetter(rune(str[idx])) {
		return 0, errors.New("invalid delimiter")
	} else if val, ok := unitMap[str[idx:]]; ok {
		num *= val
		num += val / 2 * addHalf
	} else {
		return 0, errors.New("invalid units")
	}

	return num, nil
}
