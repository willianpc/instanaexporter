package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestCanDetermineIfAttributeIsSet(t *testing.T) {
	attrMap := pcommon.NewMap()
	attrMap.InsertString("foo", "bar")
	attrMap.InsertString("fizz", "buzz")

	assert.Equal(t, false, containsAttributes(attrMap, "bingo"))
	assert.Equal(t, false, containsAttributes(attrMap, "bingo", "buzz"))

	assert.Equal(t, true, containsAttributes(attrMap, "foo"))
	assert.Equal(t, true, containsAttributes(attrMap, "foo", "fizz"))
}
