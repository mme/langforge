const path = require('path');

module.exports = {
  module: {
    rules: [
      {
        test: /\.ipynb$/,
        use: 'raw-loader'
      }
    ]
  },
  resolve: {
    modules: [
      path.resolve(__dirname, 'node_modules'),
      path.resolve(__dirname, 'notebooks')
    ]
  }
};
