const path = require('path')

const dirName = () => {
    const correct = []

    let splitDir = ""
    if (process.platform === "win32") {
        splitDir = __dirname.split('\\')
    } else {
        splitDir = __dirname.split('/')
    }

    for (let i = 0; i < splitDir.length - 1; i++) {
        correct.push(splitDir[i])
    }

    if (process.platform === "win32") {
        return correct.join("\\")
    } else {
        return correct.join("/")
    }
}

module.exports = {
    entry: ['@babel/polyfill'],
    output: {
        path: dirName() + "/dist"
    },
    resolve: {
        extensions: ['.js'],
        modules: ['node_modules', path.resolve(__dirname, './')]
    },
};