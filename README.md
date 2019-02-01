# Kratos

Kratos是[bilibili](https://www.bilibili.com)开源的一套Go微服务框架，基于“大仓(monorepo)”理念，包含大量微服务相关框架及工具。  
如：discovery(服务注册发现)、blademaster(HTTP框架)、warden(gRPC封装)、log、breaker、dapper(trace)、cache&db sdk、kratos(代码生成等工具)等等。  

我们致力于提供完整的微服务研发体验，大仓整合相关框架及工具后，微服务治理相关部分可对整体业务开发周期无感，从而更加聚焦于业务交付。对每位开发者而言，整套Kratos框架也是不错的学习仓库，可以了解和参考到[bilibili](https://www.bilibili.com)在微服务方面的技术积累和经验。

# TODOs

- [ ] log&log-agent @围城
- [ ] config @志辉
- [ ] bm @佳辉
- [ ] warden @龙虾
- [ ] naming discovery @堂辉
- [ ] cache&database @小旭
- [ ] kratos tool @普余

# issues

***抽空细分后加到issue里***

1. 需要考虑配置中心开源方式：类discovery单独 或 集成在大仓库
2. log-agent和dapper需要完整的解决方案，包含ES集群、dapperUI
3. databus&canal相关解决方案如何开源
4. proto文件相关生成工具正和到kratos工具内
