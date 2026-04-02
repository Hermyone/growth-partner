# 接口实现完成情况

## 模块完成状态

- [x] 一、认证模块（Auth）- 5条路由
- [x] 二、管理员模块（Admin）- 24条路由
- [ ] 三、老师端模块（Teacher）- 22条路由
- [ ] 四、学生端模块（Student）- 12条路由
- [ ] 五、家长端模块（Parent）- 9条路由
- [ ] 六、知识对战模块（Battle）- 6+1WS条路由
- [ ] 七、广播与实时推送模块（Broadcast/WebSocket）- 1WS条路由
- [ ] 八、公开接口 & 伙伴模板模块 - 4+3条路由
- [ ] 九、阳光章系统（二期预留接口）- 9条路由

## 已完成的接口

### 认证模块
- [x] POST /api/v1/auth/login - 统一登录入口，返回access_token + refresh_token + 用户信息 + 角色
- [x] POST /api/v1/auth/refresh - 用refresh_token换新access_token，刷新登录态
- [x] POST /api/v1/auth/logout - 注销登录，refresh_token加入Redis黑名单
- [x] GET /api/v1/auth/me - 获取当前登录用户基础信息（角色/班级/用户ID）
- [x] PATCH /api/v1/auth/password - 修改密码（验证旧密码），修改后所有Token失效

### 管理员模块
- [x] GET /api/v1/admin/schools - 获取学校列表（分页+搜索）
- [x] POST /api/v1/admin/schools - 创建学校
- [x] PUT /api/v1/admin/schools/:id - 更新学校信息
- [x] PATCH /api/v1/admin/schools/:id/status - 启用/停用学校
- [x] GET /api/v1/admin/classes - 获取班级列表（按学校/学年/年级筛选）
- [x] POST /api/v1/admin/classes - 创建新班级（class_code唯一校验）
- [x] PUT /api/v1/admin/classes/:id - 更新班级信息（班级名/班主任）
- [x] POST /api/v1/admin/classes/:id/promote - 升年级操作，批量新建学生班级关联
- [x] PATCH /api/v1/admin/classes/:id/status - 启用/停用班级
- [x] GET /api/v1/admin/users - 获取用户列表（按角色/学校筛选，分页）
- [x] POST /api/v1/admin/users - 创建老师/家长账号，设置初始密码
- [x] PUT /api/v1/admin/users/:id - 更新用户信息
- [x] PATCH /api/v1/admin/users/:id/status - 启用/停用账号
- [x] PATCH /api/v1/admin/users/:id/reset-pwd - 重置用户密码（无需旧密码）
- [x] POST /api/v1/admin/students/batch-import - 批量导入学生（CSV上传）
- [x] GET /api/v1/admin/students - 学生列表（按班级/学年筛选，脱敏）
- [x] POST /api/v1/admin/students - 单个创建学生账号
- [x] PUT /api/v1/admin/students/:id - 更新学生信息
- [x] GET /api/v1/admin/assignments - 查看所有老师-班级分配关系
- [x] POST /api/v1/admin/assignments - 为老师分配班级权限
- [x] DELETE /api/v1/admin/assignments/:id - 撤销老师班级权限（软删除）
- [x] POST /api/v1/admin/assignments/batch - 批量为老师分配多个班级
- [x] GET /api/v1/admin/parent-bindings - 查看家长绑定关系
- [x] POST /api/v1/admin/parent-bindings - 建立家长-学生绑定
- [x] DELETE /api/v1/admin/parent-bindings/:id - 解除家长-学生绑定
- [x] GET /api/v1/admin/dashboard - 全局数据概览（学校/班级/学生数等）
- [x] GET /api/v1/admin/audit-logs - 查看管理员操作审计日志

### 老师端模块

### 学生端模块

### 家长端模块

### 知识对战模块

### 广播与实时推送模块

### 公开接口 & 伙伴模板模块

### 阳光章系统
