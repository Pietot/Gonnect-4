# Grid Package

This README file is to explain the purpose of the grid package and how to use it.

## 📋 Summary

### 1. [About this package](#1---about-this-package)

### 2. [Core Components](#2---core-components)

- #### 2.1 [How the engine works ?](#21---how-the-engine-works-)

- #### 2.2 [Score system](#22---score-system)

- #### 2.3 [Bitboard Representation](#23---bitboard-representation)

- #### 2.4 [Move Sorter](#24---move-sorter)

## 1 - About this package

The grid package is the heart of the Connect-4 solver. It contains all the structs, methods and functions related to the game state representation and the core engine.

- `grid.go`: Defines the `Grid` struct and methods for manipulating game states
- `negamax.go`: Implements the negamax search algorithm with alpha-beta pruning
- `movesorter.go`: Contains the `MoveSorter` struct for sorting moves by heuristic scores

## 2 - Core Components

- ### 2.1 - How the engine works ?

  For any 2-player, zero-sum and perfect information game, the basic algorithm to use is the negamax algorithm (it's like the minimax algorithm but with a simpler implementation). All the optimizations of the solver and their explanations are described [in this README](../README.md#4---algorithms--optimizations).

- ### 2.2 - Score system

  For the negamax algorithm to work, we need to define a score system. So, we'll need to define a score for a given position. If the two player plays perfectly, the score system is defined as follows:
  - A positive score means that the current player is winning
  - A negative score means that the current player is losing
  - A score of 0 means that the position leads to a draw

  Positive (winning) score can be computed as 22 minus number of stone played by the current player at the end of the game.\
  Negative (losing) score can be computed as number of stone played by the current player at the end of the game minus 22.

  Why 22 ? Because the maximum number of moves in a Connect-4 game is 42, and each player can play at most 21 moves.
  But if the current player can win with it's last stone, the score will be 1, so 1 + 21 = 22, and 22 - 21 = 1.
  So this system means the higher the score is, the quicker the current player can win, and the lower the score is, the quicker the current player can lose.

  For example, in this position, the current player `x` can win in 2 moves with he's 5th stone, so the score is 22-5=17.

  ```text
  . . . . . . .
  . . . . . . .
  . . . . . . .
  . . . . . . .
  . . o o . . o
  . . x x . . x
  ```

  In this other example, the current player `x` can win in 7 moves (trust me), so try to guess the score.

  ```text
  . . . . . . .
  . . . . . . .
  . . . . o . .
  . . o x x . .
  . o x o x . .
  . x o x o . .
  ```

  `x`has played 6 stones and in 7 moves, when he will win the game, he will have played 6+7=13 stones, so the score is 22-13=9.

- ### 2.3 - Bitboard Representation

  Instead of using a traditional 7x6 grid representation, the bitboard uses a clever encoding where each column is represented by 7 bits (6 for the cells + 1 separator).

> [!IMPORTANT]  
> Since a binary number is "built" from right to left, each bucket of 7 bits corresponds to a column, **FROM RIGHT TO LEFT**.

  For more clarity, here is the order of bits to encode for a 7x6 board.

  ```text
  .  .  .  .  .  .  . <- extra bit
  48 41 34 27 20 13 6
  47 40 33 26 19 12 5
  46 39 32 25 18 11 4
  45 38 31 24 17 10 3
  44 37 30 23 16  9 2
  43 36 29 22 15  8 1
  ```

  Now, imagine we have the following position where x is the current player.:

  ```text
  . . . . . . .
  —————————————
  . . . . . . .
  . . . o . . .
  . . x x . . .
  . . o x . . .
  . . o o x . .
  . . o x x o .
  ```

  We can encode this position where

- `CurrentPosition`: Contains 1s for the current player's pieces and 0s for the opponent's pieces
- `Mask`: Contains 1s for all pieces on the board (both players)
- `Bottom`: Contains 1s for the bottom cell of each column (constant used for the key generation)

  ```text
      Grid         CurrentPosition        Mask
  . . . . . . .     0 0 0 0 0 0 0     0 0 0 0 0 0 0
  —————————————     —————————————     —————————————
  . . . . . . .     0 0 0 0 0 0 0     0 0 0 0 0 0 0
  . . . o . . .     0 0 0 0 0 0 0     0 0 0 1 0 0 0
  . . x x . . .     0 0 1 1 0 0 0     0 0 1 1 0 0 0
  . . o x . . .     0 0 0 1 0 0 0     0 0 1 1 0 0 0
  . . o o x . .     0 0 0 0 1 0 0     0 0 1 1 1 0 0
  . . o x x o .     0 0 0 1 1 0 0     0 0 1 1 1 1 0
  ```

  With this encoding, we can efficiently compute winning positions, possible moves, and other operations using bitwise operations. In addition, we can get a unique key for each position by adding the `CurrentPosition` and `Mask` + a constant `BOTTOM`, which is a constant to retrieve a position (CurrentPosition and mask) from the key. The key is used for the transposition table and the book.

  ```text

      Grid         CurrentPosition  +       Mask       +     BOTTOM      =       Key
  . . . . . . .     0 0 0 0 0 0 0       0 0 0 0 0 0 0     0 0 0 0 0 0 0     0 1 1 1 0 0 0
  —————————————     —————————————   +   —————————————     —————————————     —————————————
  . . . . . . .     0 0 0 0 0 0 0   +   0 0 0 0 0 0 0     0 0 0 0 0 0 0     0 0 0 0 0 0 0
  . . . o . . .     0 0 0 0 0 0 0   +   0 0 0 1 0 0 0     0 0 0 0 0 0 0     0 0 0 1 0 0 0
  . . x x . . .     0 0 1 1 0 0 0   +   0 0 1 1 0 0 0     0 0 0 0 0 0 0     0 0 0 0 0 0 0
  . . o x . . .     0 0 0 1 0 0 0   +   0 0 1 1 0 0 0     0 0 0 0 0 0 0     0 0 0 1 0 0 0
  . . o o x . .     0 0 0 0 1 0 0   +   0 0 1 1 1 0 0     0 0 0 0 0 0 0     0 0 0 0 0 0 0
  . . o x x o .     0 0 0 1 1 0 0   +   0 0 1 1 1 1 0     1 1 1 1 1 1 1     0 1 0 1 0 0 0

  Calculus:

  CurrentPosition = 0000000 1000000 1100000 1111100 1111000 0000000 0000000 (2 225 055 072 256)
  +
  Mask            = 0000000 0000000 1100000 1011000 0001000 0000000 0000000    (25 954 484 224)
  +
  BOTTOM          = 0000001 0000001 0000001 0000001 0000001 0000001 0000001 (4 432 676 798 593)
  =
  Key             = 0000001 1000010 1000010 1010110 0000001 0000001 0000001 (6 683 686 355 073)

  ```

- ### 2.4 - Move Sorter

  The `MoveSorter` is a simple priority structure that sorts moves by their heuristic score. Moves are sorted in ascending order and retrieved in descending order (best moves first), which improves alpha-beta pruning efficiency.

  The heuristic scores moves based on how many winning positions they create for the current player.

> [!NOTE]
> Yes, I could have separated the move sorter from the grid package to it's own package, but I need access to the grid package's WIDTH constant, which creates a circular import. I know there's ways to fix this, but none of them sastisfied me. So for instance, the move sorter is in the grid package, and I don't think it's that bad after all.
