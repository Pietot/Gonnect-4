This README file is to explain the purpose of the benchmark package and how to use it.

## 📋 Summary

### 1. [Methodology](#1---methodology)

### 2. [BenchmarkAnalyze](#2---benchmarkanalyze)

### 3. [BenchmarkBookCreation](#3---benchmarkbookcreation)

## 1 - Methodology

The benchmark package uses a set of test files located in the `data` directory.

Each file contains 1000 positions for a total of 6000 positions. The files originate from the work of **Pascal Pons** and are available on his [blog](http://blog.gamesolver.org/solving-connect-four/02-test-protocol/). I just reorganized/renamed them in a way that is more comprehensible and logical in my sens. Each line of each file contains a position represented as a sequence of digits, the score of the position, and the best move for that position. Each separated by a space. For example:

```
74642572132 15 C3
51756773145177 -10 C2
165746146225 -11 C7
1562527227511 13 C4
43745416472735 -13 C6
```

> [!NOTE]
> The sequence is a string of numbers representing the columns where the pieces have been played, starting from an empty grid. The columns are numbered from 1 to 7 (left to right).

The files are categorized based on the depth `d` (number of moves played) and the remaining moves `r` until a forced win.

Here are the datasets sorted from the easiest to the hardest:

| Test Set (1000 test cases each) | Number of played moves (`d`) | Number of remaining moves (`r`) |
| :-----------------------------: | :--------------------------: | :-----------------------------: |
|          test_easy_end          |          moves > 28          |         remaining < 14          |
|        test_easy_middle         |       14 < moves <= 28       |         remaining < 14          |
|         test_easy_begin         |         moves <= 14          |         remaining < 14          |
|         test_medium_end         |       14 < moves <= 28       |      14 <= remaining < 28       |
|       test_medium_middle        |         moves <= 14          |      14 <= remaining < 28       |
|         test_hard_begin         |         moves <= 14          |         remaining >= 28         |

## 2 - BenchmarkAnalyze

The `BenchmarkAnalyze` function is designed to benchmark the performance of the `Analyze`, and so, by extension, the `GetScore` and `negamax` method (the core of the solver)

The function iterates over all the position in the test files, from the easiest to the hardest, analyzes each position using the `Analyze` method, and then compute a mean of these statistics:

- **Mean Total Time**: The average time taken to analyze a position, in nanoseconds.
- **Mean Nodes**: The average number of nodes explored during the analysis of a position.
- **Mean Time per Node**: The average time taken per node explored, in nanoseconds.
- **Mean Nodes per Second**: The average number of nodes explored per second during the analysis.

Usually, you will not need to change the `BenchmarkAnalyze` function, but you can run it to see if your changes on the `GetScore` and `negamax` method brings any performance improvement or not.

> [!IMPORTANT]
> Once changes are made and benchmarked, ydon't forget to run the tests to ensure that your changes do not break the integrity of the solver. You can run the tests using the following command:
>
> ```bash
> go test -v ./...
> ```

## 3 - BenchmarkBookCreation

The `BenchmarkBookCreation` function is designed to benchmark the performance of the `CreateBook` method, which is responsible for creating a book of pre-computed positions and analysis results.

The benchmark then measures the time taken to create a book of depth 8.

> [!WARNING]
> The **BenchmarkBookCreation** is designed to create a new **.db** file called `book_benchmark.db` in the **benchmark** directory.
>
> **_DO NOT COMMIT_** `book_benchmark.db` to the repository. It is only used for benchmarking purposes and should not be part of the repository.

### Benchmark History Results

| Gonnect 4 Version |                Book Creation Technique                |      Time       |
| :---------------: | :---------------------------------------------------: | :-------------: |
|      v1.12.0      | Multi-Threading + Bbolt + Canonical Key + Early Prune | More than a day |

> [!NOTE]
> Benchmarks have been made on a 64-bit Windows 10 computer with a Ryzen 5 3600 and 16GB of RAM clocked at 3600MHz in go1.26.0 windows/amd64.
