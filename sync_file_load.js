import {
  readdirSync,
  readFileSync
} from 'fs';
import { performance } from 'perf_hooks'
import * as path from 'path';

const DIR_PATH = "./files_to_load"

function dirFilesToObj (dirPath) {
    const files =  readdirSync(dirPath);
    const fileContents = files.map(file => readFileSync(path.join(dirPath, file), "utf8"))
    return Object.assign(...files.map((k, i)=>({[k]: fileContents[i]}) ))
}


//
try {
  const args = process.argv.slice(2);
  var dirToLoad
  var logOut
  if ( args.length > 0) {
    dirToLoad = args[0]
  } else {
    dirToLoad = DIR_PATH
  }
  if ( args.length > 1) {
    logOut = args[1].toLowerCase == "true" ? true : false
  } else {
    logOut = true
  }
  var t0 = performance.now()
  const fileContents = dirFilesToObj(dirToLoad);
  var t1 = performance.now()
  if (logOut) {
    console.log(fileContents);
  }
  console.log("All files processed in " + (t1 -t0) + " milliseconds.")
} catch (err) {
  console.error(err);
}
