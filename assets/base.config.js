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
    module: {
        rules: [
            {
                test: /\.css$/i,
                exclude: /node_modules/,
                use: [
                    'style-loader',
                    {
                        loader: 'css-loader',
                        options: {
                            modules: true,
                        },
                    },
                ],
            },
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader",
                    options: {
                        "presets": [
                            [
                                "@babel/preset-env",
                                {
                                    "useBuiltIns": "entry"
                                }
                            ],
                            "@babel/preset-react"
                        ],
                        "plugins": [
                            "@babel/plugin-proposal-class-properties",
                            "@babel/plugin-proposal-export-default-from",
                            "react-hot-loader/babel"
                        ]
                    }
                }
            },
            {
                test: /\.html$/,
                use: [
                    {
                        loader: "html-loader",
                    }
                ]
            }
        ]
    },
    resolve: {
        extensions: ['.js', '.jsx'],
        modules: ['node_modules', path.resolve(__dirname, './')]
    },
};