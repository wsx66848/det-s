const webpack = require('webpack')
const path = require('path')
const config = require('../config')

module.exports = {
  entry: {
    vendor: ['vue', 'element-ui']
  },
  output: {
    filename: '[name].dll.js',
    path: config.dll.assetsRoot,
    library: '[name]_dll'
  },
  plugins: [
    new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/),
    new webpack.DllPlugin({
      path: path.join(config.dll.assetsRoot, '[name]-manifest.json'),
      // This must match the output.library option above
      name: '[name]_dll'
    }),
  ],
}