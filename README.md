# 🔗 Gonnect 4

![Location](https://img.shields.io/badge/Made_in-France-red?labelColor=blue)
![Language](https://img.shields.io/badge/Language-Go-f7d3a2?labelColor=00aed8)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/Pietot/Gonnect-4?label=Latest%20release)

Gonnect 4 is a command-line connect 4 game written in Go. It allows you to solve and/or analyze a connect 4 position.

## 📋 Summary

### 1. [Features](#1---features)

### 2. [Installation](#2---installation)

### 3. [How to use](#3---how-to-use)

### 4. [Algorithms & optimizations](#4---algorithms--optimizations)

- #### [Negamax](#41-negamax)
- #### [Alpha-Beta Pruning](#42-alpha-beta-pruning)
- #### [Move explanation order](#43-move-explanation-order)
- #### [Bitboard](#44-bitboard)
- #### [Transposition Table](#45-transposition-table)
- #### [Iterative Deepening & null window](#46-iterative-deepening--null-window)
- #### [Anticipate direct losing moves](#47-anticipate-direct-losing-moves)
- #### [Better move ordering](#48-better-move-ordering)
- #### [Lower Bound transposition table](#49-lower-bound-transposition-table)

### 5. [Benchmark](#5---benchmark)

- #### 5.1 [Testsets](#testsets)
- #### 5.2 [Results](#results)

### 6. [Contribute](#6---contribute)

### 7. [Credits](#7---credits)

## 1 - Features

- Command-line interface for easy interaction
- Ability to analyze and/or solve Connect 4 positions
- Uses a book of known positions to solve hard positions instantly

## 2 - Installation

Download the binary from the [latest releases](https://github.com/Pietot/Gonnect-4/releases/latest). If you want to build it yourself, see the [build instructions](SOURCE.md).

## 3 - How to use

Open a terminal and run the following command:

```bash
cd {path to the .exe}
gonnect4.exe
```

Then it will print you how to use the tool correctly but I will explain it here further:

- Solver

  The solver gives the score from a given position and the number of remaining moves to win. Then it shows you some statistics about the search like the total time in nanoseconds, the number of nodes evaluated, the mean time per node and the number of nodes per second.

  ```bash
  gonnect4.exe -s <sequence>
  ```

- Analyzer

  The analyzer provides insights into a given position, including the score of all possible moves, the best move to make, and the number of remaining moves until victory. It also displays statistics like the `-s` command above.

  ```bash
  gonnect4.exe -a <sequence>
  ```

- The book

  The book is a precomputed database of known positions that takes +20M nodes to **analyse** that allows the solver/analyzer to instantly provide results for those positions without performing any search.

  The book is enabled by default but if you want to disable it, you can use the following flag:

  ```bash
  gonnect4.exe -s --disable-book <sequence>
  ```

> [!NOTE]
> The sequence is a string of numbers representing the columns where the pieces have been played, starting from an empty grid. The columns are numbered from 1 to 7 (left to right).
>
> The flags can by placed in the order you want but they must be before the sequence.

## 4 - Algorithms & optimizations

- ### 4.1 Negamax

  The negamax algorithm is a variant of the minimax algorithm that simplifies the implementation of two-player games. It assumes that both players are playing optimally and tries to maximize the score for the current player while minimizing it for the opponent.

  It works by recursively exploring the game tree, evaluating the positions, and returning the best move for the current player.

- ### 4.2 Alpha-Beta Pruning

  Alpha-beta pruning is an optimization technique for the minimax algorithm that reduces the number of nodes evaluated in the game tree. It works by keeping track of the best score for both players and "pruning" branches of the tree that cannot possibly influence the final decision.

  This allows the algorithm to skip evaluating certain moves, resulting in faster search times.

- ### 4.3 Move explanation order

  Until now, the move explanation order was based from the left column to right column. A small optimization was made to prioritize moves in the center columns first, as they tend to lead to more favorable positions.

- ### 4.4 Bitboard

  For now, the game was represented using a traditional 2D array, but this representation has some drawbacks in terms of performance and memory usage.

  A bitboard is a compact representation of the game state. Each position of the board is represented by two single bits in a 64-bit integer (in practice we use a 7x6 grid plus an extra bit on top of the column so we use a 49-bit integer). We can then use bitwise operations to manipulate the board state more efficiently.

  This representation enables faster move generation, evaluation, and searching, as bitwise operations are faster than arithmetic operations.

- ### 4.5 Transposition Table

  A transposition table is a cache that stores previously evaluated board positions along with their computed scores. By referencing this table during the search, the algorithm can skip redundant calculations for positions it has already analyzed, greatly increasing efficiency.

  The effectiveness of a transposition table depends on how well it retains the most valuable positions—those that are either frequently encountered or computationally expensive to evaluate. A well-designed transposition table can dramatically reduce search time by prioritizing the storage of such positions.

  Unfortunately, the transposition table implementation is very simple and could be improved in several ways.

- ### 4.6 Iterative Deepening & null window

  Iterative deepening incrementally increases the search depth, starting shallow and deepening step by step. Each iteration uses results from previous, shallower searches stored in the transposition table, improving pruning efficiency and helping to find quick wins.

  Null window search uses a minimal search window \([alpha; alpha+1]\) to quickly determine if a position's score is above or below a threshold. This narrow window enables more pruning and faster evaluations.

  By combining both techniques, the algorithm efficiently narrows down the exact score using repeated, fast null window searches at increasing depths, leveraging early results to optimize the search process.

- ### 4.7 Anticipate direct losing moves

  This technique aims to prune the search tree by detecting and avoiding moves that would allow the opponent to win immediately on their next turn. By identifying such losing moves early, the algorithm can skip evaluating them, improving efficiency.

  This anticipation helps reduce unnecessary exploration of hopeless branches.

- ### 4.8 Better move ordering

  Move exploration order is crucial for the efficiency of the alpha-beta algorithm. Previously, moves were prioritized by column, starting from the center. This can be improved by considering moves that create alignment opportunities.

  **Ordering moves with a score function:**  
  To enhance move ordering, a score function is introduced. This function evaluates each possible move by counting the number of immediate winning positions it creates (such as open 3-alignments). Moves with higher scores are explored first. If two moves have the same score, the central-column heuristic is used as a tiebreaker. This approach prioritizes moves that are more likely to lead to a win, improving the overall search efficiency.

- ### 4.9 Lower Bound transposition table

  In negamax with alpha-beta pruning, fully explored nodes provide upper bounds while pruned nodes provide lower bounds. Traditionally, only upper bounds were stored in the transposition table, but keeping lower bounds as well can improve efficiency, even if the benefit is smaller since pruned nodes are cheaper to evaluate.

  Instead of using two separate tables, it is more efficient to store both bounds in a single table by adding a flag. This is done by shifting lower bound values by the maximum possible score, effectively doubling the score range and requiring one extra bit of storage per entry.

## 5 - Performance

- ### Testsets

  To benchmark the different algorithms, I've re-created 6 datasets from [here](http://blog.gamesolver.org/solving-connect-four/02-test-protocol/) and placed them in the <a href="tests/data/">data</a> directory.

  Each file contains 1000 positions for a total of 6000 positions. More information about the test files can be found in [this README](benchmark/README.md).

  Here are the datasets sorted from the easiest to the hardest:

  | Test Set (1000 test cases each) | Number of played moves | Number of remaining moves |
  | :-----------------------------: | :--------------------: | :-----------------------: |
  |          test_easy_end          |       moves > 28       |      remaining < 14       |
  |        test_easy_middle         |    14 < moves <= 28    |      remaining < 14       |
  |         test_easy_begin         |      moves <= 14       |      remaining < 14       |
  |         test_medium_end         |    14 < moves <= 28    |   14 <= remaining < 28    |
  |       test_medium_middle        |      moves <= 14       |   14 <= remaining < 28    |
  |         test_hard_begin         |      moves <= 14       |      remaining >= 28      |

- ### Results

  I ran my solver and anoother (from [benjaminrall](https://github.com/benjaminrall/connect-four-ai)) on these test files.

  |   Test Set    |         Solver         | Mean search time (ms) | Means Nodes | Speed (KNode/second) |
  | :-----------: | :--------------------: | :-------------------: | :---------: | :------------------: |
  |   Easy End    |       Gonnect 4        |         0.01          |     80      |        0.437         |
  |               |    connect-four-ai     |                       |             |                      |
  |               |    Gonnect 4 (book)    |                       |             |                      |
  |               | connect-four-ai (book) |                       |             |                      |
  |  Easy Middle  |       Gonnect 4        |         1.44          |     11K     |         1184         |
  |               |    connect-four-ai     |                       |             |                      |
  |               |    Gonnect 4 (book)    |                       |             |                      |
  |               | connect-four-ai (book) |                       |             |                      |
  |  Easy Begin   |       Gonnect 4        |         1 155         |    8421K    |         5016         |
  |               |    connect-four-ai     |                       |             |                      |
  |               |    Gonnect 4 (book)    |                       |             |                      |
  |               | connect-four-ai (book) |                       |             |                      |
  |  Medium End   |       Gonnect 4        |          9.7          |     75K     |         2734         |
  |               |    connect-four-ai     |                       |             |                      |
  |               |    Gonnect 4 (book)    |                       |             |                      |
  |               | connect-four-ai (book) |                       |             |                      |
  | Medium Middle |       Gonnect 4        |          986          |    7061K    |         7287         |
  |               |    connect-four-ai     |                       |             |                      |
  |               |    Gonnect 4 (book)    |                       |             |                      |
  |               | connect-four-ai (book) |                       |             |                      |
  |  Hard Begin   |       Gonnect 4        |        29 776         |  207.563M   |         7028         |
  |               |    connect-four-ai     |                       |             |                      |
  |               |    Gonnect 4 (book)    |                       |             |                      |
  |               | connect-four-ai (book) |                       |             |                      |

> [!NOTE]
> Tests have been made on a 64-bit Windows 10 computer with a Ryzen 5 3600 and 16GB of RAM clocked at 3600MHz in go1.26.0 windows/amd64.

## 6 - Contribute

If you like this project and/or want to help or improve it, you can:

- Create an issue if you find a bug or want to suggest a feature or any improvement (no matter how small it is).

- Create a pull request if you want to add a feature, fix a bug or improve the code.

- Contact me if you want to talk about the project or anything else (Discord: pietot).

> [!IMPORTANT]
> Before contributing, please make sure to read the [SOURCE.md](SOURCE.md) file for a better understanding of the project structure and how the code works. This will help you a lot to contribute more effectively and avoid any confusion.
>
> If you want to be guided/helped, you already have a file named <a href="IMPROVEMENTS.md">IMPROVEMENTS.md</a> in the project directory, where you can see all the improvements that can be made.

## 7 - Credits

- [Original source code](https://github.com/PascalPons/connect4/): The original implementation of the Connect 4 solver.
- [Pascal Pons](http://blog.gamesolver.org/solving-connect-four/12-lower-bound-transposition-table/): The online tutorial I followed to implement the algorithms.
- [Online solver](https://connect4.gamesolver.org/): A web-based Connect 4 solver I used for testing.
