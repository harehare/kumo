const path = require("path");
const webpack = require("webpack");

module.exports = {
  entry: {
    index: "./src/index.ts"
  },
  module: {
    rules: [
      {
        test: /\.ts$/,
        exclude: /node_modules/,
        use: "ts-loader"
      }
    ]
  },
  plugins: [
    new webpack.EnvironmentPlugin([
      "FIREBASE_API_KEY",
      "FIREBASE_AUTH_DOMAIN",
      "FIREBASE_PROJECT_ID",
      "FIREBASE_APP_ID"
    ])
  ],
  output: {
    path: path.resolve(__dirname, "../wwwroot/js"),
    library: "Kumo",
    libraryTarget: "umd",
    filename: "kumo.js",
    globalObject: "this"
  },
  resolve: {
    extensions: [".ts"]
  }
};
