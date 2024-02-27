package tinyelf

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"errors"
	"os"
)

var ErrNoELF = errors.New("no elf data")

type StrTab []byte

func (s *StrTab) Append(name string) (index uint32) {
	name += "\x00"
	index = uint32(len(*s))
	*s = append(*s, []byte(name)...)

	return index
}

type elf32 struct {
	Header    elf.Header32
	SymTab    []byte
	StrTab    StrTab
	ShStrTab  []byte
	Sections  []elf.Section32
	symbuf    *bytes.Buffer
	byteOrder binary.ByteOrder
}

type elf64 struct {
	Header    elf.Header64
	SymTab    []byte
	StrTab    StrTab
	ShStrTab  []byte
	Sections  []elf.Section64
	symbuf    *bytes.Buffer
	byteOrder binary.ByteOrder
}

type TinyELF struct {
	// https://github.com/golang/go/issues/48522
	elf32     *elf32
	elf64     *elf64
	byteOrder binary.ByteOrder
	filename  string
}

var IDENT32 = [16]byte{0x7f, 'E', 'L', 'F', 0x01, 0x01, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var IDENT64 = [16]byte{0x7f, 'E', 'L', 'F', 0x02, 0x01, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func New32(filename string, machine elf.Machine, flags uint32, byteOrder binary.ByteOrder) *TinyELF {
	header := elf.Header32{
		Ident:     IDENT32,
		Type:      uint16(elf.ET_REL),
		Machine:   uint16(machine),
		Version:   uint32(elf.EV_CURRENT),
		Flags:     flags,
		Ehsize:    52,
		Shentsize: 40,
		Shnum:     5,
		Shstrndx:  4,
	}

	if byteOrder == binary.BigEndian {
		header.Ident[5] = 0x2
	}

	textSection := elf.Section32{
		Name:      1,
		Type:      uint32(elf.SHT_PROGBITS),
		Flags:     uint32(elf.SHF_ALLOC | elf.SHF_EXECINSTR), // AX = 6
		Off:       uint32(header.Ehsize),
		Size:      0,
		Addralign: 4,
	}

	symtabSection := elf.Section32{
		Name:      7,
		Type:      uint32(elf.SHT_SYMTAB),
		Addralign: 4,
		Entsize:   16,
		Link:      3,
		Info:      0,
		Off:       uint32(header.Ehsize),
	}

	strtabSection := elf.Section32{
		Name:      15,
		Type:      uint32(elf.SHT_STRTAB),
		Addralign: 1,
		Off:       uint32(header.Ehsize),
	}

	shstrtab := []byte("\x00.text\x00.symtab\x00.strtab\x00.shstrtab\x00")
	shStrTabSection := elf.Section32{
		Name:      23,
		Type:      uint32(elf.SHT_STRTAB),
		Addralign: 1,
		Size:      uint32(len(shstrtab)),
		Off:       uint32(header.Ehsize),
	}

	e := &elf32{
		Header:    header,
		SymTab:    []byte{},
		StrTab:    []byte{},
		ShStrTab:  shstrtab,
		Sections:  []elf.Section32{elf.Section32{}, textSection, symtabSection, strtabSection, shStrTabSection},
		symbuf:    new(bytes.Buffer),
		byteOrder: byteOrder,
	}

	e.AddSymbol("", 0x0, 0, elf.STT_FUNC)
	e.Sections[2].Info = 1

	t := &TinyELF{
		elf32:     e,
		filename:  filename,
		byteOrder: byteOrder,
	}

	return t
}

func (e *elf32) AddSymbol(name string, value int, size int, symType elf.SymType) {
	index := e.StrTab.Append(name)

	symbol := elf.Sym32{
		Name:  index,
		Value: uint32(value),
	}

	if value != 0x0 {
		symbol.Info = uint8(elf.STB_GLOBAL<<4) + uint8(symType&0xf) // 18
		symbol.Size = uint32(size)
		symbol.Shndx = 1
		symbol.Other = 0
	}

	binary.Write(e.symbuf, e.byteOrder, symbol)
	e.SymTab = append(e.SymTab, e.symbuf.Bytes()...)

	e.Sections[2].Size += uint32(len(e.symbuf.Bytes()))
	e.Sections[3].Off += uint32(len(e.symbuf.Bytes()))

	e.Sections[3].Size = uint32(len(e.StrTab))
	e.Sections[4].Off = e.Sections[3].Off + e.Sections[3].Size

	e.symbuf.Reset()
}

func (e *elf64) AddSymbol(name string, value int, size int, symType elf.SymType) {
	index := e.StrTab.Append(name)

	symbol := elf.Sym64{
		Name:  index,
		Value: uint64(value),
	}

	if value != 0x0 {
		symbol.Info = uint8(elf.STB_GLOBAL<<4) + uint8(symType&0xf) // 18
		symbol.Size = uint64(size)
		symbol.Shndx = 1
		symbol.Other = 0
	}

	binary.Write(e.symbuf, e.byteOrder, symbol)
	e.SymTab = append(e.SymTab, e.symbuf.Bytes()...)

	e.Sections[2].Size += uint64(len(e.symbuf.Bytes()))
	e.Sections[3].Off += uint64(len(e.symbuf.Bytes()))

	e.Sections[3].Size = uint64(len(e.StrTab))
	e.Sections[4].Off = e.Sections[3].Off + e.Sections[3].Size

	e.symbuf.Reset()
}

func (t *TinyELF) AddSymbol(name string, value int, size int, symType elf.SymType) {
	if t.elf32 != nil {
		t.elf32.AddSymbol(name, value, size, symType)
		return
	}

	if t.elf64 != nil {
		t.elf64.AddSymbol(name, value, size, symType)
		return
	}
}

func New64(filename string, machine elf.Machine, flags uint32, byteOrder binary.ByteOrder) *TinyELF {
	header := elf.Header64{
		Ident:     IDENT64,
		Type:      uint16(elf.ET_REL),
		Machine:   uint16(machine),
		Version:   uint32(elf.EV_CURRENT),
		Flags:     flags,
		Ehsize:    64,
		Shentsize: 64,
		Shnum:     5,
		Shstrndx:  4,
	}

	if byteOrder == binary.BigEndian {
		header.Ident[5] = 0x2
	}

	textSection := elf.Section64{
		Name:      1,
		Type:      uint32(elf.SHT_PROGBITS),
		Flags:     uint64(elf.SHF_ALLOC | elf.SHF_EXECINSTR), // AX = 6
		Off:       uint64(header.Ehsize),
		Size:      0,
		Addralign: 4,
	}

	symtabSection := elf.Section64{
		Name:      7,
		Type:      uint32(elf.SHT_SYMTAB),
		Addralign: 8,
		Entsize:   0x18,
		Link:      3,
		Info:      0,
		Off:       uint64(header.Ehsize),
	}

	strtabSection := elf.Section64{
		Name:      15,
		Type:      uint32(elf.SHT_STRTAB),
		Addralign: 1,
		Off:       uint64(header.Ehsize),
	}

	shstrtab := []byte("\x00.text\x00.symtab\x00.strtab\x00.shstrtab\x00")
	shStrTabSection := elf.Section64{
		Name:      23,
		Type:      uint32(elf.SHT_STRTAB),
		Addralign: 1,
		Size:      uint64(len(shstrtab)),
		Off:       uint64(header.Ehsize),
	}

	e := &elf64{
		Header:    header,
		SymTab:    []byte{},
		StrTab:    []byte{},
		ShStrTab:  shstrtab,
		Sections:  []elf.Section64{elf.Section64{}, textSection, symtabSection, strtabSection, shStrTabSection},
		symbuf:    new(bytes.Buffer),
		byteOrder: byteOrder,
	}

	e.AddSymbol("", 0x0, 0, elf.STT_FUNC)
	e.Sections[2].Info = 1

	t := &TinyELF{
		elf64:     e,
		filename:  filename,
		byteOrder: byteOrder,
	}

	return t
}

func (t *TinyELF) Write() error {
	f, err := os.OpenFile(t.filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	if t.elf32 != nil {
		t.elf32.Header.Shoff = t.elf32.Sections[4].Off + t.elf32.Sections[4].Size
		binary.Write(f, t.byteOrder, t.elf32.Header)
		binary.Write(f, t.byteOrder, t.elf32.SymTab)
		binary.Write(f, t.byteOrder, t.elf32.StrTab)
		binary.Write(f, t.byteOrder, t.elf32.ShStrTab)
		binary.Write(f, t.byteOrder, t.elf32.Sections)

		return nil
	}

	if t.elf64 != nil {
		t.elf64.Header.Shoff = t.elf64.Sections[4].Off + t.elf64.Sections[4].Size
		binary.Write(f, t.byteOrder, t.elf64.Header)
		binary.Write(f, t.byteOrder, t.elf64.SymTab)
		binary.Write(f, t.byteOrder, t.elf64.StrTab)
		binary.Write(f, t.byteOrder, t.elf64.ShStrTab)
		binary.Write(f, t.byteOrder, t.elf64.Sections)

		return nil
	}

	return ErrNoELF
}
