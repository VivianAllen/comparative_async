#!/Users/vivianallen/.nvm/versions/node/v14.17.0/bin/node

const {
  readdirSync,
  readFileSync
} = require('fs')
const { performance } = require('perf_hooks')
const path = require('path')

const DIR_PATH = "./files_to_load"


function dirFilesToObj (dirPath) {
    const files = readdirSync(dirPath);
    const fileContents = files.map(file => readFileSync(path.join(dirPath, file), "utf8"))
    return Object.assign(...files.map((k, i)=>({[k]: fileContents[i]}) ))
}


function main () {
  const args = process.argv.slice(2)

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
  const fileContents = dirFilesToObj(dirToLoad);
  var t1 = performance.now()

  if (logOut) {
    console.log(fileContents);
  }
  console.log("All files processed in " + (t1 - t0) + " milliseconds.")
}


if (require.main === module) {
  main()
}
