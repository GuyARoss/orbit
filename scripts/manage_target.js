const fs = require('fs');
const { exit } = require('process')

const SUPPORTED_PLATFORMS = ['win32', 'darwin', 'linux']

async function main(platform) {
    if (!SUPPORTED_PLATFORMS.includes(platform)) {
        console.error(`unsupported platform '${platform}'`)
        exit(1)
    }


    fs.copyFile(`bin/exec/${platform}`, `bin/main`, (err) => {
        if (err) throw err;
    });
}

(async () => {
    const response = await main(process.platform, process.argv)
    console.log(response)
})()