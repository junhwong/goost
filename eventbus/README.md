
eventbus.Start(config, opts...) Done
eventbus.Register(infomer)
eventbus.Bordcast(msg, ...filters)
Infomer {
    GetIns()[]string
    Handle(context, conn, msg)
}

config{
    Mode: ServerOnly|ClientOnly|Peer
    Port:
    CAFile:
    CertFile:
    CertKey:
    User:
    Iface:
    CanbeMaster: true
    Accepts:["f", ""]
    Tags:["d"]
}

Conn{
    GetMeta()
    IsLocal() bool // 不能发送信息
    Send(context,msg) error
    Call(context,msg, func(context, conn, response)) error
}

http.paycallback{
    msg:=eventbus.New("")
    errors.Do(func()(err error){
         defer EnsureClose(request.body)
         msg.data, err =ioutil.readAll(request.body)
         return
    })
    // msg.xx=xx
    errors.Do(func()(err error){
        return eventbus.apply(ctx, msg, pay,func(ctx,repl){
            if ctx.error(){
                response.write({false})
                return
            }
            if repl.proced{
                response.write({ok})
            }
        })
    })
    errors.Ex(map[errortype]bizCode).Catch().PanicIf()
}

handle(ctx MessageContext, msg){
    eventbus.Send(msg, eventbus.WithCallback(func(ctx, msg){

    }, 10),eventbus.WithFilter(func(meta)bool{
        return true
    }),eventbus.WithOnlyLocal(),eventbus.WithoutLocal())
}