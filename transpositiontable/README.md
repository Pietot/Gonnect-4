This README file is to explain the purpose of the transpositiontable package and how to use it.

## 📋 Summary

### 1. [General Idea](#1---general-idea)

### 2. [Implementation Details](#2---implementation-details)

- #### 2.1 [Table Structure](#21---table-structure)
- #### 2.2 [Hash Collision Strategy](#22---hash-collision-strategy)

## 1 - General Idea

The transposition table is a crucial optimization for the negamax algorithm. During game tree search, the same position can be reached through different move sequences (like `4335454435` and `4545433534`). Without a transposition table, the solver would analyze the same position multiple times, wasting computation.

The transposition table acts as a cache, storing:
- **Key**: The position's unique key (a 64-bit hash)
- **Value**: The position's score

When the solver encounters a position, it first checks the transposition table. If found, it can use the cached bounds to tighten alpha-beta windows or even skip analyzing the position entirely. Unfortunately, we will never be able to store all because there are just too many positions (~7^n for depth n), so we need to use a fixed-size and some strategy to decide which positions to keep.

## 2 - Implementation Details

- ### 2.1 - Table Structure

    The transposition table is implemented as a fixed-size array of entries:

    ```go
    type TranspositionTable struct {
        table []Entry  // Fixed size: 8,388,593 entries, approximately 134MB of memory
    }

    type Entry struct {
        key   uint64
        value uint8
    }
    ```

    Each entry stores:
    - **64 bits** for the position key (not the full 64-bit key)
    - **8 bits** for the bound value

    This gives a total of 9 bytes per entry (stored as 16 bytes due to alignment).

- ### 2.2 - Hash Collision Strategy

    The table uses a simple **replacement scheme**:
    - Hash collisions overwrite the previous entry
    - No chaining or probing
    - No age or depth-based replacement policy

    This is called a "always replace" strategy, but better strategies exist (like depth-based replacement) that could be implemented in the future.