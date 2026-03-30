# What is Gonnect-4?

This markdown file contains instructions on how to build Gonnect-4 from source, as well as some informations on how the project is structured, and how the code works.

> [!IMPORTANT]
> All the instruction and commands in all the mardown files are meant to be run on Windows, but they should work on other operating systems as well, with some minor adjustments.

## 📋 Summary

### 1. [Project structure](#1---project-structure)

### 2. [Installation](#2---installation)

### 3. [Build instructions](#3---build-instructions)

- #### 3.1 [CLI tool](#31---cli-tool)
- #### 3.2 [Export the book](#32---export-the-book)
- #### 3.3 [Web assembly tool](#33---web-assembly-tool)

## 1 - Project structure

The project is structured in a way that makes it easy to understand and navigate. All package has the name of the folder they are in, except for the files in `cmd/`, but I will explain that later.

Here is a brief overview of each folder and its purpose:

- `grid/`: Contains the implementation of the game grid and with all the logic applied to it, and the chore implementation of the solver.

- `evaluation/`: Contains the struct of the evaluation and analysis of the grid, returned respectively by the Solve and Analysis methods.

- `stats/`: Contains the implementation of the statistics struct, used to keep track of some statistics during the solving / analysis process.

- `book/`: Contains the implementation of the creation of the book, which is used to store the analysis of hard positions for instant retrieval.

- `transpositiontable/`: Contains the implementation of the transposition table, which is used to cache previously computed results for performance improvements.

- `data/`: Contains some data files used by the test and benchmark package.

- `test/`: Contains the implementation of the tests.

- `benchmark/`: Contains the implementation of the benchmarking tools.

- `config/`: Not important for now, it just contains some variables.

- `utils/`: Contains some utility functions used by the other packages.

- `database/`: Contains the implementation of the database, which is used to store the book on disk when creating it.

- `progressbar/`: Contains the implementation of the progress bar, used to display the progress of the book creation.

- `cmd/`: Contains the implementation of the command line tools. All the files in this folder are in the `main` package, because they are meant to be compiled as standalone tools, and not to be imported as packages.

- `cmd/cli/`: Contains the implementation of the command line interface, which is used to solve and analyze positions from the command line.

- `cmd/export/`: Contains the implementation of the export of the book, which is used to export the book in a go file to store it directly in the binary when compiling the project.

- `cmd/web/`: Contains the implementation of the web assembly tool, which is used to create a web assembly version of the solver to be used in the browser.

> [!NOTE]
> You can also read some README.md files in some folders for more details on the implementation of the code in those folders.

## 2 - Installation

To install Gonnect-4, you need to have Go installed on your machine.

Then you have to clone the repository and navigate to the project directory:

```bash
git clone git@github.com:Pietot/Gonnect-4.git
cd Gonnect-4
```

## 3 - Build instructions

### 3.1 - CLI tool

    If you have followed the installation instructions, you should be able to build the project without any issues.

    To build the project, you can use the following command:

    ```bash
    go build -ldflags="-s -w" -o gonnect4.exe cmd/cli/main.go
    ```

### 3.2 - Export the book

    If you built your own book and want to export it, you can use the following command:

    ```bash
    go generate ./...
    ```

    Then, you can build the project as usual, and the exported book will be included in the binary.

### 3.3 - Web assembly tool

    To build the web assembly version of the solver, you can use the following command:

    ```bash
    $env:GOOS="js"; $env:GOARCH="wasm"; go build -ldflags="-s -w" -o vercel/js/gonnect4.wasm cmd/web/main.go
    ```

    This will create a `gonnect4.wasm` in the `vercel/js` directory that you can use in your web projects to solve and analyze positions in the browser.
