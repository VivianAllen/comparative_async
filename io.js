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
const {
  Worker,
  isMainThread,
  parentPort,
  workerData
} = require('worker_threads');

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


async function doFileLoadInThread(file, results) {
  // spawn worker that messages back when done
  // return promise that resolves when the message is received
  // NB worker_threads do NOT share memory!
  const worker = new Worker(__filename, { workerData: { file: file, results: results }});
  return new Promise((resolve, reject) => {
    worker.on('message', () => { resolve() })
  });
}


async function runIoHeavyMultithread() {
    console.log('Multithreaded')
    console.log('order as printed by tasks:')
    const files = await readdir(DIR_PATH);
    var results = []
    // start pool of workers, each worker does job, resolves promise then dies?
    await Promise.all(
      files.map(async file => { await doFileLoadInThread(path.join(DIR_PATH, file), results) })
    )
    printResults(results)
}


function runIoHeavyMultiprocess() {

}


async function main () {
  if (isMainThread) {
    runTimedSync(runIoHeavySync)
    await runTimedAsync(runIoHeavyAsync)
    await runTimedAsync(runIoHeavyMultithread)
  } else {
    // something something worker data?
    await loadFileContentsAsync(workerData.file, workerData.results)
    parentPort.postMessage('done')
  }
}


if (require.main === module) {
  main()
}

