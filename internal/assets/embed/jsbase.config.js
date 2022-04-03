const path = require('path')

const dirName = () => {
    const correct = []

    const splitDir = __dirname.split('/')
    for (let i = 0; i < splitDir.length - 1; i++) {
        correct.push(splitDir[i])
    }

    return correct.join("/")
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