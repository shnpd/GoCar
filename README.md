# 邻车出行
邻车出行租车小程序，前端采用wxml+typescript实现，后端基于go的微服务架构实现，包括登录服务、汽车服务、租赁服务，对象存储服务。
技术栈：gRPC、ProtoBuf、MongoDB、对象存储、JWT、RabbitMQ、Redis、Docker、k8s

## 系统架构
![](https://cdn.jsdelivr.net/gh/shnpd/blog-pic@main/javastudy/%E9%82%BB%E8%BD%A6%E5%87%BA%E8%A1%8C%E6%9E%B6%E6%9E%84%E5%9B%BE.png)
认证服务（auth）：
- 接受小程序登录请求，验证用户身份，生成JWT令牌。

对象存储服务（blob）：
- 连接云对象存储，封装对象存储相关操作

汽车服务（car）：
- 提供创建汽车、查询汽车、汽车开锁、汽车关锁、汽车状态更新等功能。
- 汽车状态更新时会将汽车状态发送到消息队列
- 通过多线程模拟汽车，通过controller管理汽车线程，从消息队列中获取汽车状态并通过go channel转发给对应的汽车线程
- 通过websocket将消息队列中的汽车位置更新实时推送到小程序。

行程服务（rental/trip）：
- 提供创建行程、查询行程、更新行程等功能。
- 接受消息队列中的汽车状态更新，更新行程状态。
- 使用redis分布式锁保证行程更新的安全性

身份服务（rental/profile）：
- 提供身份查询、提交认证、身份更新等功能。


## 项目结构
微服务内层
目录结构基本一致，只以认证服务、汽车服务为例
- server：服务端
  - auth：认证服务
    - api：proto、gateway接口定义
    - auth：服务核心实现
    - dao：数据访问层
    - token：jwt相关操作
    - wechat：微信api相关操作
  - blob：对象存储服务
    - api：proto、gateway接口定义
    - blob：服务核心实现
    - cos：云服务访问
    - dao：数据访问层
  - car：汽车服务
    - api：proto、gateway接口定义
    - car：服务核心实现
    - dao：数据访问层
    - mq：消息队列相关操作
    - sim：汽车模拟器
    - trip：行程更新的相关操作
    - ws：websocket相关操作
  - cmd：测试目录
  - gateway：网关服务
  - rental：租赁服务
    - api：proto、gateway接口定义
    - profile：身份服务
    - trip：行程服务
  - shared：共享模块、
    - auth：用户token认证
    - coolenv：三方工具
    - id：Identifier Type
    - mongo：MongoDB相关操作
    - server：grpc server相关操作
- wx/miniprogram：小程序前端


## 运行小程序
1. `cd wx/miniprogram`
2. `npm install`
3. 打开小程序开发工具，点击工具->构建npm
4. 点击编译
