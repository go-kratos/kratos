### business/warden/balancer/wrr

##### Version 1.2.1
1. 删除了netflix ribbon的权重算法，改成了平方根算法

##### Version 1.2.0
1. 实现了动态计算的调度轮询算法（使用了服务端的成功率数据，替换基于本地计算的成功率数据）

##### Version 1.1.0
1. 实现了动态计算的调度轮询算法

##### Version 1.0.0

1. 实现了带权重可以识别Color的轮询算法
