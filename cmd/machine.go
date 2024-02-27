package cmd

import (
	"debug/elf"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

type Machine struct {
	Name    string
	Value   elf.Machine
	Comment string
}

var Machines = map[string]Machine{
	"none":          {Value: elf.EM_NONE, Name: "NONE", Comment: "Unknown machine."},
	"m32":           {Value: elf.EM_M32, Name: "M32", Comment: "AT&T WE32100."},
	"sparc":         {Value: elf.EM_SPARC, Name: "SPARC", Comment: "Sun SPARC."},
	"386":           {Value: elf.EM_386, Name: "386", Comment: "Intel i386."},
	"68k":           {Value: elf.EM_68K, Name: "68K", Comment: "Motorola 68000."},
	"88k":           {Value: elf.EM_88K, Name: "88K", Comment: "Motorola 88000."},
	"860":           {Value: elf.EM_860, Name: "860", Comment: "Intel i860."},
	"mips":          {Value: elf.EM_MIPS, Name: "MIPS", Comment: "MIPS R3000 Big-Endian only."},
	"s370":          {Value: elf.EM_S370, Name: "S370", Comment: "IBM System/370."},
	"mips_rs3_le":   {Value: elf.EM_MIPS_RS3_LE, Name: "MIPS_RS3_LE", Comment: "MIPS R3000 Little-Endian."},
	"parisc":        {Value: elf.EM_PARISC, Name: "PARISC", Comment: "HP PA-RISC."},
	"vpp500":        {Value: elf.EM_VPP500, Name: "VPP500", Comment: "Fujitsu VPP500."},
	"sparc32plus":   {Value: elf.EM_SPARC32PLUS, Name: "SPARC32PLUS", Comment: "SPARC v8plus."},
	"960":           {Value: elf.EM_960, Name: "960", Comment: "Intel 80960."},
	"ppc":           {Value: elf.EM_PPC, Name: "PPC", Comment: "PowerPC 32-bit."},
	"ppc64":         {Value: elf.EM_PPC64, Name: "PPC64", Comment: "PowerPC 64-bit."},
	"s390":          {Value: elf.EM_S390, Name: "S390", Comment: "IBM System/390."},
	"v800":          {Value: elf.EM_V800, Name: "V800", Comment: "NEC V800."},
	"fr20":          {Value: elf.EM_FR20, Name: "FR20", Comment: "Fujitsu FR20."},
	"rh32":          {Value: elf.EM_RH32, Name: "RH32", Comment: "TRW RH-32."},
	"rce":           {Value: elf.EM_RCE, Name: "RCE", Comment: "Motorola RCE."},
	"arm":           {Value: elf.EM_ARM, Name: "ARM", Comment: "ARM."},
	"sh":            {Value: elf.EM_SH, Name: "SH", Comment: "Hitachi SH."},
	"sparcv9":       {Value: elf.EM_SPARCV9, Name: "SPARCV9", Comment: "SPARC v9 64-bit."},
	"tricore":       {Value: elf.EM_TRICORE, Name: "TRICORE", Comment: "Siemens TriCore embedded processor."},
	"arc":           {Value: elf.EM_ARC, Name: "ARC", Comment: "Argonaut RISC Core."},
	"h8_300":        {Value: elf.EM_H8_300, Name: "H8_300", Comment: "Hitachi H8/300."},
	"h8_300h":       {Value: elf.EM_H8_300H, Name: "H8_300H", Comment: "Hitachi H8/300H."},
	"h8s":           {Value: elf.EM_H8S, Name: "H8S", Comment: "Hitachi H8S."},
	"h8_500":        {Value: elf.EM_H8_500, Name: "H8_500", Comment: "Hitachi H8/500."},
	"ia_64":         {Value: elf.EM_IA_64, Name: "IA_64", Comment: "Intel IA-64 Processor."},
	"mips_x":        {Value: elf.EM_MIPS_X, Name: "MIPS_X", Comment: "Stanford MIPS-X."},
	"coldfire":      {Value: elf.EM_COLDFIRE, Name: "COLDFIRE", Comment: "Motorola ColdFire."},
	"68hc12":        {Value: elf.EM_68HC12, Name: "68HC12", Comment: "Motorola M68HC12."},
	"mma":           {Value: elf.EM_MMA, Name: "MMA", Comment: "Fujitsu MMA."},
	"pcp":           {Value: elf.EM_PCP, Name: "PCP", Comment: "Siemens PCP."},
	"ncpu":          {Value: elf.EM_NCPU, Name: "NCPU", Comment: "Sony nCPU."},
	"ndr1":          {Value: elf.EM_NDR1, Name: "NDR1", Comment: "Denso NDR1 microprocessor."},
	"starcore":      {Value: elf.EM_STARCORE, Name: "STARCORE", Comment: "Motorola Star*Core processor."},
	"me16":          {Value: elf.EM_ME16, Name: "ME16", Comment: "Toyota ME16 processor."},
	"st100":         {Value: elf.EM_ST100, Name: "ST100", Comment: "STMicroelectronics ST100 processor."},
	"tinyj":         {Value: elf.EM_TINYJ, Name: "TINYJ", Comment: "Advanced Logic Corp. TinyJ processor."},
	"x86_64":        {Value: elf.EM_X86_64, Name: "X86_64", Comment: "Advanced Micro Devices x86-64"},
	"pdsp":          {Value: elf.EM_PDSP, Name: "PDSP", Comment: "Sony DSP Processor"},
	"pdp10":         {Value: elf.EM_PDP10, Name: "PDP10", Comment: "Digital Equipment Corp. PDP-10"},
	"pdp11":         {Value: elf.EM_PDP11, Name: "PDP11", Comment: "Digital Equipment Corp. PDP-11"},
	"fx66":          {Value: elf.EM_FX66, Name: "FX66", Comment: "Siemens FX66 microcontroller"},
	"st9plus":       {Value: elf.EM_ST9PLUS, Name: "ST9PLUS", Comment: "STMicroelectronics ST9+ 8/16 bit microcontroller"},
	"st7":           {Value: elf.EM_ST7, Name: "ST7", Comment: "STMicroelectronics ST7 8-bit microcontroller"},
	"68hc16":        {Value: elf.EM_68HC16, Name: "68HC16", Comment: "Motorola MC68HC16 Microcontroller"},
	"68hc11":        {Value: elf.EM_68HC11, Name: "68HC11", Comment: "Motorola MC68HC11 Microcontroller"},
	"68hc08":        {Value: elf.EM_68HC08, Name: "68HC08", Comment: "Motorola MC68HC08 Microcontroller"},
	"68hc05":        {Value: elf.EM_68HC05, Name: "68HC05", Comment: "Motorola MC68HC05 Microcontroller"},
	"svx":           {Value: elf.EM_SVX, Name: "SVX", Comment: "Silicon Graphics SVx"},
	"st19":          {Value: elf.EM_ST19, Name: "ST19", Comment: "STMicroelectronics ST19 8-bit microcontroller"},
	"vax":           {Value: elf.EM_VAX, Name: "VAX", Comment: "Digital VAX"},
	"cris":          {Value: elf.EM_CRIS, Name: "CRIS", Comment: "Axis Communications 32-bit embedded processor"},
	"javelin":       {Value: elf.EM_JAVELIN, Name: "JAVELIN", Comment: "Infineon Technologies 32-bit embedded processor"},
	"firepath":      {Value: elf.EM_FIREPATH, Name: "FIREPATH", Comment: "Element 14 64-bit DSP Processor"},
	"zsp":           {Value: elf.EM_ZSP, Name: "ZSP", Comment: "LSI Logic 16-bit DSP Processor"},
	"mmix":          {Value: elf.EM_MMIX, Name: "MMIX", Comment: "Donald Knuth's educational 64-bit processor"},
	"huany":         {Value: elf.EM_HUANY, Name: "HUANY", Comment: "Harvard University machine-independent object files"},
	"prism":         {Value: elf.EM_PRISM, Name: "PRISM", Comment: "SiTera Prism"},
	"avr":           {Value: elf.EM_AVR, Name: "AVR", Comment: "Atmel AVR 8-bit microcontroller"},
	"fr30":          {Value: elf.EM_FR30, Name: "FR30", Comment: "Fujitsu FR30"},
	"d10v":          {Value: elf.EM_D10V, Name: "D10V", Comment: "Mitsubishi D10V"},
	"d30v":          {Value: elf.EM_D30V, Name: "D30V", Comment: "Mitsubishi D30V"},
	"v850":          {Value: elf.EM_V850, Name: "V850", Comment: "NEC v850"},
	"m32r":          {Value: elf.EM_M32R, Name: "M32R", Comment: "Mitsubishi M32R"},
	"mn10300":       {Value: elf.EM_MN10300, Name: "MN10300", Comment: "Matsushita MN10300"},
	"mn10200":       {Value: elf.EM_MN10200, Name: "MN10200", Comment: "Matsushita MN10200"},
	"pj":            {Value: elf.EM_PJ, Name: "PJ", Comment: "picoJava"},
	"openrisc":      {Value: elf.EM_OPENRISC, Name: "OPENRISC", Comment: "OpenRISC 32-bit embedded processor"},
	"arc_compact":   {Value: elf.EM_ARC_COMPACT, Name: "ARC_COMPACT", Comment: "ARC International ARCompact processor (old spelling/synonym: EM_ARC_A5)"},
	"xtensa":        {Value: elf.EM_XTENSA, Name: "XTENSA", Comment: "Tensilica Xtensa Architecture"},
	"videocore":     {Value: elf.EM_VIDEOCORE, Name: "VIDEOCORE", Comment: "Alphamosaic VideoCore processor"},
	"tmm_gpp":       {Value: elf.EM_TMM_GPP, Name: "TMM_GPP", Comment: "Thompson Multimedia General Purpose Processor"},
	"ns32k":         {Value: elf.EM_NS32K, Name: "NS32K", Comment: "National Semiconductor 32000 series"},
	"tpc":           {Value: elf.EM_TPC, Name: "TPC", Comment: "Tenor Network TPC processor"},
	"snp1k":         {Value: elf.EM_SNP1K, Name: "SNP1K", Comment: "Trebia SNP 1000 processor"},
	"st200":         {Value: elf.EM_ST200, Name: "ST200", Comment: "STMicroelectronics (www.st.com) ST200 microcontroller"},
	"ip2k":          {Value: elf.EM_IP2K, Name: "IP2K", Comment: "Ubicom IP2xxx microcontroller family"},
	"max":           {Value: elf.EM_MAX, Name: "MAX", Comment: "MAX Processor"},
	"cr":            {Value: elf.EM_CR, Name: "CR", Comment: "National Semiconductor CompactRISC microprocessor"},
	"f2mc16":        {Value: elf.EM_F2MC16, Name: "F2MC16", Comment: "Fujitsu F2MC16"},
	"msp430":        {Value: elf.EM_MSP430, Name: "MSP430", Comment: "Texas Instruments embedded microcontroller msp430"},
	"blackfin":      {Value: elf.EM_BLACKFIN, Name: "BLACKFIN", Comment: "Analog Devices Blackfin (DSP) processor"},
	"se_c33":        {Value: elf.EM_SE_C33, Name: "SE_C33", Comment: "S1C33 Family of Seiko Epson processors"},
	"sep":           {Value: elf.EM_SEP, Name: "SEP", Comment: "Sharp embedded microprocessor"},
	"arca":          {Value: elf.EM_ARCA, Name: "ARCA", Comment: "Arca RISC Microprocessor"},
	"unicore":       {Value: elf.EM_UNICORE, Name: "UNICORE", Comment: "Microprocessor series from PKU-Unity Ltd. and MPRC of Peking University"},
	"excess":        {Value: elf.EM_EXCESS, Name: "EXCESS", Comment: "eXcess: 16/32/64-bit configurable embedded CPU"},
	"dxp":           {Value: elf.EM_DXP, Name: "DXP", Comment: "Icera Semiconductor Inc. Deep Execution Processor"},
	"altera_nios2":  {Value: elf.EM_ALTERA_NIOS2, Name: "ALTERA_NIOS2", Comment: "Altera Nios II soft-core processor"},
	"crx":           {Value: elf.EM_CRX, Name: "CRX", Comment: "National Semiconductor CompactRISC CRX microprocessor"},
	"xgate":         {Value: elf.EM_XGATE, Name: "XGATE", Comment: "Motorola XGATE embedded processor"},
	"c166":          {Value: elf.EM_C166, Name: "C166", Comment: "Infineon C16x/XC16x processor"},
	"m16c":          {Value: elf.EM_M16C, Name: "M16C", Comment: "Renesas M16C series microprocessors"},
	"dspic30f":      {Value: elf.EM_DSPIC30F, Name: "DSPIC30F", Comment: "Microchip Technology dsPIC30F Digital Signal Controller"},
	"ce":            {Value: elf.EM_CE, Name: "CE", Comment: "Freescale Communication Engine RISC core"},
	"m32c":          {Value: elf.EM_M32C, Name: "M32C", Comment: "Renesas M32C series microprocessors"},
	"tsk3000":       {Value: elf.EM_TSK3000, Name: "TSK3000", Comment: "Altium TSK3000 core"},
	"rs08":          {Value: elf.EM_RS08, Name: "RS08", Comment: "Freescale RS08 embedded processor"},
	"sharc":         {Value: elf.EM_SHARC, Name: "SHARC", Comment: "Analog Devices SHARC family of 32-bit DSP processors"},
	"ecog2":         {Value: elf.EM_ECOG2, Name: "ECOG2", Comment: "Cyan Technology eCOG2 microprocessor"},
	"score7":        {Value: elf.EM_SCORE7, Name: "SCORE7", Comment: "Sunplus S+core7 RISC processor"},
	"dsp24":         {Value: elf.EM_DSP24, Name: "DSP24", Comment: "New Japan Radio (NJR) 24-bit DSP Processor"},
	"videocore3":    {Value: elf.EM_VIDEOCORE3, Name: "VIDEOCORE3", Comment: "Broadcom VideoCore III processor"},
	"latticemico32": {Value: elf.EM_LATTICEMICO32, Name: "LATTICEMICO32", Comment: "RISC processor for Lattice FPGA architecture"},
	"se_c17":        {Value: elf.EM_SE_C17, Name: "SE_C17", Comment: "Seiko Epson C17 family"},
	"ti_c6000":      {Value: elf.EM_TI_C6000, Name: "TI_C6000", Comment: "The Texas Instruments TMS320C6000 DSP family"},
	"ti_c2000":      {Value: elf.EM_TI_C2000, Name: "TI_C2000", Comment: "The Texas Instruments TMS320C2000 DSP family"},
	"ti_c5500":      {Value: elf.EM_TI_C5500, Name: "TI_C5500", Comment: "The Texas Instruments TMS320C55x DSP family"},
	"ti_arp32":      {Value: elf.EM_TI_ARP32, Name: "TI_ARP32", Comment: "Texas Instruments Application Specific RISC Processor, 32bit fetch"},
	"ti_pru":        {Value: elf.EM_TI_PRU, Name: "TI_PRU", Comment: "Texas Instruments Programmable Realtime Unit"},
	"mmdsp_plus":    {Value: elf.EM_MMDSP_PLUS, Name: "MMDSP_PLUS", Comment: "STMicroelectronics 64bit VLIW Data Signal Processor"},
	"cypress_m8c":   {Value: elf.EM_CYPRESS_M8C, Name: "CYPRESS_M8C", Comment: "Cypress M8C microprocessor"},
	"r32c":          {Value: elf.EM_R32C, Name: "R32C", Comment: "Renesas R32C series microprocessors"},
	"trimedia":      {Value: elf.EM_TRIMEDIA, Name: "TRIMEDIA", Comment: "NXP Semiconductors TriMedia architecture family"},
	"qdsp6":         {Value: elf.EM_QDSP6, Name: "QDSP6", Comment: "QUALCOMM DSP6 Processor"},
	"8051":          {Value: elf.EM_8051, Name: "8051", Comment: "Intel 8051 and variants"},
	"stxp7x":        {Value: elf.EM_STXP7X, Name: "STXP7X", Comment: "STMicroelectronics STxP7x family of configurable and extensible RISC processors"},
	"nds32":         {Value: elf.EM_NDS32, Name: "NDS32", Comment: "Andes Technology compact code size embedded RISC processor family"},
	"ecog1":         {Value: elf.EM_ECOG1, Name: "ECOG1", Comment: "Cyan Technology eCOG1X family"},
	"ecog1x":        {Value: elf.EM_ECOG1X, Name: "ECOG1X", Comment: "Cyan Technology eCOG1X family"},
	"maxq30":        {Value: elf.EM_MAXQ30, Name: "MAXQ30", Comment: "Dallas Semiconductor MAXQ30 Core Micro-controllers"},
	"ximo16":        {Value: elf.EM_XIMO16, Name: "XIMO16", Comment: "New Japan Radio (NJR) 16-bit DSP Processor"},
	"manik":         {Value: elf.EM_MANIK, Name: "MANIK", Comment: "M2000 Reconfigurable RISC Microprocessor"},
	"craynv2":       {Value: elf.EM_CRAYNV2, Name: "CRAYNV2", Comment: "Cray Inc. NV2 vector architecture"},
	"rx":            {Value: elf.EM_RX, Name: "RX", Comment: "Renesas RX family"},
	"metag":         {Value: elf.EM_METAG, Name: "METAG", Comment: "Imagination Technologies META processor architecture"},
	"mcst_elbrus":   {Value: elf.EM_MCST_ELBRUS, Name: "MCST_ELBRUS", Comment: "MCST Elbrus general purpose hardware architecture"},
	"ecog16":        {Value: elf.EM_ECOG16, Name: "ECOG16", Comment: "Cyan Technology eCOG16 family"},
	"cr16":          {Value: elf.EM_CR16, Name: "CR16", Comment: "National Semiconductor CompactRISC CR16 16-bit microprocessor"},
	"etpu":          {Value: elf.EM_ETPU, Name: "ETPU", Comment: "Freescale Extended Time Processing Unit"},
	"sle9x":         {Value: elf.EM_SLE9X, Name: "SLE9X", Comment: "Infineon Technologies SLE9X core"},
	"l10m":          {Value: elf.EM_L10M, Name: "L10M", Comment: "Intel L10M"},
	"k10m":          {Value: elf.EM_K10M, Name: "K10M", Comment: "Intel K10M"},
	"aarch64":       {Value: elf.EM_AARCH64, Name: "AARCH64", Comment: "ARM 64-bit Architecture (AArch64)"},
	"avr32":         {Value: elf.EM_AVR32, Name: "AVR32", Comment: "Atmel Corporation 32-bit microprocessor family"},
	"stm8":          {Value: elf.EM_STM8, Name: "STM8", Comment: "STMicroeletronics STM8 8-bit microcontroller"},
	"tile64":        {Value: elf.EM_TILE64, Name: "TILE64", Comment: "Tilera TILE64 multicore architecture family"},
	"tilepro":       {Value: elf.EM_TILEPRO, Name: "TILEPRO", Comment: "Tilera TILEPro multicore architecture family"},
	"microblaze":    {Value: elf.EM_MICROBLAZE, Name: "MICROBLAZE", Comment: "Xilinx MicroBlaze 32-bit RISC soft processor core"},
	"cuda":          {Value: elf.EM_CUDA, Name: "CUDA", Comment: "NVIDIA CUDA architecture"},
	"tilegx":        {Value: elf.EM_TILEGX, Name: "TILEGX", Comment: "Tilera TILE-Gx multicore architecture family"},
	"cloudshield":   {Value: elf.EM_CLOUDSHIELD, Name: "CLOUDSHIELD", Comment: "CloudShield architecture family"},
	"corea_1st":     {Value: elf.EM_COREA_1ST, Name: "COREA_1ST", Comment: "KIPO-KAIST Core-A 1st generation processor family"},
	"corea_2nd":     {Value: elf.EM_COREA_2ND, Name: "COREA_2ND", Comment: "KIPO-KAIST Core-A 2nd generation processor family"},
	"arc_compact2":  {Value: elf.EM_ARC_COMPACT2, Name: "ARC_COMPACT2", Comment: "Synopsys ARCompact V2"},
	"open8":         {Value: elf.EM_OPEN8, Name: "OPEN8", Comment: "Open8 8-bit RISC soft processor core"},
	"rl78":          {Value: elf.EM_RL78, Name: "RL78", Comment: "Renesas RL78 family"},
	"videocore5":    {Value: elf.EM_VIDEOCORE5, Name: "VIDEOCORE5", Comment: "Broadcom VideoCore V processor"},
	"78kor":         {Value: elf.EM_78KOR, Name: "78KOR", Comment: "Renesas 78KOR family"},
	"56800ex":       {Value: elf.EM_56800EX, Name: "56800EX", Comment: "Freescale 56800EX Digital Signal Controller (DSC)"},
	"ba1":           {Value: elf.EM_BA1, Name: "BA1", Comment: "Beyond BA1 CPU architecture"},
	"ba2":           {Value: elf.EM_BA2, Name: "BA2", Comment: "Beyond BA2 CPU architecture"},
	"xcore":         {Value: elf.EM_XCORE, Name: "XCORE", Comment: "XMOS xCORE processor family"},
	"mchp_pic":      {Value: elf.EM_MCHP_PIC, Name: "MCHP_PIC", Comment: "Microchip 8-bit PIC(r) family"},
	"intel205":      {Value: elf.EM_INTEL205, Name: "INTEL205", Comment: "Reserved by Intel"},
	"intel206":      {Value: elf.EM_INTEL206, Name: "INTEL206", Comment: "Reserved by Intel"},
	"intel207":      {Value: elf.EM_INTEL207, Name: "INTEL207", Comment: "Reserved by Intel"},
	"intel208":      {Value: elf.EM_INTEL208, Name: "INTEL208", Comment: "Reserved by Intel"},
	"intel209":      {Value: elf.EM_INTEL209, Name: "INTEL209", Comment: "Reserved by Intel"},
	"km32":          {Value: elf.EM_KM32, Name: "KM32", Comment: "KM211 KM32 32-bit processor"},
	"kmx32":         {Value: elf.EM_KMX32, Name: "KMX32", Comment: "KM211 KMX32 32-bit processor"},
	"kmx16":         {Value: elf.EM_KMX16, Name: "KMX16", Comment: "KM211 KMX16 16-bit processor"},
	"kmx8":          {Value: elf.EM_KMX8, Name: "KMX8", Comment: "KM211 KMX8 8-bit processor"},
	"kvarc":         {Value: elf.EM_KVARC, Name: "KVARC", Comment: "KM211 KVARC processor"},
	"cdp":           {Value: elf.EM_CDP, Name: "CDP", Comment: "Paneve CDP architecture family"},
	"coge":          {Value: elf.EM_COGE, Name: "COGE", Comment: "Cognitive Smart Memory Processor"},
	"cool":          {Value: elf.EM_COOL, Name: "COOL", Comment: "Bluechip Systems CoolEngine"},
	"norc":          {Value: elf.EM_NORC, Name: "NORC", Comment: "Nanoradio Optimized RISC"},
	"csr_kalimba":   {Value: elf.EM_CSR_KALIMBA, Name: "CSR_KALIMBA", Comment: "CSR Kalimba architecture family"},
	"z80":           {Value: elf.EM_Z80, Name: "Z80", Comment: "Zilog Z80"},
	"visium":        {Value: elf.EM_VISIUM, Name: "VISIUM", Comment: "Controls and Data Services VISIUMcore processor"},
	"ft32":          {Value: elf.EM_FT32, Name: "FT32", Comment: "FTDI Chip FT32 high performance 32-bit RISC architecture"},
	"moxie":         {Value: elf.EM_MOXIE, Name: "MOXIE", Comment: "Moxie processor family"},
	"amdgpu":        {Value: elf.EM_AMDGPU, Name: "AMDGPU", Comment: "AMD GPU architecture"},
	"riscv":         {Value: elf.EM_RISCV, Name: "RISCV", Comment: "RISC-V"},
	"lanai":         {Value: elf.EM_LANAI, Name: "LANAI", Comment: "Lanai 32-bit processor"},
	"bpf":           {Value: elf.EM_BPF, Name: "BPF", Comment: "Linux BPF – in-kernel virtual machine"},
	"loongarch":     {Value: elf.EM_LOONGARCH, Name: "LOONGARCH", Comment: "LoongArch"},
	"486":           {Value: elf.EM_486, Name: "486", Comment: "Intel i486."},
	"mips_rs4_be":   {Value: elf.EM_MIPS_RS4_BE, Name: "MIPS_RS4_BE", Comment: "MIPS R4000 Big-Endian"},
	"alpha_std":     {Value: elf.EM_ALPHA_STD, Name: "ALPHA_STD", Comment: "Digital Alpha (standard value)."},
	"alpha":         {Value: elf.EM_ALPHA, Name: "ALPHA", Comment: "Alpha (written in the absence of an ABI)"},
}

var MachinesByID = map[elf.Machine]Machine{
	elf.EM_NONE:          {Value: elf.EM_NONE, Name: "NONE", Comment: "Unknown machine."},
	elf.EM_M32:           {Value: elf.EM_M32, Name: "M32", Comment: "AT&T WE32100."},
	elf.EM_SPARC:         {Value: elf.EM_SPARC, Name: "SPARC", Comment: "Sun SPARC."},
	elf.EM_386:           {Value: elf.EM_386, Name: "386", Comment: "Intel i386."},
	elf.EM_68K:           {Value: elf.EM_68K, Name: "68K", Comment: "Motorola 68000."},
	elf.EM_88K:           {Value: elf.EM_88K, Name: "88K", Comment: "Motorola 88000."},
	elf.EM_860:           {Value: elf.EM_860, Name: "860", Comment: "Intel i860."},
	elf.EM_MIPS:          {Value: elf.EM_MIPS, Name: "MIPS", Comment: "MIPS R3000 Big-Endian only."},
	elf.EM_S370:          {Value: elf.EM_S370, Name: "S370", Comment: "IBM System/370."},
	elf.EM_PARISC:        {Value: elf.EM_PARISC, Name: "PARISC", Comment: "HP PA-RISC."},
	elf.EM_VPP500:        {Value: elf.EM_VPP500, Name: "VPP500", Comment: "Fujitsu VPP500."},
	elf.EM_SPARC32PLUS:   {Value: elf.EM_SPARC32PLUS, Name: "SPARC32PLUS", Comment: "SPARC v8plus."},
	elf.EM_960:           {Value: elf.EM_960, Name: "960", Comment: "Intel 80960."},
	elf.EM_PPC:           {Value: elf.EM_PPC, Name: "PPC", Comment: "PowerPC 32-bit."},
	elf.EM_PPC64:         {Value: elf.EM_PPC64, Name: "PPC64", Comment: "PowerPC 64-bit."},
	elf.EM_S390:          {Value: elf.EM_S390, Name: "S390", Comment: "IBM System/390."},
	elf.EM_V800:          {Value: elf.EM_V800, Name: "V800", Comment: "NEC V800."},
	elf.EM_FR20:          {Value: elf.EM_FR20, Name: "FR20", Comment: "Fujitsu FR20."},
	elf.EM_RH32:          {Value: elf.EM_RH32, Name: "RH32", Comment: "TRW RH-32."},
	elf.EM_RCE:           {Value: elf.EM_RCE, Name: "RCE", Comment: "Motorola RCE."},
	elf.EM_ARM:           {Value: elf.EM_ARM, Name: "ARM", Comment: "ARM."},
	elf.EM_SH:            {Value: elf.EM_SH, Name: "SH", Comment: "Hitachi SH."},
	elf.EM_SPARCV9:       {Value: elf.EM_SPARCV9, Name: "SPARCV9", Comment: "SPARC v9 64-bit."},
	elf.EM_TRICORE:       {Value: elf.EM_TRICORE, Name: "TRICORE", Comment: "Siemens TriCore embedded processor."},
	elf.EM_ARC:           {Value: elf.EM_ARC, Name: "ARC", Comment: "Argonaut RISC Core."},
	elf.EM_H8_300:        {Value: elf.EM_H8_300, Name: "H8_300", Comment: "Hitachi H8/300."},
	elf.EM_H8_300H:       {Value: elf.EM_H8_300H, Name: "H8_300H", Comment: "Hitachi H8/300H."},
	elf.EM_H8S:           {Value: elf.EM_H8S, Name: "H8S", Comment: "Hitachi H8S."},
	elf.EM_H8_500:        {Value: elf.EM_H8_500, Name: "H8_500", Comment: "Hitachi H8/500."},
	elf.EM_IA_64:         {Value: elf.EM_IA_64, Name: "IA_64", Comment: "Intel IA-64 Processor."},
	elf.EM_MIPS_X:        {Value: elf.EM_MIPS_X, Name: "MIPS_X", Comment: "Stanford MIPS-X."},
	elf.EM_COLDFIRE:      {Value: elf.EM_COLDFIRE, Name: "COLDFIRE", Comment: "Motorola ColdFire."},
	elf.EM_68HC12:        {Value: elf.EM_68HC12, Name: "68HC12", Comment: "Motorola M68HC12."},
	elf.EM_MMA:           {Value: elf.EM_MMA, Name: "MMA", Comment: "Fujitsu MMA."},
	elf.EM_PCP:           {Value: elf.EM_PCP, Name: "PCP", Comment: "Siemens PCP."},
	elf.EM_NCPU:          {Value: elf.EM_NCPU, Name: "NCPU", Comment: "Sony nCPU."},
	elf.EM_NDR1:          {Value: elf.EM_NDR1, Name: "NDR1", Comment: "Denso NDR1 microprocessor."},
	elf.EM_STARCORE:      {Value: elf.EM_STARCORE, Name: "STARCORE", Comment: "Motorola Star*Core processor."},
	elf.EM_ME16:          {Value: elf.EM_ME16, Name: "ME16", Comment: "Toyota ME16 processor."},
	elf.EM_ST100:         {Value: elf.EM_ST100, Name: "ST100", Comment: "STMicroelectronics ST100 processor."},
	elf.EM_TINYJ:         {Value: elf.EM_TINYJ, Name: "TINYJ", Comment: "Advanced Logic Corp. TinyJ processor."},
	elf.EM_X86_64:        {Value: elf.EM_X86_64, Name: "X86_64", Comment: "Advanced Micro Devices x86-64"},
	elf.EM_PDSP:          {Value: elf.EM_PDSP, Name: "PDSP", Comment: "Sony DSP Processor"},
	elf.EM_PDP10:         {Value: elf.EM_PDP10, Name: "PDP10", Comment: "Digital Equipment Corp. PDP-10"},
	elf.EM_PDP11:         {Value: elf.EM_PDP11, Name: "PDP11", Comment: "Digital Equipment Corp. PDP-11"},
	elf.EM_FX66:          {Value: elf.EM_FX66, Name: "FX66", Comment: "Siemens FX66 microcontroller"},
	elf.EM_ST9PLUS:       {Value: elf.EM_ST9PLUS, Name: "ST9PLUS", Comment: "STMicroelectronics ST9+ 8/16 bit microcontroller"},
	elf.EM_ST7:           {Value: elf.EM_ST7, Name: "ST7", Comment: "STMicroelectronics ST7 8-bit microcontroller"},
	elf.EM_68HC16:        {Value: elf.EM_68HC16, Name: "68HC16", Comment: "Motorola MC68HC16 Microcontroller"},
	elf.EM_68HC11:        {Value: elf.EM_68HC11, Name: "68HC11", Comment: "Motorola MC68HC11 Microcontroller"},
	elf.EM_68HC08:        {Value: elf.EM_68HC08, Name: "68HC08", Comment: "Motorola MC68HC08 Microcontroller"},
	elf.EM_68HC05:        {Value: elf.EM_68HC05, Name: "68HC05", Comment: "Motorola MC68HC05 Microcontroller"},
	elf.EM_SVX:           {Value: elf.EM_SVX, Name: "SVX", Comment: "Silicon Graphics SVx"},
	elf.EM_ST19:          {Value: elf.EM_ST19, Name: "ST19", Comment: "STMicroelectronics ST19 8-bit microcontroller"},
	elf.EM_VAX:           {Value: elf.EM_VAX, Name: "VAX", Comment: "Digital VAX"},
	elf.EM_CRIS:          {Value: elf.EM_CRIS, Name: "CRIS", Comment: "Axis Communications 32-bit embedded processor"},
	elf.EM_JAVELIN:       {Value: elf.EM_JAVELIN, Name: "JAVELIN", Comment: "Infineon Technologies 32-bit embedded processor"},
	elf.EM_FIREPATH:      {Value: elf.EM_FIREPATH, Name: "FIREPATH", Comment: "Element 14 64-bit DSP Processor"},
	elf.EM_ZSP:           {Value: elf.EM_ZSP, Name: "ZSP", Comment: "LSI Logic 16-bit DSP Processor"},
	elf.EM_MMIX:          {Value: elf.EM_MMIX, Name: "MMIX", Comment: "Donald Knuth's educational 64-bit processor"},
	elf.EM_HUANY:         {Value: elf.EM_HUANY, Name: "HUANY", Comment: "Harvard University machine-independent object files"},
	elf.EM_PRISM:         {Value: elf.EM_PRISM, Name: "PRISM", Comment: "SiTera Prism"},
	elf.EM_AVR:           {Value: elf.EM_AVR, Name: "AVR", Comment: "Atmel AVR 8-bit microcontroller"},
	elf.EM_FR30:          {Value: elf.EM_FR30, Name: "FR30", Comment: "Fujitsu FR30"},
	elf.EM_D10V:          {Value: elf.EM_D10V, Name: "D10V", Comment: "Mitsubishi D10V"},
	elf.EM_D30V:          {Value: elf.EM_D30V, Name: "D30V", Comment: "Mitsubishi D30V"},
	elf.EM_V850:          {Value: elf.EM_V850, Name: "V850", Comment: "NEC v850"},
	elf.EM_M32R:          {Value: elf.EM_M32R, Name: "M32R", Comment: "Mitsubishi M32R"},
	elf.EM_MN10300:       {Value: elf.EM_MN10300, Name: "MN10300", Comment: "Matsushita MN10300"},
	elf.EM_MN10200:       {Value: elf.EM_MN10200, Name: "MN10200", Comment: "Matsushita MN10200"},
	elf.EM_PJ:            {Value: elf.EM_PJ, Name: "PJ", Comment: "picoJava"},
	elf.EM_OPENRISC:      {Value: elf.EM_OPENRISC, Name: "OPENRISC", Comment: "OpenRISC 32-bit embedded processor"},
	elf.EM_ARC_COMPACT:   {Value: elf.EM_ARC_COMPACT, Name: "ARC_COMPACT", Comment: "ARC International ARCompact processor (old spelling/synonym: EM_ARC_A5)"},
	elf.EM_XTENSA:        {Value: elf.EM_XTENSA, Name: "XTENSA", Comment: "Tensilica Xtensa Architecture"},
	elf.EM_VIDEOCORE:     {Value: elf.EM_VIDEOCORE, Name: "VIDEOCORE", Comment: "Alphamosaic VideoCore processor"},
	elf.EM_TMM_GPP:       {Value: elf.EM_TMM_GPP, Name: "TMM_GPP", Comment: "Thompson Multimedia General Purpose Processor"},
	elf.EM_NS32K:         {Value: elf.EM_NS32K, Name: "NS32K", Comment: "National Semiconductor 32000 series"},
	elf.EM_TPC:           {Value: elf.EM_TPC, Name: "TPC", Comment: "Tenor Network TPC processor"},
	elf.EM_SNP1K:         {Value: elf.EM_SNP1K, Name: "SNP1K", Comment: "Trebia SNP 1000 processor"},
	elf.EM_ST200:         {Value: elf.EM_ST200, Name: "ST200", Comment: "STMicroelectronics (www.st.com) ST200 microcontroller"},
	elf.EM_IP2K:          {Value: elf.EM_IP2K, Name: "IP2K", Comment: "Ubicom IP2xxx microcontroller family"},
	elf.EM_MAX:           {Value: elf.EM_MAX, Name: "MAX", Comment: "MAX Processor"},
	elf.EM_CR:            {Value: elf.EM_CR, Name: "CR", Comment: "National Semiconductor CompactRISC microprocessor"},
	elf.EM_F2MC16:        {Value: elf.EM_F2MC16, Name: "F2MC16", Comment: "Fujitsu F2MC16"},
	elf.EM_MSP430:        {Value: elf.EM_MSP430, Name: "MSP430", Comment: "Texas Instruments embedded microcontroller msp430"},
	elf.EM_BLACKFIN:      {Value: elf.EM_BLACKFIN, Name: "BLACKFIN", Comment: "Analog Devices Blackfin (DSP) processor"},
	elf.EM_SE_C33:        {Value: elf.EM_SE_C33, Name: "SE_C33", Comment: "S1C33 Family of Seiko Epson processors"},
	elf.EM_SEP:           {Value: elf.EM_SEP, Name: "SEP", Comment: "Sharp embedded microprocessor"},
	elf.EM_ARCA:          {Value: elf.EM_ARCA, Name: "ARCA", Comment: "Arca RISC Microprocessor"},
	elf.EM_UNICORE:       {Value: elf.EM_UNICORE, Name: "UNICORE", Comment: "Microprocessor series from PKU-Unity Ltd. and MPRC of Peking University"},
	elf.EM_EXCESS:        {Value: elf.EM_EXCESS, Name: "EXCESS", Comment: "eXcess: 16/32/64-bit configurable embedded CPU"},
	elf.EM_DXP:           {Value: elf.EM_DXP, Name: "DXP", Comment: "Icera Semiconductor Inc. Deep Execution Processor"},
	elf.EM_ALTERA_NIOS2:  {Value: elf.EM_ALTERA_NIOS2, Name: "ALTERA_NIOS2", Comment: "Altera Nios II soft-core processor"},
	elf.EM_CRX:           {Value: elf.EM_CRX, Name: "CRX", Comment: "National Semiconductor CompactRISC CRX microprocessor"},
	elf.EM_XGATE:         {Value: elf.EM_XGATE, Name: "XGATE", Comment: "Motorola XGATE embedded processor"},
	elf.EM_C166:          {Value: elf.EM_C166, Name: "C166", Comment: "Infineon C16x/XC16x processor"},
	elf.EM_M16C:          {Value: elf.EM_M16C, Name: "M16C", Comment: "Renesas M16C series microprocessors"},
	elf.EM_DSPIC30F:      {Value: elf.EM_DSPIC30F, Name: "DSPIC30F", Comment: "Microchip Technology dsPIC30F Digital Signal Controller"},
	elf.EM_CE:            {Value: elf.EM_CE, Name: "CE", Comment: "Freescale Communication Engine RISC core"},
	elf.EM_M32C:          {Value: elf.EM_M32C, Name: "M32C", Comment: "Renesas M32C series microprocessors"},
	elf.EM_TSK3000:       {Value: elf.EM_TSK3000, Name: "TSK3000", Comment: "Altium TSK3000 core"},
	elf.EM_RS08:          {Value: elf.EM_RS08, Name: "RS08", Comment: "Freescale RS08 embedded processor"},
	elf.EM_SHARC:         {Value: elf.EM_SHARC, Name: "SHARC", Comment: "Analog Devices SHARC family of 32-bit DSP processors"},
	elf.EM_ECOG2:         {Value: elf.EM_ECOG2, Name: "ECOG2", Comment: "Cyan Technology eCOG2 microprocessor"},
	elf.EM_SCORE7:        {Value: elf.EM_SCORE7, Name: "SCORE7", Comment: "Sunplus S+core7 RISC processor"},
	elf.EM_DSP24:         {Value: elf.EM_DSP24, Name: "DSP24", Comment: "New Japan Radio (NJR) 24-bit DSP Processor"},
	elf.EM_VIDEOCORE3:    {Value: elf.EM_VIDEOCORE3, Name: "VIDEOCORE3", Comment: "Broadcom VideoCore III processor"},
	elf.EM_LATTICEMICO32: {Value: elf.EM_LATTICEMICO32, Name: "LATTICEMICO32", Comment: "RISC processor for Lattice FPGA architecture"},
	elf.EM_SE_C17:        {Value: elf.EM_SE_C17, Name: "SE_C17", Comment: "Seiko Epson C17 family"},
	elf.EM_TI_C6000:      {Value: elf.EM_TI_C6000, Name: "TI_C6000", Comment: "The Texas Instruments TMS320C6000 DSP family"},
	elf.EM_TI_C2000:      {Value: elf.EM_TI_C2000, Name: "TI_C2000", Comment: "The Texas Instruments TMS320C2000 DSP family"},
	elf.EM_TI_C5500:      {Value: elf.EM_TI_C5500, Name: "TI_C5500", Comment: "The Texas Instruments TMS320C55x DSP family"},
	elf.EM_TI_ARP32:      {Value: elf.EM_TI_ARP32, Name: "TI_ARP32", Comment: "Texas Instruments Application Specific RISC Processor, 32bit fetch"},
	elf.EM_TI_PRU:        {Value: elf.EM_TI_PRU, Name: "TI_PRU", Comment: "Texas Instruments Programmable Realtime Unit"},
	elf.EM_MMDSP_PLUS:    {Value: elf.EM_MMDSP_PLUS, Name: "MMDSP_PLUS", Comment: "STMicroelectronics 64bit VLIW Data Signal Processor"},
	elf.EM_CYPRESS_M8C:   {Value: elf.EM_CYPRESS_M8C, Name: "CYPRESS_M8C", Comment: "Cypress M8C microprocessor"},
	elf.EM_R32C:          {Value: elf.EM_R32C, Name: "R32C", Comment: "Renesas R32C series microprocessors"},
	elf.EM_TRIMEDIA:      {Value: elf.EM_TRIMEDIA, Name: "TRIMEDIA", Comment: "NXP Semiconductors TriMedia architecture family"},
	elf.EM_QDSP6:         {Value: elf.EM_QDSP6, Name: "QDSP6", Comment: "QUALCOMM DSP6 Processor"},
	elf.EM_8051:          {Value: elf.EM_8051, Name: "8051", Comment: "Intel 8051 and variants"},
	elf.EM_STXP7X:        {Value: elf.EM_STXP7X, Name: "STXP7X", Comment: "STMicroelectronics STxP7x family of configurable and extensible RISC processors"},
	elf.EM_NDS32:         {Value: elf.EM_NDS32, Name: "NDS32", Comment: "Andes Technology compact code size embedded RISC processor family"},
	elf.EM_ECOG1X:        {Value: elf.EM_ECOG1X, Name: "ECOG1X", Comment: "Cyan Technology eCOG1X family"},
	elf.EM_MAXQ30:        {Value: elf.EM_MAXQ30, Name: "MAXQ30", Comment: "Dallas Semiconductor MAXQ30 Core Micro-controllers"},
	elf.EM_XIMO16:        {Value: elf.EM_XIMO16, Name: "XIMO16", Comment: "New Japan Radio (NJR) 16-bit DSP Processor"},
	elf.EM_MANIK:         {Value: elf.EM_MANIK, Name: "MANIK", Comment: "M2000 Reconfigurable RISC Microprocessor"},
	elf.EM_CRAYNV2:       {Value: elf.EM_CRAYNV2, Name: "CRAYNV2", Comment: "Cray Inc. NV2 vector architecture"},
	elf.EM_RX:            {Value: elf.EM_RX, Name: "RX", Comment: "Renesas RX family"},
	elf.EM_METAG:         {Value: elf.EM_METAG, Name: "METAG", Comment: "Imagination Technologies META processor architecture"},
	elf.EM_MCST_ELBRUS:   {Value: elf.EM_MCST_ELBRUS, Name: "MCST_ELBRUS", Comment: "MCST Elbrus general purpose hardware architecture"},
	elf.EM_ECOG16:        {Value: elf.EM_ECOG16, Name: "ECOG16", Comment: "Cyan Technology eCOG16 family"},
	elf.EM_CR16:          {Value: elf.EM_CR16, Name: "CR16", Comment: "National Semiconductor CompactRISC CR16 16-bit microprocessor"},
	elf.EM_ETPU:          {Value: elf.EM_ETPU, Name: "ETPU", Comment: "Freescale Extended Time Processing Unit"},
	elf.EM_SLE9X:         {Value: elf.EM_SLE9X, Name: "SLE9X", Comment: "Infineon Technologies SLE9X core"},
	elf.EM_L10M:          {Value: elf.EM_L10M, Name: "L10M", Comment: "Intel L10M"},
	elf.EM_K10M:          {Value: elf.EM_K10M, Name: "K10M", Comment: "Intel K10M"},
	elf.EM_AARCH64:       {Value: elf.EM_AARCH64, Name: "AARCH64", Comment: "ARM 64-bit Architecture (AArch64)"},
	elf.EM_AVR32:         {Value: elf.EM_AVR32, Name: "AVR32", Comment: "Atmel Corporation 32-bit microprocessor family"},
	elf.EM_STM8:          {Value: elf.EM_STM8, Name: "STM8", Comment: "STMicroeletronics STM8 8-bit microcontroller"},
	elf.EM_TILE64:        {Value: elf.EM_TILE64, Name: "TILE64", Comment: "Tilera TILE64 multicore architecture family"},
	elf.EM_TILEPRO:       {Value: elf.EM_TILEPRO, Name: "TILEPRO", Comment: "Tilera TILEPro multicore architecture family"},
	elf.EM_MICROBLAZE:    {Value: elf.EM_MICROBLAZE, Name: "MICROBLAZE", Comment: "Xilinx MicroBlaze 32-bit RISC soft processor core"},
	elf.EM_CUDA:          {Value: elf.EM_CUDA, Name: "CUDA", Comment: "NVIDIA CUDA architecture"},
	elf.EM_TILEGX:        {Value: elf.EM_TILEGX, Name: "TILEGX", Comment: "Tilera TILE-Gx multicore architecture family"},
	elf.EM_CLOUDSHIELD:   {Value: elf.EM_CLOUDSHIELD, Name: "CLOUDSHIELD", Comment: "CloudShield architecture family"},
	elf.EM_COREA_1ST:     {Value: elf.EM_COREA_1ST, Name: "COREA_1ST", Comment: "KIPO-KAIST Core-A 1st generation processor family"},
	elf.EM_COREA_2ND:     {Value: elf.EM_COREA_2ND, Name: "COREA_2ND", Comment: "KIPO-KAIST Core-A 2nd generation processor family"},
	elf.EM_ARC_COMPACT2:  {Value: elf.EM_ARC_COMPACT2, Name: "ARC_COMPACT2", Comment: "Synopsys ARCompact V2"},
	elf.EM_OPEN8:         {Value: elf.EM_OPEN8, Name: "OPEN8", Comment: "Open8 8-bit RISC soft processor core"},
	elf.EM_RL78:          {Value: elf.EM_RL78, Name: "RL78", Comment: "Renesas RL78 family"},
	elf.EM_VIDEOCORE5:    {Value: elf.EM_VIDEOCORE5, Name: "VIDEOCORE5", Comment: "Broadcom VideoCore V processor"},
	elf.EM_78KOR:         {Value: elf.EM_78KOR, Name: "78KOR", Comment: "Renesas 78KOR family"},
	elf.EM_56800EX:       {Value: elf.EM_56800EX, Name: "56800EX", Comment: "Freescale 56800EX Digital Signal Controller (DSC)"},
	elf.EM_BA1:           {Value: elf.EM_BA1, Name: "BA1", Comment: "Beyond BA1 CPU architecture"},
	elf.EM_BA2:           {Value: elf.EM_BA2, Name: "BA2", Comment: "Beyond BA2 CPU architecture"},
	elf.EM_XCORE:         {Value: elf.EM_XCORE, Name: "XCORE", Comment: "XMOS xCORE processor family"},
	elf.EM_MCHP_PIC:      {Value: elf.EM_MCHP_PIC, Name: "MCHP_PIC", Comment: "Microchip 8-bit PIC(r) family"},
	elf.EM_INTEL205:      {Value: elf.EM_INTEL205, Name: "INTEL205", Comment: "Reserved by Intel"},
	elf.EM_INTEL206:      {Value: elf.EM_INTEL206, Name: "INTEL206", Comment: "Reserved by Intel"},
	elf.EM_INTEL207:      {Value: elf.EM_INTEL207, Name: "INTEL207", Comment: "Reserved by Intel"},
	elf.EM_INTEL208:      {Value: elf.EM_INTEL208, Name: "INTEL208", Comment: "Reserved by Intel"},
	elf.EM_INTEL209:      {Value: elf.EM_INTEL209, Name: "INTEL209", Comment: "Reserved by Intel"},
	elf.EM_KM32:          {Value: elf.EM_KM32, Name: "KM32", Comment: "KM211 KM32 32-bit processor"},
	elf.EM_KMX32:         {Value: elf.EM_KMX32, Name: "KMX32", Comment: "KM211 KMX32 32-bit processor"},
	elf.EM_KMX16:         {Value: elf.EM_KMX16, Name: "KMX16", Comment: "KM211 KMX16 16-bit processor"},
	elf.EM_KMX8:          {Value: elf.EM_KMX8, Name: "KMX8", Comment: "KM211 KMX8 8-bit processor"},
	elf.EM_KVARC:         {Value: elf.EM_KVARC, Name: "KVARC", Comment: "KM211 KVARC processor"},
	elf.EM_CDP:           {Value: elf.EM_CDP, Name: "CDP", Comment: "Paneve CDP architecture family"},
	elf.EM_COGE:          {Value: elf.EM_COGE, Name: "COGE", Comment: "Cognitive Smart Memory Processor"},
	elf.EM_COOL:          {Value: elf.EM_COOL, Name: "COOL", Comment: "Bluechip Systems CoolEngine"},
	elf.EM_NORC:          {Value: elf.EM_NORC, Name: "NORC", Comment: "Nanoradio Optimized RISC"},
	elf.EM_CSR_KALIMBA:   {Value: elf.EM_CSR_KALIMBA, Name: "CSR_KALIMBA", Comment: "CSR Kalimba architecture family"},
	elf.EM_Z80:           {Value: elf.EM_Z80, Name: "Z80", Comment: "Zilog Z80"},
	elf.EM_VISIUM:        {Value: elf.EM_VISIUM, Name: "VISIUM", Comment: "Controls and Data Services VISIUMcore processor"},
	elf.EM_FT32:          {Value: elf.EM_FT32, Name: "FT32", Comment: "FTDI Chip FT32 high performance 32-bit RISC architecture"},
	elf.EM_MOXIE:         {Value: elf.EM_MOXIE, Name: "MOXIE", Comment: "Moxie processor family"},
	elf.EM_AMDGPU:        {Value: elf.EM_AMDGPU, Name: "AMDGPU", Comment: "AMD GPU architecture"},
	elf.EM_RISCV:         {Value: elf.EM_RISCV, Name: "RISCV", Comment: "RISC-V"},
	elf.EM_LANAI:         {Value: elf.EM_LANAI, Name: "LANAI", Comment: "Lanai 32-bit processor"},
	elf.EM_BPF:           {Value: elf.EM_BPF, Name: "BPF", Comment: "Linux BPF – in-kernel virtual machine"},
	elf.EM_LOONGARCH:     {Value: elf.EM_LOONGARCH, Name: "LOONGARCH", Comment: "LoongArch"},
	elf.EM_486:           {Value: elf.EM_486, Name: "486", Comment: "Intel i486."},
	elf.EM_MIPS_RS4_BE:   {Value: elf.EM_MIPS_RS4_BE, Name: "MIPS_RS4_BE", Comment: "MIPS R4000 Big-Endian"},
	elf.EM_ALPHA_STD:     {Value: elf.EM_ALPHA_STD, Name: "ALPHA_STD", Comment: "Digital Alpha (standard value)."},
	elf.EM_ALPHA:         {Value: elf.EM_ALPHA, Name: "ALPHA", Comment: "Alpha (written in the absence of an ABI)"},
}

func ListMachines() {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.TabIndent)
	keys := make([]string, 0, len(Machines))
	for k := range Machines {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return int(Machines[keys[i]].Value) < int(Machines[keys[j]].Value)
	})

	fmt.Fprintln(w, "name\tcomment\tid\t")
	fmt.Fprintln(w, "====\t=======\t==\t")
	for _, v := range keys {
		//nolint:gosimple
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%d\t", v, Machines[v].Comment, int(Machines[v].Value)))
	}
	w.Flush()
}