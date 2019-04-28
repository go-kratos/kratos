package model

//缓存key的常量
const (
	CacheKeyOrderList  = "order_ls"  //订单列表缓存
	CacheKeyOrderMn    = "order_mn"  //order_main缓存
	CacheKeyOrderCnt   = "order_cnt" //订单数缓存
	CacheKeyOrderDt    = "order_dt"  //order_detail缓存
	CacheKeyOrderSKU   = "order_sku" //order_sku缓存
	CacheKeyOrderPayCh = "order_ch"  //order_pay_charge缓存

	CacheKeyStock   = "stock:%d"    // 库存数 redis key 前缀
	CacheKeyStockL  = "locked:%d"   // 锁定库存数 redisKey 前缀
	CacheKeySku     = "sku:%d"      // skuId => sku redis key 前缀
	CacheKeyItemSku = "sku.item:%d" // itemId => sku redis key 前缀

	// 票相关 key
	CacheKeyScreenSales      = "ticket:screen.sales"      // hash {sid:cnt}  各场次总的出票数
	CacheKeyScreenDailySales = "ticket:screen.daily"      // hash {sid:cnt}  各场次的当日销量
	CacheKeyUserBuyScreen    = "ticket:user.screen:%d"    // hash {sid:cnt}  用户购买各场次票数量
	CacheKeyOrderTickets     = "ticket:order.tks:%d"      // 一笔订单下所有电子票信息
	CacheKeyScreenTickets    = "ticket:screens.tks:%d:%d" // 一个场次下用户电子票信息
	CacheKeyTicketQr         = "ticket:qr.tk:%s"          // 一个二维码对应电子票信息
	CacheKeyTicket           = "ticket:tk:%d"             // 单张电子票信息
	CacheKeyTicketPool       = "ticket:pool:%d"           // sku票池
	CacheKeyTicketSend       = "ticket:send:%d"           // 票的赠送信息 send_tid => ticket_send
	CacheKeyTicketRecv       = "ticket:recv:%d"           // 票的赠送信息 recv_tid => ticket_send

	RedisExpireStock    = 120  // 库存量缓存过期时间
	RedisExpireStockTmp = 2    // 库存量缓存过期时间
	RedisExpireSku      = 1800 // SKU信息缓存过期时间
	RedisExpireSkuTmp   = 2    // SKU信息缓存过期时间

	RedisExpireTenMin    = 600   // 过期时间 10 分钟
	RedisExpireTenMinTmp = 2     // 过期时间 10 分钟
	RedisExpireOneDay    = 86400 // 过期时间 1 天
	RedisExpireOneDayTmp = 2     // 过期时间 1 天 兼容版

	CacheKeyPromo          = "%d:promotion:sales"        //拼团活动缓存
	CacheKeyPromoGroup     = "%d:promotion:group:sales"  //团缓存
	CacheKeyPromoOrder     = "%d:promotion:order:sales"  //拼团订单缓存
	CacheKeyPromoOrders    = "%d:promotion:orders:sales" //拼团团订单缓存
	RedisExpirePromo       = 1                           // 过期时间5分钟 to do
	RedisExpirePromoGroup  = 1                           // 过期时间5分钟 to do
	RedisExpirePromoOrder  = 1                           // 过期时间5分钟 to do
	RedisExpirePromoOrders = 1                           // 过期时间5分钟 to do
)
