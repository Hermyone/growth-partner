***

## 身份声明与严苛要求

我是本项目的高级全栈架构师，负责「成长伙伴」教育系统的整体技术架构设计（Golang + Gin + PostgreSQL + Redis + WebSocket，Clean Architecture 三层分层）。

现在我需要你以**同等级别的高级后端工程师**身份，对本项目已开发完成的全部接口进行**系统性、生产级、无死角的验证体系设计**。

**对你的严苛要求（不可妥协）：**

- 输出必须结构化，每个模块独立成章，使用 Markdown + 表格 + 代码块
- 验证逻辑必须覆盖正常流、边界流、异常流、安全流全部场景
- 所有测试代码/脚本必须可直接复制执行，不得有占位符
- 安全验证必须涵盖 OWASP Top 10 相关场景
- WebSocket 验证必须覆盖并发、心跳、断线重连完整生命周期
- 最终必须输出一份可执行的「接口验证手册」，含优先级与自动化建议
- 不允许输出任何废话、重复内容、无法执行的伪代码

***

## 项目背景

**项目名称**：成长伙伴（Growth Partner）\
**技术栈**：Golang 1.23 · Gin · GORM · PostgreSQL 16 · Redis 7 · Gorilla WebSocket\
**认证方案**：JWT Bearer Token（Access Token 24h + Refresh Token 7d）\
**统一响应格式**：`{ "code": 0, "message": "success", "data": {}, "error": null, "timestamp": 1234567890 }`\
**API 基础路径**：`http://localhost:8080/api/v1`\
**角色体系**：`admin` / `teacher` / `parent` / `student`\
**权限核心**：老师必须通过 `teacher_class_assignments` 表校验对班级的操作权限

***

## 全模块验证任务

请对以下 **9 大模块、共约 72 条路由**逐一设计验证逻辑，每个模块必须严格按照【七大验证维度】输出。

***

### 【七大验证维度】（每个模块必须全部覆盖）

| 维度            | 要求                                                     |
| :------------ | :----------------------------------------------------- |
| ① 功能正确性       | 正常流程、各参数边界值、枚举值合法性、返回字段完整性                             |
| ② 安全性         | JWT鉴权、角色越权、水平越权（访问他人数据）、SQL注入、速率限制、Token黑名单            |
| ③ 性能与并发       | Redis缓存命中率、数据库慢查询、WebSocket并发连接数                       |
| ④ 数据一致性       | 事务回滚验证、Redis与PostgreSQL数据同步、幂等性                        |
| ⑤ 异常与容错       | 全错误码覆盖、网络超时、数据库宕机降级、日志记录完整性                            |
| ⑥ WebSocket专属 | 连接建立握手、Ping/Pong心跳、断线重连、消息顺序、并发广播                      |
| ⑦ 工具与脚本       | Postman Collection JSON / Newman CLI命令 / Go测试代码 / 压测脚本 |

***

### 一、认证模块（Auth）— 5 条路由

**路由清单：**

```
POST   /api/v1/auth/login          # 统一登录
POST   /api/v1/auth/refresh        # Token刷新
POST   /api/v1/auth/logout         # 注销登录
GET    /api/v1/auth/me             # 获取当前用户信息
PATCH  /api/v1/auth/password       # 修改密码
```

**验证重点指令（必须全部覆盖）：**

① **功能正确性**：

- 四种角色（admin/teacher/parent/student）分别登录，验证返回的 JWT Claims 中 role/class\_id/child\_id 字段是否正确
- 密码错误时返回 401，连续错误 5 次后账号是否被锁定（或返回 429）
- refresh\_token 过期后刷新是否返回正确错误码
- 修改密码后旧 access\_token 是否立即失效（Redis 黑名单验证）

② **安全性**：

- 暴力破解防护：10秒内超过5次登录失败，IP或账号维度的限流
- Token 伪造：修改 JWT payload 后重新 base64 编码，验证服务器是否拒绝
- 并发登录：同一账号同时发起 100 个登录请求，验证 Token 一致性
- logout 接口：同一 token 注销两次，第二次应返回 401

③ **数据一致性**：

- 修改密码后，Redis 中 `password_version` key 是否正确更新
- refresh\_token 使用一次后是否失效（防重放）

**输出格式要求**：

- 完整的 Postman Collection JSON（含环境变量设置）
- Newman 一键执行命令
- Go 测试代码（`auth_test.go`，使用 `httptest` 包）
- 压测脚本（使用 `hey` 或 `wrk`）

***

### 二、管理员模块（Admin）— 24 条路由

**路由清单：**

```
GET    /api/v1/admin/schools
POST   /api/v1/admin/schools
PUT    /api/v1/admin/schools/:id
PATCH  /api/v1/admin/schools/:id/status
GET    /api/v1/admin/classes
POST   /api/v1/admin/classes
PUT    /api/v1/admin/classes/:id
POST   /api/v1/admin/classes/:id/promote       # 升年级（高风险）
PATCH  /api/v1/admin/classes/:id/status
GET    /api/v1/admin/users
POST   /api/v1/admin/users
PUT    /api/v1/admin/users/:id
PATCH  /api/v1/admin/users/:id/status
PATCH  /api/v1/admin/users/:id/reset-pwd
POST   /api/v1/admin/students/batch-import
GET    /api/v1/admin/students
POST   /api/v1/admin/students
PUT    /api/v1/admin/students/:id
GET    /api/v1/admin/assignments
POST   /api/v1/admin/assignments
DELETE /api/v1/admin/assignments/:id
POST   /api/v1/admin/assignments/batch
GET    /api/v1/admin/parent-bindings
POST   /api/v1/admin/parent-bindings
DELETE /api/v1/admin/parent-bindings/:id
GET    /api/v1/admin/dashboard
GET    /api/v1/admin/audit-logs
```

**验证重点指令：**

① **升年级接口专项验证（`/classes/:id/promote`）**——这是最高风险接口：

- 对同一班级重复调用两次，第二次必须返回错误（幂等保护，通过 school\_year 唯一索引）
- 事务验证：模拟数据库在 `class_enrollments` 批量插入到一半时失败，验证 classes 表的 grade 字段是否也回滚
- 验证影响的学生数量是否正确返回
- 使用 teacher/student 角色调用此接口，必须返回 403

② **批量导入学生（`/students/batch-import`）**：

- 上传格式错误的 CSV（编码不对、字段缺失、超大文件），验证各种错误响应
- 部分数据合法、部分非法时，是整体回滚还是跳过非法行（验证设计一致性）
- 同一学号重复导入，验证唯一索引冲突处理

③ **权限越权验证**：

- 用 teacher/parent/student 的 Token 调用任意 admin 接口，必须返回 403
- 不带 Token 调用，必须返回 401

④ **数据一致性**：

- 创建老师账号 + 分配班级（assignments），再撤销分配，验证 teacher\_class\_assignments.is\_active 变化
- 升年级后验证：classes.grade 变化、class\_enrollments 新增记录、旧记录保留

**输出格式**：按路由分组的 Postman Collection + 关键接口的 Go 集成测试代码

***

### 三、老师端模块（Teacher）— 22 条路由

**路由清单：**

```
GET    /api/v1/teacher/my-classes
GET    /api/v1/teacher/classes/:classId/overview
GET    /api/v1/teacher/classes/:classId/students
POST   /api/v1/teacher/behaviors                    # 【核心】打分
GET    /api/v1/teacher/behaviors
GET    /api/v1/teacher/behaviors/:id
DELETE /api/v1/teacher/behaviors/:id                # 24h内撤销
POST   /api/v1/teacher/behaviors/batch              # 批量打分
GET    /api/v1/teacher/broadcasts
POST   /api/v1/teacher/broadcasts
DELETE /api/v1/teacher/broadcasts/:id
GET    /api/v1/teacher/challenges
POST   /api/v1/teacher/challenges
PATCH  /api/v1/teacher/challenges/:id/complete
GET    /api/v1/teacher/questions
POST   /api/v1/teacher/questions
PUT    /api/v1/teacher/questions/:id
DELETE /api/v1/teacher/questions/:id
POST   /api/v1/teacher/questions/batch-import
GET    /api/v1/teacher/blindbox/pool
POST   /api/v1/teacher/blindbox/pool
PUT    /api/v1/teacher/blindbox/pool/:id
DELETE /api/v1/teacher/blindbox/pool/:id
POST   /api/v1/teacher/blindbox/draw/:childId
PATCH  /api/v1/teacher/blindbox/draws/:drawId/redeem
POST   /api/v1/teacher/reports/weekly
GET    /api/v1/teacher/reports/weekly
GET    /api/v1/teacher/reports/weekly/:id/download
```

**验证重点指令：**

① **打分接口事务完整性验证（最核心）**：

```
POST /api/v1/teacher/behaviors
{
  "child_id": 1,
  "class_id": 1,
  "dimension": "virtue",
  "description": "主动帮助同学",
  "growth_value": 3
}
```

验证步骤：

- 调用成功后，查询 `behavior_records` 新增 1 条
- 查询 `growth_records` 新增 1 条，`delta=3`，`after_points` 正确
- 查询 `partners` 表 `growth_points` 增加了 3
- 若触发进化，`partners.current_stage` 正确变化，`growth_records.is_evolution_trigger=true`
- 模拟数据库中间失败（使用 mock），验证三表全部回滚

② **班级权限越权验证**：

- 老师 A（仅授权班级1）调用班级2的打分接口 → 必须返回 403
- `can_score=false` 的老师尝试打分 → 必须返回 403
- 老师调用其他老师班级的广播接口（`can_broadcast=false`）→ 403

③ **24小时撤销规则**：

- 对 25 小时前创建的行为记录调用 DELETE → 返回业务错误（非404）
- 撤销成功后验证成长值是否正确回滚（delta 写负值到 growth\_records）

④ **批量打分并发验证**：

- 同时对同一学生发起 50 个批量打分请求，验证最终 growth\_points 是否等于预期累加值（无并发丢失）

⑤ **盲盒超卖防护**：

- 奖励池 stock=1，并发发起 10 个 draw 请求，最终只有 1 个成功，验证 `SELECT FOR UPDATE` 行锁有效性

***

### 四、学生端模块（Student）— 12 条路由

**路由清单：**

```
GET    /api/v1/student/partner
GET    /api/v1/student/partners
POST   /api/v1/student/partner                    # 选择伙伴
PATCH  /api/v1/student/partner/nickname
GET    /api/v1/student/partner/growth-history
GET    /api/v1/student/partner/templates
GET    /api/v1/student/behaviors
GET    /api/v1/student/behaviors/stats
GET    /api/v1/student/broadcasts
PATCH  /api/v1/student/broadcasts/:id/read
POST   /api/v1/student/broadcasts/read-all
GET    /api/v1/student/growth-calendar/months
GET    /api/v1/student/growth-calendar/months/:month
GET    /api/v1/student/growth-calendar/annual/:year
GET    /api/v1/student/milestones
GET    /api/v1/student/blindbox/my-draws
```

**验证重点指令：**

① **选择伙伴接口约束验证**：

- 首次选择：正常流程，验证 `partners` 表新增记录，`status=active`，`sequence_no=1`
- 无解锁权限时选择第二只：必须返回业务错误（查 `partner_unlock_logs`）
- 已有 `status=active` 伙伴时再次调用：返回业务错误
- 数据库层局部唯一索引验证：直接 SQL 尝试插入两条 `status=active` 的记录，验证数据库拒绝

② **水平越权验证**：

- 学生 A 尝试查询学生 B 的 `partner`/`behaviors`/`broadcasts` → 必须返回 403 或 404
- 学生尝试调用 teacher/admin 接口 → 必须返回 403

③ **数据隔离验证**：

- 确认所有 student 接口的响应中，绝对不包含 `real_name_enc`、`phone_enc`、`student_no_enc` 字段
- 确认不暴露其他学生的任何信息

***

### 五、家长端模块（Parent）— 9 条路由

**路由清单：**

```
GET    /api/v1/parent/children
GET    /api/v1/parent/children/:childId/partner
GET    /api/v1/parent/children/:childId/partners
GET    /api/v1/parent/children/:childId/behaviors
GET    /api/v1/parent/children/:childId/broadcasts
GET    /api/v1/parent/children/:childId/milestones
GET    /api/v1/parent/children/:childId/monthly-card
GET    /api/v1/parent/children/:childId/annual-report
GET    /api/v1/parent/children/:childId/battles
```

**验证重点指令：**

① **核心越权测试**：

- 家长 A（绑定 child\_id=1）尝试访问 `children/2/partner` → 必须返回 403
- 未绑定任何孩子的家长账号调用所有接口 → 返回空数据或 404，不报 500

② **对战记录验证**：

- 确认 `/battles` 接口返回数据中**不包含** `score`、`is_winner`、`player_b_score` 等胜负字段
- 仅展示参与次数和时间，保护正向理念

***

### 六、知识对战模块（Battle）— 6 条 HTTP + 1 条 WebSocket

**HTTP 路由清单：**

```
GET    /api/v1/battle/subjects
POST   /api/v1/battle/rooms
POST   /api/v1/battle/rooms/:roomCode/join
GET    /api/v1/battle/rooms/:roomCode
GET    /api/v1/battle/history
GET    /api/v1/battle/history/:roomId/review
```

**WebSocket：**

```
WS     /api/v1/battle/ws?room_code=XXX
```

**验证重点指令：**

① **完整对战流程端到端验证**：

```
Step 1: 学生A POST /battle/rooms → 获取 room_code
Step 2: 学生B POST /battle/rooms/:roomCode/join → 双方收到 battle:ready
Step 3: 服务端推送 battle:question（10道题）
Step 4: 双方 WS 发送 battle:answer，验证每题收到 battle:result
Step 5: 所有题完成后，验证 battle:finish 推送（含 growth_gained，不含对方分数）
Step 6: 查询 battle_participants 验证双方均有 growth_gained=3（或配置值）
Step 7: 查询 growth_records 验证两条成长值流水记录
```

② **WebSocket 专属验证**：

- 心跳测试：连接后 30 秒不发送任何消息，验证服务端是否在 40 秒内断开连接
- 断线重连：学生 A 在对战中途断开 WebSocket，验证：房间状态变化 + 学生 B 收到 `battle:error` 消息
- 并发连接：100 个 WebSocket 同时连接（使用 `gorilla/websocket` 测试客户端），验证服务端内存和 Hub 管理正常
- 消息幂等：对战结束后重复触发成长值发放，验证 `growth_gained` 只写入一次（幂等保护）

③ **题目安全验证**：

- 确认 `battle:question` 消息中**不包含** `answer` 字段
- 客户端直接发送含 answer 的 WS 消息尝试作弊，服务端必须忽略并以自身校验为准

**验证工具**：

```bash
# WebSocket 测试（wscat）
npm install -g wscat
wscat -c "ws://localhost:8080/api/v1/battle/ws?room_code=ABC123" \
  -H "Authorization: Bearer <token>"

# 并发 WebSocket 压测（Go 代码）
# 请输出完整的 Go 并发测试代码，同时建立 100 条 WS 连接
```

***

### 七、广播与实时推送模块（Broadcast/WebSocket）— 1 条 WebSocket

**WebSocket：**

```
WS     /api/v1/ws    # 主广播通道
```

**验证重点指令：**

① **完整推送链路验证**：

```
Step 1: 学生连接 /api/v1/ws
Step 2: 老师调用 POST /teacher/behaviors 为该学生打分
Step 3: 验证学生 WS 收到 notify:behavior 消息（含 partner_message）
Step 4: 若触发进化，验证紧接着收到 notify:evolution 消息
Step 5: 验证消息延迟 < 500ms（实时性要求）
```

② **Redis Pub/Sub 验证**：

- 直接用 `redis-cli PUBLISH broadcast:class:1 '{"type":"test"}'`，验证该班级所有在线学生收到消息
- Redis 重启后，服务端是否自动重新订阅（断线重连机制）

③ **多端登录处理**：

- 同一学生账号同时建立两条 WS 连接，第二次连接建立时第一条必须被踢下线

④ **离线消息补推验证**：

- 学生离线时老师打分，学生上线后验证是否收到补推的 `notify:behavior` 消息
- 超过 5 分钟未推送的消息，验证是否放弃推送并记录日志

***

### 八、公开接口 & 伙伴模板模块 — 7 条路由

```
GET    /api/v1/partner-templates
GET    /api/v1/partner-templates/:id
GET    /health
GET    /api/v1/config/client
POST   /api/v1/admin/partner-templates
PUT    /api/v1/admin/partner-templates/:id
POST   /api/v1/admin/partner-templates/seed
```

**验证重点指令：**

① **种子数据接口幂等验证**：

- `POST /admin/partner-templates/seed` 调用两次，第二次必须返回"已初始化"提示，不重复插入
- 验证 30 个模板全部存在，`code` 字段唯一

② **缓存验证**：

- 第一次请求 `/partner-templates` 后，Redis 是否缓存了结果
- 第二次请求是否走缓存（响应时间对比，且 DB 查询日志无新增）
- 管理员更新模板后，缓存是否被正确清除（Cache-Aside 模式验证）

***

### 九、阳光章系统（二期预留接口）— 9 条路由

```
GET    /api/v1/admin/sunshine/colors
POST   /api/v1/admin/sunshine/colors
PUT    /api/v1/admin/sunshine/colors/:id
POST   /api/v1/teacher/sunshine/stamp
GET    /api/v1/teacher/sunshine/stamps
GET    /api/v1/student/sunshine/my-stamps
POST   /api/v1/teacher/sunshine/awards/evaluate
GET    /api/v1/teacher/sunshine/awards
GET    /api/v1/student/sunshine/my-awards
GET    /api/v1/parent/children/:childId/sunshine
```

**验证重点指令：**

① **状态验证（二期接口需返回固定占位响应）**：

- 所有二期接口当前应返回 `{"code": 0, "message": "功能即将上线", "data": null}`
- 验证接口已存在（不返回 404），且需要正确的权限（不是公开接口）

② **表结构预验证**：

- 验证数据库中以下表已存在且结构正确：`sunshine_colors`、`sunshine_stamps`、`sunshine_star_awards`
- 提供 SQL 验证脚本

***

## 最终输出要求：「接口验证执行手册」

你必须在全部模块验证逻辑输出后，生成一份完整的\*\*「接口验证执行手册」\*\*，包含：

### 验证优先级矩阵

| 优先级 | 模块                  | 理由          |
| :-- | :------------------ | :---------- |
| P0  | 认证模块 + 打分接口 + WS主连接 | 核心链路，其他功能依赖 |
| P1  | 管理员模块（升年级/权限分配）     | 高风险写操作      |
| P2  | 老师端（广播/盲盒/题库）       | 主要功能        |
| P3  | 学生/家长端              | 只读为主        |
| P4  | 阳光章（二期）             | 占位验证        |

### 测试执行顺序

```bash
# 1. 环境准备
docker compose up -d
sleep 10  # 等待服务就绪

# 2. 健康检查
curl http://localhost:8080/health

# 3. 初始化种子数据
curl -X POST http://localhost:8080/api/v1/admin/partner-templates/seed \
  -H "Authorization: Bearer <admin_token>"

# 4. 按优先级执行 Newman 测试集合
newman run auth_collection.json -e env.json
newman run admin_collection.json -e env.json
newman run teacher_collection.json -e env.json
newman run student_collection.json -e env.json
newman run parent_collection.json -e env.json
newman run battle_collection.json -e env.json

# 5. WebSocket 专项测试
go test ./test/websocket/... -v -timeout 120s

# 6. 压力测试
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/partner-templates
```

### 自动化建议

- CI/CD 集成：将 Newman 测试集成到 GitHub Actions，每次 PR 自动运行
- 覆盖率目标：功能测试 100% 路由覆盖，安全测试覆盖所有角色越权场景
- 监控告警：集成 Prometheus + Grafana，对打分接口的 P99 延迟设置告警阈值 < 200ms
- 定期回归：WebSocket 并发测试每周定时运行，验证连接数上限未退化

***

