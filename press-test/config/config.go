package config

const (

	//企业微信webhook地址
	WEBHOOK_URL = ""
	//'WEBHOOK_URL':'',

	//ssh-user
	User = ""
	Key  = "/Users/koko/.ssh/id_rsa"

	KubeconfigPATH = "/Users/koko/.ssh/188-kubeconfig"

	//cassandra-template path
	TemplatePath = "/Users/koko/go/src/BOC/press-test/schedule_script/cassandra-pv.yaml"

	K8sNamespace = "yuga1"

	Client = "gaohuitest"

	//# k8s namespace cassandra名称
	K8sNamespaceCassandra = "boc"

	//# fp base url
	FpBaseUrl = ""

	//# cassandra或yugabytedb pod数量, 每个表示一组压测
	DB_POD_RANGE = "1"

	FP_POD_RANGE = "1"

	//# fp pod初始数量, 首轮压测使用
	FpPodDefault = 1

	//# 随着cassandra或yugabytedb pod数量的增加，fp pod增加个数
	FpPodAdd = 1

	//# connect初始数
	//# 'CONNECT_NUM_ADD': 30,
	ConnectNumAdd = 4

	//#feature
	//# 'feature': [20, 40, 60, 80, 100],
	Feature = "1"

	//#数据维度
	DataRange = "7"

	Host = "10.1.9.188, 10.1.9.189"

	//# cassandra-data path
	CassandraDataPath  = "\\/home\\/datavisor\\/gaohui\\/boc1\\/cassandra\\/"
	CassandraData1Path = "/home/datavisor/gaohui/boc1/cassandra/"

	//#fp最终数量
	TotalFPNum = 1

	//#cassandra最终数量
	TotalCassandraNum = 1

	//# featureID List
	FeatureList = "1"
	//# best: 8t48c
	//# connect增加数
	ConnectNumDefault = 8

	//# wrk 压测脚本路径
	WrkScript = "/Users/koko/go/src/BOC/press-test/wrk_script/test.sh"
	Logpath   = "/Users/koko/go/src/press-test/logs"
	//# 压测结果保存位置
	CsvFilePath      = ""
	CsvFileName      = "wrk_result.csv"
	WrkRawlogPath    = ""
	WrkRawlog        = "wrk_rawlog"
	SchedulerLogPath = ""
	WrkLogPath       = ""

	DEPLOY          = "deploy"
	FP              = "fp"
	STS             = "sts"
	YB              = "yb-tserver"
	KAFKA           = "kafka"
	CASSANDRA       = "cassandra"
	ACTIONDELETE    = "delete"
	ACTIONCREATE    = "create"
	CASSANDRAPREFIX = "cassandra-data"
	YBPREFIX        = "yb-"
	SSHPASSWD       = ""
	SSHUSER         = ""
	SSHPORT         = "22"
	IsUpdate        = true
)
