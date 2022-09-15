# Yet-Another-RPC-Framework
电子科技大学2021秋季学期《分布式系统》课程设计：Yet-Another RPC Framework。
## 目录说明
client_test，regusrtry_test，server_test这三个文件夹中的是测试程序，分别包含client、server、registry的源码。
YA-RPC文件夹中包含RPC框架源码。
## 特性
- 主要参考了Dubbo架构
- At-Least-Once语义
- 使用gob序列化库，支持所有数据类型的传输
- 对关键全局变量访问控制，保证多线程安全
- 带有心跳机制和负载均衡功能的注册中心
