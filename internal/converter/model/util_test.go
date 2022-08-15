package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestCanConvertSpanId(t *testing.T) {
	bytes := [8]byte{1, 2, 3, 4, 10, 11, 12, 13}

	assert.Equal(t, "010203040a0b0c0d", convertSpanId(pcommon.NewSpanID(bytes)))
}
