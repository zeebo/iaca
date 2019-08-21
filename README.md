# iaca

This package helps one use Intel's IACA tool with go programs. It does this by providing some exported functions that get rewritten into the appropriate markers by the `iacaify` tool. Specifically, the `iaca.Start()` and `iaca.Stop()` function calls. To demonstate, first build an example binary (one exists at `github.com/zeebo/iaca/example`) which just contains 

```go
package main

import "github.com/zeebo/iaca"

var x [4]uint64

func main() {
	iaca.Start()
	x[0]++
	x[1]++
	x[2]++
	x[3]++
	x[0]++
	x[1]++
	x[2]++
	x[3]++
	iaca.Stop()
}
```

One can use the objdump tool to inspect the generated code. In it, the function calls are still visible, and have some `NOPL` instructions proceeding it thanks to mid-stack inlining.

```
$ go tool objdump -s main.main example
TEXT main.main(SB) /home/jeff/go/src/github.com/zeebo/iaca/example/main.go
  main.go:7   0x44f950  64488b0c25f8ffffff  MOVQ FS:0xfffffff8, CX
  main.go:7   0x44f959  483b6110            CMPQ 0x10(CX), SP
  main.go:7   0x44f95d  0f868e000000        JBE 0x44f9f1
  main.go:7   0x44f963  4883ec08            SUBQ $0x8, SP
  main.go:7   0x44f967  48892c24            MOVQ BP, 0(SP)
  main.go:7   0x44f96b  488d2c24            LEAQ 0(SP), BP

  main.go:8   0x44f96f  90                  NOPL
  iaca.go:8   0x44f970  90                  NOPL
  iaca.go:6   0x44f971  90                  NOPL
  iaca.go:5   0x44f972  e8b9ffffff          CALL github.com/zeebo/iaca.padStart(SB)

  main.go:9   0x44f977  488b05e2e40800      MOVQ main.x(SB), AX
  main.go:9   0x44f97e  488d4801            LEAQ 0x1(AX), CX
  main.go:9   0x44f982  48890dd7e40800      MOVQ CX, main.x(SB)
  main.go:10  0x44f989  488b0dd8e40800      MOVQ main.x+8(SB), CX
  main.go:10  0x44f990  488d5101            LEAQ 0x1(CX), DX
  main.go:10  0x44f994  488915cde40800      MOVQ DX, main.x+8(SB)
  main.go:11  0x44f99b  488b15cee40800      MOVQ main.x+16(SB), DX
  main.go:11  0x44f9a2  488d5a01            LEAQ 0x1(DX), BX
  main.go:11  0x44f9a6  48891dc3e40800      MOVQ BX, main.x+16(SB)
  main.go:12  0x44f9ad  488b1dc4e40800      MOVQ main.x+24(SB), BX
  main.go:13  0x44f9b4  4883c002            ADDQ $0x2, AX
  main.go:13  0x44f9b8  488905a1e40800      MOVQ AX, main.x(SB)
  main.go:14  0x44f9bf  488d4102            LEAQ 0x2(CX), AX
  main.go:14  0x44f9c3  4889059ee40800      MOVQ AX, main.x+8(SB)
  main.go:15  0x44f9ca  488d4202            LEAQ 0x2(DX), AX
  main.go:15  0x44f9ce  4889059be40800      MOVQ AX, main.x+16(SB)
  main.go:16  0x44f9d5  488d4302            LEAQ 0x2(BX), AX
  main.go:16  0x44f9d9  48890598e40800      MOVQ AX, main.x+24(SB)

  main.go:17  0x44f9e0  90                  NOPL
  iaca.go:15  0x44f9e1  90                  NOPL
  iaca.go:13  0x44f9e2  90                  NOPL
  iaca.go:12  0x44f9e3  e858ffffff          CALL github.com/zeebo/iaca.padStop(SB)

  iaca.go:12  0x44f9e8  488b2c24            MOVQ 0(SP), BP
  iaca.go:12  0x44f9ec  4883c408            ADDQ $0x8, SP
  iaca.go:12  0x44f9f0  c3                  RET
  main.go:7   0x44f9f1  e84a7effff          CALL runtime.morestack_noctxt(SB)
  main.go:7   0x44f9f6  e955ffffff          JMP main.main(SB)
```

Then, if you run `iacaify` (import path `github.com/zeebo/iaca/iacaify`) on the binary and objdump the resulting output, you can see that the calls to the iaca code were replaced with the appropriate markers.

```
$ iacaify example > example2
$ go tool objdump -s main.main out2
TEXT main.main(SB) /home/jeff/go/src/github.com/zeebo/iaca/example/main.go
  main.go:7   0x44f950  64488b0c25f8ffffff  MOVQ FS:0xfffffff8, CX
  main.go:7   0x44f959  483b6110            CMPQ 0x10(CX), SP
  main.go:7   0x44f95d  0f868e000000        JBE 0x44f9f1
  main.go:7   0x44f963  4883ec08            SUBQ $0x8, SP
  main.go:7   0x44f967  48892c24            MOVQ BP, 0(SP)
  main.go:7   0x44f96b  488d2c24            LEAQ 0(SP), BP

  main.go:8   0x44f96f  bb6f000000          MOVL $0x6f, BX
  iaca.go:5   0x44f974  646790              NOPL

  main.go:9   0x44f977  488b05e2e40800      MOVQ main.x(SB), AX
  main.go:9   0x44f97e  488d4801            LEAQ 0x1(AX), CX
  main.go:9   0x44f982  48890dd7e40800      MOVQ CX, main.x(SB)
  main.go:10  0x44f989  488b0dd8e40800      MOVQ main.x+8(SB), CX
  main.go:10  0x44f990  488d5101            LEAQ 0x1(CX), DX
  main.go:10  0x44f994  488915cde40800      MOVQ DX, main.x+8(SB)
  main.go:11  0x44f99b  488b15cee40800      MOVQ main.x+16(SB), DX
  main.go:11  0x44f9a2  488d5a01            LEAQ 0x1(DX), BX
  main.go:11  0x44f9a6  48891dc3e40800      MOVQ BX, main.x+16(SB)
  main.go:12  0x44f9ad  488b1dc4e40800      MOVQ main.x+24(SB), BX
  main.go:13  0x44f9b4  4883c002            ADDQ $0x2, AX
  main.go:13  0x44f9b8  488905a1e40800      MOVQ AX, main.x(SB)
  main.go:14  0x44f9bf  488d4102            LEAQ 0x2(CX), AX
  main.go:14  0x44f9c3  4889059ee40800      MOVQ AX, main.x+8(SB)
  main.go:15  0x44f9ca  488d4202            LEAQ 0x2(DX), AX
  main.go:15  0x44f9ce  4889059be40800      MOVQ AX, main.x+16(SB)
  main.go:16  0x44f9d5  488d4302            LEAQ 0x2(BX), AX
  main.go:16  0x44f9d9  48890598e40800      MOVQ AX, main.x+24(SB)

  main.go:17  0x44f9e0  bbde000000          MOVL $0xde, BX
  iaca.go:12  0x44f9e5  646790              NOPL

  iaca.go:12  0x44f9e8  488b2c24            MOVQ 0(SP), BP
  iaca.go:12  0x44f9ec  4883c408            ADDQ $0x8, SP
  iaca.go:12  0x44f9f0  c3                  RET
  main.go:7   0x44f9f1  e84a7effff          CALL runtime.morestack_noctxt(SB)
  main.go:7   0x44f9f6  e955ffffff          JMP main.main(SB)
```

Then, Intel's iaca tool works like expected.

```
$ iaca ./example
COULD NOT FIND START_MARKER

$ iaca <( iacaify ./example )
Intel(R) Architecture Code Analyzer Version -  v3.0-28-g1ba2cbb build date: 2017-10-23;16:42:45
Analyzed File -  /dev/fd/63
Binary Format - 64Bit
Architecture  -  SKL
Analysis Type - Throughput

Throughput Analysis Report
--------------------------
Block Throughput: 7.00 Cycles       Throughput Bottleneck: Backend
Loop Count:  22
Port Binding In Cycles Per Iteration:
--------------------------------------------------------------------------------------------------
|  Port  |   0   -  DV   |   1   |   2   -  D    |   3   -  D    |   4   |   5   |   6   |   7   |
--------------------------------------------------------------------------------------------------
| Cycles |  0.5     0.0  |  3.0  |  3.7     2.7  |  3.6     1.3  |  7.0  |  3.0  |  0.5  |  3.7  |
--------------------------------------------------------------------------------------------------

DV - Divider pipe (on port 0)
D - Data fetch pipe (on ports 2 and 3)
F - Macro Fusion with the previous instruction occurred
* - instruction micro-ops not bound to a port
^ - Micro Fusion occurred
# - ESP Tracking sync uop was issued
@ - SSE instruction followed an AVX256/AVX512 instruction, dozens of cycles penalty is expected
X - instruction not supported, was not accounted in Analysis

| Num Of   |                    Ports pressure in cycles                         |      |
|  Uops    |  0  - DV    |  1   |  2  -  D    |  3  -  D    |  4   |  5   |  6   |  7   |
-----------------------------------------------------------------------------------------
|   1      |             |      | 0.7     0.7 | 0.3     0.3 |      |      |      |      | mov rax, qword ptr [rip+0x8e4e2]
|   1      |             | 1.0  |             |             |      |      |      |      | lea rcx, ptr [rax+0x1]
|   2^     |             |      |             | 0.3         | 1.0  |      |      | 0.7  | mov qword ptr [rip+0x8e4d7], rcx
|   1      |             |      | 0.6     0.6 | 0.4     0.4 |      |      |      |      | mov rcx, qword ptr [rip+0x8e4d8]
|   1      |             |      |             |             |      | 1.0  |      |      | lea rdx, ptr [rcx+0x1]
|   2^     |             |      |             | 0.3         | 1.0  |      |      | 0.7  | mov qword ptr [rip+0x8e4cd], rdx
|   1      |             |      | 0.7     0.7 | 0.3     0.3 |      |      |      |      | mov rdx, qword ptr [rip+0x8e4ce]
|   1      |             | 1.0  |             |             |      |      |      |      | lea rbx, ptr [rdx+0x1]
|   2^     |             |      |             | 0.4         | 1.0  |      |      | 0.6  | mov qword ptr [rip+0x8e4c3], rbx
|   1      |             |      | 0.7     0.7 | 0.3     0.3 |      |      |      |      | mov rbx, qword ptr [rip+0x8e4c4]
|   1      | 0.5         |      |             |             |      |      | 0.5  |      | add rax, 0x2
|   2^     |             |      |             | 0.3         | 1.0  |      |      | 0.7  | mov qword ptr [rip+0x8e4a1], rax
|   1      |             |      |             |             |      | 1.0  |      |      | lea rax, ptr [rcx+0x2]
|   2^     |             |      | 0.3         | 0.4         | 1.0  |      |      | 0.3  | mov qword ptr [rip+0x8e49e], rax
|   1      |             | 1.0  |             |             |      |      |      |      | lea rax, ptr [rdx+0x2]
|   2^     |             |      | 0.3         | 0.3         | 1.0  |      |      | 0.4  | mov qword ptr [rip+0x8e49b], rax
|   1      |             |      |             |             |      | 1.0  |      |      | lea rax, ptr [rbx+0x2]
|   2^     |             |      | 0.4         | 0.3         | 1.0  |      |      | 0.3  | mov qword ptr [rip+0x8e498], rax
Total Num Of Uops: 25
Analysis Notes:
Backend allocation was stalled due to unavailable allocation resources.
```