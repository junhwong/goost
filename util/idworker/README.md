# idworker

idworker用于生成分布式环境下全局唯一的长整形ID。设计目标：全局唯一、时间序、扩展良好(针对持久化存储，如：分库分表)。

参考：

- https://developer.twitter.com/en/docs/basics/twitter-ids.html
- https://github.com/twitter-archive/snowflake
- https://gitee.com/mayanjun/idworker
- https://segmentfault.com/a/1190000011282426?utm_source=tag-newest
- https://www.petalstopicots.com/grandma-jennies-snowflake-patterns-two-versions/
