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

### 6. [Improve the project](#6---improve-the-project)

### 7. [Credits](#7---credits)

## 1 - Features

- Command-line interface for easy interaction
- Ability to analyze and/or solve Connect 4 positions

## 2 - Installation

Download the binary from the [releases page](https://github.com/Pietot/Gonnect-4/releases/latest).

## 3 - How to use

Open a terminal and run the following command:

```bash
cd {path to the .exe}
gonnect_4.exe
```

Then it will print you how to use the tool correctly but I will explain it here further:

- Solver

  The solver gives the score from a given position and the number of remaining moves to win. Then it shows you some statistics about the search like the total time in nanoseconds, the number of nodes evaluated, the mean time per node and the number of nodes per second.

  ```bash
  gonnect_4.exe --solve | --s <sequence>
  ```

- Analyzer

  The analyzer provides insights into a given position, including the score of all possible moves, the best move to make, and the number of remaining moves until victory.

  ```bash
  gonnect_4.exe --analyze | --a <sequence>
  ```

> **Note**: The sequence is a string of numbers representing the columns where the pieces have been played, starting from an empty grid. The columns are numbered from 1 to 7 (left to right).

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

## 5 - Benchmark

- ### Testsets

  To benchmark the different algorithms, I've re-created 6 datasets from [here](http://blog.gamesolver.org/solving-connect-four/02-test-protocol/) and placed them in the <a href="tests/data/">data</a> directory.

  Here are the datasets sorted from the easiest to the hardest:

  | Test Set (1000 test cases each) | Number of played moves | Number of remaining moves |
  | :-----------------------------: | :--------------------: | :-----------------------: |
  |          test_easy_end          |       28 < moves       |      remaining < 14       |
  |        test_easy_middle         |    14 < moves <= 28    |      remaining < 14       |
  |         test_easy_begin         |      moves <= 14       |      remaining < 14       |
  |         test_medium_end         |    14 < moves <= 28    |   14 <= remaining < 28    |
  |       test_medium_middle        |      moves <= 14       |   14 <= remaining < 28    |
  |         test_hard_begin         |      moves <= 14       |      28 <= remaining      |

- ### Results

  I ran the tests on each algorithm and collected the results in a CSV file and computed the **mean** for each category.

  Each algorithm had around 12 hours to solve all positions from a test set. If a test was not solved within this time frame, it was marked as failed and the other tests were not tested.

  Here are the algorithms ranked from the most efficient to the least efficient:

  | Rank |             Algorithm             | Test passed | Search time (ms) | Number of nodes | Time/node (ns) | Node/second |
  | :--: | :-------------------------------: | :---------: | :--------------: | :-------------: | :------------: | :---------: |
  |  1   |  Lower bound transposition table  |     6/6     |      1_572       |   11_079_843    |       99       |  3_607_353  |
  |  2   |       Better Move Ordering        |     6/6     |      1_640       |   12_105_638    |      110       |  3_880_170  |
  |  3   |     Skip direct losing moves      |     5/6     |      1_327       |   23_831_884    |       64       |  6_614_864  |
  |  4   | Iterative deepening & null window |     5/6     |      2_673       |   41_519_065    |       77       |  6_562_169  |
  |  5   |        Transposition table        |     5/6     |      7_279       |   103_134_312   |      112       |  5_421_597  |
  |  6   |             Bitboards             |     2/6     |        33        |     591_676     |       74       |  2_860_704  |
  |  7   |  Better column order exploration  |     2/6     |       109        |     591_676     |      227       |  1_169_145  |
  |  8   |        Alpha-Beta prunning        |     0/6     |        -         |        -        |       -        |      -      |
  |  9   |              Negamax              |     0/6     |        -         |        -        |       -        |      -      |

<p align="center">
  <a href="csv/benchmark_results.csv">Download csv here</a>
</p>

> **Note**: Tests have been made on a 64-bit Windows 10 computer with a Ryzen 5 3600 and 16GB of RAM clocked at 3600MHz in go1.24.4 windows/amd64.
>
> In the csv, for all the test with a low mean search time (few ms), the **time per node** and **nodes per second** are not relevant and can vary a lot between runs. This is because when a position is solved in less than a nanosecond, the elapsed time is effectively mesured as zero, leading to misleading statistics.

## 6 - Improve the project

If you like this project and/or want to help or improve it, you can:

- Create an issue if you find a bug or want to suggest a feature or any improvement (no matter how small it is).

- Create a pull request if you want to add a feature, fix a bug or improve the code.

- Contact me if you want to talk about the project or anything else (Discord: pietot).

> **Note**: If you want to be guided/helped, you already have a file named <a href="IMPROVEMENTS.md">IMPROVEMENTS.md</a> in the project directory, where you can see all the improvements that can be made.

## 7 - Credits

- [Original source code](https://github.com/PascalPons/connect4/): The original implementation of the Connect 4 solver.
- [Pascal Pons](http://blog.gamesolver.org/solving-connect-four/12-lower-bound-transposition-table/): The online tutorial I followed to implement the algorithms.
- [Online solver](https://connect4.gamesolver.org/): A web-based Connect 4 solver I used for testing.
