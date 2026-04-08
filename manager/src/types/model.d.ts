import { TableType } from '/@/types/global';

// 树形节点数据
export interface TreeNodeData {
	name: string;
	id: number;
	parentId: number;
	children?: TreeNodeData[];
}

export interface MenuDataTree {
	id: number;
	pid: number;
	title: string;
	children?: MenuDataTree[];
}

// 系统配置模型
export interface TableConfigColumn {
	configId: number;
	configName: string;
	configKey: string;
	configValue: string;
	configType: number;
	remark: string;
	createdAt: string;
}

export interface TableConfigModel extends TableType {
	tableData: Array<TableConfigColumn>;
	param: {
		pageNum: number;
		pageSize: number;
		configName: string;
		configKey: string;
		configType: string;
		dateRange: string[];
	};
}

// 部门模型
export interface TableDeptColumn {
	deptId: number;
	parentId: number;
	deptName: string;
	status: number;
	orderNum: number;
	createdAt: string;
	children?: TableDeptColumn[];
}

export interface TableDeptModel extends TableType {
	tableData: Array<TableDataRow>;
	param: {
		pageNum: number;
		pageSize: number;
		dateRange: [];
		deptName: string;
		status: string;
	};
}

export interface RuleFormDeptState {
	deptId: number;
	parentId: number;
	deptName: string;
	orderNum: number;
	leader: string;
	phone: string | number;
	email: string;
	status: number;
}

export interface DeptSate {
	isShowDialog: boolean;
	ruleForm: RuleFormDeptState;
	deptData: Array<TreeNodeData>;
	rules: object;
}

// 数据字典模型
export interface TableDictColumn {
	dictId: number;
	dictName: string;
	dictType: string;
	status: number;
	remark: string;
	createdAt: string;
}

export interface TableDictModel extends TableType{
	tableData: Array<TableDictColumn>;
	param: {
		pageNum: number;
		pageSize: number;
		dictName: string;
		dictType: string;
		status: string;
		dateRange: string[];
	};
}

export interface TableDictTypeColumn {
	dictCode: number;
	dictSort: number;
	dictLabel: string;
	dictValue: string;
	dictType: string;
	status: number;
	remark: string;
	createdAt: string;
}

export interface TableDictTypeModel extends TableType {
	tableData: Array<TableDictTypeColumn>;
	param: {
		pageNum: number;
		pageSize: number;
		dictType: string;
		dictLabel: string;
		status: string;
	};
}

export interface RuleFormDictState {
	dictCode: number;
	dictLabel: string;
	dictValue: string;
	dictColor: string;
	dictSort: number;
	isDefault: number;
	status: number;
	remark: string;
	dictType: string;
}

export interface DictState {
	isShowDialog: boolean;
	ruleForm: RuleFormDictState;
	rules: {};
}

export interface RuleFormDictStateEx {
	dictId: number;
	dictName: string;
	dictType: string;
	status: number;
	remark: string;
}

export interface DictStateEx {
	isShowDialog: boolean;
	ruleForm: RuleFormDictStateEx;
	rules: {};
}

// 岗位模型
export interface TablePostColumn {
	postId: number;
	postCode: string;
	postName: string;
	postSort: number;
	status: number;
	remark: string;
	createdAt: string;
}

export interface TablePostModel extends TableType {
	tableData: Array<TablePostColumn>;
	param: {
		postName: string;
		status: string;
		postCode: string;
		pageNum: number;
		pageSize: number;
	};
}

// 角色模型
export interface TableRoleColumn {
	id: number;
	status: number;
	listOrder: number;
	name: string;
	remark: string;
	dataScope: number;
	createdAt: string;
}

export interface TableRoleModel extends TableType {
	tableData: Array<TableRoleColumn>;
	param: {
		roleName: string;
		roleStatus: string;
		pageNum: number;
		pageSize: number;
	};
}

export interface RoleDataRow {
	id: number;
	name: string;
	status: number;
	listOrder: number;
	remark: string;
	menuIds: Array<number>;
}

export interface RoleState {
	loading: boolean;
	isShowDialog: boolean;
	formData: RoleDataRow;
	menuData: Array<MenuDataTree>;
	menuExpand: boolean;
	menuNodeAll: boolean;
	menuCheckStrictly: boolean;
	menuProps: {
		children: string;
		label: string;
	};
	rules: object;
}

// 用户模型
export interface TableUserModel extends TableType {
	deptProps: {};
	deptData: any[];
	tableData: any[];
	param: {
		pageNum: number;
		pageSize: number;
		deptId: string;
		mobile: string;
		status: string;
		keyWords: string;
		dateRange: string[];
	};
}

// 登录日志模型
export interface TableLoginLogsColumn {
	infoId: number;
	loginName: string;
	ipaddr: string;
	loginLocation: string;
	browser: string;
	os: string;
	status: number;
	msg: string;
	loginTime: string;
	module: string;
}

export interface TableLoginLogsModel extends TableType {
	tableData: Array<TableDataRow>;
	param: {
		pageNum: number;
		pageSize: number;
		dateRange: string[];
		status: string;
		ipaddr: string;
		loginLocation: string;
		userName: string;
	};
}

// 操作日志模型
export interface LinkedSysOperLogSysDept {
	deptId:number|undefined;        // 部门id
	deptName:string|undefined;      // 部门名称
}

export interface TableOperLogsColumn {
	operId:number;          // 日志编号
	title:string;           // 系统模块
	requestMethod:string;   // 请求方式
	operName:string;        // 操作人员
	deptName:string;        // 部门名称
	operUrl:string;         // 请求URL
	operIp:string;          // 主机地址
	operLocation:string;    // 操作地点
	operParam:string;       // 请求参数
	status:number;          // 操作状态（0正常 1异常）
	operTime:string;        // 操作时间
	linkedSysOperLogSysDept:LinkedSysOperLogSysDept;
}

export interface TableOperLogInfoData {
	operId:number|undefined;        // 日志编号
	title:string|undefined;         // 系统模块
	businessType:number|undefined;  // 操作类型
	method:string|undefined;        // 操作方法
	requestMethod:string|undefined; // 请求方式
	operatorType:number|undefined;  // 操作类别
	operName:string|undefined;      // 操作人员
	deptName:string|undefined;      // 部门名称
	operUrl:string|undefined;       // 请求URL
	operIp:string|undefined;        // 主机地址
	operLocation:string|undefined;  // 操作地点
	operParam:string|undefined;     // 请求参数
	jsonResult:string|undefined;    // 返回参数
	status:boolean;                 // 操作状态（0正常 1异常）
	errorMsg:string|undefined;      // 错误消息
	operTime:string|undefined;      // 操作时间
	linkedSysOperLogSysDept:LinkedSysOperLogSysDept;
}

export interface TableOperLogEditState{
	loading:boolean;
	isShowDialog: boolean;
	formData:TableOperLogInfoData;
	rules: object;
}

export interface TableOperLogsModel extends TableType {
	operIds:any[];
	tableData: Array<TableOperLogsColumn>;
	param: {
		pageNum: number;
		pageSize: number;
		title: string|undefined;
		requestMethod: string|undefined;
		operName: string|undefined;
		status: number|undefined;
		dateRange: string[];
	};
}

// 在线用户模型
export interface TableOnlineColumn {
	id: number;
	uuid: string;
	token: string;
	createTime: string;
	userName: string;
	ip: string;
	explorer: string;
	os: string;
}

export interface TableOnlinedModel extends TableType{
	tableData: Array<TableOnlineColumn>;
	param: {
		ipaddr: string;
		userName: string;
		pageNum: number;
		pageSize: number;
	};
}

