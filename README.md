# ReMindful - 智能间隔重复学习系统

ReMindful是一个基于间隔重复算法（Spaced Repetition）的智能学习卡片系统，帮助用户高效记忆和复习知识点。

## 功能特性

### 🎯 核心功能
- **智能复习算法**: 基于SuperMemo2算法的间隔重复系统
- **学习卡片管理**: 支持多种卡片类型（基础、填空、问答）
- **标签系统**: 灵活的标签分类和管理
- **复习统计**: 详细的学习进度和复习数据分析
- **用户系统**: 完整的用户注册、登录和个人信息管理

### 📊 数据分析
- **学习进度跟踪**: 实时监控学习进度和掌握情况
- **复习热力图**: 可视化复习活动分布
- **性能统计**: 复习质量和时间分析
- **连续学习**: 学习连续天数统计

### 🔧 技术特性
- **RESTful API**: 完整的REST API接口
- **JWT认证**: 安全的用户认证系统
- **Redis缓存**: 高性能数据缓存
- **MySQL数据库**: 可靠的数据存储
- **Swagger文档**: 完整的API文档

## 技术栈

- **后端框架**: Go + Gin
- **数据库**: MySQL + Redis
- **认证**: JWT
- **文档**: Swagger
- **算法**: SuperMemo2间隔重复算法

## 快速开始

### 环境要求
- Go 1.24+
- MySQL 8.0+
- Redis 6.0+

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd ReMindful
```

2. **安装依赖**
```bash
go mod download
```

3. **配置数据库**
```bash
# 创建MySQL数据库
mysql -u root -p
CREATE DATABASE remindful CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. **配置文件**
复制并修改配置文件：
```bash
cp config.yaml.example config.yaml
# 编辑config.yaml，配置数据库连接信息
```

5. **生成API文档**
```bash
chmod +x scripts/generate-docs.sh
./scripts/generate-docs.sh
```

6. **启动服务**
```bash
go run cmd/server/main.go
```

服务启动后访问：
- API服务: http://localhost:8080
- API文档: http://localhost:8080/swagger/index.html

## API接口

### 用户管理
- `POST /api/v1/send-code` - 发送验证码
- `POST /api/v1/register` - 用户注册
- `POST /api/v1/login` - 用户登录
- `GET /api/v1/user` - 获取用户信息
- `PUT /api/v1/user` - 更新用户信息

### 学习卡片
- `POST /api/v1/learning-cards` - 创建卡片
- `GET /api/v1/learning-cards` - 获取卡片列表
- `GET /api/v1/learning-cards/review` - 获取需要复习的卡片
- `GET /api/v1/learning-cards/:id` - 获取单个卡片
- `PUT /api/v1/learning-cards/:id` - 更新卡片
- `DELETE /api/v1/learning-cards/:id` - 删除卡片
- `POST /api/v1/learning-cards/:id/review` - 复习卡片

### 标签管理
- `POST /api/v1/tags` - 创建标签
- `GET /api/v1/tags` - 获取标签列表
- `GET /api/v1/tags/:id` - 获取单个标签
- `PUT /api/v1/tags/:id` - 更新标签
- `DELETE /api/v1/tags/:id` - 删除标签

### 复习日志
- `GET /api/v1/review-logs` - 获取复习日志
- `GET /api/v1/review-logs/stats` - 获取复习统计
- `GET /api/v1/review-logs/progress` - 获取学习进度
- `GET /api/v1/review-logs/heatmap` - 获取复习热力图

## 配置说明

### config.yaml
```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: remindful
  charset: utf8mb4

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your_jwt_secret
  expiration: 24h

email:
  host: smtp.163.com
  port: 465
  username: your_email@163.com
  password: your_email_password
  from: your_email@163.com
  TLS: false
  SSL: true
```

## 间隔重复算法

本系统采用SuperMemo2算法，根据用户的复习表现动态调整复习间隔：

- **质量评分**: 0-5分，表示复习质量
- **间隔计算**: 基于难度系数和复习次数
- **自适应调整**: 根据用户表现调整难度系数

## 项目结构

```
ReMindful/
├── cmd/server/          # 应用入口
├── internal/            # 内部模块
│   ├── config/         # 配置管理
│   ├── handler/        # HTTP处理器
│   ├── middleware/     # 中间件
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   ├── router/         # 路由配置
│   └── service/        # 业务逻辑层
├── pkg/                # 公共包
│   ├── algorithm/      # 算法实现
│   ├── database/       # 数据库工具
│   ├── jwt/           # JWT工具
│   └── utils/         # 工具函数
├── docs/              # API文档
├── scripts/           # 脚本文件
└── config.yaml        # 配置文件
```

## 开发指南

### 添加新功能
1. 在`internal/model/`中定义数据模型
2. 在`internal/repository/`中实现数据访问
3. 在`internal/service/`中实现业务逻辑
4. 在`internal/handler/`中实现HTTP处理
5. 在`internal/router/`中添加路由

### 代码规范
- 使用Go标准代码格式
- 添加适当的注释和文档
- 遵循RESTful API设计原则
- 使用Swagger注解生成API文档

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交Issue或联系开发团队。 