decompelf
=========

Produces tiny elf file with functions and objects provided by [decomp2dbg Ghidra plugin](https://github.com/mahaloz/decomp2dbg/tree/main/decompilers/d2d_ghidra)

### Build:

```shell
make build
```

### Usage:

```shell
Usage of ./decompelf:
  -arch int
    	32 or 64 bit
  -byteorder string
    	l - little endian, b - big endian (default "l")
  -elftype int
    	https://pkg.go.dev/debug/elf#Type (default 1)
  -flags string
    	ELF flags, ex. 0x0
  -l	list all machines
  -machine string
    	ex. X86_64
  -out string
    	 (default "/tmp/tinyelf")
  -url string
    	decomp2dbg server url (default "http://localhost:3662/RPC2")
```

Command-line options take priority over decomp2dbg-provided values.

<details>
<summary>
Example:
</summary>

```shell
./decompelf --machine arm --flags 0x5000000 --arch 32 --byteorder l --out /tmp/tinyelf
2024/02/25 19:20:59 INFO new tinyelf filename=test machine_id=40 machine_name=ARM machine_comment=ARM. is_32bit=true flags=0x5000400 byteorder=LittleEndian image_base=0x10000
2024/02/25 19:20:59 INFO got function headers functions=2
2024/02/25 19:21:01 INFO got global vars objects=1
2024/02/25 19:21:01 INFO done

$ readelf -a --wide /tmp/tinyelf 
ELF Header:
  Magic:   7f 45 4c 46 01 01 01 00 00 00 00 00 00 00 00 00 
  Class:                             ELF32
  Data:                              2's complement, little endian
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI Version:                       0
  Type:                              REL (Relocatable file)
  Machine:                           ARM
  Version:                           0x1
  Entry point address:               0x0
  Start of program headers:          0 (bytes into file)
  Start of section headers:          171 (bytes into file)
  Flags:                             0x5000400, Version5 EABI, hard-float ABI
  Size of this header:               52 (bytes)
  Size of program headers:           0 (bytes)
  Number of program headers:         0
  Size of section headers:           40 (bytes)
  Number of section headers:         5
  Section header string table index: 4

Section Headers:
  [Nr] Name              Type            Addr     Off    Size   ES Flg Lk Inf Al
  [ 0]                   NULL            00000000 000000 000000 00      0   0  0
  [ 1] .text             PROGBITS        00000000 000034 000000 00  AX  0   0  4
  [ 2] .symtab           SYMTAB          00000000 000034 000040 10      3   1  4
  [ 3] .strtab           STRTAB          00000000 000074 000016 00      0   0  1
  [ 4] .shstrtab         STRTAB          00000000 00008a 000021 00      0   0  1
Key to Flags:
  W (write), A (alloc), X (execute), M (merge), S (strings), I (info),
  L (link order), O (extra OS processing required), G (group), T (TLS),
  C (compressed), x (unknown), o (OS specific), E (exclude),
  D (mbind), y (purecode), p (processor specific)

There are no section groups in this file.

There are no program headers in this file.

There is no dynamic section in this file.

There are no relocations in this file.

There are no unwind sections in this file.

Symbol table '.symtab' contains 4 entries:
   Num:    Value  Size Type    Bind   Vis      Ndx Name
     0: 00000000     0 NOTYPE  LOCAL  DEFAULT  UND 
     1: 0019aa58    24 FUNC    GLOBAL DEFAULT    1 hahahaha
     2: 0019ba59    24 OBJECT  GLOBAL DEFAULT    1 hohoho
     3: 0019ba79    24 FUNC    GLOBAL DEFAULT    1 hehe
```

</details>


### Using with gdb

Example gdb script (put into .gdbinit):

```shell
define decomp
  dont-repeat
  shell decompelf
  symbol-file /tmp/tinyelf
end
```

Example gdb script with custom elf parameters :

```shell
define decomp
  dont-repeat
  shell decompelf --machine arm --flags 0x5000000 --arch 32 --byteorder b
  symbol-file /tmp/tinyelf
end
```

Run:

```shell
(gdb) decomp
2024/02/26 17:27:20 INFO new tinyelf filename=test machine_id=40 machine_name=ARM machine_comment=ARM. is_32bit=true flags=0x5000000 byteorder=BigEndian image_base=0x10000
2024/02/26 17:27:21 INFO got function headers functions=200
2024/02/26 17:27:23 INFO got global vars objects=10
2024/02/26 17:27:23 INFO done

(gdb) pipe info functions | head -n7
All defined functions:

Non-debugging symbols:
0x00000001  Elf32_Ehdr_00010000.e_ident_magic_str
0x0000000a  Elf32_Ehdr_00010000.e_ident_pad[1]
0x00000034  Elf32_Phdr_ARRAY_00010025
0x00000100  Elf32_Phdr_ARRAY_00010028[6].p_paddr
```
