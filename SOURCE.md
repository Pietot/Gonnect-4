# What is Gonnect-4?

This markdown file contains instructions on how to build Gonnect-4 from source, as well as some informations on how the project is structured, and how the code works.

## 📋 Summary

### 1. [Project structure](#1---project-structure)

### 2. [Installation](#2---installation)

### 3. [Build instructions](#3---build-instructions)

## 1 - Project structure

The project is structured in a way that makes it easy to understand and navigate. All package has the name of the folder they are in, except fot `export/main.go`, but I will explain that later.

Here is a brief overview of each folder and its purpose:

- `grid/`: Contains the implementation of the game grid and with all the logic applied to it, and the chore implementation of the solver.

- `evaluation/`: Contains the struct of the evaluation and analysis of the grid, returned respectively by the Solve and Analysis methods.

- `stats/`: Contains the implementation of the statistics struct, used to keep track of some statistics during the solving / analysis process.

- `book/`: Contains the implementation of the creation of the book, which is used to store the analysis of hard positions for instant retrieval.

- `transpositiontable/`: Contains the implementation of the transposition table, which is used to cache previously computed results for performance improvements.

- `data/`: Contains some data files used by the test and benchmark package.

- `test/`: Contains the implementation of the tests.

- `benchmark/`: Contains the implementation of the benchmarking tools.

- `config/`: Not important for now, it just contains a variable to store if the book is enabled or not.

- `utils/`: Contains some utility functions used by the other packages.

- `database/`: Contains the implementation of the database, which is used to store the book on disk when creating it.

- `export/`: Contains the implementation of the export of the book, which is used to export the book in go file to store it directly in the binary when compiling the project.

> **[!NOTE]**
> You can also read some README.md files in some folders for more details on the implementation of the code in those folders.

## 2 - Installation

To install Gonnect-4, you need to have Go installed on your machine.

Then you have to clone the repository and navigate to the project directory:

```bash
git clone git@github.com:Pietot/Gonnect-4.git
cd Gonnect-4
```

After that, you need to download the database for the book in the [latest releases](https://github.com/Pietot/Gonnect-4/releases/latest), and extract the `book.db` file in the `database/` folder.

> **[!NOTE]**
> You can skip the last step if you don't want to generate/use the book, but it will make the analysis of hard positions much slower and you will need to manually update some code (in the next section).

## 3 - Build instructions

If you have followed the installation instructions, you should be able to build the project without any issues.

To build the project, you can use the following command:

```bash
go generate ./...
go build
```

> **[!IMPORTANT]**
> If you want to build the project without the book, you will to comment mutiple lines of code in order to compile the code:

- Comment the lines related to the book in `grid/negamax.go`:

  Line 8 & 12:

  ```go
  "github.com/Pietot/Gonnect-4/database"
  "github.com/Pietot/Gonnect-4/utils"
  ```

  Line 48-60:

  ```go
  scores, found := utils.GetScores(&database.ExportedBook, grid.Key(), grid.MirrorKey())
  if found {
  	score, _ := utils.GetBestScoreAndMove(scores)
  	return evaluation.Evaluation{
  			Score:          &score,
  			RemainingMoves: GetRemainingMoves(score, grid.nbMoves),
  		}, stats.Stats{
  			TotalTimeNanoseconds: 0,
  			NodeCount:            0,
  			MeanTimePerNode:      0,
  			NodesPerSecond:       0,
  		}
  }
  ```

  Line 99-111:

  ```go
  	sc, found := utils.GetScores(&database.ExportedBook, grid.Key(), grid.MirrorKey())
  	if found {
  		scores.Scores = sc
  		_, bestMove = utils.GetBestScoreAndMove(sc)
  		scores.BestMove = &bestMove
  		scores.RemainingMoves = GetRemainingMoves(*sc[bestMove], grid.nbMoves)
  		return scores, stats.Stats{
  			TotalTimeNanoseconds: 0,
  			NodeCount:            0,
  			MeanTimePerNode:      0,
  			NodesPerSecond:       0,
  		}
  	}
  ```

Then you can build the project with the following command:

```bash
go build
```
