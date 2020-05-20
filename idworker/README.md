
简单方便，可扩展的分布式ID生成方式，默认实现 twitter snowflake 无依赖算法。

twitter snowflake算法是兼容Mysql ID并随时间递增无重复的高效率算法，它的基本逻辑是将int64位切分位3-4个部分：
符号 时间戳 机房 节点 序列号
然而，大多数公司或项目没有这么高的流量并且是js不安全，所以提供简短的方式：
符号 时间戳 节点 序列号


workerGen(){
    redis:=...
    redis.lock("woker")
    id:=redis.get("woker","uuid")
    if id==0{
        id=redis.len("workers")
        id+=1
        redis.ladd("workers","uuid",1)
    }
    redis.unlock("woker")


}



一些实现参考：
https://www.cnblogs.com/relucent/p/4955340.html
https://blog.csdn.net/ycb1689/article/details/89331634?depth_1-utm_source=distribute.pc_relevant.none-task&utm_source=distribute.pc_relevant.none-task
https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Language_Resources
https://github.com/twitter-archive/snowflake
https://github.com/twitter/twitter-server
http://www.ecma-international.org/ecma-262/6.0/index.html#sec-ecmascript-language-types-number-type
https://blog.csdn.net/u012488504/article/details/82194495
https://www.cnblogs.com/Hollson/p/9116218.html
https://tech.meituan.com/2017/04/21/mt-leaf.html
https://blog.csdn.net/X5fnncxzq4/article/details/79549514
https://github.com/twitter-archive/snowflake/blob/snowflake-2010/src/main/scala/com/twitter/service/snowflake/IdWorker.scala

1 587 801 157
2 147 483 647
   35 796 685

365*999=364 635
      1 048 576
         65 536

31|
19=524288
12|
10=1024
2=4
