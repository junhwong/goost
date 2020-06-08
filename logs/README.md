设计目标：兼容[Syslog RFC5424](https://tools.ietf.org/html/rfc5424)协议规范，结构化，快速(高性能)，简单(接口使用方便)，能配合现在主流的日志存储和分析方案，如：ELK。 [go-boot/log]()主要参考了下面几个库(感谢他们的贡献)。
  * log(go)
  * [zap](https://github.com/uber-go/zap)
  * [logrus](https://github.com/sirupsen/logrus)
https://opentracing.io/
https://www.w3.org/TR/trace-context/
http://www.chinaw3c.org/
https://opentelemetry.io/

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
https://github.com/openzipkin/zipkin-go/blob/master/model/traceid.go
type TraceID struct {
	High uint64
	Low  uint64
}
func (t *randomTimestamped) TraceID() (id model.TraceID) {
	seededIDLock.Lock()
	id = model.TraceID{
		High: uint64(time.Now().Unix()<<32) + uint64(seededIDGen.Int31()),
		Low:  uint64(seededIDGen.Int63()),
	}
	seededIDLock.Unlock()
	return
}
func (t TraceID) String() string {
	if t.High == 0 {
		return fmt.Sprintf("%016x", t.Low)
	}
	return fmt.Sprintf("%016x%016x", t.High, t.Low)
}
func TraceIDFromHex(h string) (t TraceID, err error) {
	if len(h) > 16 {
		if t.High, err = strconv.ParseUint(h[0:len(h)-16], 16, 64); err != nil {
			return
		}
		t.Low, err = strconv.ParseUint(h[len(h)-16:], 16, 64)
		return
	}
	t.Low, err = strconv.ParseUint(h, 16, 64)
	return
}
func (t *randomTimestamped) SpanID(traceID model.TraceID) (id model.ID) {
	if !traceID.Empty() {
		return model.ID(traceID.Low)
	}
	seededIDLock.Lock()
	id = model.ID(seededIDGen.Int63())
	seededIDLock.Unlock()
	return
}
https://github.com/openzipkin/zipkin-go/blob/master/model/span_id.go
// ID type
type ID uint64

// String outputs the 64-bit ID as hex string.
func (i ID) String() string {
	return fmt.Sprintf("%016x", uint64(i))
}

// MarshalJSON serializes an ID type (SpanID, ParentSpanID) to HEX.
func (i ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", i.String())), nil
}

// UnmarshalJSON deserializes an ID type (SpanID, ParentSpanID) from HEX.
func (i *ID) UnmarshalJSON(b []byte) (err error) {
	var id uint64
	if len(b) < 3 {
		return nil
	}
	id, err = strconv.ParseUint(string(b[1:len(b)-1]), 16, 64)
	*i = ID(id)
	return err
}


http{

func(Request)
  log:=logs.Start(ctx, path)
  defer log.End()
  log=logs.FromContext(ctx) // trace{traceid=1,name=request,time=...,msg=start}
  defer log.End() // // trace{traceid=1,name=request,time=...,msg=end,du=10ms}
  log.SetName("")
  log.debug()
  subcall(log.ctx,12324343)

}

subcall(ctx,...){
  log=logs.From(ctx,ops)


}
jeager
https://github.com/jaegertracing/jaeger
https://www.elastic.co/cn/apm
10.1.1.1 - - [09/Feb/2019:22:41:28 +0800] "GET /index.html HTTP/1.1" 200 612 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"


\[%{SYSLOGTIMESTAMP:timestamp}\]\t\[%{IP:ipaddress}:%{INT}\]\t\[%{USERNAME:verb}\]\t?%{BASE16NUM:hex}?.*