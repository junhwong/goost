设计目标：兼容[Syslog RFC5424](https://tools.ietf.org/html/rfc5424)协议规范，结构化，快速(高性能)，简单(接口使用方便)，能配合现在主流的日志存储和分析方案，如：ELK。 [go-boot/log]()主要参考了下面几个库(感谢他们的贡献)。
  * log(go)
  * [zap](https://github.com/uber-go/zap)
  * [logrus](https://github.com/sirupsen/logrus)


```go
// 标准接口
Logger
  .WithContext(Context)                     Loger // 根据上下文获取 Prefix Trace Level 参数
  .WithPrefix(string, joinParent...bool)    Loger //
  .WithTrace(Trace)    Loger  //
  .WithLevel(Level)    Loger  // 默认级别：一般情况下仅有 Print 方法生效
//  .WithCallStack(bool) Loger  // 每一个
  .Print(string,  ...interface{})
  .Debug(string,  ...interface{})
  .Info(string,   ...interface{})
  .Warn(string,   ...interface{})
  .Error(string,  ...interface{})
  .Fatal(string,  ...interface{})
  .Trace(...Field)
```