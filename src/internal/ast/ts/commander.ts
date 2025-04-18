import { program } from 'commander';
import { getSymbolName } from './index';

program
  .name('ast-util')
  .description('A utility for gathering information about porjects')
  .version('0.0.1')

program.command('symbols')
  .description('Retrieve symbols from a file')
  .argument('<path to file>', 'File to retrieve symbols from')
  .action((file: string) => {
    console.log(`Retrieving symbols from ${file}`);

    const symbols = getSymbolName(file);
    if (symbols.length > 0) {
      console.log(`Found ${symbols.length} symbols:`);
      symbols.forEach(symbol => {
        console.log(`- ${symbol}`);
      });
    }
  });
  
program.parse();
