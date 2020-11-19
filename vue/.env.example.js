module.exports = {
  proxys: {
    '/upload': {
      // target: 'http://10.3.220.200:8003',
      target: 'http://127.0.0.1:8003',
      changeOrigin: true
    }
  }
}
