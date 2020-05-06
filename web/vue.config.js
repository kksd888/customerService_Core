// 增加环境变量
process.env.VUE_APP_VERSION = require('./package.json').version;

module.exports = {
    publicPath: process.env.NODE_ENV === "production"
        ? "/static"
        : "/",
    parallel: require('os').cpus().length > 1,
};
