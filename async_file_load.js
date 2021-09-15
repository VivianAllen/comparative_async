import {
  readdir,
  readFile
} from 'fs/promises';
import * as path from 'path';

const DIR_PATH = "./files_to_load"

async function dirFilesToObj (dirPath = DIR_PATH) {
    const files = await readdir(dirPath);
    const fileContents = await Promise.all(
      files.map(file => readFile(path.join(dirPath, file), "utf8"))
    )
    return Object.assign(...files.map((k, i)=>({[k]: fileContents[i]}) ))
}



// get list of files to load
function filesToObj (dirPath) {
  /*
  Scan dirPath for files and run all through appendFileDataToObj to collate fileData,
  then return as object.
  */
  var output = {}
  const filesInDir = fs.readdirSync(dirPath)
  filesInDir.forEach(fileName => appendFileDataToObj(output, fileName, dirPath))
  return output
}

// spawn async processes for each file
// load contents to single array
// display
try {
  const files = await dirFilesToObj();
  console.log(files);
} catch (err) {
  console.error(err);
}
