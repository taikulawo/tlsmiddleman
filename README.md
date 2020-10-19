一个 `HTTPS` 流量解析器，用于演示用途 :)

`go build .`会编译出tlsmiddleman的可执行文件。设置HTTP代理到本地的8000端口，然后访问网站就可以解密出TLS流量

`TLS` 生成证书部分参考自 `lantern`

> https://github.com/Dukou007/lantern/blob/fd97e0e047d1fd491b5739e9501419b8ea78b7dd/src/github.com/getlantern/keyman/keyman.go#L150

有一个新的想法，借助现有的HTTPS加密解密可以做成一个工具，方面我们控制流量传输，比如HTTP chunked，可以观察chunked对前端页面的渲染UI影响

之前读过Netty的代码，ChannelPipeline设计很不错，我可以参考自己也搞一个，将中间件挂载上去流式处理

https://netty.io/4.0/api/io/netty/channel/ChannelPipeline.html

```
                                                 I/O Request
                                            via Channel or
                                        ChannelHandlerContext
                                                      |
  +---------------------------------------------------+---------------+
  |                           ChannelPipeline         |               |
  |                                                  \|/              |
  |    +---------------------+            +-----------+----------+    |
  |    | Inbound Handler  N  |            | Outbound Handler  1  |    |
  |    +----------+----------+            +-----------+----------+    |
  |              /|\                                  |               |
  |               |                                  \|/              |
  |    +----------+----------+            +-----------+----------+    |
  |    | Inbound Handler N-1 |            | Outbound Handler  2  |    |
  |    +----------+----------+            +-----------+----------+    |
  |              /|\                                  .               |
  |               .                                   .               |
  | ChannelHandlerContext.fireIN_EVT() ChannelHandlerContext.OUT_EVT()|
  |        [ method call]                       [method call]         |
  |               .                                   .               |
  |               .                                  \|/              |
  |    +----------+----------+            +-----------+----------+    |
  |    | Inbound Handler  2  |            | Outbound Handler M-1 |    |
  |    +----------+----------+            +-----------+----------+    |
  |              /|\                                  |               |
  |               |                                  \|/              |
  |    +----------+----------+            +-----------+----------+    |
  |    | Inbound Handler  1  |            | Outbound Handler  M  |    |
  |    +----------+----------+            +-----------+----------+    |
  |              /|\                                  |               |
  +---------------+-----------------------------------+---------------+
                  |                                  \|/
  +---------------+-----------------------------------+---------------+
  |               |                                   |               |
  |       [ Socket.read() ]                    [ Socket.write() ]     |
  |                                                                   |
  |  Netty Internal I/O Threads (Transport Implementation)            |
  +-------------------------------------------------------------------+
```
