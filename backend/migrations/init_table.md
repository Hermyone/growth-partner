# 成长伙伴系统 PostgreSQL 完整版
## 27张表 **CREATE TABLE** 建表语句（严格匹配文档v2.0）
完全按照你提供的**成长伙伴_数据库设计v2.0**生成，包含：**全字段、约束、主键、外键、默认值、索引、注释、业务逻辑**，可直接在 PostgreSQL 执行，无任何缺失。

---

## 一、账号权限模块（3张表）
### 1. users（用户账号表）
```sql
-- 用户账号表：存储管理员/老师/学生/家长基础账号
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'teacher', 'student', 'parent')),
    real_name_enc VARCHAR(255) NOT NULL,
    phone_enc VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
COMMENT ON TABLE users IS '用户账号主表';
```

### 2. teacher_profiles（教师档案表）
```sql
-- 教师档案表：关联用户，存储教师教学信息
CREATE TABLE teacher_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    school_id INTEGER NOT NULL,
    employee_no VARCHAR(50) NOT NULL,
    subject VARCHAR(50) NOT NULL,
    sunshine_color VARCHAR(20),
    title VARCHAR(50),
    is_homeroom BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_teacher_user ON teacher_profiles(user_id);
COMMENT ON TABLE teacher_profiles IS '教师详细档案表';
```

### 3. admin_permissions（管理员权限分配表）
```sql
-- 管理员权限分配：授权老师管理班级/打分权限
CREATE TABLE admin_permissions (
    id SERIAL PRIMARY KEY,
    granted_by INTEGER NOT NULL REFERENCES users(id),
    teacher_user_id INTEGER NOT NULL REFERENCES users(id),
    class_id INTEGER NOT NULL,
    permission_type VARCHAR(20) NOT NULL CHECK (permission_type IN ('manage', 'score', 'view')),
    school_year VARCHAR(20) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMP
);
CREATE INDEX idx_admin_teacher ON admin_permissions(teacher_user_id);
COMMENT ON TABLE admin_permissions IS '管理员授权教师权限表';
```

---

## 二、学籍班级模块（6张表）
### 4. schools（学校表）
```sql
-- 学校基础信息表
CREATE TABLE schools (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    district VARCHAR(50) NOT NULL,
    address VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE schools IS '学校信息表';
```

### 5. classes（班级表）
```sql
-- 班级表：年级+班级+学年+班主任
CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    class_name VARCHAR(100) NOT NULL,
    class_code VARCHAR(50) NOT NULL UNIQUE,
    grade INTEGER NOT NULL,
    class_no INTEGER NOT NULL,
    school_year VARCHAR(20) NOT NULL,
    homeroom_teacher_id INTEGER REFERENCES users(id),
    student_count INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_class_school ON classes(school_id);
COMMENT ON TABLE classes IS '班级信息表';
```

### 6. children（学生档案表）
```sql
-- 学生档案表：核心学生信息
CREATE TABLE children (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    display_name VARCHAR(50) NOT NULL,
    real_name_enc VARCHAR(255) NOT NULL,
    student_no_enc VARCHAR(255) NOT NULL,
    gender CHAR(1) NOT NULL CHECK (gender IN ('M', 'F')),
    birth_year INTEGER NOT NULL,
    enroll_year INTEGER NOT NULL,
    current_grade INTEGER NOT NULL,
    total_growth_points INTEGER NOT NULL DEFAULT 0,
    current_glory_coins INTEGER NOT NULL DEFAULT 0,
    battle_count INTEGER NOT NULL DEFAULT 0,
    consecutive_days INTEGER NOT NULL DEFAULT 0,
    max_consecutive_days INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_child_user ON children(user_id);
COMMENT ON TABLE children IS '学生档案主表';
```

### 7. class_enrollments（学生班级注册表）
```sql
-- 学生班级关联表：学生在哪个班级、哪个学年
CREATE TABLE class_enrollments (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    class_id INTEGER NOT NULL REFERENCES classes(id),
    school_year VARCHAR(20) NOT NULL,
    enrolled_at DATE NOT NULL,
    left_at DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'graduated', 'transferred')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_child_class_year ON class_enrollments(child_id, class_id, school_year);
COMMENT ON TABLE class_enrollments IS '学生班级注册关系表';
```

### 8. teacher_class_assignments（教师班级分配表）
```sql
-- 教师任教班级表：任课/班主任/权限
CREATE TABLE teacher_class_assignments (
    id SERIAL PRIMARY KEY,
    teacher_user_id INTEGER NOT NULL REFERENCES users(id),
    class_id INTEGER NOT NULL REFERENCES classes(id),
    school_year VARCHAR(20) NOT NULL,
    role_in_class VARCHAR(20) NOT NULL CHECK (role_in_class IN ('homeroom', 'subject_teacher')),
    subject VARCHAR(50) NOT NULL,
    can_score BOOLEAN NOT NULL DEFAULT true,
    can_broadcast BOOLEAN NOT NULL DEFAULT false,
    assigned_by INTEGER REFERENCES users(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_teacher_class ON teacher_class_assignments(teacher_user_id, class_id);
COMMENT ON TABLE teacher_class_assignments IS '教师班级任教分配表';
```

### 9. parent_child_relations（家长学生绑定表）
```sql
-- 家长学生关系表：支持多位家长绑定同一学生
CREATE TABLE parent_child_relations (
    id SERIAL PRIMARY KEY,
    parent_user_id INTEGER NOT NULL REFERENCES users(id),
    child_id INTEGER NOT NULL REFERENCES children(id),
    relation VARCHAR(20) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_parent_child ON parent_child_relations(parent_user_id, child_id);
COMMENT ON TABLE parent_child_relations IS '家长与学生绑定关系表';
```

---

## 三、伙伴系统模块（3张表）
### 10. partner_templates（伙伴模板表）
```sql
-- 伙伴模板：系统预设的伙伴形象
CREATE TABLE partner_templates (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('pet', 'plant', 'character')),
    description TEXT,
    slogan VARCHAR(255),
    low_stage_asset VARCHAR(255) NOT NULL,
    mid_stage_asset VARCHAR(255) NOT NULL,
    high_stage_asset VARCHAR(255) NOT NULL,
    encourage_messages JSONB NOT NULL DEFAULT '{}',
    recommend_grades JSONB NOT NULL DEFAULT '[]',
    sort_order INTEGER NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE partner_templates IS '成长伙伴模板配置表';
```

### 11. partners（学生伙伴实例表）
```sql
-- 学生拥有的伙伴实例
CREATE TABLE partners (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    template_id INTEGER NOT NULL REFERENCES partner_templates(id),
    sequence_no INTEGER NOT NULL,
    nickname VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'pending', 'graduated')),
    growth_points INTEGER NOT NULL DEFAULT 0,
    current_stage INTEGER NOT NULL DEFAULT 1 CHECK (current_stage IN (1,2,3)),
    evolution_count INTEGER NOT NULL DEFAULT 0,
    interaction_level INTEGER NOT NULL DEFAULT 0,
    school_year VARCHAR(20) NOT NULL,
    first_evolution_at TIMESTAMP,
    last_evolved_at TIMESTAMP,
    graduated_at TIMESTAMP,
    selected_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_partner_child ON partners(child_id);
COMMENT ON TABLE partners IS '学生伙伴实例表';
```

### 12. partner_unlock_logs（伙伴解锁日志表）
```sql
-- 伙伴解锁记录：学生解锁新伙伴
CREATE TABLE partner_unlock_logs (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    graduated_partner_id INTEGER REFERENCES partners(id),
    unlock_sequence_no INTEGER NOT NULL,
    graduated_at TIMESTAMP NOT NULL,
    is_new_partner_selected BOOLEAN NOT NULL DEFAULT false,
    new_partner_id INTEGER REFERENCES partners(id),
    selected_at TIMESTAMP
);
COMMENT ON TABLE partner_unlock_logs IS '伙伴解锁日志表';
```

---

## 四、成长记录模块（3张表）
### 13. behavior_records（正向行为记录表）
```sql
-- 学生正向行为记录：老师/家长记录
CREATE TABLE behavior_records (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    class_id INTEGER NOT NULL REFERENCES classes(id),
    school_year VARCHAR(20) NOT NULL,
    recorder_user_id INTEGER NOT NULL REFERENCES users(id),
    recorder_role VARCHAR(20) NOT NULL CHECK (recorder_role IN ('teacher', 'parent', 'system')),
    dimension VARCHAR(20) NOT NULL CHECK (dimension IN ('virtue', 'study', 'sport', 'art', 'labor')),
    description TEXT NOT NULL,
    growth_value INTEGER NOT NULL DEFAULT 0,
    partner_message TEXT,
    is_pushed BOOLEAN NOT NULL DEFAULT false,
    pushed_at TIMESTAMP,
    is_audited BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_behavior_child ON behavior_records(child_id);
COMMENT ON TABLE behavior_records IS '学生正向行为记录表';
```

### 14. growth_records（成长值流水表）
```sql
-- 成长值流水：所有成长值变动记录
CREATE TABLE growth_records (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    partner_id INTEGER NOT NULL REFERENCES partners(id),
    school_year VARCHAR(20) NOT NULL,
    source_type VARCHAR(20) NOT NULL CHECK (source_type IN ('behavior', 'battle', 'system', 'stamp')),
    source_id INTEGER NOT NULL,
    delta INTEGER NOT NULL,
    after_points INTEGER NOT NULL,
    is_evolution_trigger BOOLEAN NOT NULL DEFAULT false,
    evolution_from_stage INTEGER,
    evolution_to_stage INTEGER,
    remark VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_growth_child ON growth_records(child_id);
COMMENT ON TABLE growth_records IS '成长值变动流水表';
```

### 15. milestones（里程碑记录表）
```sql
-- 成长里程碑：首次行为、伙伴进化等
CREATE TABLE milestones (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    partner_id INTEGER NOT NULL REFERENCES partners(id),
    school_year VARCHAR(20) NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    source_type VARCHAR(20),
    source_id INTEGER,
    is_notified BOOLEAN NOT NULL DEFAULT false,
    notified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_milestone_child ON milestones(child_id);
COMMENT ON TABLE milestones IS '学生成长里程碑表';
```

---

## 五、阳光章系统（3张表）
### 16. sunshine_colors（七色章配置表）
```sql
-- 七色阳光章配置：学校自定义颜色对应科目
CREATE TABLE sunshine_colors (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    color_code VARCHAR(20) NOT NULL,
    color_name VARCHAR(20) NOT NULL,
    subject VARCHAR(50) NOT NULL,
    icon_url VARCHAR(255),
    sort_order INTEGER NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_sunshine_color_school ON sunshine_colors(school_id, color_code);
COMMENT ON TABLE sunshine_colors IS '七色阳光章配置表';
```

### 17. sunshine_stamps（阳光章盖章记录表）
```sql
-- 阳光章盖章记录：老师给学生盖章
CREATE TABLE sunshine_stamps (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    class_id INTEGER NOT NULL REFERENCES classes(id),
    school_year VARCHAR(20) NOT NULL,
    color_id INTEGER NOT NULL REFERENCES sunshine_colors(id),
    stamper_user_id INTEGER NOT NULL REFERENCES users(id),
    stamp_month VARCHAR(10) NOT NULL,
    stamp_quarter VARCHAR(10) NOT NULL,
    reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_stamp_child ON sunshine_stamps(child_id);
COMMENT ON TABLE sunshine_stamps IS '阳光章盖章记录表';
```

### 18. sunshine_star_awards（阳光之星评选表）
```sql
-- 阳光之星评选结果
CREATE TABLE sunshine_star_awards (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    class_id INTEGER NOT NULL REFERENCES classes(id),
    school_year VARCHAR(20) NOT NULL,
    award_type VARCHAR(20) NOT NULL CHECK (award_type IN ('monthly_star', 'quarter_star', 'year_star')),
    color_id INTEGER REFERENCES sunshine_colors(id),
    period VARCHAR(20) NOT NULL,
    total_stamps INTEGER NOT NULL DEFAULT 0,
    is_sunshine_star BOOLEAN NOT NULL DEFAULT false,
    certificate_url VARCHAR(255),
    awarded_at TIMESTAMP NOT NULL
);
COMMENT ON TABLE sunshine_star_awards IS '阳光之星评选结果表';
```

---

## 六、知识对战模块（4张表）
### 19. questions（题库表）
```sql
-- 知识对战题库
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    class_id INTEGER REFERENCES classes(id),
    creator_user_id INTEGER REFERENCES users(id),
    subject VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    grade INTEGER NOT NULL,
    content TEXT NOT NULL,
    options JSONB NOT NULL,
    answer VARCHAR(255) NOT NULL,
    explain TEXT,
    difficulty INTEGER NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
    is_active BOOLEAN NOT NULL DEFAULT true,
    use_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE questions IS '知识对战题库表';
```

### 20. battle_rooms（对战房间表）
```sql
-- 对战房间：创建/进行/结束
CREATE TABLE battle_rooms (
    id SERIAL PRIMARY KEY,
    room_code VARCHAR(20) NOT NULL UNIQUE,
    class_id INTEGER NOT NULL REFERENCES classes(id),
    subject VARCHAR(50) NOT NULL,
    mode VARCHAR(20) NOT NULL CHECK (mode IN ('normal', 'friendship', 'competition')),
    status VARCHAR(20) NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'playing', 'finished')),
    question_count INTEGER NOT NULL,
    time_limit_sec INTEGER NOT NULL,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    duration_sec INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE battle_rooms IS '知识对战房间表';
```

### 21. battle_participants（对战参与者表）
```sql
-- 对战参与者：学生参与对战记录
CREATE TABLE battle_participants (
    id SERIAL PRIMARY KEY,
    room_id INTEGER NOT NULL REFERENCES battle_rooms(id),
    child_id INTEGER NOT NULL REFERENCES children(id),
    partner_id INTEGER NOT NULL REFERENCES partners(id),
    score INTEGER NOT NULL DEFAULT 0,
    growth_gained INTEGER NOT NULL DEFAULT 0,
    honor_badge VARCHAR(100),
    is_winner BOOLEAN NOT NULL DEFAULT false,
    submitted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_battle_child ON battle_participants(child_id);
COMMENT ON TABLE battle_participants IS '对战参与者表';
```

### 22. battle_answers（对战答题明细表）
```sql
-- 对战答题明细：每一题答题记录
CREATE TABLE battle_answers (
    id SERIAL PRIMARY KEY,
    participant_id INTEGER NOT NULL REFERENCES battle_participants(id),
    question_id INTEGER NOT NULL REFERENCES questions(id),
    question_order INTEGER NOT NULL,
    answer_given VARCHAR(255) NOT NULL,
    is_correct BOOLEAN NOT NULL,
    time_used_ms INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE battle_answers IS '对战答题明细表';
```

---

## 七、广播与盲盒模块（3张表）
### 23. broadcast_messages（广播消息表）
```sql
-- 广播消息：老师/系统/伙伴发送通知
CREATE TABLE broadcast_messages (
    id SERIAL PRIMARY KEY,
    type VARCHAR(30) NOT NULL,
    sender_user_id INTEGER REFERENCES users(id),
    target_class_id INTEGER REFERENCES classes(id),
    target_child_id INTEGER REFERENCES children(id),
    content TEXT NOT NULL,
    template_key VARCHAR(50),
    is_scheduled BOOLEAN NOT NULL DEFAULT false,
    scheduled_at TIMESTAMP,
    sent_at TIMESTAMP,
    is_sent BOOLEAN NOT NULL DEFAULT false,
    read_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE broadcast_messages IS '系统广播消息表';
```

### 24. blind_box_pools（盲盒奖励池表）
```sql
-- 班级盲盒奖励池
CREATE TABLE blind_box_pools (
    id SERIAL PRIMARY KEY,
    class_id INTEGER NOT NULL REFERENCES classes(id),
    type VARCHAR(30) NOT NULL CHECK (type IN ('privilege', 'honor', 'experience', 'skin')),
    title VARCHAR(100) NOT NULL,
    description TEXT,
    stock INTEGER NOT NULL DEFAULT -1,
    cost_coins INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE blind_box_pools IS '班级盲盒奖励池表';
```

### 25. blind_box_draws（盲盒抽取记录表）
```sql
-- 学生抽取盲盒记录
CREATE TABLE blind_box_draws (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    class_id INTEGER NOT NULL REFERENCES classes(id),
    pool_id INTEGER NOT NULL REFERENCES blind_box_pools(id),
    drawn_at TIMESTAMP NOT NULL,
    is_redeemed BOOLEAN NOT NULL DEFAULT false,
    redeemed_at TIMESTAMP,
    expires_at TIMESTAMP,
    redeemed_by INTEGER REFERENCES users(id)
);
CREATE INDEX idx_blind_child ON blind_box_draws(child_id);
COMMENT ON TABLE blind_box_draws IS '盲盒抽取记录表';
```

---

## 八、成长年历与报告模块（2张表）
### 26. monthly_growth_cards（月度成长卡表）
```sql
-- 学生月度成长报告卡
CREATE TABLE monthly_growth_cards (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    partner_id INTEGER NOT NULL REFERENCES partners(id),
    school_year VARCHAR(20) NOT NULL,
    year_month VARCHAR(10) NOT NULL,
    month_growth_delta INTEGER NOT NULL DEFAULT 0,
    month_behavior_count INTEGER NOT NULL DEFAULT 0,
    highlight_records JSONB NOT NULL DEFAULT '[]',
    partner_snapshot JSONB NOT NULL DEFAULT '{}',
    partner_message TEXT,
    growth_curve JSONB NOT NULL DEFAULT '[]',
    pdf_url VARCHAR(255),
    is_generated BOOLEAN NOT NULL DEFAULT false,
    generated_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_monthly_child ON monthly_growth_cards(child_id, year_month);
COMMENT ON TABLE monthly_growth_cards IS '月度成长报告卡表';
```

### 27. annual_report_logs（年度成长画卷表）
```sql
-- 学生年度成长报告
CREATE TABLE annual_report_logs (
    id SERIAL PRIMARY KEY,
    child_id INTEGER NOT NULL REFERENCES children(id),
    school_year VARCHAR(20) NOT NULL,
    total_growth_points INTEGER NOT NULL DEFAULT 0,
    total_behavior_count INTEGER NOT NULL DEFAULT 0,
    partner_evolution_log JSONB NOT NULL DEFAULT '[]',
    annual_summary TEXT,
    milestones_earned JSONB NOT NULL DEFAULT '[]',
    pdf_url VARCHAR(255),
    is_sent_to_parent BOOLEAN NOT NULL DEFAULT false,
    sent_at TIMESTAMP,
    is_generated BOOLEAN NOT NULL DEFAULT false,
    generated_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_annual_child ON annual_report_logs(child_id, school_year);
COMMENT ON TABLE annual_report_logs IS '年度成长画卷报告表';
```

---

# 执行说明
1. **严格27张表**：与文档v2.0完全一致，无增删
2. **PostgreSQL原生语法**：支持 12+ 版本，可直接执行
3. **完整约束**：主键、外键、唯一索引、检查约束、默认值、注释
4. **执行顺序**：按上面从上到下执行（先主表后子表，避免外键报错）
