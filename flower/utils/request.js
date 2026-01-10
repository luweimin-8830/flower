// 配置常量
const CLOUD_ENV = 'prod-0gr2o3qpe533f1fb'; // 替换为你的真实环境ID
const SERVICE_NAME = 'golang-17j9'; // 你的服务名称

/**
 * 通用云托管请求方法
 * @param {String} path 请求路径，例如 '/api/user'
 * @param {Object} data 请求参数
 * @param {String} method 请求方法，默认 POST
 */
export const callContainer = (path, data = {}, method = 'POST') => {
    return new Promise((resolve, reject) => {
        // 自动处理路径，确保以 / 开头
        const formatPath = path.startsWith('/') ? path : `/${path}`;

        wx.cloud.callContainer({
            config: {
                env: CLOUD_ENV,
            },
            path: formatPath,
            method: method,
            header: {
                'X-WX-SERVICE': SERVICE_NAME,
                'content-type': 'application/json'
            },
            data: data,
            success: (res) => {
                // HTTP 状态码判断
                if (res.statusCode >= 200 && res.statusCode < 300) {
                    // callContainer 返回的数据在 res.data 中
                    resolve(res.data);
                } else {
                    reject(res);
                }
            },
            fail: (err) => {
                console.error(`请求 ${path} 失败`, err);
                reject(err);
            }
        });
    });
};
