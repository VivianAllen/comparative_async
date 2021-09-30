#!/Users/vivianallen/.nvm/versions/node/v14.17.0/bin/node

const {
  readdirSync,
  readFileSync
} = require('fs')
const {
  readdir,
  readFile
} = require('fs/promises')
const { performance } = require('perf_hooks')
const path = require('path')

const DIR_PATH = "./files_to_load"


function loadFileContentsSync (filePath, results) {
  readFileSync(filePath, "utf8")
  resultsStr = path.basename(filePath)
  process.stdout.write(resultsStr + ', ')
  results.push(resultsStr)
}


async function loadFileContentsAsync (filePath, results) {
  await readFile(filePath, "utf8")
  resultsStr = path.basename(filePath)
  process.stdout.write(resultsStr + ', ')
  results.push(resultsStr)
}


function printResults (results) {
  process.stdout.write('\norder as reported in shared object:\n')
  console.log(results.join(', '))
}


function runTimedSync(func) {
  var start_execution_time_ms = performance.now()
  func();
  var execution_time_s = (performance.now() - start_execution_time_ms) / 1000
  console.log(`Function ${func.name} executed in ${execution_time_s}s\n`)
}


async function runTimedAsync(func) {
  var start_execution_time_ms = performance.now()
  await func();
  var execution_time_s = (performance.now() - start_execution_time_ms) / 1000
  console.log(`Function ${func.name} executed in ${execution_time_s}s\n`)
}


function runIoHeavySync () {
  console.log('Synchronous')
  console.log('order as printed by tasks:')
  const files = readdirSync(DIR_PATH);
  var results = []
  files.map(file => { loadFileContentsSync(path.join(DIR_PATH, file), results) })
  printResults(results)
}


async function runIoHeavyAsync () {
    console.log('Asynchronous')
    console.log('order as printed by tasks:')
    const files = await readdir(DIR_PATH);
    var results = []
    await Promise.all(
      files.map(async file => { await loadFileContentsAsync(path.join(DIR_PATH, file), results) })
    )
    printResults(results)
}


async function main () {
  runTimedSync(runIoHeavySync)
  runTimedAsync(runIoHeavyAsync)
}


if (require.main === module) {
  main()
}

