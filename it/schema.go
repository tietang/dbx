package it

import "fmt"

// 测试表结构：每次 init 时自动执行（CREATE TABLE IF NOT EXISTS，可重复运行）
var itTables = []string{
	// red_envelope_goods 表
	"CREATE TABLE IF NOT EXISTS `red_envelope_goods`" +
		"(" +
		"  `id`              bigint(20)     NOT NULL AUTO_INCREMENT COMMENT '自增ID'," +
		"  `envelope_no`     varchar(32)    NOT NULL COMMENT '红包编号,红包唯一标识 '," +
		"  `envelope_type`   tinyint(4)     NOT NULL COMMENT '红包类型：普通红包，碰运气红包,过期红包'," +
		"  `username`        varchar(64)             DEFAULT NULL COMMENT '用户名称'," +
		"  `user_id`         varchar(40)    NOT NULL COMMENT '用户编号, 红包所属用户 '," +
		"  `blessing`        varchar(64)             DEFAULT NULL COMMENT '祝福语'," +
		"  `amount`          decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '红包总金额'," +
		"  `amount_one`      decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '单个红包金额，碰运气红包无效'," +
		"  `quantity`        int(10)        NOT NULL COMMENT '红包总数量 '," +
		"  `remain_amount`   decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '红包剩余金额额'," +
		"  `remain_quantity` int(10)        NOT NULL COMMENT '红包剩余数量 '," +
		"  `expired_at`      datetime(3)    NOT NULL COMMENT '过期时间'," +
		"  `status`          tinyint(4)     NOT NULL COMMENT '红包/订单状态：0 创建、1 发布启用、2过期、3失效'," +
		"  `order_type`      tinyint(4)     NOT NULL COMMENT '订单类型：发布单、退款单 '," +
		"  `pay_status`      tinyint(4)     NOT NULL COMMENT '支付状态：未支付，支付中，已支付，支付失败 '," +
		"  `created_at`      datetime(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间'," +
		"  `updated_at`      datetime(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间'," +
		"  PRIMARY KEY (`id`) USING BTREE," +
		"  UNIQUE KEY `envelope_no_idx` (`envelope_no`) USING BTREE," +
		"  KEY `id_user_idx` (`user_id`) USING BTREE" +
		") ENGINE = InnoDB" +
		"  AUTO_INCREMENT = 101" +
		"  DEFAULT CHARSET = utf8" +
		"  ROW_FORMAT = DYNAMIC;",
	// account 表（Account 实体对应）
	"CREATE TABLE IF NOT EXISTS `account`" +
		"(" +
		"  `id`            bigint(20)     NOT NULL AUTO_INCREMENT COMMENT '账户ID'," +
		"  `account_no`    varchar(45)    NOT NULL COMMENT '账户编号,账户唯一标识'," +
		"  `account_name`  varchar(64)             DEFAULT NULL COMMENT '账户名称'," +
		"  `account_type`  tinyint(4)     NOT NULL DEFAULT '0' COMMENT '账户类型：积分账户、会员卡账户、钱包账户、红包账户'," +
		"  `currency_code` varchar(10)    NOT NULL DEFAULT 'CNY' COMMENT '货币类型编码：CNY人民币，EUR欧元，USD美元'," +
		"  `user_id`       varchar(40)    NOT NULL COMMENT '用户编号, 账户所属用户'," +
		"  `username`      varchar(64)             DEFAULT NULL COMMENT '用户名称'," +
		"  `balance`       decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '账户可用余额'," +
		"  `status`        tinyint(4)     NOT NULL DEFAULT '0' COMMENT '账户状态：0账户初始化，1启用，2停用'," +
		"  `created_at`    datetime(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间'," +
		"  `updated_at`    datetime(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间'," +
		"  PRIMARY KEY (`id`) USING BTREE," +
		"  UNIQUE KEY `account_no_idx` (`account_no`) USING BTREE" +
		") ENGINE = InnoDB" +
		"  DEFAULT CHARSET = utf8" +
		"  ROW_FORMAT = DYNAMIC;",
}

// createTables 自动创建全部测试表（IF NOT EXISTS，可重复执行）
func createTables() error {
	for _, ddl := range itTables {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("exec ddl failed: %v\n%s", err, ddl)
		}
	}
	return nil
}
