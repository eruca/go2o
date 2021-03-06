package glob

import (
	"com/domain/interface/member"
	"com/domain/interface/partner"
	"com/domain/interface/promotion"
	"com/domain/interface/sale"
	"com/domain/interface/shopping"
	"com/ording/entity"
	"com/share/variable"
	"fmt"
	"ops/cf"
	"ops/cf/db"
	"ops/cf/log"
	"ops/cf/web"
)

func getDb(c *cf.Config, l log.ILogger) db.Connector {
	//数据库连接字符串
	//root@tcp(127.0.0.1:3306)/foodording?charset=utf8
	var connStr string
	driver := c.Get(variable.DbDriver)
	dbCharset := c.Get(variable.DbCharset)
	if dbCharset == "" {
		dbCharset = "utf8"
	}
	connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		c.Get(variable.DbUsr),
		c.Get(variable.DbPwd),
		c.Get(variable.DbServer),
		c.Get(variable.DbPort),
		c.Get(variable.DbName),
		dbCharset,
	)
	connector := db.NewCommonConnector(driver, connStr, l, 30)

	//table mapping
	orm := connector.GetOrm()

	orm.CreateTableMap(entity.Partner{}, "pt_partner")
	orm.CreateTableMap(entity.Shop{}, "pt_shop")
	//orm.CreateTableMap(entity.Order{}, "pt_order")
	//orm.CreateTableMap(entity.SiteConf{}, "pt_siteconf")
	//orm.CreateTableMap(entity.SaleConf{}, "pt_saleconf")
	orm.CreateTableMap(entity.Category{}, "it_category")
	orm.CreateTableMap(entity.FoodItem{}, "it_item")
	//orm.CreateTableMap(entity.Member{}, "mm_member")
	//orm.CreateTableMap(entity.MemberRelation{}, "mm_relation")
	//orm.CreateTableMap(entity.MemberAccount{}, "mm_account")
	//orm.CreateTableMap(entity.BankInfo{}, "mm.bank")
	//orm.CreateTableMap(entity.IncomeLog{}, "mm_income_log")
	//orm.CreateTableMap(entity.IntegralLog{}, "mm_integral_log")
	orm.CreateTableMap(entity.DeliverAddress{}, "mm_deliver_addr")

	/** new **/
	orm.CreateTableMap(member.ValueMember{}, "mm_member")
	orm.CreateTableMap(member.IncomeLog{}, "mm_income_log")
	orm.CreateTableMap(member.IntegralLog{}, "mm_integral_log")
	orm.CreateTableMap(member.Account{}, "mm_account")
	orm.CreateTableMap(member.MemberLevel{}, "conf_member_level")
	orm.CreateTableMap(member.DeliverAddress{}, "mm_deliver_addr")
	orm.CreateTableMap(member.MemberRelation{}, "mm_relation")
	orm.CreateTableMap(shopping.ValueOrder{}, "pt_order")
	orm.CreateTableMap(shopping.OrderCoupon{}, "pt_order_coupon")

	orm.CreateTableMap(sale.ValueProduct{}, "it_item")
	orm.CreateTableMap(partner.ValuePartner{}, "pt_partner")
	orm.CreateTableMap(partner.SiteConf{}, "pt_siteconf")
	orm.CreateTableMap(partner.SaleConf{}, "pt_saleconf")

	/** 促销 **/
	orm.CreateTableMap(promotion.ValueCoupon{}, "pm_coupon")
	orm.CreateTableMap(promotion.ValueCouponBind{}, "pm_coupon_bind")
	orm.CreateTableMap(promotion.ValueCouponTake{}, "pm_coupon_take")

	return connector
}

func initTemplate(cfg *cf.Config) *web.TemplateWrapper {
	return &web.TemplateWrapper{
		Init: func(m *map[string]interface{}) {
			(*m)["static_serv"] = cfg.Get(variable.StaticServer)
			(*m)["img_serv"] = cfg.Get(variable.ImageServer)
			(*m)["version"] = cfg.Get(variable.Version)
		},
	}
}

//	MasterToken: crypto.EncodeUsrPwd("master", "123456"),
//		deviceDir := filepath.Dir(options.DevicePatchDir)
//		fi,err := os.Stat(options.DevicePatchDir)
//
//		if fi == nil || err == os.ErrNotExist {
//			os.MkdirAll(deviceDir,os.ModePerm)
//		}else if !fi.IsDir() {
//			os.Remove(options.DevicePatchDir)
//			os.MkdirAll(deviceDir,os.ModePerm)
//}
//func initDbMap(db *sql.DB) *gorp.DbMap {
//	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
//	dbmap.AddTableWithName(entity.FoodItem{}, "it_item")
//	return dbmap
//}
