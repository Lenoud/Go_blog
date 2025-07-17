# Blog System

一个基于Go语言开发的博客系统

## 项目结构

```
.
├── config/             # 配置文件
├── controllers/        # 控制器
├── middleware/         # 中间件
├── models/            # 数据模型
├── routes/            # 路由
├── services/          # 业务逻辑
├── storage/           # 数据存储
│   ├── mysql/        # MySQL相关
│   ├── mongodb/      # MongoDB相关
│   └── redis/        # Redis相关
└── utils/            # 工具函数
```

## 技术栈

- Go 1.24
- Gin Web框架
- MySQL 8
- MongoDB 7
- Redis 集群
- JWT认证

## 功能特性

- 用户认证（登录/注册）
- 文章管理（CRUD）
- 分类管理
- 缓存支持 