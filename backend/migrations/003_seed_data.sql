-- 成长伙伴系统 - 种子数据初始化
-- 包含：学校、班级、测试用户（教师、学生、家长）、伙伴模板、阳光章颜色等

-- ============================================
-- 1. 初始化测试用户 - 教师
-- ============================================
INSERT INTO users (username, password_hash, role, real_name_enc, phone_enc, avatar_url, is_active, last_login_at, created_at, updated_at)
VALUES 
    ('teacher001', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEKjNvl7Y8qY8QY8Cu', 'teacher', '张老师', '13800138001', 'https://api.dicebear.com/7.x/avataaars/svg?seed=teacher001', true, NOW(), NOW(), NOW()),
    ('teacher002', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEKjNvl7Y8qY8QY8Cu', 'teacher', '李老师', '13800138002', 'https://api.dicebear.com/7.x/avataaars/svg?seed=teacher002', true, NOW(), NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- 密码都是: 123456

-- ============================================
-- 2. 初始化测试用户 - 学生
-- ============================================
INSERT INTO users (username, password_hash, role, real_name_enc, phone_enc, avatar_url, is_active, last_login_at, created_at, updated_at)
VALUES 
    ('student001', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEKjNvl7Y8qY8QY8Cu', 'student', '王小明', '13900139001', 'https://api.dicebear.com/7.x/avataaars/svg?seed=student001', true, NOW(), NOW(), NOW()),
    ('student002', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEKjNvl7Y8qY8QY8Cu', 'student', '李小红', '13900139002', 'https://api.dicebear.com/7.x/avataaars/svg?seed=student002', true, NOW(), NOW(), NOW()),
    ('student003', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEKjNvl7Y8qY8QY8Cu', 'student', '张小华', '13900139003', 'https://api.dicebear.com/7.x/avataaars/svg?seed=student003', true, NOW(), NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- ============================================
-- 3. 初始化测试用户 - 家长
-- ============================================
INSERT INTO users (username, password_hash, role, real_name_enc, phone_enc, avatar_url, is_active, last_login_at, created_at, updated_at)
VALUES 
    ('parent001', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEKjNvl7Y8qY8QY8Cu', 'parent', '王爸爸', '13700137001', 'https://api.dicebear.com/7.x/avataaars/svg?seed=parent001', true, NOW(), NOW(), NOW()),
    ('parent002', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEKjNvl7Y8qY8QY8Cu', 'parent', '李妈妈', '13700137002', 'https://api.dicebear.com/7.x/avataaars/svg?seed=parent002', true, NOW(), NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- ============================================
-- 4. 初始化伙伴模板（30个预设伙伴）
-- ============================================
INSERT INTO partner_templates (code, name, type, description, slogan, low_stage_asset, mid_stage_asset, high_stage_asset, encourage_messages, recommend_grades, sort_order, is_active, created_at, updated_at)
VALUES 
    -- 普通伙伴 (N)
    ('rabbit001', '小兔子跳跳', 'animal', '活泼可爱的小兔子，喜欢蹦蹦跳跳', '跳跳跳，快乐每一天！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=rabbit', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=rabbit2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=rabbit3', '["加油！", "你真棒！", "继续保持！"]', '["一年级", "二年级"]', 1, true, NOW(), NOW()),
    ('cat001', '小猫咪咪咪', 'animal', '温顺的小猫咪，喜欢晒太阳', '喵喵喵，今天也要开心哦！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=cat', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=cat2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=cat3', '["做得很好！", "你是最棒的！", "我相信你！"]', '["一年级", "二年级"]', 2, true, NOW(), NOW()),
    ('dog001', '小狗旺财', 'animal', '忠诚的小狗，人类最好的朋友', '汪汪汪，我们一起加油！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=dog', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=dog2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=dog3', '["太厉害了！", "你是我的骄傲！", "继续加油哦！"]', '["一年级", "二年级"]', 3, true, NOW(), NOW()),
    ('panda001', '小熊猫团团', 'animal', '爱吃竹子的小熊猫，圆滚滚的很可爱', '团团转，快乐学习！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=panda', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=panda2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=panda3', '["你真聪明！", "进步很大！", "继续保持！"]', '["一年级", "二年级", "三年级"]', 4, true, NOW(), NOW()),
    ('squirrel001', '小松鼠果果', 'animal', '勤劳的小松鼠，喜欢收集坚果', '储存知识，收获成长！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=squirrel', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=squirrel2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=squirrel3', '["好样的！", "你真勤奋！", "知识就是力量！"]', '["二年级", "三年级"]', 5, true, NOW(), NOW()),
    ('frog001', '小青蛙呱呱', 'animal', '爱唱歌的小青蛙，雨后叫得最欢', '呱呱呱，学习顶呱呱！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=frog', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=frog2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=frog3', '["顶呱呱！", "你真棒！", "学习顶呱呱！"]', '["二年级", "三年级"]', 6, true, NOW(), NOW()),
    
    -- 稀有伙伴 (R)
    ('fox001', '小狐狸阿狸', 'animal', '聪明的小狐狸，有很多鬼点子', '聪明伶俐，学习有方！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=fox', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=fox2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=fox3', '["聪明绝顶！", "机智过人！", "点子真多！"]', '["三年级", "四年级"]', 10, true, NOW(), NOW()),
    ('deer001', '小鹿斑斑', 'animal', '优雅的小鹿，在森林里自由奔跑', '优雅前行，步步高升！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=deer', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=deer2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=deer3', '["优雅大方！", "你真出色！", "前途无量！"]', '["三年级", "四年级"]', 11, true, NOW(), NOW()),
    ('hedgehog001', '小刺猬球球', 'animal', '外表有刺但内心柔软的小刺猬', '保护自己，温暖他人！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=hedgehog', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=hedgehog2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=hedgehog3', '["坚强勇敢！", "内心温暖！", "细心体贴！"]', '["四年级", "五年级"]', 12, true, NOW(), NOW()),
    ('penguin001', '小企鹅冰冰', 'animal', '来自南极的小企鹅，不怕寒冷', '勇敢前行，无所畏惧！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=penguin', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=penguin2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=penguin3', '["耐寒坚强！", "团结友爱！", "你真棒！"]', '["四年级", "五年级"]', 13, true, NOW(), NOW()),
    ('dolphin001', '小海豚蓝蓝', 'animal', '聪明的小海豚，海洋里的精灵', '聪明灵动，快乐学习！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=dolphin', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=dolphin2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=dolphin3', '["聪明伶俐！", "海洋精灵！", "你真出色！"]', '["五年级", "六年级"]', 14, true, NOW(), NOW()),
    
    -- 史诗伙伴 (SR)
    ('lion001', '小狮子雷雷', 'animal', '勇敢的小狮子，草原之王的后代', '勇敢无畏，王者风范！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=lion', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=lion2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=lion3', '["勇敢霸气！", "王者风范！", "正义凛然！"]', '["五年级", "六年级"]', 20, true, NOW(), NOW()),
    ('tiger001', '小老虎威威', 'animal', '威风凛凛的小老虎，百兽之王', '威风凛凛，勇往直前！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=tiger', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=tiger2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=tiger3', '["威风八面！", "勇敢无畏！", "独立自强！"]', '["五年级", "六年级"]', 21, true, NOW(), NOW()),
    ('unicorn001', '小独角兽闪闪', 'fantasy', '传说中的独角兽，拥有魔法力量', '魔法力量，梦想成真！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=unicorn', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=unicorn2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=unicorn3', '["神奇魔法！", "纯洁高贵！", "梦想成真！"]', '["六年级"]', 22, true, NOW(), NOW()),
    ('phoenix001', '小凤凰焰焰', 'fantasy', '不死鸟凤凰，象征着重生和希望', '浴火重生，希望永存！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=phoenix', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=phoenix2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=phoenix3', '["坚韧不拔！", "热情洋溢！", "希望之光！"]', '["六年级"]', 23, true, NOW(), NOW()),
    ('dragon001', '小神龙腾腾', 'fantasy', '东方神龙，带来好运和福气', '龙腾虎跃，好运连连！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=dragon', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=dragon2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=dragon3', '["神秘强大！", "祥瑞之兆！", "龙腾虎跃！"]', '["六年级"]', 24, true, NOW(), NOW()),
    
    -- 传说伙伴 (SSR)
    ('elf001', '小精灵露露', 'fantasy', '森林精灵，守护自然的小仙子', '守护自然，快乐成长！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=elf', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=elf2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=elf3', '["灵动可爱！", "善良纯真！", "守护之力！"]', '["全年级"]', 30, true, NOW(), NOW()),
    ('angel001', '小天使安琪', 'fantasy', '来自天堂的小天使，传播爱与光明', '爱与光明，温暖人心！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=angel', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=angel2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=angel3', '["善良纯洁！", "光明使者！", "爱的化身！"]', '["全年级"]', 31, true, NOW(), NOW()),
    ('devil001', '小恶魔皮皮', 'fantasy', '调皮的小恶魔，其实内心很善良', '调皮可爱，心地善良！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=devil', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=devil2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=devil3', '["调皮聪明！", "心地善良！", "鬼点子多！"]', '["全年级"]', 32, true, NOW(), NOW()),
    ('robot001', '小机器人豆豆', 'robot', '来自未来的智能机器人', '智能科技，未来已来！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=robot', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=robot2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=robot3', '["智能可靠！", "科技先锋！", "未来之星！"]', '["全年级"]', 33, true, NOW(), NOW()),
    
    -- 更多普通伙伴
    ('duck001', '小鸭子嘎嘎', 'animal', '爱游泳的小鸭子，排队走路很可爱', '嘎嘎嘎，快乐学习！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=duck', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=duck2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=duck3', '["活泼可爱！", "团结友爱！", "排队整齐！"]', '["一年级", "二年级"]', 7, true, NOW(), NOW()),
    ('chick001', '小鸡叽叽', 'animal', '毛茸茸的小鸡，总是跟着鸡妈妈', '叽叽叽，快乐成长！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=chick', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=chick2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=chick3', '["乖巧可爱！", "活泼好动！", "毛茸茸的！"]', '["一年级"]', 8, true, NOW(), NOW()),
    ('sheep001', '小绵羊绵绵', 'animal', '温顺的小绵羊，毛茸茸像棉花糖', '绵绵软软，温暖人心！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=sheep', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=sheep2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=sheep3', '["温顺乖巧！", "合群友善！", "像棉花糖！"]', '["一年级", "二年级"]', 9, true, NOW(), NOW()),
    ('cow001', '小牛牛牛', 'animal', '勤劳的小牛，耕地的好帮手', '勤勤恳恳，脚踏实地！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=cow', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=cow2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=cow3', '["勤劳踏实！", "可靠稳重！", "勤勤恳恳！"]', '["二年级", "三年级"]', 15, true, NOW(), NOW()),
    ('horse001', '小马奔奔', 'animal', '爱奔跑的小马，草原上的风', '奔腾不息，勇往直前！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=horse', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=horse2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=horse3', '["自由奔放！", "活力四射！", "奔腾不息！"]', '["三年级", "四年级"]', 16, true, NOW(), NOW()),
    
    -- 更多稀有伙伴
    ('koala001', '小考拉懒懒', 'animal', '爱睡觉的小考拉，一天睡20小时', '好好休息，好好学习！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=koala', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=koala2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=koala3', '["慵懒可爱！", "安静乖巧！", "好好休息！"]', '["三年级", "四年级"]', 17, true, NOW(), NOW()),
    ('kangaroo001', '小袋鼠跳跳', 'animal', '口袋里装着小宝宝的小袋鼠', '跳跃前进，茁壮成长！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=kangaroo', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=kangaroo2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=kangaroo3', '["母爱伟大！", "强壮有力！", "跳跃高手！"]', '["四年级", "五年级"]', 18, true, NOW(), NOW()),
    ('raccoon001', '小浣熊洗洗', 'animal', '爱干净的小浣熊，吃东西前要洗洗', '爱干净，讲卫生！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=raccoon', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=raccoon2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=raccoon3', '["爱干净！", "机灵聪明！", "讲卫生！"]', '["四年级", "五年级"]', 19, true, NOW(), NOW()),
    ('owl001', '小猫头鹰夜夜', 'animal', '夜间活动的小猫头鹰，智慧的象征', '智慧之光，照亮前方！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=owl', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=owl2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=owl3', '["智慧过人！", "警觉敏锐！", "神秘莫测！"]', '["五年级", "六年级"]', 25, true, NOW(), NOW()),
    ('bee001', '小蜜蜂嗡嗡', 'animal', '勤劳的小蜜蜂，采蜜忙不停', '勤劳采蜜，收获甜蜜！', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=bee', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=bee2', 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=bee3', '["勤劳奉献！", "团结协作！", "采蜜忙！"]', '["全年级"]', 26, true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- 种子数据初始化完成
SELECT '种子数据初始化完成！' as status;
SELECT '测试账号:' as info;
SELECT '  教师: teacher001 / 123456' as account;
SELECT '  学生: student001 / 123456' as account;
SELECT '  家长: parent001 / 123456' as account;
