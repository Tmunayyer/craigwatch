const fs = require('fs');

// webpack.config.js
const VueLoaderPlugin = require('vue-loader/lib/plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

// a custom plugin so there arent 1000 bundles when working on frontend
class TinyTidyPlugin {
    constructor(options) {
        this.options = options;
    }

    apply = (compiler) => {
        if (compiler.options.mode !== "development") return;

        const callback = this._doSomething.bind(this, compiler);
        compiler.hooks.beforeCompile.tapAsync('TinyTidyPlugin', callback);
    };

    _doSomething = function (compiler, details, done) {
        // remove old bundles according to filename scheme
        console.log("removing old bundles...");
        const files = fs.readdirSync(compiler.outputPath);
        const bundleName = compiler.options.output.filename.split(".");
        const prefix = bundleName[0];

        for (let i = 0; i < files.length; i++) {
            if (files[i].startsWith(prefix) && files[i].endsWith('.js')) {
                fs.unlinkSync(compiler.outputPath + "/" + files[i]);
            }
        }

        done();
    };
};

module.exports = {
    entry: './src/index.js',
    output: {
        filename: 'bundle.[contenthash].js'
    },
    mode: process.env.MODE || "development",
    module: {
        rules: [
            {
                test: /\.vue$/,
                loader: 'vue-loader'
            },
            // this will apply to both plain `.js` files
            // AND `<script>` blocks in `.vue` files
            {
                test: /\.js$/,
                loader: 'babel-loader'
            },
            // this will apply to both plain `.css` files
            // AND `<style>` blocks in `.vue` files
            {
                test: /\.css$/,
                use: [
                    'vue-style-loader',
                    'css-loader'
                ]
            }
        ]
    },
    plugins: [
        // make sure to include the plugin for the magic
        new VueLoaderPlugin(),
        new HtmlWebpackPlugin({
            template: './src/index.html'
        }),
        new TinyTidyPlugin(),
    ]
};