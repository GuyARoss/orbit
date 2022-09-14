const path = require('path');

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
    entry: './index.js',
    output: {
        path: dirName() + "/dist"
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /(node_modules|bower_components)/,
                use: {
                    loader: 'swc-loader',
                    options: {
                        jsc: {
                            target: "es5",
                            parser: {
                                syntax: "ecmascript",
                                jsx: true,
                                numericSeparator: false,
                                classPrivateProperty: false,
                                privateMethod: false,
                                classProperty: false,
                                functionBind: false,
                                decorators: false,
                                decoratorsBeforeExport: false
                            },
                            transform: {
                                react: {
                                    pragma: "React.createElement",
                                    pragmaFrag: "React.Fragment",
                                    throwIfNamespace: true,
                                    development: true,
                                    useBuiltins: false
                                },
                                optimizer: {
                                    globals: {
                                        vars: {
                                            __DEBUG__: "true"
                                        }
                                    }
                                }
                            }
                        },
                        module: {
                            type: "es6"
                        },
                        minify: false
                    }
                }
            },
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
                test: /\.html$/,
                use: [
                    {
                        loader: "html-loader",
                    }
                ]
            }

        ],
    },
    resolve: {
        extensions: ['.js', '.jsx'],
        modules: ['node_modules', path.resolve(__dirname, './')]
    },
};