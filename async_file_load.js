#!/Users/vivianallen/.nvm/versions/node/v14.17.0/bin/node

const {
readdir,
  readFile
} = require('fs/promises')
const { performance } = require('perf_hooks')
const path = require('path')

const DIR_PATH = "./files_to_load"


async function dirFilesToObj (dirPath) {
    const files = await readdir(dirPath);
    const fileContents = await Promise.all(
      files.map(file => readFile(path.join(dirPath, file), "utf8"))
    )
    return Object.assign(...files.map((k, i)=>({[k]: fileContents[i]}) ))
}


async function main () {
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
  const fileContents = await dirFilesToObj(dirToLoad);
  var t1 = performance.now()

  if (logOut) {
    console.log(fileContents);
  }
  console.log("All files processed in " + (t1 - t0) + " milliseconds.")
}


if (require.main === module) {
  main()
}

