# Zhuji

A bottom-up compiler project for learning compiler construction. 

## What it does

Compiles a simple RISC-V-like bytecode to x86-64 native exectables. The pipeline goes from bytecode through codegen to assembly, then uses as and ld to produce an executable.

## Current state

The VM has 32 registers and supports arithmetic (ADD, SUB, MUL, DIV, MOD, ADDI), memory operations (LW, SW), and branches (BEQ, BLT, BNE, BGE). The x86-64 codegen handles all the arithmetic and branch instructions but not memory yet.

## Building
To make an execuwtable you need to produce the .s file and writing the output to a text file, say, `output.s` in the proj dir, then call make asm as it expects the output.s file to be there and will turn it into an executable program, read the Makefile.

```sh
make test    # run tests
make asm     # assemble/link output.s
```

## Structure

```
internal/
  vm/         - bytecode interpreter
  codegen/    - x86-64 code generator
  assembler/  - RISC-V text to bytecode
  registers/  - register file
  memory/     - byte-addressable RAM
  opcodes/    - instruction definitions
```

## Next

LW/SW codegen, then loop programs.
