//
// Copyright 2019 Archon Technologies, Inc.
//

package bytez

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	var size Size
	var err error

	size, err = Parse("")
	require.Error(t, err, "parsing empty string returns error")

	size, err = Parse("4321")
	require.NoError(t, err, "parsing exact bytes succeeds")
	require.Equal(t, 4321, int(size), "parsing exact bytes succeeds")

	size, err = Parse("4k")
	require.NoError(t, err, "error parsing string")
	require.Equal(t, 4*Kilobyte, size, "incorrect value")
}
