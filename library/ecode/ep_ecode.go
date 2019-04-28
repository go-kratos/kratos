package ecode

// ep ecode interval is [0,990000]
var (
	//merlin/paas
	MerlinInvalidClusterErr = New(80001) // 集群不合法
	MerlinPaasRequestErr    = New(80002) // Paas 请求错误

	//merlin/tree
	MerlinGetUserTreeFailed        = New(80010) //获取 tree 节点失败
	MerlinTreeResponseErr          = New(80011) //请求 Tree 失败
	MerlinShouldTreeFullPath       = New(80012) //节点不合法
	MerlinTreeRequestErr           = New(80013) //服务树请求错误
	MerlinLoseTreeContainerNodeErr = New(80014) //当前节点下不存在 dev/container 子节点
	MerlinUserNoAccessTreeNode     = New(80034) //用户没有该服务树节点权限

	//merlin
	MerlinDuplicateMachineNameErr = New(81001) // 机器名称重复
	MerlinInvalidMachineAmountErr = New(81002) //机器数量不合法
	MerlinInvalidNodeAmountErr    = New(81003) //挂载节点数必须大于0且不大于10
	MerlinUpdateNodeErr           = New(81004) //更新节点失败
	MerlinCanNotBeDel             = New(81040) //创建状态机器无法删除

	//other
	MerlinIllegalPageNumErr  = New(89001) //分页页码不合法
	MerlinIllegalPageSizeErr = New(89002) //分页大小不合法

	MerlinDelayMachineErr                     = New(89010) //机器自主延期失败
	MerlinApplyMachineErr                     = New(89011) //机器申请延期失败
	MerlinCancelMachineErr                    = New(89012) //机器取消延期失败
	MerlinAuditMachineErr                     = New(89013) //机器审核延期失败
	MerlinApplyMachineByApplyEndTimeMore3MErr = New(89014) //机器申请延期失败

	//hubbili
	MerlinHubRequestErr           = New(89015) //请求bilihub失败
	MerlinHubNoRight              = New(89016) //没有权限执行
	MerlinImagePullErr            = New(89017) //下载镜像失败
	MerlinImagePushErr            = New(89018) //上传镜像失败
	MerlinImageTagErr             = New(89019) //Tag镜像失败
	MerlinSnapshotInDoingErr      = New(89024) //快照进行中
	MerlinNoHubAccount            = New(89026) //该用户没有Hub账号
	MerlinDuplicateImageNameErr   = New(89028) //镜像名称重复
	MerlinMachine2ImageInDoingErr = New(89029) //机器转镜像进行中
	MerlinMachineImageNotSameErr  = New(89030) //镜像名称不一致

	MerlinDeviceNotBind              = New(89020) //设备未绑定
	MerlinDeviceFarmErr              = New(89021) //DeviceFarm Error
	MerlinDeviceFarmMachineStatusErr = New(89025) //Merlin Device Farm Machine StatusErr
	MerlinDeviceIsNotRealMachineErr  = New(89031) //该操作只支持真机
	MerlinDeviceIsLendOut            = New(89032) //真机需在机架上才能操作
	MerlinDeviceNoRight              = New(89033) //无权限操作

	// user
	MerlinUserNotExist = New(89034) //用户不存在

	MartheBuglyErr       = New(89022) //MartheBuglyErr
	MartheTapdErr        = New(89023) //Tapd请求错误
	MartheTapdResDataErr = New(89050) //Tapd返回数据错误
	MartheTaskInRunning  = New(89047) //有任务正在执行

	MartheNoProjectInfo     = New(89045)
	MartheBugTaskInRunning  = New(89041)
	MartheDuplicateErr      = New(81044) // 名称重复
	MartheNoCookie          = New(89042)
	MartheCookieExpired     = New(89043) //cookie过期
	MartheFilterSqlError    = New(89046) //过滤sql error.
	MartheTimeConflictError = New(89048) //时间冲突.

	//melloi/PaaS
	MelloiPaasRequestErr        = New(60002) // Paas 请求错误
	MeilloiIllegalPageNumErr    = New(60004) //分页页码不合法
	MeilloillegalPageSizeErr    = New(60005) //分页大小不合法
	MelloiTreeRequestErr        = New(60001) //Tree 请求错误
	MelloiAdminExist            = New(60003) //管理员存在
	MelloiUpdateUserErr         = New(60006) //更新用户权限
	MelloiApplyRequestErr       = New(60008) //申请请求错误
	MelloiLabelRelationNotExist = New(60009) //标签关系存在
	MelloiLabelCountErr         = New(60010) //Label 数量超过2
	MelloiLabelExistErr         = New(60011) //Label数量存在
	MelloiRunNotInTime          = New(60012) //非压测时间段
	MelloiJmeterGenerateErr     = New(60013) //Jmeter脚本生成失败
	MelloiProtoFileNotUploaded  = New(60014) // proto文件没有上传
	MelloiProtocError           = New(60015) // protoc 编译失败
	MelloiProtoJavaPluginError  = New(60016) // protoc java插件编译失败
	MelloiJavacCompileError     = New(60017) // javac 编译失败
	MelloiJarError              = New(60018) // Jar打包失败
	MelloiUrlParseError         = New(60019) // URL 解析错误
	MelloiBeyondFileSize        = New(60020) // 文件太大
	MelloiCopyFileErr           = New(60021) // 文件复制错误
)
