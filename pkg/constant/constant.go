package constant

const OPERATE_YOP_TEST_USER = 77 // 操作特殊账号

const OPERATE_PARTITION_TARGET_GMV = 78 // 操作分区的目标gmv

const (
	IN_VALID_ORG_ID = 41
)

const (
	SnapshotTaskTime = 1
	TaskTime         = 5
	SnapshotDeleted  = 0
	SnapshotWaitExec = 1
	Processing       = 2
	SnapshotFailed   = 3
	SnapshotSuccess  = 4
)

const (
	ManualInput       = 1
	UseSnapshot       = 2
	ImportUserList    = 3
	UserListSeparator = ","
)

// 材料预留，需要忽略的藏品id
// 1019327是正式环境pass卡id，1019894是测试环境pass卡id
var IGNORE_RESERVE_COLLECTION = []uint64{1019327, 1019894}

const REDIS_IP_WHITELIST_KEY = "new_backend:ip_whitelist"

const (
	COMPENSATE_STATUS_CLOSE = 0 // 未开启补偿状态
	COMPENSATE_STATUS_OPEN  = 1 // 开启补偿状态
)

const REDIS_HD_PARTITION_TARGET_GMV_SWITCH_KEY = "hotbox:hd_partition_target_gmv"
