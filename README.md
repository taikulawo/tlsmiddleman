一个**完全可用**的 `HTTPS` 流量解析器

使用方式
1. 运行 tlsmiddleman 会在当前目录生成证书文件，将证书文件添加到OS的证书信任列表
2. 设置HTTP代理到本地的8000端口，使用http 8000代理访问网站，控制台则会打印解密后的https内容

`TLS` 生成证书部分参考自 `lantern`

> https://github.com/Dukou007/lantern/blob/fd97e0e047d1fd491b5739e9501419b8ea78b7dd/src/github.com/getlantern/keyman/keyman.go#L150

