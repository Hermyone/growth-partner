const {
  Document,
  Packer,
  Paragraph,
  TextRun,
  Table,
  TableRow,
  TableCell,
  HeadingLevel,
  AlignmentType,
  BorderStyle,
  WidthType,
  ShadingType,
  VerticalAlign,
  LevelFormat,
  PageBreak,
} = require("docx");
const fs = require("fs");

// ─── 颜色 ──────────────────────────────────────────────────────
const C = {
  navy: "1F4E79",
  blue: "2E75B6",
  lightBlue: "D6E4F0",
  deepBlue: "1A3A5C",
  green: "1E6B3C",
  lightGreen: "D8F0E3",
  orange: "C45911",
  lightOrange: "FCE9D8",
  purple: "5C2D91",
  lightPurple: "EDE7F6",
  teal: "006064",
  lightTeal: "E0F7FA",
  red: "C62828",
  lightRed: "FFEBEE",
  gray: "404040",
  lightGray: "F5F5F5",
  midGray: "E0E0E0",
  white: "FFFFFF",
  rowAlt: "EBF3FB",
  border: "BFBFBF",
};

// ─── 边框 ──────────────────────────────────────────────────────
const B = (color = C.border) => ({ style: BorderStyle.SINGLE, size: 4, color });
const allBorders = (c = C.border) => ({
  top: B(c),
  bottom: B(c),
  left: B(c),
  right: B(c),
});
const noBorder = { style: BorderStyle.NONE, size: 0, color: "FFFFFF" };
const noBorders = {
  top: noBorder,
  bottom: noBorder,
  left: noBorder,
  right: noBorder,
};

// ─── 工具函数 ─────────────────────────────────────────────────
const run = (text, opts = {}) =>
  new TextRun({
    text,
    font: "Consolas",
    size: opts.size || 18,
    bold: opts.bold || false,
    color: opts.color || "1A1A1A",
    italics: opts.italic || false,
    highlight: opts.highlight || undefined,
  });

const p = (text, opts = {}) =>
  new Paragraph({
    children: [
      new TextRun({
        text,
        font: opts.mono ? "Consolas" : "Arial",
        size: opts.size || 20,
        bold: opts.bold || false,
        color: opts.color || "1A1A1A",
        italics: opts.italic || false,
      }),
    ],
    spacing: { before: opts.before || 60, after: opts.after || 60 },
    alignment: opts.align || AlignmentType.LEFT,
    indent: opts.indent ? { left: opts.indent } : undefined,
    shading: opts.bg ? { fill: opts.bg, type: ShadingType.CLEAR } : undefined,
  });

const pEmpty = (s = 80) =>
  new Paragraph({
    children: [new TextRun("")],
    spacing: { before: s, after: s },
  });

const pBullet = (text, color = "1A1A1A", indent = 360) =>
  new Paragraph({
    numbering: { reference: "bullets", level: 0 },
    children: [new TextRun({ text, font: "Arial", size: 19, color })],
    spacing: { before: 40, after: 40 },
  });

const pNum = (text, color = "1A1A1A") =>
  new Paragraph({
    numbering: { reference: "numbers", level: 0 },
    children: [new TextRun({ text, font: "Arial", size: 19, color })],
    spacing: { before: 40, after: 40 },
  });

// ─── 标题 ──────────────────────────────────────────────────────
const h1 = (text, bg = C.navy) =>
  new Paragraph({
    heading: HeadingLevel.HEADING_1,
    children: [
      new TextRun({
        text,
        font: "Arial",
        size: 38,
        bold: true,
        color: C.white,
      }),
    ],
    shading: { fill: bg, type: ShadingType.CLEAR },
    spacing: { before: 360, after: 200 },
    indent: { left: 200 },
  });

const h2 = (text, bg = C.blue) =>
  new Paragraph({
    heading: HeadingLevel.HEADING_2,
    children: [
      new TextRun({
        text,
        font: "Arial",
        size: 28,
        bold: true,
        color: C.white,
      }),
    ],
    shading: { fill: bg, type: ShadingType.CLEAR },
    spacing: { before: 280, after: 160 },
    indent: { left: 160 },
  });

const h3 = (text, color = C.navy) =>
  new Paragraph({
    heading: HeadingLevel.HEADING_3,
    children: [
      new TextRun({ text, font: "Arial", size: 24, bold: true, color }),
    ],
    spacing: { before: 220, after: 100 },
    border: { bottom: { style: BorderStyle.SINGLE, size: 6, color: C.blue } },
  });

const h4 = (text, color = C.green) =>
  new Paragraph({
    children: [
      new TextRun({
        text: `▶  ${text}`,
        font: "Arial",
        size: 21,
        bold: true,
        color,
      }),
    ],
    spacing: { before: 160, after: 60 },
  });

// ─── 注解框 ───────────────────────────────────────────────────
const note = (text, icon = "💡", bg = C.lightBlue, color = C.navy) =>
  new Paragraph({
    children: [
      new TextRun({
        text: `${icon}  ${text}`,
        font: "Arial",
        size: 18,
        color,
        italics: true,
      }),
    ],
    shading: { fill: bg, type: ShadingType.CLEAR },
    spacing: { before: 80, after: 80 },
    indent: { left: 280, right: 280 },
  });

const warn = (text) => note(text, "⚠️", C.lightOrange, C.orange);
const success = (text) => note(text, "✅", C.lightGreen, C.green);
const critical = (text) => note(text, "🔒", C.lightRed, C.red);

// ─── 路由行辅助 ───────────────────────────────────────────────
// 路由表列宽：Method(900) + 路径(3200) + 权限(1000) + 说明(4260) = 9360
const RT_WIDTHS = [900, 3200, 1000, 4260];
const METHOD_COLORS = {
  GET: { bg: "006064", text: "FFFFFF" },
  POST: { bg: "1E6B3C", text: "FFFFFF" },
  PUT: { bg: "C45911", text: "FFFFFF" },
  PATCH: { bg: "5C2D91", text: "FFFFFF" },
  DELETE: { bg: "C62828", text: "FFFFFF" },
  WS: { bg: "1A3A5C", text: "FFFFFF" },
};

const routeTableHeader = () =>
  new TableRow({
    tableHeader: true,
    children: ["Method", "路由路径", "权限", "功能说明"].map(
      (t, i) =>
        new TableCell({
          borders: allBorders(C.navy),
          width: { size: RT_WIDTHS[i], type: WidthType.DXA },
          shading: { fill: C.navy, type: ShadingType.CLEAR },
          margins: { top: 100, bottom: 100, left: 120, right: 120 },
          verticalAlign: VerticalAlign.CENTER,
          children: [
            new Paragraph({
              alignment: AlignmentType.CENTER,
              children: [
                new TextRun({
                  text: t,
                  font: "Arial",
                  size: 18,
                  bold: true,
                  color: C.white,
                }),
              ],
            }),
          ],
        })
    ),
  });

const routeRow = (method, path, role, desc, isAlt = false) => {
  const mc = METHOD_COLORS[method] || { bg: C.gray, text: C.white };
  return new TableRow({
    children: [
      new TableCell({
        borders: allBorders(C.border),
        width: { size: RT_WIDTHS[0], type: WidthType.DXA },
        shading: { fill: mc.bg, type: ShadingType.CLEAR },
        margins: { top: 80, bottom: 80, left: 80, right: 80 },
        verticalAlign: VerticalAlign.CENTER,
        children: [
          new Paragraph({
            alignment: AlignmentType.CENTER,
            children: [
              new TextRun({
                text: method,
                font: "Consolas",
                size: 16,
                bold: true,
                color: mc.text,
              }),
            ],
          }),
        ],
      }),
      new TableCell({
        borders: allBorders(C.border),
        width: { size: RT_WIDTHS[1], type: WidthType.DXA },
        shading: {
          fill: isAlt ? C.lightGray : C.white,
          type: ShadingType.CLEAR,
        },
        margins: { top: 80, bottom: 80, left: 120, right: 120 },
        verticalAlign: VerticalAlign.CENTER,
        children: [
          new Paragraph({
            children: [
              new TextRun({
                text: path,
                font: "Consolas",
                size: 17,
                color: C.deepBlue,
              }),
            ],
          }),
        ],
      }),
      new TableCell({
        borders: allBorders(C.border),
        width: { size: RT_WIDTHS[2], type: WidthType.DXA },
        shading: {
          fill: isAlt ? C.lightGray : C.white,
          type: ShadingType.CLEAR,
        },
        margins: { top: 80, bottom: 80, left: 80, right: 80 },
        verticalAlign: VerticalAlign.CENTER,
        children: [
          new Paragraph({
            alignment: AlignmentType.CENTER,
            children: [
              new TextRun({
                text: role,
                font: "Arial",
                size: 16,
                color: C.gray,
              }),
            ],
          }),
        ],
      }),
      new TableCell({
        borders: allBorders(C.border),
        width: { size: RT_WIDTHS[3], type: WidthType.DXA },
        shading: {
          fill: isAlt ? C.lightGray : C.white,
          type: ShadingType.CLEAR,
        },
        margins: { top: 80, bottom: 80, left: 120, right: 120 },
        verticalAlign: VerticalAlign.CENTER,
        children: [
          new Paragraph({
            children: [
              new TextRun({
                text: desc,
                font: "Arial",
                size: 18,
                color: "1A1A1A",
              }),
            ],
          }),
        ],
      }),
    ],
  });
};

const routeTable = (rows) =>
  new Table({
    width: { size: 9360, type: WidthType.DXA },
    columnWidths: RT_WIDTHS,
    rows: [
      routeTableHeader(),
      ...rows.map((r, i) => routeRow(r[0], r[1], r[2], r[3], i % 2 === 1)),
    ],
  });

// ─── 文件说明表 ───────────────────────────────────────────────
const FILE_WIDTHS = [2800, 1600, 4960];
const fileTable = (rows) =>
  new Table({
    width: { size: 9360, type: WidthType.DXA },
    columnWidths: FILE_WIDTHS,
    rows: [
      new TableRow({
        tableHeader: true,
        children: ["文件路径", "对应层", "职责说明"].map(
          (t, i) =>
            new TableCell({
              borders: allBorders(C.teal),
              width: { size: FILE_WIDTHS[i], type: WidthType.DXA },
              shading: { fill: C.teal, type: ShadingType.CLEAR },
              margins: { top: 100, bottom: 100, left: 120, right: 120 },
              children: [
                new Paragraph({
                  alignment: AlignmentType.CENTER,
                  children: [
                    new TextRun({
                      text: t,
                      font: "Arial",
                      size: 18,
                      bold: true,
                      color: C.white,
                    }),
                  ],
                }),
              ],
            })
        ),
      }),
      ...rows.map(
        (r, i) =>
          new TableRow({
            children: r.map(
              (t, j) =>
                new TableCell({
                  borders: allBorders(C.border),
                  width: { size: FILE_WIDTHS[j], type: WidthType.DXA },
                  shading: {
                    fill: i % 2 === 1 ? C.lightTeal : C.white,
                    type: ShadingType.CLEAR,
                  },
                  margins: { top: 80, bottom: 80, left: 120, right: 120 },
                  children: [
                    new Paragraph({
                      children: [
                        new TextRun({
                          text: t,
                          font: j === 0 ? "Consolas" : "Arial",
                          size: j === 0 ? 16 : 18,
                          color: j === 0 ? C.teal : "1A1A1A",
                        }),
                      ],
                    }),
                  ],
                })
            ),
          })
      ),
    ],
  });

// ─── 模块说明表（两列） ───────────────────────────────────────
const MOD_WIDTHS = [2200, 7160];
const modTable = (rows, hbg = C.purple) =>
  new Table({
    width: { size: 9360, type: WidthType.DXA },
    columnWidths: MOD_WIDTHS,
    rows: [
      new TableRow({
        tableHeader: true,
        children: ["要点", "详细说明"].map(
          (t, i) =>
            new TableCell({
              borders: allBorders(hbg),
              width: { size: MOD_WIDTHS[i], type: WidthType.DXA },
              shading: { fill: hbg, type: ShadingType.CLEAR },
              margins: { top: 100, bottom: 100, left: 120, right: 120 },
              children: [
                new Paragraph({
                  alignment: AlignmentType.CENTER,
                  children: [
                    new TextRun({
                      text: t,
                      font: "Arial",
                      size: 18,
                      bold: true,
                      color: C.white,
                    }),
                  ],
                }),
              ],
            })
        ),
      }),
      ...rows.map(
        (r, i) =>
          new TableRow({
            children: [
              new TableCell({
                borders: allBorders(C.border),
                width: { size: MOD_WIDTHS[0], type: WidthType.DXA },
                shading: {
                  fill: i % 2 === 1 ? C.lightPurple : C.white,
                  type: ShadingType.CLEAR,
                },
                margins: { top: 80, bottom: 80, left: 120, right: 120 },
                children: [
                  new Paragraph({
                    children: [
                      new TextRun({
                        text: r[0],
                        font: "Arial",
                        size: 18,
                        bold: true,
                        color: C.purple,
                      }),
                    ],
                  }),
                ],
              }),
              new TableCell({
                borders: allBorders(C.border),
                width: { size: MOD_WIDTHS[1], type: WidthType.DXA },
                shading: {
                  fill: i % 2 === 1 ? C.lightPurple : C.white,
                  type: ShadingType.CLEAR,
                },
                margins: { top: 80, bottom: 80, left: 120, right: 120 },
                children: [
                  new Paragraph({
                    children: [
                      new TextRun({
                        text: r[1],
                        font: "Arial",
                        size: 18,
                        color: "1A1A1A",
                      }),
                    ],
                  }),
                ],
              }),
            ],
          })
      ),
    ],
  });

// ═══════════════════════════════════════════════════════════════════
// 正文内容构建
// ═══════════════════════════════════════════════════════════════════
const children = [];

// ─── 封面 ──────────────────────────────────────────────────────
children.push(pEmpty(400));
children.push(
  new Paragraph({
    children: [
      new TextRun({
        text: "成长伙伴",
        font: "Arial",
        size: 80,
        bold: true,
        color: C.navy,
      }),
    ],
    alignment: AlignmentType.CENTER,
    spacing: { before: 0, after: 100 },
  })
);
children.push(
  new Paragraph({
    children: [
      new TextRun({
        text: "Growth Partner System",
        font: "Arial",
        size: 30,
        color: C.blue,
        italics: true,
      }),
    ],
    alignment: AlignmentType.CENTER,
    spacing: { before: 0, after: 200 },
  })
);
children.push(
  new Paragraph({
    children: [
      new TextRun({
        text: "后端 API 完整功能清单 & 路由设计",
        font: "Arial",
        size: 46,
        bold: true,
        color: C.gray,
      }),
    ],
    alignment: AlignmentType.CENTER,
    spacing: { before: 0, after: 80 },
  })
);
children.push(
  new Paragraph({
    children: [
      new TextRun({
        text: "Backend API Specification  ·  架构师审查版  ·  v1.0",
        font: "Arial",
        size: 22,
        color: C.gray,
      }),
    ],
    alignment: AlignmentType.CENTER,
    spacing: { before: 0, after: 500 },
  })
);

// 封面元信息表
const META_W = [2800, 3760];
children.push(
  new Table({
    width: { size: 6560, type: WidthType.DXA },
    columnWidths: META_W,
    rows: [
      [
        "技术栈",
        "Golang 1.23 · Gin · GORM · PostgreSQL 16 · Redis 7 · WebSocket",
      ],
      ["架构模式", "Handler → Service → Repository 三层架构，接口驱动"],
      ["认证方案", "JWT Bearer Token，三角色（学生/老师/家长/管理员）"],
      ["API 版本", "/api/v1  全部路由前缀"],
      ["响应格式", "统一 JSON：{ code, message, data, error, timestamp }"],
      ["总路由数", "约 72 条（含 WebSocket 2 条）"],
      ["总功能模块", "10 大模块，27 张数据表"],
    ].map(
      (r, i) =>
        new TableRow({
          children: [
            new TableCell({
              borders: allBorders(C.lightBlue),
              width: { size: META_W[0], type: WidthType.DXA },
              shading: { fill: C.lightBlue, type: ShadingType.CLEAR },
              margins: { top: 100, bottom: 100, left: 200, right: 200 },
              children: [
                new Paragraph({
                  alignment: AlignmentType.RIGHT,
                  children: [
                    new TextRun({
                      text: r[0],
                      font: "Arial",
                      size: 19,
                      bold: true,
                      color: C.navy,
                    }),
                  ],
                }),
              ],
            }),
            new TableCell({
              borders: allBorders(C.lightBlue),
              width: { size: META_W[1], type: WidthType.DXA },
              margins: { top: 100, bottom: 100, left: 200, right: 200 },
              children: [
                new Paragraph({
                  children: [
                    new TextRun({
                      text: r[1],
                      font: "Arial",
                      size: 19,
                      color: "1A1A1A",
                    }),
                  ],
                }),
              ],
            }),
          ],
        })
    ),
  })
);
children.push(pEmpty(200));

children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第一章：工程结构 & 分层约定
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第一章  工程结构 & 分层约定"));
children.push(pEmpty(40));
children.push(
  note(
    "每个功能模块都必须严格按照 Handler → Service → Repository 三层实现，任何层不得跨层调用。Handler 只做参数绑定+响应，Service 含全部业务逻辑，Repository 只做数据库/Redis操作。"
  )
);
children.push(pEmpty(40));

children.push(h3("1.1  文件结构总览"));
children.push(
  fileTable([
    ["backend/main.go", "入口", "依赖注入，启动 HTTP 服务，优雅关闭"],
    [
      "backend/config/config.go",
      "配置层",
      "从环境变量加载所有配置（DB/Redis/JWT/AES），无配置文件",
    ],
    [
      "backend/internal/router/router.go",
      "路由层",
      "注册所有路由，挂载中间件，按角色分组",
    ],
    [
      "backend/internal/middleware/*.go",
      "中间件层",
      "JWT认证 / 角色鉴权 / 统一响应 / CORS / 请求日志 / RequestID",
    ],
    [
      "backend/internal/model/*.go",
      "模型层",
      "GORM Struct定义，27张表的Go结构体",
    ],
    [
      "backend/internal/repository/*.go",
      "Repository层",
      "纯数据库/Redis操作，返回领域对象，禁止含业务逻辑",
    ],
    [
      "backend/internal/service/*.go",
      "Service层",
      "全部业务逻辑，事务管理，跨Repository协调",
    ],
    [
      "backend/internal/handler/*.go",
      "Handler层",
      "参数绑定+校验，调Service，写统一响应，禁止含业务逻辑",
    ],
    [
      "backend/internal/websocket/*.go",
      "WebSocket层",
      "Hub连接管理，客户端心跳，消息协议",
    ],
    ["backend/pkg/jwt/jwt.go", "工具包", "JWT生成/验证，Claims定义"],
    ["backend/pkg/encrypt/aes.go", "工具包", "AES-256-GCM 加解密（敏感字段）"],
    [
      "backend/pkg/validator/validator.go",
      "工具包",
      "Gin参数校验器注册，中文错误信息",
    ],
    [
      "backend/migrations/*.sql",
      "数据库迁移",
      "按序号命名的SQL文件，初始化表结构和种子数据",
    ],
    [
      "backend/scripts/seed_templates.go",
      "数据初始化",
      "30个伙伴模板初始数据写入",
    ],
  ])
);

children.push(pEmpty(60));

children.push(h3("1.2  中间件执行顺序（全局）"));
children.push(p("所有请求按以下顺序经过中间件：", { bold: true }));
children.push(
  pNum(
    "Recovery()  →  请求日志(Logger)  →  RequestID注入  →  CORS  →  [路由匹配]  →  JWT认证(Auth)  →  角色鉴权(RequireRole)  →  Handler"
  )
);
children.push(pEmpty(40));

children.push(h3("1.3  权限角色说明"));
children.push(
  new Table({
    width: { size: 9360, type: WidthType.DXA },
    columnWidths: [1200, 1600, 6560],
    rows: [
      new TableRow({
        tableHeader: true,
        children: ["角色值", "角色名称", "权限说明"].map(
          (t, i) =>
            new TableCell({
              borders: allBorders(C.navy),
              width: { size: [1200, 1600, 6560][i], type: WidthType.DXA },
              shading: { fill: C.navy, type: ShadingType.CLEAR },
              margins: { top: 80, bottom: 80, left: 120, right: 120 },
              children: [
                new Paragraph({
                  alignment: AlignmentType.CENTER,
                  children: [
                    new TextRun({
                      text: t,
                      font: "Arial",
                      size: 18,
                      bold: true,
                      color: C.white,
                    }),
                  ],
                }),
              ],
            })
        ),
      }),
      ...[
        [
          "admin",
          "系统管理员",
          "最高权限：创建学校/班级、分配老师权限、管理所有数据、查看所有报表",
        ],
        [
          "teacher",
          "老师/班主任/教练",
          "管理被授权班级：打分、发广播、管理盲盒、查看班级概况、生成报告",
        ],
        [
          "parent",
          "家长",
          "只读：查看自己孩子的伙伴/行为/广播，不可跨孩子查询",
        ],
        [
          "student",
          "学生",
          "只读+互动：查看自己伙伴、参与对战、查看广播、选择伙伴",
        ],
      ].map(
        (r, i) =>
          new TableRow({
            children: r.map(
              (t, j) =>
                new TableCell({
                  borders: allBorders(C.border),
                  width: { size: [1200, 1600, 6560][j], type: WidthType.DXA },
                  shading: {
                    fill: i % 2 === 1 ? C.rowAlt : C.white,
                    type: ShadingType.CLEAR,
                  },
                  margins: { top: 80, bottom: 80, left: 120, right: 120 },
                  children: [
                    new Paragraph({
                      children: [
                        new TextRun({
                          text: t,
                          font: j === 0 ? "Consolas" : "Arial",
                          size: j === 0 ? 17 : 18,
                          bold: j === 0,
                          color: j === 0 ? C.orange : "1A1A1A",
                        }),
                      ],
                    }),
                  ],
                })
            ),
          })
      ),
    ],
  })
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第二章：认证模块
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第二章  认证模块（Auth）", C.navy));
children.push(pEmpty(40));
children.push(
  note(
    "所有角色共用同一套登录接口，通过 role 字段区分。Token 分 Access（24h）和 Refresh（7d）两种，前端用 Refresh Token 静默续期，不强制重登。"
  )
);
children.push(pEmpty(40));

children.push(h3("2.1  路由清单"));
children.push(
  routeTable([
    [
      "POST",
      "/api/v1/auth/login",
      "公开",
      "统一登录入口，返回 access_token + refresh_token + 用户信息 + 角色",
    ],
    [
      "POST",
      "/api/v1/auth/refresh",
      "公开",
      "用 refresh_token 换新 access_token，刷新登录态，不需重新输密码",
    ],
    [
      "POST",
      "/api/v1/auth/logout",
      "JWT",
      "注销登录，将 refresh_token 加入 Redis 黑名单（TTL=token剩余有效期）",
    ],
    [
      "GET",
      "/api/v1/auth/me",
      "JWT",
      "获取当前登录用户的基础信息（角色/班级/用户ID）",
    ],
    [
      "PATCH",
      "/api/v1/auth/password",
      "JWT",
      "修改密码（需验证旧密码），修改后强制所有 Token 失效",
    ],
  ])
);

children.push(pEmpty(40));
children.push(h3("2.2  功能要点"));
children.push(
  modTable(
    [
      [
        "登录逻辑",
        "验证 username+password → 查 users 表 → bcrypt 校验 → 生成 JWT Claims（含 user_id/role/class_id/child_id）→ 写 last_login_at → 返回双 Token",
      ],
      [
        "Token 黑名单",
        "注销时将 jti（Token唯一ID）写入 Redis，TTL = Token 剩余有效期；Auth中间件每次验证都检查黑名单",
      ],
      [
        "密码策略",
        "最少8位，必须含字母+数字；bcrypt cost=12；修改密码后写 Redis 标记（password_version），让旧 Token 失效",
      ],
      [
        "敏感信息",
        "login 接口返回 display_name（非真实姓名），不返回加密字段，不返回 password_hash",
      ],
    ],
    C.navy
  )
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第三章：管理员模块
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第三章  管理员模块（Admin）", C.red));
children.push(pEmpty(40));
children.push(
  critical("所有 /api/v1/admin/* 路由必须同时验证 JWT + role=admin，双重校验。")
);
children.push(pEmpty(40));

children.push(h3("3.1  学校管理"));
children.push(
  routeTable([
    ["GET", "/api/v1/admin/schools", "admin", "获取学校列表（分页+搜索）"],
    ["POST", "/api/v1/admin/schools", "admin", "创建学校"],
    ["PUT", "/api/v1/admin/schools/:id", "admin", "更新学校信息"],
    ["PATCH", "/api/v1/admin/schools/:id/status", "admin", "启用/停用学校"],
  ])
);
children.push(pEmpty(40));

children.push(h3("3.2  班级管理"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/admin/classes",
      "admin",
      "获取班级列表（可按学校/学年/年级筛选）",
    ],
    [
      "POST",
      "/api/v1/admin/classes",
      "admin",
      "创建新班级（写入 class_code 唯一校验）",
    ],
    [
      "PUT",
      "/api/v1/admin/classes/:id",
      "admin",
      "更新班级信息（班级名/班主任）",
    ],
    [
      "POST",
      "/api/v1/admin/classes/:id/promote",
      "admin",
      "升年级操作：grade+1，更新 school_year，批量为该班学生新建 class_enrollments",
    ],
    ["PATCH", "/api/v1/admin/classes/:id/status", "admin", "启用/停用班级"],
  ])
);
children.push(pEmpty(40));

children.push(h3("3.3  用户管理（老师/家长账号）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/admin/users",
      "admin",
      "获取用户列表（按角色/学校筛选，分页）",
    ],
    [
      "POST",
      "/api/v1/admin/users",
      "admin",
      "创建用户账号（老师/家长），初始密码由管理员设置",
    ],
    ["PUT", "/api/v1/admin/users/:id", "admin", "更新用户信息"],
    ["PATCH", "/api/v1/admin/users/:id/status", "admin", "启用/停用账号"],
    [
      "PATCH",
      "/api/v1/admin/users/:id/reset-pwd",
      "admin",
      "重置用户密码（管理员操作，无需旧密码）",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("3.4  学生账号批量管理"));
children.push(
  routeTable([
    [
      "POST",
      "/api/v1/admin/students/batch-import",
      "admin",
      "批量导入学生（CSV上传）：创建 users + children + class_enrollments",
    ],
    [
      "GET",
      "/api/v1/admin/students",
      "admin",
      "学生列表（按班级/学年筛选，脱敏展示）",
    ],
    ["POST", "/api/v1/admin/students", "admin", "单个创建学生账号"],
    ["PUT", "/api/v1/admin/students/:id", "admin", "更新学生信息"],
  ])
);
children.push(pEmpty(40));

children.push(h3("3.5  老师班级权限分配（核心：张老师管N个班）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/admin/assignments",
      "admin",
      "查看所有老师-班级分配关系（按学校/学年筛选）",
    ],
    [
      "POST",
      "/api/v1/admin/assignments",
      "admin",
      "为老师分配班级权限（写 teacher_class_assignments）：指定 teacher_user_id/class_id/role_in_class/subject/can_score/can_broadcast",
    ],
    [
      "DELETE",
      "/api/v1/admin/assignments/:id",
      "admin",
      "撤销老师的某个班级权限（软删除，is_active=false）",
    ],
    [
      "POST",
      "/api/v1/admin/assignments/batch",
      "admin",
      "批量为一个老师分配多个班级（如张老师一次分配10个班）",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("3.6  家长-学生绑定管理"));
children.push(
  routeTable([
    ["GET", "/api/v1/admin/parent-bindings", "admin", "查看家长绑定关系"],
    [
      "POST",
      "/api/v1/admin/parent-bindings",
      "admin",
      "建立家长-学生绑定（写 parent_child_relations）",
    ],
    [
      "DELETE",
      "/api/v1/admin/parent-bindings/:id",
      "admin",
      "解除家长-学生绑定",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("3.7  数据概览（管理员仪表盘）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/admin/dashboard",
      "admin",
      "全局数据概览：学校数/班级数/学生数/今日行为数/本周进化数",
    ],
    [
      "GET",
      "/api/v1/admin/audit-logs",
      "admin",
      "操作审计日志（管理员操作记录，合规要求）",
    ],
  ])
);

children.push(
  warn(
    "升年级接口（/classes/:id/promote）是高风险写操作，必须：1) 使用数据库事务；2) 防止重复升年级（检查 school_year 是否已存在）；3) 返回操作影响的学生数量。"
  )
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第四章：老师端（园长）模块
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第四章  老师端模块（Teacher）", C.green));
children.push(pEmpty(40));
children.push(
  note(
    '老师登录后，一次性获取"我的班级列表"，前端通过下拉/Tab 切换班级，所有打分/查询接口都携带 class_id 参数，后端验证老师对该班级是否有权限（查 teacher_class_assignments），禁止越权操作。'
  )
);
children.push(pEmpty(40));

children.push(h3("4.1  我的班级列表（核心：一师多班）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/teacher/my-classes",
      "teacher",
      "获取当前老师被授权的所有班级（含权限类型：can_score/can_broadcast）。前端用此接口渲染班级下拉列表，支持快速切换。",
    ],
    [
      "GET",
      "/api/v1/teacher/classes/:classId/overview",
      "teacher",
      "获取某班级概览：学生总数/本周行为数/平均成长值分布/伙伴进化概况（不含排名）",
    ],
    [
      "GET",
      "/api/v1/teacher/classes/:classId/students",
      "teacher",
      "获取某班级学生列表（含每人当前伙伴阶段/成长值区间，不含绝对数值排名）",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("4.2  正向行为打分（核心功能）"));
children.push(
  routeTable([
    [
      "POST",
      "/api/v1/teacher/behaviors",
      "teacher",
      "【核心】为学生添加正向行为记录：child_id/class_id/dimension/description/growth_value，自动触发成长值计算+伙伴进化检测+WebSocket推送给学生",
    ],
    [
      "GET",
      "/api/v1/teacher/behaviors",
      "teacher",
      "查看班级行为记录列表（按班级/学年/维度/时间筛选，分页）",
    ],
    ["GET", "/api/v1/teacher/behaviors/:id", "teacher", "查看单条行为记录详情"],
    [
      "DELETE",
      "/api/v1/teacher/behaviors/:id",
      "teacher",
      "撤销一条行为记录（软删除，仅限24小时内，同步扣减成长值）",
    ],
    [
      "POST",
      "/api/v1/teacher/behaviors/batch",
      "teacher",
      "批量打分：同一维度/描述，一次为多个学生添加行为记录",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("4.3  广播发送（园长广播）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/teacher/broadcasts",
      "teacher",
      "查看我发送的广播列表（含已发/定时待发）",
    ],
    [
      "POST",
      "/api/v1/teacher/broadcasts",
      "teacher",
      "发送广播（立即发送或定时发送）：target_class_id/content/type/scheduled_at",
    ],
    [
      "DELETE",
      "/api/v1/teacher/broadcasts/:id",
      "teacher",
      "取消一条定时广播（仅未发送的可取消）",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("4.4  集体挑战管理"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/teacher/challenges",
      "teacher",
      "查看班级当前进行中的集体挑战",
    ],
    [
      "POST",
      "/api/v1/teacher/challenges",
      "teacher",
      '创建集体挑战（如"连续3天全班作业全交，所有伙伴获成长礼包"）：条件配置+奖励值',
    ],
    [
      "PATCH",
      "/api/v1/teacher/challenges/:id/complete",
      "teacher",
      "手动标记挑战完成，触发批量成长值发放+广播通知",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("4.5  题库管理（班级专属）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/teacher/questions",
      "teacher",
      "查看班级题库（含系统公共题库+本班专属题目）",
    ],
    ["POST", "/api/v1/teacher/questions", "teacher", "添加班级专属题目"],
    ["PUT", "/api/v1/teacher/questions/:id", "teacher", "编辑题目"],
    [
      "DELETE",
      "/api/v1/teacher/questions/:id",
      "teacher",
      "删除题目（软删除）",
    ],
    [
      "POST",
      "/api/v1/teacher/questions/batch-import",
      "teacher",
      "批量导入题目（CSV格式）",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("4.6  盲盒奖励池管理"));
children.push(
  routeTable([
    ["GET", "/api/v1/teacher/blindbox/pool", "teacher", "查看本班盲盒奖励池"],
    [
      "POST",
      "/api/v1/teacher/blindbox/pool",
      "teacher",
      "向奖励池添加奖励（type/title/description/stock）",
    ],
    ["PUT", "/api/v1/teacher/blindbox/pool/:id", "teacher", "编辑奖励配置"],
    [
      "DELETE",
      "/api/v1/teacher/blindbox/pool/:id",
      "teacher",
      "下架某个奖励（软删除）",
    ],
    [
      "POST",
      "/api/v1/teacher/blindbox/draw/:childId",
      "teacher",
      "为某个学生触发抽盲盒（写入 blind_box_draws）",
    ],
    [
      "PATCH",
      "/api/v1/teacher/blindbox/draws/:drawId/redeem",
      "teacher",
      "确认兑换学生已获得的盲盒奖励",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("4.7  周报 PDF 生成"));
children.push(
  routeTable([
    [
      "POST",
      "/api/v1/teacher/reports/weekly",
      "teacher",
      "触发生成本班本周正能量周报 PDF（异步，chromedp渲染）",
    ],
    [
      "GET",
      "/api/v1/teacher/reports/weekly",
      "teacher",
      "查看历史周报列表（含下载链接）",
    ],
    [
      "GET",
      "/api/v1/teacher/reports/weekly/:id/download",
      "teacher",
      "下载指定周报 PDF 文件流",
    ],
  ])
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第五章：学生端模块
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第五章  学生端模块（Student）", C.teal));
children.push(pEmpty(40));
children.push(
  note(
    '学生端所有数据只返回"自己的"数据，任何接口都不暴露其他学生信息，后端必须严格校验 child_id 归属。'
  )
);
children.push(pEmpty(40));

children.push(h3("5.1  伙伴系统（核心）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/student/partner",
      "student",
      "获取当前活跃伙伴详情（template信息/成长值/进化阶段/昵称/动画资源URL）",
    ],
    [
      "GET",
      "/api/v1/student/partners",
      "student",
      "获取所有历史伙伴列表（含已毕业的，按 sequence_no 排序）",
    ],
    [
      "POST",
      "/api/v1/student/partner",
      "student",
      "选择新伙伴（首次选择，或满级后选下一只）：template_id+nickname。后端校验：是否有解锁权限（查 partner_unlock_logs）",
    ],
    [
      "PATCH",
      "/api/v1/student/partner/nickname",
      "student",
      "修改当前伙伴昵称",
    ],
    [
      "GET",
      "/api/v1/student/partner/growth-history",
      "student",
      "获取当前伙伴成长值流水（分页，含来源类型/来源描述/delta/时间）",
    ],
    [
      "GET",
      "/api/v1/student/partner/templates",
      "student",
      "获取可供选择的伙伴模板列表（含三种类型30个模板）",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("5.2  行为记录查看"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/student/behaviors",
      "student",
      "查看自己的行为记录（按学年/维度筛选，分页，含伙伴鼓励话）",
    ],
    [
      "GET",
      "/api/v1/student/behaviors/stats",
      "student",
      "行为统计：各维度累计次数/本月次数，用于前端雷达图展示",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("5.3  广播中心"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/student/broadcasts",
      "student",
      "获取收到的广播消息列表（伙伴广播+园长广播，分页）",
    ],
    [
      "PATCH",
      "/api/v1/student/broadcasts/:id/read",
      "student",
      "标记广播为已读",
    ],
    ["POST", "/api/v1/student/broadcasts/read-all", "student", "一键全部已读"],
  ])
);
children.push(pEmpty(40));

children.push(h3("5.4  成长年历"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/student/growth-calendar/months",
      "student",
      "获取全部月度成长卡列表（按学年分组）",
    ],
    [
      "GET",
      "/api/v1/student/growth-calendar/months/:month",
      "student",
      "获取指定月份的成长卡详情（伙伴快照/三件高光/曲线图数据/寄语）",
    ],
    [
      "GET",
      "/api/v1/student/growth-calendar/annual/:year",
      "student",
      "获取年度成长画卷数据（12张月卡汇总+总结语+里程碑时间线）",
    ],
    [
      "GET",
      "/api/v1/student/milestones",
      "student",
      "获取里程碑列表（勋章墙）",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("5.5  盲盒（学生侧查看）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/student/blindbox/my-draws",
      "student",
      "查看自己已抽到的盲盒奖励（含未兑换/已兑换/已过期）",
    ],
  ])
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第六章：家长端模块
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第六章  家长端模块（Parent）", C.purple));
children.push(pEmpty(40));
children.push(
  note(
    '家长端只能查看"已绑定的自己孩子"的数据，不能查看其他孩子。后端每个接口都要验证：请求的 child_id 必须在 parent_child_relations 中与当前登录家长绑定。'
  )
);
children.push(pEmpty(40));

children.push(
  routeTable([
    [
      "GET",
      "/api/v1/parent/children",
      "parent",
      "获取自己绑定的孩子列表（一个家长可绑多个孩子：兄弟姐妹）",
    ],
    [
      "GET",
      "/api/v1/parent/children/:childId/partner",
      "parent",
      "查看孩子当前伙伴状态（形态/成长值/阶段/昵称）",
    ],
    [
      "GET",
      "/api/v1/parent/children/:childId/partners",
      "parent",
      "查看孩子的历史伙伴列表",
    ],
    [
      "GET",
      "/api/v1/parent/children/:childId/behaviors",
      "parent",
      "查看孩子的正向行为记录（时间线，按学年/维度筛选）",
    ],
    [
      "GET",
      "/api/v1/parent/children/:childId/broadcasts",
      "parent",
      "查看孩子收到的伙伴鼓励广播（家长可收藏）",
    ],
    [
      "GET",
      "/api/v1/parent/children/:childId/milestones",
      "parent",
      "查看孩子的里程碑贴纸",
    ],
    [
      "GET",
      "/api/v1/parent/children/:childId/monthly-card",
      "parent",
      "查看孩子本月/历史月度成长卡",
    ],
    [
      "GET",
      "/api/v1/parent/children/:childId/annual-report",
      "parent",
      "查看孩子年度成长画卷（含PDF下载链接）",
    ],
    [
      "GET",
      "/api/v1/parent/children/:childId/battles",
      "parent",
      '查看孩子的对战参与记录（仅展示"参与了N次"，不显示胜负）',
    ],
  ])
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第七章：知识对战模块（Battle）
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第七章  知识对战模块（Battle）", C.orange));
children.push(pEmpty(40));
children.push(
  note(
    "对战核心依赖 WebSocket 实时通信。HTTP 接口负责房间创建/查询，WebSocket 负责实时题目同步/倒计时/结果推送。对战结束后异步发放双方成长值（正向，不论胜负）。"
  )
);
children.push(pEmpty(40));

children.push(h3("7.1  房间管理（HTTP 接口）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/battle/subjects",
      "student",
      "获取可用对战科目列表（含该班级题库数量）",
    ],
    [
      "POST",
      "/api/v1/battle/rooms",
      "student",
      "创建对战房间：subject/mode/class_id，返回 room_code（6位邀请码）",
    ],
    [
      "POST",
      "/api/v1/battle/rooms/:roomCode/join",
      "student",
      "通过邀请码加入房间（触发双方 WebSocket 通知：对手已加入，准备开始）",
    ],
    [
      "GET",
      "/api/v1/battle/rooms/:roomCode",
      "student",
      "查询房间状态（等待中/对战中/已结束）",
    ],
    [
      "GET",
      "/api/v1/battle/history",
      "student",
      "查看自己的对战历史（仅展示场次和参与时间，不暴露胜负排名）",
    ],
    [
      "GET",
      "/api/v1/battle/history/:roomId/review",
      "student",
      "对战复盘：查看自己的答题明细（每道题答对/答错/耗时），不显示对手成绩",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("7.2  WebSocket 实时通信"));
children.push(
  routeTable([
    [
      "WS",
      "/api/v1/battle/ws?room_code=XXX",
      "student",
      "对战实时通道：连接后进入房间，接收题目推送/倒计时/对手状态/结果通知",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("7.3  WebSocket 消息协议（Battle 专用）"));
children.push(
  new Table({
    width: { size: 9360, type: WidthType.DXA },
    columnWidths: [2400, 1400, 5560],
    rows: [
      new TableRow({
        tableHeader: true,
        children: ["消息类型", "方向", "数据结构/说明"].map(
          (t, i) =>
            new TableCell({
              borders: allBorders(C.orange),
              width: { size: [2400, 1400, 5560][i], type: WidthType.DXA },
              shading: { fill: C.orange, type: ShadingType.CLEAR },
              margins: { top: 80, bottom: 80, left: 120, right: 120 },
              children: [
                new Paragraph({
                  alignment: AlignmentType.CENTER,
                  children: [
                    new TextRun({
                      text: t,
                      font: "Arial",
                      size: 18,
                      bold: true,
                      color: C.white,
                    }),
                  ],
                }),
              ],
            })
        ),
      }),
      ...[
        ["battle:ready", "服务→客户端", "对手加入，准备开始倒计时（3/2/1）"],
        [
          "battle:question",
          "服务→客户端",
          "推送题目：{question_id, content, options, time_limit_sec, question_order}，答案不含在内",
        ],
        [
          "battle:answer",
          "客户端→服务",
          "提交答案：{question_id, answer_given}，服务端校验并记录",
        ],
        [
          "battle:result",
          "服务→客户端",
          "单题结果：{is_correct, correct_answer, time_used_ms}",
        ],
        [
          "battle:progress",
          "服务→客户端",
          '对手已完成本题（不透露答案），给出"等待对手"提示',
        ],
        [
          "battle:finish",
          "服务→客户端",
          "对战结束：{my_score, growth_gained, honor_badge（胜方有）}，不含对手分数",
        ],
        [
          "battle:ping",
          "客户端→服务",
          "心跳包，服务端响应 pong，30s无心跳断开连接",
        ],
        ["battle:error", "服务→客户端", "异常通知：超时/对手断线等"],
      ].map(
        (r, i) =>
          new TableRow({
            children: r.map(
              (t, j) =>
                new TableCell({
                  borders: allBorders(C.border),
                  width: { size: [2400, 1400, 5560][j], type: WidthType.DXA },
                  shading: {
                    fill: i % 2 === 1 ? C.lightOrange : C.white,
                    type: ShadingType.CLEAR,
                  },
                  margins: { top: 80, bottom: 80, left: 120, right: 120 },
                  children: [
                    new Paragraph({
                      children: [
                        new TextRun({
                          text: t,
                          font: j === 0 ? "Consolas" : "Arial",
                          size: j === 0 ? 16 : 18,
                          color: j === 0 ? C.orange : "1A1A1A",
                        }),
                      ],
                    }),
                  ],
                })
            ),
          })
      ),
    ],
  })
);

children.push(pEmpty(40));
children.push(
  warn(
    "对战结束后的成长值发放必须是幂等操作：如果网络断线重连导致重复触发，要通过 battle_participants 表的 growth_gained 字段检查是否已发放，防止重复发分。"
  )
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第八章：广播系统模块（WebSocket 核心）
// ═══════════════════════════════════════════════════════════════════
children.push(
  h1("第八章  广播与实时推送模块（Broadcast / WebSocket）", C.deepBlue)
);
children.push(pEmpty(40));
children.push(
  note(
    "广播系统基于 Redis Pub/Sub + WebSocket Hub 实现。学生客户端保持一个长连接，接收所有实时推送（行为打分通知/伙伴进化/园长广播/里程碑）。同伴互评消息纯内存广播，不落库。"
  )
);
children.push(pEmpty(40));

children.push(h3("8.1  WebSocket 主连接"));
children.push(
  routeTable([
    [
      "WS",
      "/api/v1/ws",
      "JWT",
      "学生/家长主 WebSocket 连接。连接建立后，服务端将该连接注册到 Hub（按 user_id 索引）。所有推送（打分/进化/广播）都通过此连接送达。",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("8.2  WebSocket Hub 设计要点"));
children.push(
  modTable(
    [
      [
        "连接管理",
        "Hub 维护 map[uint64]*Client（user_id → 连接）。一个用户同时只有一条活跃连接（多端登录时踢旧连接）。",
      ],
      [
        "心跳机制",
        "服务端每30秒发送 ping，客户端必须在10秒内回 pong，否则标记连接断开并从 Hub 移除。",
      ],
      [
        "断线重连",
        "客户端断线后，未读广播由数据库保存，重连后自动拉取 is_pushed=false 的记录补推。",
      ],
      [
        "Redis订阅",
        "每个班级对应一个 Redis Channel（broadcast:class:{class_id}），园长发广播写入 Redis，所有该班级在线连接收到推送。",
      ],
      [
        "同伴互评",
        "纯内存广播（不写数据库）：赠送者→接收者的 WebSocket 直接推送，接收者不在线则消息丢弃，保护隐私。",
      ],
      ["消息协议", "JSON格式：{type, data, timestamp}，type 枚举见8.3节"],
    ],
    C.deepBlue
  )
);
children.push(pEmpty(40));

children.push(h3("8.3  推送消息类型（服务端→客户端）"));
children.push(
  new Table({
    width: { size: 9360, type: WidthType.DXA },
    columnWidths: [2600, 6760],
    rows: [
      new TableRow({
        tableHeader: true,
        children: ["消息 type", "触发时机 & 数据内容"].map(
          (t, i) =>
            new TableCell({
              borders: allBorders(C.deepBlue),
              width: { size: [2600, 6760][i], type: WidthType.DXA },
              shading: { fill: C.deepBlue, type: ShadingType.CLEAR },
              margins: { top: 80, bottom: 80, left: 120, right: 120 },
              children: [
                new Paragraph({
                  alignment: AlignmentType.CENTER,
                  children: [
                    new TextRun({
                      text: t,
                      font: "Arial",
                      size: 18,
                      bold: true,
                      color: C.white,
                    }),
                  ],
                }),
              ],
            })
        ),
      }),
      ...[
        [
          "notify:behavior",
          "老师为学生打分后立即推送 → 含 dimension/growth_value/partner_message（伙伴鼓励话）",
        ],
        [
          "notify:evolution",
          "伙伴进化时推送 → 含 from_stage/to_stage/evolution_message/新形象资源URL",
        ],
        [
          "notify:milestone",
          "触发里程碑时推送 → 含 milestone_type/title/content",
        ],
        [
          "broadcast:class",
          "园长发班级广播时推送给全班在线学生 → 含 content/sender_name",
        ],
        [
          "broadcast:partner",
          "系统定时发送伙伴对主人的鼓励广播 → 含 content（早安/晚安/自定义）",
        ],
        ["notify:blindbox", "学生抽到盲盒后推送 → 含 reward_title/reward_type"],
        [
          "notify:peer_gift",
          "同伴互评（纯内存，不落库）→ 含 from_display_name/tag_label（美好时刻标签）",
        ],
        [
          "notify:challenge",
          "班级集体挑战完成时全班推送 → 含 challenge_title/reward_message",
        ],
        ["system:ping", "服务端心跳，客户端收到后回 pong，维持连接活跃"],
      ].map(
        (r, i) =>
          new TableRow({
            children: r.map(
              (t, j) =>
                new TableCell({
                  borders: allBorders(C.border),
                  width: { size: [2600, 6760][j], type: WidthType.DXA },
                  shading: {
                    fill: i % 2 === 1 ? C.lightBlue : C.white,
                    type: ShadingType.CLEAR,
                  },
                  margins: { top: 80, bottom: 80, left: 120, right: 120 },
                  children: [
                    new Paragraph({
                      children: [
                        new TextRun({
                          text: t,
                          font: j === 0 ? "Consolas" : "Arial",
                          size: j === 0 ? 16 : 18,
                          color: j === 0 ? C.deepBlue : "1A1A1A",
                        }),
                      ],
                    }),
                  ],
                })
            ),
          })
      ),
    ],
  })
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第九章：公开接口 & 伙伴模板
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第九章  公开接口 & 伙伴模板模块", C.gray));
children.push(pEmpty(40));

children.push(h3("9.1  公开接口（无需登录）"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/partner-templates",
      "公开",
      "获取全部激活的伙伴模板列表（30个），含三种类型、名称、描述、各阶段资源路径",
    ],
    ["GET", "/api/v1/partner-templates/:id", "公开", "获取单个模板详情"],
    [
      "GET",
      "/health",
      "公开",
      "健康检查，Docker healthcheck 使用，返回 {status:ok}",
    ],
    [
      "GET",
      "/api/v1/config/client",
      "公开",
      "返回前端所需的全局配置（成长值阈值/维度列表/七色列表），避免硬编码到客户端",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("9.2  伙伴模板管理（管理员）"));
children.push(
  routeTable([
    [
      "POST",
      "/api/v1/admin/partner-templates",
      "admin",
      "新增伙伴模板（初始化30个模板用）",
    ],
    [
      "PUT",
      "/api/v1/admin/partner-templates/:id",
      "admin",
      "更新模板信息（资源URL/鼓励语/启停）",
    ],
    [
      "POST",
      "/api/v1/admin/partner-templates/seed",
      "admin",
      "一键初始化30个预设模板（只能执行一次，幂等保护）",
    ],
  ])
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第十章：阳光章系统（二期接口预留）
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第十章  阳光章系统（二期，接口预留）", C.orange));
children.push(pEmpty(40));
children.push(
  warn(
    "二期实现，但路由和接口签名现在就要定义好，前端预留对接位置。表结构已在第二轮 v2.0 设计完毕。"
  )
);
children.push(pEmpty(40));

children.push(h3("10.1  七色配置管理"));
children.push(
  routeTable([
    [
      "GET",
      "/api/v1/admin/sunshine/colors",
      "admin",
      "【二期】获取学校七色-科目配置",
    ],
    [
      "POST",
      "/api/v1/admin/sunshine/colors",
      "admin",
      "【二期】配置七色-科目映射",
    ],
    [
      "PUT",
      "/api/v1/admin/sunshine/colors/:id",
      "admin",
      "【二期】更新颜色配置",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("10.2  盖章操作"));
children.push(
  routeTable([
    [
      "POST",
      "/api/v1/teacher/sunshine/stamp",
      "teacher",
      "【二期】老师为学生盖章：child_id/color_id/reason，自动记录 stamp_month/stamp_quarter",
    ],
    [
      "GET",
      "/api/v1/teacher/sunshine/stamps",
      "teacher",
      "【二期】查看班级盖章记录（按月/颜色筛选）",
    ],
    [
      "GET",
      "/api/v1/student/sunshine/my-stamps",
      "student",
      "【二期】学生查看自己的七色盖章概况（各色数量/当月情况）",
    ],
  ])
);
children.push(pEmpty(40));

children.push(h3("10.3  之星评选"));
children.push(
  routeTable([
    [
      "POST",
      "/api/v1/teacher/sunshine/awards/evaluate",
      "teacher",
      "【二期】触发月度/季度/年度评选（系统根据盖章数自动评出各色之星）",
    ],
    [
      "GET",
      "/api/v1/teacher/sunshine/awards",
      "teacher",
      "【二期】查看评选结果列表",
    ],
    [
      "GET",
      "/api/v1/student/sunshine/my-awards",
      "student",
      "【二期】学生查看自己获得的之星称号",
    ],
    [
      "GET",
      "/api/v1/parent/children/:childId/sunshine",
      "parent",
      "【二期】家长查看孩子阳光章情况",
    ],
  ])
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第十一章：核心业务流程 & 关键逻辑
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第十一章  核心业务流程 & 关键逻辑规范", C.red));
children.push(pEmpty(40));

children.push(h3("11.1  打分 → 成长值 → 进化 完整流程（最重要）"));
children.push(
  modTable(
    [
      [
        "Step 1",
        "老师调用 POST /teacher/behaviors，参数校验（维度/成长值范围）",
      ],
      [
        "Step 2",
        "检查老师对该班级是否有 can_score 权限（查 teacher_class_assignments）",
      ],
      [
        "Step 3",
        "开启数据库事务：写 behavior_records + 更新 partners.growth_points + 写 growth_records",
      ],
      [
        "Step 4",
        "进化检测：if new_points >= 阶段阈值 AND new_stage > current_stage，then 更新 partners.current_stage，写进化快照到 growth_records.is_evolution_trigger=true",
      ],
      [
        "Step 5",
        "满级检测：if current_stage=high AND growth_points>=300，then partners.status=graduated + 写 partner_unlock_logs（解锁选下一只权限）",
      ],
      [
        "Step 6",
        "里程碑检测：异步检查是否触发里程碑（goroutine），不阻塞主流程",
      ],
      ["Step 7", "事务提交"],
      [
        "Step 8",
        "异步推送 WebSocket（goroutine）：notify:behavior + notify:evolution（如有进化）+ notify:milestone（如有里程碑）",
      ],
      [
        "Step 9",
        "行为记录写 behavior_records.partner_message（从模板库取对应维度的鼓励话）",
      ],
    ],
    C.red
  )
);
children.push(pEmpty(40));

children.push(h3("11.2  多伙伴选择流程"));
children.push(
  modTable(
    [
      [
        "前置检查",
        "调用 POST /student/partner 时：查 partner_unlock_logs，确认有可用的未使用解锁次数（is_new_partner_selected=false）",
      ],
      [
        "创建伙伴",
        "计算 sequence_no = 历史伙伴数+1，创建新 partners 记录（status=active）",
      ],
      [
        "更新解锁日志",
        "将对应 partner_unlock_logs 记录的 is_new_partner_selected=true，new_partner_id=新伙伴ID",
      ],
      [
        "唯一性保证",
        "局部唯一索引：UNIQUE(child_id) WHERE status='active'，数据库层保证同一时刻只有一只活跃伙伴",
      ],
    ],
    C.teal
  )
);
children.push(pEmpty(40));

children.push(h3("11.3  老师多班权限校验（每个接口必须执行）"));
children.push(
  p(
    "每个 teacher 端接口收到 class_id 参数后，必须执行以下校验（可封装为 Service 层公共方法）：",
    { bold: true }
  )
);
children.push(pNum("从 JWT Claims 取 teacher_user_id"));
children.push(
  pNum(
    "查 teacher_class_assignments WHERE teacher_user_id=? AND class_id=? AND school_year=当前学年 AND is_active=true"
  )
);
children.push(
  pNum("检查对应权限位（打分需 can_score=true，广播需 can_broadcast=true）")
);
children.push(pNum("未通过则返回 403 FORBIDDEN，不执行任何业务逻辑"));
children.push(pEmpty(40));

children.push(h3("11.4  数据加密规范"));
children.push(
  modTable(
    [
      [
        "加密字段",
        "users.real_name_enc / users.phone_enc / children.real_name_enc / children.student_no_enc",
      ],
      ["算法", "AES-256-GCM，密钥从 APP_AES_KEY 环境变量读取（必须32字节）"],
      ["存储格式", "base64(iv + ciphertext + tag)，存入 VARCHAR(256) 字段"],
      [
        "解密时机",
        "仅后台管理员接口/审计接口需要解密明文，所有对外接口只返回脱敏展示值",
      ],
      [
        "禁止行为",
        "绝对禁止将真实姓名/学号拼入日志、错误信息、响应体的任何位置",
      ],
    ],
    C.purple
  )
);
children.push(pEmpty(40));

children.push(h3("11.5  定时任务（Cron）清单"));
children.push(
  new Table({
    width: { size: 9360, type: WidthType.DXA },
    columnWidths: [2000, 1600, 5760],
    rows: [
      new TableRow({
        tableHeader: true,
        children: ["任务名", "执行周期", "功能说明"].map(
          (t, i) =>
            new TableCell({
              borders: allBorders(C.navy),
              width: { size: [2000, 1600, 5760][i], type: WidthType.DXA },
              shading: { fill: C.navy, type: ShadingType.CLEAR },
              margins: { top: 80, bottom: 80, left: 120, right: 120 },
              children: [
                new Paragraph({
                  alignment: AlignmentType.CENTER,
                  children: [
                    new TextRun({
                      text: t,
                      font: "Arial",
                      size: 18,
                      bold: true,
                      color: C.white,
                    }),
                  ],
                }),
              ],
            })
        ),
      }),
      ...[
        [
          "定时广播发送",
          "每分钟",
          "查询 broadcast_messages WHERE is_scheduled=true AND is_sent=false AND scheduled_at<=NOW()，执行发送",
        ],
        ["早安广播", "每天7:30", "触发系统早安广播，向所有在线学生推送"],
        ["晚安广播", "每天21:00", "触发系统晚安广播"],
        [
          "月度成长卡生成",
          "每月1日0:00",
          "为所有学生生成上月月度成长卡（异步逐个生成，写 monthly_growth_cards）",
        ],
        [
          "连续打卡检测",
          "每天23:55",
          "检查今日是否有行为记录，无则 consecutive_days 重置为0",
        ],
        [
          "年度画卷生成",
          "每年12月25日",
          "触发生成全年学生的年度成长画卷（学年末）",
        ],
        ["Token黑名单清理", "每小时", "清理 Redis 中已过期的 Token 黑名单 Key"],
        [
          "离线消息补推",
          "每5分钟",
          "检查 behavior_records WHERE is_pushed=false AND created_at<5分钟前，对在线用户补推，超时丢弃记录日志",
        ],
      ].map(
        (r, i) =>
          new TableRow({
            children: r.map(
              (t, j) =>
                new TableCell({
                  borders: allBorders(C.border),
                  width: { size: [2000, 1600, 5760][j], type: WidthType.DXA },
                  shading: {
                    fill: i % 2 === 1 ? C.rowAlt : C.white,
                    type: ShadingType.CLEAR,
                  },
                  margins: { top: 80, bottom: 80, left: 120, right: 120 },
                  children: [
                    new Paragraph({
                      children: [
                        new TextRun({
                          text: t,
                          font: j === 0 ? "Consolas" : "Arial",
                          size: j === 0 ? 16 : 18,
                          color: "1A1A1A",
                        }),
                      ],
                    }),
                  ],
                })
            ),
          })
      ),
    ],
  })
);

children.push(pEmpty(80));
children.push(new Paragraph({ children: [new PageBreak()] }));

// ═══════════════════════════════════════════════════════════════════
// 第十二章：架构师检查要点
// ═══════════════════════════════════════════════════════════════════
children.push(h1("第十二章  架构师 Code Review 检查清单", C.red));
children.push(pEmpty(40));
children.push(critical("以下每一条都是硬性要求，代码提交前必须全部通过。"));
children.push(pEmpty(40));

children.push(h3("12.1  安全性"));
children.push(
  pBullet("所有 /admin/* 路由必须同时校验 JWT + role=admin，不能仅依赖路由分组")
);
children.push(
  pBullet(
    "所有 /teacher/* 路由必须验证 teacher_user_id 对 class_id 的权限，防越权"
  )
);
children.push(
  pBullet("所有 /parent/* 路由必须验证 child_id 属于该 parent，防越权")
);
children.push(
  pBullet("所有 /student/* 路由必须验证操作的资源属于该 student，防越权")
);
children.push(
  pBullet(
    "敏感字段（real_name/phone/student_no）绝对不出现在 JSON 响应和日志中"
  )
);
children.push(
  pBullet(
    "SQL 查询必须使用参数化查询，禁止字符串拼接（GORM 默认安全，不要绕过）"
  )
);
children.push(
  pBullet("文件上传（CSV/图片）必须校验文件类型和大小，防止任意文件上传")
);

children.push(h3("12.2  数据一致性"));
children.push(pBullet("打分→成长值→进化 三步必须在一个数据库事务内完成"));
children.push(
  pBullet(
    "成长值只增不减，任何减少成长值的操作必须有专项审批流（当前版本无此功能）"
  )
);
children.push(
  pBullet(
    "撤销行为记录（24h内）必须同步扣回成长值，并写 growth_records（delta为负）"
  )
);
children.push(
  pBullet(
    "升年级操作必须是幂等的（重复执行结果相同），通过 school_year 唯一索引保证"
  )
);
children.push(
  pBullet(
    "盲盒抽取必须扣减 stock，并使用数据库行锁防止超卖（SELECT FOR UPDATE）"
  )
);

children.push(h3("12.3  性能"));
children.push(
  pBullet("所有列表接口必须分页（默认20条，最大100条），禁止全表扫描返回")
);
children.push(
  pBullet(
    "班级概览接口必须使用 Redis 缓存（TTL=5分钟），避免每次请求都聚合查询"
  )
);
children.push(
  pBullet("WebSocket 消息推送必须是异步 goroutine，不阻塞 HTTP 响应")
);
children.push(
  pBullet("成长值统计（total_growth_points等）使用冗余字段，不能每次实时 SUM")
);
children.push(
  pBullet(
    "月度成长卡生成是 chromedp 重操作，必须走异步队列（Redis List），有并发限制"
  )
);

children.push(h3("12.4  代码规范"));
children.push(pBullet("Handler 层禁止出现任何 SQL/Redis 操作，只能调 Service"));
children.push(pBullet("Repository 层禁止出现任何业务判断逻辑，只做 CRUD"));
children.push(pBullet("Service 层的每个公共方法都必须有 context.Context 参数"));
children.push(
  pBullet("所有 error 必须被处理（不允许 _ = err），用 fmt.Errorf 包装上下文")
);
children.push(
  pBullet(
    "所有对外 JSON 字段用 snake_case，Go 结构体字段用 CamelCase + json tag"
  )
);
children.push(
  pBullet(
    "时间字段统一用 TIMESTAMPTZ（带时区），Go 侧用 time.Time，前端统一用 ISO 8601"
  )
);

children.push(pEmpty(80));

// ─── 结尾 ─────────────────────────────────────────────────────
children.push(
  new Paragraph({
    shading: { fill: C.navy, type: ShadingType.CLEAR },
    spacing: { before: 200, after: 200 },
    indent: { left: 200 },
    children: [
      new TextRun({
        text: "架构师签发  ·  成长伙伴后端功能清单 v1.0  ·  供 AI 代码生成使用",
        font: "Arial",
        size: 20,
        color: C.white,
      }),
    ],
  })
);

// ═══════════════════════════════════════════════════════════════════
// 生成文档
// ═══════════════════════════════════════════════════════════════════
const doc = new Document({
  styles: {
    default: { document: { run: { font: "Arial", size: 20 } } },
    paragraphStyles: [
      {
        id: "Heading1",
        name: "Heading 1",
        basedOn: "Normal",
        next: "Normal",
        quickFormat: true,
        run: { size: 38, bold: true, font: "Arial", color: C.white },
        paragraph: { spacing: { before: 360, after: 200 }, outlineLevel: 0 },
      },
      {
        id: "Heading2",
        name: "Heading 2",
        basedOn: "Normal",
        next: "Normal",
        quickFormat: true,
        run: { size: 28, bold: true, font: "Arial", color: C.white },
        paragraph: { spacing: { before: 280, after: 160 }, outlineLevel: 1 },
      },
      {
        id: "Heading3",
        name: "Heading 3",
        basedOn: "Normal",
        next: "Normal",
        quickFormat: true,
        run: { size: 24, bold: true, font: "Arial", color: C.navy },
        paragraph: { spacing: { before: 220, after: 100 }, outlineLevel: 2 },
      },
    ],
  },
  numbering: {
    config: [
      {
        reference: "bullets",
        levels: [
          {
            level: 0,
            format: LevelFormat.BULLET,
            text: "•",
            alignment: AlignmentType.LEFT,
            style: { paragraph: { indent: { left: 560, hanging: 280 } } },
          },
        ],
      },
      {
        reference: "numbers",
        levels: [
          {
            level: 0,
            format: LevelFormat.DECIMAL,
            text: "%1.",
            alignment: AlignmentType.LEFT,
            style: { paragraph: { indent: { left: 560, hanging: 280 } } },
          },
        ],
      },
    ],
  },
  sections: [
    {
      properties: {
        page: {
          size: { width: 11906, height: 16838 },
          margin: { top: 900, right: 900, bottom: 900, left: 900 },
        },
      },
      children,
    },
  ],
});

Packer.toBuffer(doc)
  .then((buf) => {
    fs.writeFileSync(
      "/mnt/user-data/outputs/成长伙伴_后端API功能清单v1.0.docx",
      buf
    );
    console.log("✅ 后端功能清单文档生成成功！");
  })
  .catch((err) => {
    console.error("❌ 生成失败:", err.message);
    process.exit(1);
  });
