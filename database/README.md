# Database Package

This README file is to explain the purpose of the database package and how to use it.

## 📋 Summary

### 1. [About this package](#1---about-this-package)

### 2. [Database Structure](#2---database-structure)

- #### 2.1 [BadgerDB](#21---badgerdb)

- #### 2.2 [Key Prefixes](#22---key-prefixes)

### 3. [Exporting the Book](#3---exporting-the-book)

## 1 - About this package

The database package provides persistent storage for the book creation process and the final book itself. It uses **BadgerDB**, a fast, embedded key-value database written in pure Go, to store:

- **Positions to analyze** (queue)
- **Analysis results** (scores for each column)
- **Pending positions** (positions currently in the queue)

## 2 - Database Structure

- ### 2.1 - BadgerDB

  BadgerDB was chosen for its performance characteristics:
  - LSM tree-based storage optimized for write-heavy workloads
  - Embedded database
  - Built-in compression and memory management
  - Pure Go implementation

- ### 2.2 - Key Prefixes

  To organize different types of data, the database uses three prefixes:

  | Prefix | Purpose | Key Format                                 | Value                              |
  | ------ | ------- | ------------------------------------------ | ---------------------------------- |
  | `R:`   | Results | `R:[8 bytes: position key]`                | `[7]int8` (scores for each column) |
  | `Q:`   | Queue   | `Q:[1 byte: depth][8 bytes: position key]` | `[1 byte: depth]`                  |
  | `P:`   | Pending | `P:[8 bytes: position key]`                | `[1 byte: depth]`                  |

  - **Results** (`R:`): Stores the analysis results for positions that have been solved
  - **Queue** (`Q:`): Stores positions waiting to be analyzed, sorted by depth for efficient BFS
  - **Pending** (`P:`): Fast lookup to check if a position is already in the queue (avoids duplicates)

  > [!NOTE]
  > The queue uses a composite key with depth as the first byte to ensure positions are processed in breadth-first order. This allows the database iterator to naturally retrieve positions level by level.

## 3 - Exporting the Book

To use the book in the solver, one way is to directly use the database for lookups during solving. However, this include to copy the database file with the project, and adds more complexity to the project.

The other solution I found is to export the book as a Go file containing a map of positions and their analysis results. This way, the book can be directly included in the binary when compiling the project, and it can be easily used during solving without any additional complexity.

If you created your own book and want to export it, use `go generate ./...` before building the project.
