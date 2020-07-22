const fs = require('fs');

// webpack.config.js
const VueLoaderPlugin = require('vue-loader/lib/plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

// a custom plugin so there arent 1000 bundles when working on frontend
class TinyTidyPlugin {
    constructor(options) {
        this.path;
        this.firstRun = true;
        this.scrubbing = false;
        this.options = options;
    }

    apply = (compiler) => {
        if (compiler.options.mode !== "development") return;
        this.path = compiler.options.output.path;

        compiler.hooks.done.tapAsync("TinyTidyPlugin", this.removeOldFiles);
    };

    removeOldFiles = ({ compilation }, done) => {
        if (this.firstRun) {
            this.firstRun = false;
            return done();
        }

        const builtFiles = Object.keys(compilation.assets);
        this.backgroundScrub(builtFiles);
        done();
    };

    backgroundScrub = (builtFiles) => {
        this.srubbing = true;
        const existingFiles = fs.readdirSync(this.path);

        let currentBundle;
        for (let i = 0; i < builtFiles.length; i++) {
            if (builtFiles[i].startsWith("bundle")) currentBundle = builtFiles[i];
        }

        for (let i = 0; i < existingFiles.length; i++) {
            const file = existingFiles[i];
            if (!file.startsWith("bundle")) continue;
            if (!file.endsWith("js")) continue;

            if (file === currentBundle) continue;

            fs.unlinkSync(this.path + "/" + file);
        }

        this.scrubbing = false;
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