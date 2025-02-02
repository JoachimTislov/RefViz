# Project-Visualizer

Uses references to map out a graph of the code base, maps are flexible and can be for a single file or folder.

## About and recommended usage

Creating maps requires extracting symbols and their references which a json file will keep track of. The cache will be checked each time a new map is generated to boost preformance, I recommend scanning the whole code base and fill up the cache to prevent long generation time. The cache will each time check if the modification value of the file differs, and if they do, scan the file again. 

## Dependencies

- [Gopls CLI](https://github.com/golang/tools/blob/master/gopls/doc/command-line.md)

## Supported Languages

- [Golang](https://go.dev/doc/)

## Unsupported languages which I want to support in the future

- [Python](https://docs.python.org/3/)
- [C](https://devdocs.io/c/)/[C++](https://devdocs.io/cpp/)
- [C#](https://learn.microsoft.com/en-us/dotnet/csharp/tour-of-csharp/)
- [Rust](https://doc.rust-lang.org/beta/)
- [Typescript](https://www.typescriptlang.org/fr/docs/)/[Javascript](https://devdocs.io/javascript/)
- [Java](https://docs.oracle.com/en/java/)

Im uncertain if theres CLI tools such as Gopls, offering simple `References` and `symbols` commands.

## Supported graph types

- [Graphviz](https://graphviz.org/documentation/)

### Potential new graph types

- [D3.js](https://d3js.org/)
- [Plotly.js](https://plotly.com/javascript/)

## Potential tools to add

- Python: [Pyright](https://github.com/microsoft/pyright/blob/main/docs/command-line.md)
- C/C++: [clangd](https://clangd.llvm.org/)
- Typescript/Javascript: [tsserver](https://github.com/typescript-language-server/typescript-language-server)
- Rust: [rust-analyzer](https://rust-analyzer.github.io/)
- C#: [omnisharp](https://www.omnisharp.net/)
- Java: [jdtls](https://github.com/eclipse-jdtls/eclipse.jdt.ls)
  
## End goal

Create a system which is capable of visualizing any project's code base. Realistically support the more popular languages listed above. 

### Subgoals

- Finish the project for just Go code, since Gopls offer a fantastic CLI to extract references for symbols which can also be extracted :)
- Write readable and structured code
- Document the process/code 
- Implement various CLI tools for finding references.
- Optionally write generic code which can find references for any language, but thats not my priority, not sure if its even possible
- Establish communication with LSP servers for unsupported languages, especially when libraries lack CLI support.
