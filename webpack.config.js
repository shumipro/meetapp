var webpack = require('webpack');

module.exports = {
  cache: true,

  watch: true,

  entry: {
    'main': ['./public/js/ma/main.js']
  },

  output: {
    filename: '[name].js'
  },

  devtool: 'inline-source-map',

  module: {
    loaders: [
      { test: /\.js$|\.jsx$/, exclude: /node_modules|build/, loader: 'babel-loader'}
    ]
  },

  plugins: [
    // set jquery as a global variable for waypoints
    new webpack.ProvidePlugin({
      $: "jquery",
      jQuery: "jquery",
      "window.jQuery": "jquery"
    })
  ],

  resolve: {
    root: __dirname,
    alias: {
    },
    extensions: ['', '.js', '.jsx']
  }
};
