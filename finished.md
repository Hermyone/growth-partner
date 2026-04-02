# 接口实现完成情况

## 模块完成状态

- [x] 一、认证模块（Auth）- 5条路由
- [x] 二、管理员模块（Admin）- 24条路由
- [x] 三、老师端模块（Teacher）- 22条路由
- [x] 四、学生端模块（Student）- 12条路由
- [x] 五、家长端模块（Parent）- 9条路由
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
- [x] GET /api/v1/teacher/my-classes - 获取当前老师被授权的所有班级（含权限类型）
- [x] GET /api/v1/teacher/classes/:classId/overview - 获取班级概览（学生数/行为数/成长值等）
- [x] GET /api/v1/teacher/classes/:classId/students - 获取班级学生列表（含伙伴/成长值信息）
- [x] POST /api/v1/teacher/behaviors - 为学生添加正向行为记录，触发成长值/伙伴进化
- [x] GET /api/v1/teacher/behaviors - 查看班级行为记录列表（多条件筛选，分页）
- [x] GET /api/v1/teacher/behaviors/:id - 获取单条行为记录详情
- [x] DELETE /api/v1/teacher/behaviors/:id - 撤销行为记录（24小时内，扣减成长值）
- [x] POST /api/v1/teacher/behaviors/batch - 批量为多个学生打分
- [x] GET /api/v1/teacher/broadcasts - 查看自己发送的广播列表（已发/定时待发）
- [x] POST /api/v1/teacher/broadcasts - 发送广播（立即/定时）
- [x] DELETE /api/v1/teacher/broadcasts/:id - 取消定时广播（仅未发送）
- [x] GET /api/v1/teacher/challenges - 查看班级当前进行中的集体挑战
- [x] POST /api/v1/teacher/challenges - 创建集体挑战，配置条件+奖励
- [x] PATCH /api/v1/teacher/challenges/:id/complete - 手动标记挑战完成，批量发放成长值
- [x] GET /api/v1/teacher/questions - 查看班级题库（公共+专属）
- [x] POST /api/v1/teacher/questions - 添加班级专属题目
- [x] PUT /api/v1/teacher/questions/:id - 编辑题目
- [x] DELETE /api/v1/teacher/questions/:id - 删除题目（软删除）
- [x] POST /api/v1/teacher/questions/batch-import - 批量导入题目（CSV）
- [x] GET /api/v1/teacher/blindbox/pool - 查看本班盲盒奖励池
- [x] POST /api/v1/teacher/blindbox/pool - 向奖励池添加奖励
- [x] PUT /api/v1/teacher/blindbox/pool/:id - 编辑奖励配置
- [x] DELETE /api/v1/teacher/blindbox/pool/:id - 下架奖励（软删除）
- [x] POST /api/v1/teacher/blindbox/draw/:childId - 为学生触发抽盲盒
- [x] PATCH /api/v1/teacher/blindbox/draws/:drawId/redeem - 确认兑换学生盲盒奖励
- [x] POST /api/v1/teacher/reports/weekly - 触发生成本班本周正能量周报PDF（异步）
- [x] GET /api/v1/teacher/reports/weekly - 查看历史周报列表（含下载链接）
- [x] GET /api/v1/teacher/reports/weekly/:id/download - 下载指定周报PDF

### 学生端模块
- [x] GET /api/v1/student/partner - 获取当前活跃伙伴详情
- [x] GET /api/v1/student/partners - 获取所有历史伙伴列表
- [x] POST /api/v1/student/partner - 选择新伙伴（首次/满级后），校验解锁权限
- [x] PATCH /api/v1/student/partner/nickname - 修改当前伙伴昵称
- [x] GET /api/v1/student/partner/growth-history - 获取伙伴成长值流水（分页）
- [x] GET /api/v1/student/partner/templates - 获取可供选择的伙伴模板列表
- [x] GET /api/v1/student/behaviors - 查看自己的行为记录（多条件筛选，分页）
- [x] GET /api/v1/student/behaviors/stats - 行为统计，用于前端雷达图展示
- [x] GET /api/v1/student/broadcasts - 获取收到的广播消息列表（伙伴+园长）
- [x] PATCH /api/v1/student/broadcasts/:id/read - 标记广播为已读
- [x] POST /api/v1/student/broadcasts/read-all - 一键标记所有广播为已读
- [x] GET /api/v1/student/growth-calendar/months - 获取全部月度成长卡列表（按学年分组）
- [x] GET /api/v1/student/growth-calendar/months/:month - 获取指定月份成长卡详情
- [x] GET /api/v1/student/growth-calendar/annual/:year - 获取年度成长画卷数据
- [x] GET /api/v1/student/milestones - 获取里程碑列表（勋章墙）
- [x] GET /api/v1/student/blindbox/my-draws - 查看自己已抽到的盲盒奖励（含状态）

### 家长端模块
- [x] GET /api/v1/parent/children - 获取自己绑定的孩子列表
- [x] GET /api/v1/parent/children/:childId/partner - 查看孩子当前伙伴状态
- [x] GET /api/v1/parent/children/:childId/partners - 查看孩子的历史伙伴列表
- [x] GET /api/v1/parent/children/:childId/behaviors - 查看孩子的正向行为记录
- [x] GET /api/v1/parent/children/:childId/broadcasts - 查看孩子收到的伙伴鼓励广播
- [x] GET /api/v1/parent/children/:childId/milestones - 查看孩子的里程碑贴纸
- [x] GET /api/v1/parent/children/:childId/monthly-card - 查看孩子本月/历史月度成长卡
- [x] GET /api/v1/parent/children/:childId/annual-report - 查看孩子年度成长画卷（含PDF下载）
- [x] GET /api/v1/parent/children/:childId/battles - 查看孩子的对战参与记录

### 知识对战模块

### 广播与实时推送模块

### 公开接口 & 伙伴模板模块

### 阳光章系统
