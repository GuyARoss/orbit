const { exec } = require('child_process')
const { exit } = require('process')
const { spawn } = require('child_process');
var path = require("path");

const SUPPORTED_PLATFORMS = ['win32', 'darwin', 'linux']

async function main(platform, argv) {
    const forwardArgv = argv.splice(2, process.argv.length)
    if (!SUPPORTED_PLATFORMS.includes(platform)) {
        console.error(`unsupported platform '${platform}'`)
        exit(1)
    }

    return new Promise((res, rej) => {
        exec(`bash ./exec/${platform} ${forwardArgv.join(' ')}`, (error, stdout, stderr) => {
            if (error) {
                rej(error)
            }

            if (stderr) {
                res(stderr)
            }

            res(stdout)
        })
    })
}
(async () => {
    const response = await main(process.platform, process.argv)
    console.log(response)
})()
// const pth = path.resolve("./exec/win32")

// console.log(pth)

// const ls = spawn(pth, process.argv.splice(2, process.argv.length));


// ls.stdout.on('data', (data) => {
//     console.log(`stdout: ${data}`);
// });

// ls.stderr.on('data', (data) => {
//     console.log(`stderr: ${data}`);
// });

// ls.on('close', (code) => {
//     console.log(`child process exited with code ${code}`);
// });

// ls.on('error', (err) => {
//     console.log(err)
// })