package cmd

import (
	"debug/elf"
	"decompelf/src/decomp2dbg/client"
	"decompelf/src/tinyelf"
	"encoding/binary"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func Start() {
	var url string
	var filename string
	var fMachine string
	var flags string
	var fByteOrder string
	var arch int
	var list bool
	flag.StringVar(&url, "url", "http://localhost:3662/RPC2", "decomp2dbg server url")
	flag.StringVar(&filename, "out", "/tmp/tinyelf", "")
	flag.StringVar(&fMachine, "machine", "", "ex. X86_64")
	flag.StringVar(&flags, "flags", "", "ELF flags, ex. 0x0")
	flag.StringVar(&fByteOrder, "byteorder", "l", "l - little endian, b - big endian")
	flag.IntVar(&arch, "arch", 0, "32 or 64 bit")
	flag.BoolVar(&list, "l", false, "list all machines")
	flag.Parse()

	if list {
		ListMachines()
		os.Exit(0)
	}

	c := client.Client{URL: url}

	reply, err := c.Ping()
	if err != nil {
		slog.Error("failed to ping decomp2dbg server", "error", err.Error())
		os.Exit(1)
	}

	if reply.Params.Param.Value.Boolean != "1" {
		slog.Error("decomp2dbg server ping reply is false")
		os.Exit(1)
	}

	elfInfo, err := c.ElfInfo()
	if err != nil {
		slog.Error("failed to get elf info from decomp2dbg", "error", err.Error())
		os.Exit(1)
	}

	var mach Machine
	var ok bool
	if fMachine != "" {
		if mach, ok = Machines[strings.ToLower(fMachine)]; !ok {
			slog.Error("invalid machine, call decompelf -l to list all machines", "machine", fMachine)
			os.Exit(1)
		}
	} else {
		if mach, ok = MachinesByID[elfInfo.Machine]; !ok {
			slog.Error("invalid machine from decomp2dbg", "machine", fMachine)
			os.Exit(1)
		}
	}

	var flagsInt int64
	if flags == "" {
		flagsInt = int64(elfInfo.Flags)
	} else {
		flags, _ = strings.CutPrefix(flags, "0x")
		flagsInt, err = strconv.ParseInt(flags, 16, 32)
		if err != nil {
			slog.Error("failed to parse flags into hex", "flag", flags, "using flags", flagsInt)
		}
	}

	var byteOrder binary.ByteOrder = binary.LittleEndian
	if fByteOrder == "" {
		if elfInfo.IsBigEndian {
			byteOrder = binary.BigEndian
		}
	} else {
		if fByteOrder == "b" {
			byteOrder = binary.BigEndian
		}
	}

	var t *tinyelf.TinyELF
	var is32 bool
	if arch == 0 {
		is32 = elfInfo.Is32Bit
	} else {
		is32 = arch == 32
	}

	slog.Info("new tinyelf", "filename", elfInfo.Name, "machine_id", int(mach.Value), "machine_name", mach.Name, "machine_comment", mach.Comment,
		"is_32bit", is32, "flags", fmt.Sprintf("0x%02x", flagsInt), "byteorder", byteOrder, "image_base", fmt.Sprintf("0x%02x", elfInfo.ImageBase))

	if is32 {
		t = tinyelf.New32(filename, mach.Value, uint32(flagsInt), byteOrder)
	} else {
		t = tinyelf.New64(filename, mach.Value, uint32(flagsInt), byteOrder)
	}

	fh, err := c.FunctionHeaders()
	if err != nil {
		slog.Error("failed to get function headers", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("function headers", "total", len(fh))

	gv, err := c.GlobalVars()
	if err != nil {
		slog.Error("failed to get global vars", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("global vars", "total", len(gv))

	for _, s := range fh {
		t.AddSymbol(s.Name, s.Value-elfInfo.ImageBase, s.Size, elf.STT_FUNC)
	}

	for _, s := range gv {
		t.AddSymbol(s.Name, s.Value, 8, elf.STT_OBJECT)
	}

	if err = t.Write(); err != nil {
		slog.Error("failed to save tiny elf", "filename", filename, "error", err.Error())
		os.Exit(1)
	}

	slog.Info("done")
}
