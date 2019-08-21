# iaca

This package helps one use Intel's IACA tool with go programs. It does this by processing the output of `go tool objdump`. To demonstate, first build an example binary (one exists at `github.com/zeebo/iaca/example`) which just contains

```go
package main

var x [4]uint64

func main() {
	x[0]++
	x[1]++
	x[2]++
	x[3]++
	x[0]++
	x[1]++
	x[2]++
	x[3]++
}
```

One can use the objdump tool to inspect the generated code.

```
$ go tool objdump -s main.main example
TEXT main.main(SB) /home/jeff/go/src/github.com/zeebo/iaca/example/main.go
  main.go:6   0x44f930  488b0529e50800  MOVQ main.x(SB), AX
  main.go:6   0x44f937  488d4801        LEAQ 0x1(AX), CX
  main.go:6   0x44f93b  48890d1ee50800  MOVQ CX, main.x(SB)
  main.go:7   0x44f942  488b0d1fe50800  MOVQ main.x+8(SB), CX
  main.go:7   0x44f949  488d5101        LEAQ 0x1(CX), DX
  main.go:7   0x44f94d  48891514e50800  MOVQ DX, main.x+8(SB)
  main.go:8   0x44f954  488b1515e50800  MOVQ main.x+16(SB), DX
  main.go:8   0x44f95b  488d5a01        LEAQ 0x1(DX), BX
  main.go:8   0x44f95f  48891d0ae50800  MOVQ BX, main.x+16(SB)
  main.go:9   0x44f966  488b1d0be50800  MOVQ main.x+24(SB), BX
  main.go:10  0x44f96d  4883c002        ADDQ $0x2, AX
  main.go:10  0x44f971  488905e8e40800  MOVQ AX, main.x(SB)
  main.go:11  0x44f978  488d4102        LEAQ 0x2(CX), AX
  main.go:11  0x44f97c  488905e5e40800  MOVQ AX, main.x+8(SB)
  main.go:12  0x44f983  488d4202        LEAQ 0x2(DX), AX
  main.go:12  0x44f987  488905e2e40800  MOVQ AX, main.x+16(SB)
  main.go:13  0x44f98e  488d4302        LEAQ 0x2(BX), AX
  main.go:13  0x44f992  488905dfe40800  MOVQ AX, main.x+24(SB)
  main.go:14  0x44f999  c3              RET
```

Running the output of this package over that produces the binary data with the iaca markers inserted.

```
$ go tool objdump -s main.main a.out | go run github.com/zeebo/iaca | xxd
00000000: bb6f 0000 0064 6790 488b 0529 e508 0048  .o...dg.H..)...H
00000010: 8d48 0148 890d 1ee5 0800 488b 0d1f e508  .H.H......H.....
00000020: 0048 8d51 0148 8915 14e5 0800 488b 1515  .H.Q.H......H...
00000030: e508 0048 8d5a 0148 891d 0ae5 0800 488b  ...H.Z.H......H.
00000040: 1d0b e508 0048 83c0 0248 8905 e8e4 0800  .....H...H......
00000050: 488d 4102 4889 05e5 e408 0048 8d42 0248  H.A.H......H.B.H
00000060: 8905 e2e4 0800 488d 4302 4889 05df e408  ......H.C.H.....
00000070: 00c3 bbde 0000 0064 6790                 .......dg.
```

The starting marker is the byte sequence `bb6f000000646790` and the ending marker is the byte sequence `bbde000000646790`, and in between is the same byte sequence as what was produced by `go tool objdump`. Then, Intel's iaca tool works like expected.

```
$ iaca <( go tool objdump -s main.main a.out | go run github.com/zeebo/iaca )
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
| Cycles |  0.5     0.0  |  3.0  |  4.0     3.0  |  4.0     2.0  |  7.0  |  3.0  |  0.5  |  4.0  |
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
|   1      |             |      |             | 1.0     1.0 |      |      |      |      | mov rax, qword ptr [rip+0x8e529]
|   1      |             | 1.0  |             |             |      |      |      |      | lea rcx, ptr [rax+0x1]
|   2^     |             |      |             |             | 1.0  |      |      | 1.0  | mov qword ptr [rip+0x8e51e], rcx
|   1      |             |      | 1.0     1.0 |             |      |      |      |      | mov rcx, qword ptr [rip+0x8e51f]
|   1      |             |      |             |             |      | 1.0  |      |      | lea rdx, ptr [rcx+0x1]
|   2^     |             |      |             |             | 1.0  |      |      | 1.0  | mov qword ptr [rip+0x8e514], rdx
|   1      |             |      |             | 1.0     1.0 |      |      |      |      | mov rdx, qword ptr [rip+0x8e515]
|   1      |             | 1.0  |             |             |      |      |      |      | lea rbx, ptr [rdx+0x1]
|   2^     |             |      |             |             | 1.0  |      |      | 1.0  | mov qword ptr [rip+0x8e50a], rbx
|   1      |             |      | 1.0     1.0 |             |      |      |      |      | mov rbx, qword ptr [rip+0x8e50b]
|   1      | 0.5         |      |             |             |      |      | 0.5  |      | add rax, 0x2
|   2^     |             |      |             | 1.0         | 1.0  |      |      |      | mov qword ptr [rip+0x8e4e8], rax
|   1      |             |      |             |             |      | 1.0  |      |      | lea rax, ptr [rcx+0x2]
|   2^     |             |      |             |             | 1.0  |      |      | 1.0  | mov qword ptr [rip+0x8e4e5], rax
|   1      |             | 1.0  |             |             |      |      |      |      | lea rax, ptr [rdx+0x2]
|   2^     |             |      | 1.0         |             | 1.0  |      |      |      | mov qword ptr [rip+0x8e4e2], rax
|   1      |             |      |             |             |      | 1.0  |      |      | lea rax, ptr [rbx+0x2]
|   2^     |             |      |             | 1.0         | 1.0  |      |      |      | mov qword ptr [rip+0x8e4df], rax
|   3^     |             |      | 1.0     1.0 |             |      |      |      |      | ret
Total Num Of Uops: 28
Analysis Notes:
Backend allocation was stalled due to unavailable allocation resources.
```

If you just want to check some hot kernel inside of some function, just trim the `go tool objdump` output to the portion you care about.
