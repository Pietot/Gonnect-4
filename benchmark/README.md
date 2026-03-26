This README file is to explain the purpose of the benchmark package and how to use it.

## 📋 Summary

### 1. [Methodology](#1---methodology)

### 2. [Benchmark Evaluation](#2---benchmark-evaluation)

- #### 2.1 [BenchmarkSolve History Results](#21---benchmarksolve-history-results-no-book)

- #### 2.2 [BenchmarkAnalyze History Results](#22---benchmarkanalyze-history-results-no-book)

### 3. [BenchmarkBookCreation](#3---benchmarkbookcreation)

- #### 3.1 [BenchmarkBookCreation History Results](#31---benchmarkbookcreation-history-results)

### 4. [How to benchmark](#4---how-to-benchmark)

## 1 - Methodology

The benchmark package uses a set of test files located in the `data` directory.

Each file contains 1000 positions for a total of 6000 positions. The files originate from the work of **Pascal Pons** and are available on his [blog](http://blog.gamesolver.org/solving-connect-four/02-test-protocol/). I just reorganized/renamed them in a way that is more comprehensible and logical in my sens. Files are just an array of json objects in this format:

```json
[
  {
    "sequence": "274552224131661",
    "score": 0,
    "analysis": [-9, -11, -12, 0, -11, -11, -11]
  },
  {
    "sequence": "5455174361263362",
    "score": -1,
    "analysis": [-12, -1, -12, -13, -12, -12, -12]
  },
  {
    "sequence": "2531276566711153",
    "score": 2,
    "analysis": [-2, 2, -1, 0, -1, 0, -12]
  },
  {
    "sequence": "37313333717124171162542",
    "score": 3,
    "analysis": [null, -7, null, 3, -8, -2, -7]
  },
  ...
]
```

> [!NOTE]
> The sequence is a string of numbers representing the columns where the pieces have been played, starting from an empty grid. The columns are numbered from 1 to 7 (left to right). null values in the analysis array indicate that it's column is full, so it is not possible to play in that column.

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

## 2 - Benchmark Evaluation

The `BenchmarkSolve` function is designed to benchmark the performance of the `Solve` method, which is responsible for solving a given position and returning its score and remaining until a forced win.

The `BenchmarkAnalyze` function is designed to benchmark the performance of the `Analyze`, and so, by extension, the `GetScore` and `negamax` method (the core of the solver)

These functions iterates over all the position in the test files, from the easiest to the hardest, solve/analyzes each position, and then compute a mean of these statistics:

- **Mean Total Time**: The average time taken to solve/analyze a position, in nanoseconds.
- **Mean Nodes**: The average number of nodes explored during the solving/analysis of a position.
- **Mean Time per Node**: The average time taken per node explored, in nanoseconds.
- **Mean Nodes per Second**: The average number of nodes explored per second during the solving/analysis.

Usually, you will not need to change the `runBenchmark` function, but you can run it to see if your changes on the `GetScore` and `negamax` method brings any performance improvement or not.

> [!IMPORTANT]
> Once changes are made and benchmarked, don't forget to run the tests to ensure that your changes do not break the integrity of the solver. You can run the tests using the following command:

```bash
go test -v ./...
```

### 2.1 - BenchmarkSolve History Results (no book)

| Gonnect 4 Version | Improvements | Total time |
| :---------------: | :----------: | :--------: |
|      v1.14.0      |      -       |    2h52    |

### 2.2 - BenchmarkAnalyze History Results (no book)

| Gonnect 4 Version | Improvements | Total time |
| :---------------: | :----------: | :--------: |
|      v1.14.0      |      -       |    9h44    |

## 3 - BenchmarkBookCreation

The `BenchmarkBookCreation` function is designed to benchmark the performance of the `CreateBook` method, which is responsible for creating a book of pre-computed positions and analysis results.

The benchmark then measures the time taken to create a book of depth 8.

> [!WARNING]
> The **BenchmarkBookCreation** is designed to create a new **db** folder called `gonnect4_benchmark_db` in the **benchmark** directory.

### 3.1 - BenchmarkBookCreation History Results

| Gonnect 4 Version |                Book Creation Techniques                 | Depth reached - Progress | Time |
| :---------------: | :-----------------------------------------------------: | :----------------------: | :--: |
|      v1.12.0      | Multi-Threading + Bbolt + Canonical Key + early pruning |          4-5/8           | 11h  |
|      v1.14.0      | Mono-Threading + Badger + Canonical Key + early pruning |  4/8 - 231/568 (40.67%)  | 11h  |

> [!NOTE]
> Benchmarks have been made on a 64-bit Windows 10 computer with a Ryzen 5 3600 and 16GB of RAM clocked at 3600MHz in go1.26.0 windows/amd64.

## 4 - How to benchmark

To run the benchmarks, go to the `main.go` file. Then remove all the code in the `main` function and replace it with the benchmark functions you want to run. For example, if you want to run all the benchmarks, replace the `main` function with this code:

```go
func main() {
  benchmark.BenchmarkSolve()
  benchmark.BenchmarkAnalyze()
  benchmark.BenchmarkBookCreation()
}
```

> [!IMPORTANT]
> Don't forget to clean the imports but if you have a decent IDE, it should be done automatically.

Then, run the program by using the following command:

```bash
go run .
```
