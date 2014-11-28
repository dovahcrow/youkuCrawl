package main

import (
	"path/filepath"
	"testing"
)

func TestFilePath(t *testing.T) {
	t.Log(filepath.EvalSymlinks("~/"))
}
