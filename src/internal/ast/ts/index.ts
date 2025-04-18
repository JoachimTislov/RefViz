import { Project, SyntaxKind } from "ts-morph";

// Solution to transport data to the go project
// 1. Transfer via json - normal stragegy
// 2. Transfer via protobuf ?
// 3. Setup a server and respond to requests - gRPC ?
// 4. Use the standard output to transfer data - console log and os.stdin

const project = new Project({
  tsConfigFilePath:
    "C:/Users/joach/OneDrive/Dokumenter/GitHub/Project-Visualizer/sample-code/quickfeed/public/tsconfig.json",
});

export function getSymbolName(filepath: string) {
  const sourceFile = project.getSourceFileOrThrow(filepath);
  const symbols: string[] = [];

  sourceFile.forEachDescendant((node) => {
    const symbol = node.getSymbol();
    if (symbol) {
      symbols.push(symbol.getName() + " (" + node.getKindName() + ")");
    }
  });

  return symbols;
}

function f() {
  // Analyze all files in project
  project.getSourceFiles().forEach((sourceFile) => {
    console.log(`File: ${sourceFile.getFilePath()}`);

    sourceFile.getClasses().forEach((cls) => {
      const name = cls.getName() || "<anonymous>";
      const refs = cls.findReferences();

      refs.forEach((ref) => {});
    });

    sourceFile.getFunctions().forEach((fn) => {
      const name = fn.getName() || "<anonymous>";
      console.log(`\nFunction: ${name}`);

      // Outgoing calls
      console.log("  Calls:");
      fn.forEachDescendant((node) => {
        if (node.getKind() === SyntaxKind.CallExpression) {
          const callExpr = node.asKindOrThrow(SyntaxKind.CallExpression);
          const exprSymbol = callExpr.getExpression().getSymbol();
          if (exprSymbol) {
            console.log("   ➔", exprSymbol.getName());
          }
        }
      });

      // Incoming references
      console.log("  Referenced by:");
      const refSymbols = fn.findReferences();
      refSymbols.forEach((refSymbol) => {
        const definition = refSymbol.getDefinition();
        const references = refSymbol.getReferences();

        references.forEach((ref) => {
          const r = ref.getNode();
        });
      });
    });
  });
}
