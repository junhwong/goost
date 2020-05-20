参考

https://segmentfault.com/a/1190000009906317?utm_source=tag-newest
https://blog.csdn.net/zhangqiang_accp/article/details/85316910
https://www.jianshu.com/p/fe292b15a06a
https://docs.spring.io/spring-statemachine/docs/3.0.0.M2/reference/


解决大量if ifelse else switch 等条件分支
1. 订单状态之间的流转，如：代付款到待发货的转换，中间夹杂付款操作，支付超时等触发事件, 现实的订单业务往往不是简单的FSM而是复杂交错的业务流程，“状态”也只是沦为客户(端)的某个操作指导
2. 订单事件
3. 主动触发，如：关闭订单操作，这类操作依赖特定状态和其它条件才能操作，如关闭订单需要退还之前的支付

例子：
初始创建
待付款：[付款，关闭]
待发货：[发货，取消]
已发货(运输中)：[]
已发货(已签收)
已完成
已取消
已关闭


sm.WithState(代付款).Target(action, gurd)

sm.WithState(代付款)

sm.Operation(发货，参数)结果

sm.Start(stateobj)

```go
//order
func BuildOrderStateMachine(){
    builder:=NewStateMachineBuilder()
    builder.Config(OrderState.Create,type(OrderCreateParams),OrderState.WaitPrepay,type(OrderCreateParams),func(ctx,interface{})(interface{},error))
    builder.Config(OrderState.WaitPrepay,type(Order),[OrderState.WaitPrepay,OrderState.Other...],type(OrderCreateParams),func(ctx,interface{})(interface{},error))
    // 待预付订单金额
    builder.Source(OrderState.WaitPrepay, type(OrderCreateParams))
        .Target(action, gurds...)
        .Target(OrderState.WaitPrepay)
        .Target(OrderState.WaitPrepay)
        .Target(OrderState.WaitPrepay)

    return builder.build()
}
//service
func CreateOrder(ctx, params){
    sm:=BuildOrderStateMachine()
    derfer sm.Done()
    order=orderRepo.Create(params)
    sm.Process(order)
    return sm.SendEvent(OrderState.Create, order)
}
//web api
func orderCreate(http.Request){
    params:=...
    order,err:=orderservice.CreateOrder(prams)
    if err!=nil{

    }
}
//state,event,action,guard
interface StateObject{
    GetState() string
    GetVersion() int64 // 毫秒时间戳
}

type interface State{
    StateDesc()string
    //总是安全的
    StateValue()interface{}
    StateString() string
    StateInt() (uint64,error)
}

type orderInitHandler{}
func(orderInitHandler) Handle(ctx StateContext){
    //1.处理预付
    //2.改变预付状态：成功-待发货，失败-待预付
    //时间：超时-预付超时-关闭
    id:=conv.Int(ctx.Value("order_id"))
    ver:=conv.Int(ctx.Value("updated_time"))
    order:=repo.GetById(id)
    if ver!=order.UpdatedTime{
        return
    }

    return OrderState.WaitPrepay, order,nil
}

```