package util

import (
	"bytes"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
)

// BufferOut parser ast node to SQL string
func BufferOut(node ast.Node) (string, error) {
	out := new(bytes.Buffer)
	restore := format.NewRestoreCtx(format.DefaultRestoreFlags, out)
	err := node.Restore(restore)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
