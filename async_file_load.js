#!/Users/vivianallen/.nvm/versions/node/v14.17.0/bin/node

import {
  readdir,
  readFile
} from 'fs/promises';
import { performance } from 'perf_hooks'
import * as path from 'path';

const DIR_PATH = "./files_to_load"

async function dirFilesToObj (dirPath) {
    const files = await readdir(dirPath);
    const fileContents = await Promise.all(
      files.map(file => readFile(path.join(dirPath, file), "utf8"))
    )
    return Object.assign(...files.map((k, i)=>({[k]: fileContents[i]}) ))
}


// main
const args = process.argv.slice(2);

var logOut
if ( args.length > 0) {
  logOut = args[0].toLowerCase() === "true" ? true : false
} else {
  logOut = false
}

var dirToLoad
if ( args.length > 1) {
  dirToLoad = args[1]
} else {
  dirToLoad = DIR_PATH
}

var t0 = performance.now()
const fileContents = await dirFilesToObj(dirToLoad);
var t1 = performance.now()

if (logOut) {
  console.log(fileContents);
}

console.log("All files processed in " + (t1 - t0) + " milliseconds.")


