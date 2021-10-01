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
const { fork } = require('child_process');

const DIR_PATH = "./files_to_load"


function runTimedSync(func) {
  const start_execution_time_ms = performance.now()
  const results = func();
  const execution_time_s = (performance.now() - start_execution_time_ms) / 1000
  console.log(`Function ${func.name} executed in ${execution_time_s}s`)
  console.log(results)
  console.log('')
}


async function runTimedAsync(func) {
  const start_execution_time_ms = performance.now()
  const results = await func();
  const execution_time_s = (performance.now() - start_execution_time_ms) / 1000
  console.log(`Function ${func.name} executed in ${execution_time_s}s`)
  console.log(results)
  console.log('')
}


function loadFileContentsSync (filePath) {
  const file_contents = readFileSync(filePath, "utf8")
  return {
    'filename': path.basename(filePath),
    'file_contents_length': file_contents.length,
    'processing_time': new Date().toISOString()
  }
}


async function loadFileContentsAsync (filePath, results) {
  var file_contents = await readFile(filePath, "utf8")
  return {
    'filename': path.basename(filePath),
    'file_contents_length': file_contents.length,
    'processing_time': new Date().toISOString()
  }
}


function ioHeavySync () {
  const files = readdirSync(DIR_PATH);
  const results = files.map(file => { return loadFileContentsSync(path.join(DIR_PATH, file)) })
  return results
}


async function ioHeavyAsync () {
    const files = await readdir(DIR_PATH);
    var results = await Promise.all(
      files.map(async file => { return await loadFileContentsAsync(path.join(DIR_PATH, file)) })
    )
    return results
}


async function doFileLoadInThread(file) {
  // spawn worker that messages back when done
  // return promise that resolves when the message is received
  // NB worker_threads do NOT share memory!
  const workerThread = new Worker(__filename, { workerData: { file: file }});
  return new Promise((resolve, reject) => {
    workerThread.on('message', (msg) => { resolve(msg) })
  });
}


async function ioHeavyMultithread() {
    const files = await readdir(DIR_PATH);
    var results = await Promise.all(
      files.map(async file => { return await doFileLoadInThread(path.join(DIR_PATH, file)) })
    )
    return results
}


async function doFileLoadInProcess(file) {
  // spawn worker that messages back when done
  // return promise that resolves when the message is received
  // NB worker processes do NOT share memory!
  const workerProcess = fork(__filename, ['child']);
  return new Promise((resolve, reject) => {
    workerProcess.send(file)
    workerProcess.on('message', (msg) => { resolve(msg) })
  });
}


async function ioHeavyMultiprocess() {
  const files = await readdir(DIR_PATH);
  var results = await Promise.all(
    files.map(async file => { return await doFileLoadInProcess(path.join(DIR_PATH, file)) })
  )
  return results
}


async function main () {
  const isMainProcess = !(process.argv.length == 3 && process.argv[2] === 'child')
  if (isMainThread && isMainProcess)  {
    runTimedSync(ioHeavySync)
    await runTimedAsync(ioHeavyAsync)
    await runTimedAsync(ioHeavyMultithread)
    await runTimedAsync(ioHeavyMultiprocess)
  } else if (isMainProcess) {
    // something something worker data?
    var results = await loadFileContentsAsync(workerData.file)
    parentPort.postMessage(results)
    process.exit()
  } else if (isMainThread) {
    process.on('message', async (msg) => {
      var results = await loadFileContentsAsync(msg)
      process.send(results)
      process.exit()
    })
  }
}


if (require.main === module) {
  main()
}

