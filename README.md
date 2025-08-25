# 🔗 Gonnect 4

![Location](https://img.shields.io/badge/Made_in-France-red?labelColor=blue)
![Language](https://img.shields.io/badge/Language-Go-f7d3a2?labelColor=00aed8)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/Pietot/Gonnect-4?label=Latest%20release)

Gonnect 4 is a command-line connect 4 game written in Go. It allows you to play against an unbeatable AI. The game is played in a terminal and uses ASCII characters to represent the board.

## 📋 Summary

### 1. [Features](#1---features)

### 2. [Installation](#2---installation)

### 3. [How to use](#3---how-to-use)

### 4. [Algorithms & optimizations](#4---algorithms--optimizations)

### 5. [Benchmark](#5---benchmark)

- #### 5.1 [Testsets](#5---testsets)
- #### 5.2 [Results](#5---results)

### 6. [Improve the project](#6---improve-the-project)

### 7. [Credits](#7---credits)

## 1 - Features

## 2 - Installation

## 3 - How to use

## 4 - Algorithms & optimizations

- ### Negamax

  The negamax algorithm is a variant of the minimax algorithm that simplifies the implementation of two-player games. It assumes that both players are playing optimally and tries to maximize the score for the current player while minimizing it for the opponent.

  It works by recursively exploring the game tree, evaluating the positions, and returning the best move for the current player.

- ### Alpha-Beta Pruning

  Finished

- ### Move explanation order

  Finished

- ### Bitboard

  Finished

- ### Transposition Table

  Finished

- ### Iterative Deepening

  Finished

- ### Anticipate direct losing moves

  Finished

### Better move ordering

  In progress...

## 5 - Benchmark

- ### Testsets

  To benchmark the different algorithms, I've created 6 datasets from [here](http://blog.gamesolver.org/solving-connect-four/02-test-protocol/) and placed them in the <a href="tests/data/">data</a> directory.

  Here are the datasets:

  | Test Set (1000 test cases each) | Test Set name |     nb moves      |  nb remaining moves   |
  | :-----------------------------: | :-----------: | :---------------: | :-------------------: |
  |           Test_L1_S1            |  Easy-Begin   |    moves <= 11    |    remaining <= 11    |
  |           Test_L1_S2            |  Easy-Middle  | 12 <= moves <= 21 |    remaining <= 11    |
  |           Test_L1_S3            |   Easy-End    | 22 <= moves <= 31 |    remaining <= 11    |
  |           Test_L2_S2            | Medium-Begin  |    moves <= 11    | 12 <= remaining <= 30 |
  |           Test_L2_S3            | Medium-Middle | 12 <= moves <= 21 | 12 <= remaining <= 30 |
  |           Test_L3_S1            |  Hard-Begin   |    moves <= 11    | 31 <= remaining <= 42 |

- ### Results

  I ran the tests on each algorithm and collected the results in a CSV file and computed the **mean** for each category.

  Here are the algorithms ranked from the most efficient to the least efficient:

  | Rank | Algorithm | Search time | Number of positions | Time per position | Positions per second |
  | :--: | :-------: | :---------: | :-----------------: | :---------------: | :------------------: |
  |  1   |           |             |                     |                   |                      |
  |  2   |           |             |                     |                   |                      |
  |  3   |           |             |                     |                   |                      |
  |  4   |           |             |                     |                   |                      |
  |  5   |           |             |                     |                   |                      |
  |  6   |           |             |                     |                   |                      |
  |  7   |  Negamax  |             |                     |                   |                      |

<br>

<!-- image -->

<p align="center">
  <a href="assets/csv/">Download csv here</a>
</p>

> **Note**: Tests have been made on a 64-bit Windows 10 computer with a Ryzen 5 3600 and 16GB of RAM clocked at 3600MHz in go1.22.5 windows/amd64.

## 6 - Improve the project

If you like this project and/or want to help or improve it, you can:

- Create an issue if you find a bug or want to suggest a feature or any improvement (no matter how small it is).

- Create a pull request if you want to add a feature, fix a bug or improve the code.

- Contact me if you want to talk about the project or anything else (Discord: pietot).

> **Note**: If you want to be guided/helped, you already have a file named <a href="improvements.txt">improvements.txt</a> in the project directory, where you can see all the improvements that can be made.

## 7 - Credits
