//go:build mage

// Magefile for go-lucky specific tasks
package main

import (
	"github.com/magefile/mage/sh"
)

// TestQuick runs fast unit tests excluding performance tests
func TestQuick() error {
	return sh.RunV("go", "test", "-short", "./...")
}
