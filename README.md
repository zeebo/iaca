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

One can use the objdump tool to inspect the generated code. In it, there's extra stores to some memory locations with some specific contants that the tooling will look for.

```
$ go tool objdump -s main.main example
TEXT main.main(SB) /home/jeff/go/src/github.com/zeebo/iaca/example/main.go
  main.go:8   0x44f930  90                    NOPL
  main.go:8   0x44f931  48b845170a2acde2524a  MOVQ $0x4a52e2cd2a0a1745, AX
  iaca.go:5   0x44f93b  488905eee30800        MOVQ AX, github.com/zeebo/iaca.global(SB)

  main.go:9   0x44f942  488b0517e50800        MOVQ main.x(SB), AX
  main.go:9   0x44f949  488d4801              LEAQ 0x1(AX), CX
  main.go:9   0x44f94d  48890d0ce50800        MOVQ CX, main.x(SB)
  main.go:10  0x44f954  488b0d0de50800        MOVQ main.x+8(SB), CX
  main.go:10  0x44f95b  488d5101              LEAQ 0x1(CX), DX
  main.go:10  0x44f95f  48891502e50800        MOVQ DX, main.x+8(SB)
  main.go:11  0x44f966  488b1503e50800        MOVQ main.x+16(SB), DX
  main.go:11  0x44f96d  488d5a01              LEAQ 0x1(DX), BX
  main.go:11  0x44f971  48891df8e40800        MOVQ BX, main.x+16(SB)
  main.go:12  0x44f978  488b1df9e40800        MOVQ main.x+24(SB), BX
  main.go:13  0x44f97f  4883c002              ADDQ $0x2, AX
  main.go:13  0x44f983  488905d6e40800        MOVQ AX, main.x(SB)
  main.go:14  0x44f98a  488d4102              LEAQ 0x2(CX), AX
  main.go:14  0x44f98e  488905d3e40800        MOVQ AX, main.x+8(SB)
  main.go:15  0x44f995  488d4202              LEAQ 0x2(DX), AX
  main.go:15  0x44f999  488905d0e40800        MOVQ AX, main.x+16(SB)
  main.go:16  0x44f9a0  488d4302              LEAQ 0x2(BX), AX
  main.go:16  0x44f9a4  488905cde40800        MOVQ AX, main.x+24(SB)

  main.go:17  0x44f9ab  90                    NOPL
  iaca.go:7   0x44f9ac  48b8a0283d2f8d21ac47  MOVQ $0x47ac218d2f3d28a0, AX
  iaca.go:7   0x44f9b6  48890573e30800        MOVQ AX, github.com/zeebo/iaca.global(SB)

  iaca.go:7   0x44f9bd  c3                    RET
```

Then, if you run `iacaify` (import path `github.com/zeebo/iaca/iacaify`) on the binary and objdump the resulting output, you can see that the calls to the iaca code were replaced with the appropriate markers.

```
$ iacaify example > example2
$ go tool objdump -s main.main out2
TEXT main.main(SB) /home/jeff/go/src/github.com/zeebo/iaca/example/main.go
  main.go:8   0x44f930  90              NOPL
  main.go:8   0x44f931  90              NOPL
  main.go:8   0x44f932  90              NOPL
  main.go:8   0x44f933  90              NOPL
  main.go:8   0x44f934  90              NOPL
  main.go:8   0x44f935  90              NOPL
  main.go:8   0x44f936  90              NOPL
  main.go:8   0x44f937  90              NOPL
  main.go:8   0x44f938  90              NOPL
  main.go:8   0x44f939  90              NOPL
  main.go:8   0x44f93a  bb6f000000      MOVL $0x6f, BX
  iaca.go:5   0x44f93f  646790          NOPL

  main.go:9   0x44f942  488b0517e50800  MOVQ main.x(SB), AX
  main.go:9   0x44f949  488d4801        LEAQ 0x1(AX), CX
  main.go:9   0x44f94d  48890d0ce50800  MOVQ CX, main.x(SB)
  main.go:10  0x44f954  488b0d0de50800  MOVQ main.x+8(SB), CX
  main.go:10  0x44f95b  488d5101        LEAQ 0x1(CX), DX
  main.go:10  0x44f95f  48891502e50800  MOVQ DX, main.x+8(SB)
  main.go:11  0x44f966  488b1503e50800  MOVQ main.x+16(SB), DX
  main.go:11  0x44f96d  488d5a01        LEAQ 0x1(DX), BX
  main.go:11  0x44f971  48891df8e40800  MOVQ BX, main.x+16(SB)
  main.go:12  0x44f978  488b1df9e40800  MOVQ main.x+24(SB), BX
  main.go:13  0x44f97f  4883c002        ADDQ $0x2, AX
  main.go:13  0x44f983  488905d6e40800  MOVQ AX, main.x(SB)
  main.go:14  0x44f98a  488d4102        LEAQ 0x2(CX), AX
  main.go:14  0x44f98e  488905d3e40800  MOVQ AX, main.x+8(SB)
  main.go:15  0x44f995  488d4202        LEAQ 0x2(DX), AX
  main.go:15  0x44f999  488905d0e40800  MOVQ AX, main.x+16(SB)
  main.go:16  0x44f9a0  488d4302        LEAQ 0x2(BX), AX
  main.go:16  0x44f9a4  488905cde40800  MOVQ AX, main.x+24(SB)

  main.go:17  0x44f9ab  bbde000000      MOVL $0xde, BX
  iaca.go:7   0x44f9b0  646790          NOPL
  iaca.go:7   0x44f9b3  90              NOPL
  iaca.go:7   0x44f9b4  90              NOPL
  iaca.go:7   0x44f9b5  90              NOPL
  iaca.go:7   0x44f9b6  90              NOPL
  iaca.go:7   0x44f9b7  90              NOPL
  iaca.go:7   0x44f9b8  90              NOPL
  iaca.go:7   0x44f9b9  90              NOPL
  iaca.go:7   0x44f9ba  90              NOPL
  iaca.go:7   0x44f9bb  90              NOPL
  iaca.go:7   0x44f9bc  90              NOPL

  iaca.go:7   0x44f9bd  c3              RET
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
