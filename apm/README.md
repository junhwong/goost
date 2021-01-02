# APM documents

简化并实现应用程序可观测性。

```go

func main(){
  ch := runtime.ListenTerminate()
  apm.SetHandler(custom) //
  go apm.RunOrDie(settings, ch)

  //其它初始化...
}


counter:=apm.Counter("name", fields...)


func biz(){
  options.WithType()
  options.WithName()
  options.WithBadges(...)
  options.WithFields(...)
  ctx, span:=apm.FromContex(ctx,...options)
  defer span.End()
  apm.Debug("called") // global loging
  span.Fail() // 标记该事务执行失败，但不会结束事务执行
  childCall(ctx,...){
    ctx, span:=apm.FromContex(ctx,...options)
    defer span.End()
    span.Debug() // loging
    counter.Inc(ctx, fields...) // metrics
  }
}

```

## Tracing

## Logging

## Metrics
