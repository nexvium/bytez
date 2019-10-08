// bytez package encapsulates functionality for working with large byte sizes in a human-friendly
// way. The 'z' in the name is not an attempt to be cutesy but to create a package name that is
// short yet unique.
//
// Unfortunately, the history is computing is one of inconsistency and ambiguity when it comes to
// specifying large byte sizes. For example, does 1KB mean 1000 bytes or 1024? It has been used to
// mean both. The units KiB to indicate binary (powers of two) units were introduced in 1998 but
// are not universally used.
//
// This package handles the ambiguity by using the letter case of the first letter of the units to
// determine the base: lowercase indicates base 10 and uppercase indicates base 2. (The "smaller"
// letter represents the smaller units). Thus, the following units all represent 4000 bytes: 4k,
// 4kb, 4kB. And the following units all represent 4096 bytes: 4K, 4KB, 4Kb, 4KiB
//
// When converting from numbers to strings, this package uses the two-letter lowercase units (e.g.
// "mb") for powers of 10 and the three-letter mixed case (e.g. "MiB") for powers of 2.
package bytez

import (
	"errors"
	"strings"
)

type Size uint64

// Decimal (SI) constants
const (
	Kilobyte Size = 1000
	Megabyte      = Kilobyte * 1000
	Gigabyte      = Megabyte * 1000
	Terabyte      = Gigabyte * 1000
	Petabyte      = Terabyte * 1000
	Exabyte       = Petabyte * 1000
)

// Binary (ISO/IEC) constants
const (
	Kibibyte Size = 1024
	Mebibyte      = Kibibyte * 1024
	Gibibyte      = Mebibyte * 1024
	Tebibyte      = Gibibyte * 1024
	Pebibyte      = Tebibyte * 1024
	Exbibyte      = Pebibyte * 1024
)

var unitMap = map[string]Size{
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

func (sz *Size) MarshalText() ([]byte, error) {
	if str, err := sz.AsString(); err != nil {
		return nil, err
	} else {
		return []byte(str), nil
	}
}

// As String returns the byte size as a string.
func (sz *Size) AsString() (string, error) {
	return "", nil
}

func Parse(str string) (Size, error) {
	var num Size

	var i int
	str = strings.Trim(str, " \t\r\n")
	for i = 0; i < len(str); i++ {
		if str[i] < '0' || str[i] > '9' {
			break
		} else {
			num = num*10 + Size(str[i]-'0')
		}
	}

	if i == 0 {
		return 0, errors.New("no number in units")
	}

	// If the number has no units label, it is an exact number of bytes.
	if i == len(str) {
		return Size(num), nil
	}

	// Special case: allow '.5' to specify half units like 2.5GiB.
	divisor := 1
	if i < len(str) && str[i] == '.' {
		if str[i:i+1] == ".5" {
			divisor = 2
			i += 2
		} else {
			return 0, errors.New("invalid fractional part")
		}
	}

	if val, ok := unitMap[str[i:]]; ok {
		num *= val
	} else {
		return 0, errors.New("invalid units")
	}

	return num / Size(divisor), nil
}
