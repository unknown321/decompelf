//go:build ignore
// +build ignore

package main

import (
	"debug/elf"
	"decompelf/src/tinyelf"
	"encoding/binary"
)

func main() {
	filename := "/tmp/tinyelf"
	t := tinyelf.New32(filename, elf.EM_ARM, 0x5000000, binary.LittleEndian)

	t.AddSymbol("hahahaha", 0x0019aa58, 24, elf.STT_FUNC) // ghidra - 0x10000
	t.AddSymbol("hohoho", 0x0019ba59, 24, elf.STT_OBJECT)
	t.AddSymbol("hehe", 0x0019ba79, 24, elf.STT_FUNC)

	err := t.Write()
	if err != nil {
		panic(err)
	}
}
