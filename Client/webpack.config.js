const path = require('path');
const webpack = require('webpack');

module.exports = {
  entry: './GBA.js',
  output: { path: __dirname, filename: 'GBA.min.js' },
  watch: true,

  module: {
    loaders: [
      {
        test: /.jsx?$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
        query: {
          presets: ['es2015', 'react'],
        },
      },
    ],
  },
};
