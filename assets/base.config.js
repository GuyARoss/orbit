const dirName = () => {
    const correct = []

    const splitDir = __dirname.split('\\')
    for (let i = 0; i < splitDir.length - 1; i++) {
        correct.push(splitDir[i])
    }

    return correct.join("\\")
}

module.exports = {
    entry: ['@babel/polyfill'],
    output: {
        path: dirName() + "\\dist"
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: [
                    'style-loader',
                    'css-loader'
                ]
            },
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader"
                }
            },
            {
                test: /\.html$/,
                use: [
                    {
                        loader: "html-loader"
                    }
                ]
            }
        ]
    },
};