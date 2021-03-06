package dao

import (
	"com/ording"
	"com/ording/entity"
	"com/share/variable"
	"database/sql"
	"errors"
	"fmt"
	"ops/cf/app"
	"ops/cf/db"
	"regexp"
	"time"
)

var (
	commHostRegexp *regexp.Regexp
)

type partnerDao struct {
	db.Connector
	app.Context
}

func (this *partnerDao) getHostRegexp() *regexp.Regexp {
	if commHostRegexp == nil {
		commHostRegexp = regexp.MustCompile("([^.]+)." +
			this.Context.Config().Get(variable.ServerDomain))
	}
	return commHostRegexp
}

//检查用户名和密码是否正确，如果正确则返回ID,否则返回-1
func (this *partnerDao) Verify(usr, pwd string) int {

	var id int
	id = -1
	epwd := ording.EncodePartnerPwd(usr, pwd)

	f := func(row *sql.Row) {
		row.Scan(&id)
	}

	this.Connector.QueryRow("SELECT id FROM pt_partner WHERE usr=? AND pwd=?", f, usr, epwd)

	return id
}

//根据ID获取Partner信息
func (this *partnerDao) GetPartnerById(id int) *entity.Partner {
	e := new(entity.Partner)
	if this.Context.Db().GetOrm().Get(e, id) != nil {
		return nil
	}
	return e
	//	pt := new(entity.Partner)
	//	f := func(row *sql.Row) {
	//		var expTime, joinTime, updateTime string
	//
	//		row.Scan(&pt.Id, &pt.User,
	//			&pt.Pwd, &expTime,
	//			&pt.Name, &pt.Tel,
	//			&pt.Address, &joinTime,
	//			&updateTime)
	//
	//		var err error
	//		//转换成时间
	//		pt.Expires, err = time.Parse("2006-01-02 15:04:05", expTime)
	//		pt.JoinTime, _ = time.Parse("2006-01-02 15:04:05", joinTime)
	//		pt.UpdateTime, _ = time.Parse("2006-01-02 15:04:05", updateTime)
	//
	//		if err != nil {
	//			fmt.Println(err.Error())
	//		}
	//	}
	//
	//	this.Connector.QueryRow("SELECT * FROM partners WHERE id=?", f, id)
	//
	//	return pt
}

//根据用户名获取合作商
func (this *partnerDao) GetPartnerByUsr(usr string) *entity.Partner {
	e := new(entity.Partner)
	if this.Context.Db().GetOrm().GetBy(e, "usr='"+usr+"'") != nil {
		return nil
	}
	return e
}

//保存合作商
func (this *partnerDao) SavePartner(partner *entity.Partner) (int, error) {
	if partner.Id <= 0 {
		//		//多行字符用``
		//		_, id, err := this.Connector.Exec(`
		//			INSERT INTO partners
		//			(
		//			user,
		//			pwd,
		//			expires,
		//			name,
		//			tel,
		//			address,
		//			jointime,
		//			updatetime)
		//			VALUES
		//			(?,?,?,?,?,?,?,?)`, partner.User, partner.Pwd, partner.Expires, partner.Name,
		//			partner.Tel, partner.Address, partner.JoinTime, partner.UpdateTime)
		//		return id, err

		orm := this.Connector.GetOrm()
		_, _, err := orm.Save(nil, partner)
		var newPtId int
		err = this.Connector.ExecScalar(`SELECT MAX(id) FROM pt_partner`, &newPtId)
		if err != nil {
			return 0, err
		}

		var conf entity.SiteConf = entity.SiteConf{
			PtId:       newPtId,
			Host:       "",
			Logo:       "nologo.gif",
			IndexTitle: partner.Name,
			SubTitle:   partner.Name,
			State:      1,
			StateHtml:  "",
		}
		orm.Save(nil, conf)

		return newPtId, err

	} else {
		//		_, _, err := this.Connector.Exec(`UPDATE partners
		//			SET
		//			user = ?,
		//			pwd = ?,
		//			expires = ?,
		//			name = ?,
		//			tel = ?,
		//			address = ?,
		//			jointime = ?,
		//			updatetime = ?
		//			WHERE id = ?
		//			`, partner.User, partner.Pwd, partner.Expires, partner.Name, partner.Tel,
		//			partner.Address, partner.JoinTime, partner.UpdateTime,
		//			partner.Id)
		//		return partner.Id, err

		_, _, err := this.Connector.GetOrm().Save(partner.Id, partner)
		return partner.Id, err

	}
}



//删除合作商户
func (this *partnerDao) DelPartner(partnerId int) error {
	return errors.New("not allow")
	return this.Connector.GetOrm().DeleteByPk(entity.Partner{}, partnerId)
}

//获取合作商域名
func (this *partnerDao) GetHostById(partnerId int) (host string) {
	this.Connector.ExecScalar(`SELECT host FROM pt_siteconf WHERE pt_id=? LIMIT 0,1`,
		&host, partnerId)
	if host == "" {
		var usr string
		this.Connector.ExecScalar(
			`SELECT usr FROM pt_partner WHERE id=?`,
			&usr, partnerId)
		host = fmt.Sprintf("%s.%s", usr, this.Context.Config().Get(variable.ServerDomain))
	}
	return host
}

//根据host获取合作商
func (this *partnerDao) GetPartnerByHost(host string) (partner *entity.Partner, err error) {
	//  $ 获取合作商ID
	// $ hostname : 域名
	// *.wdian.net  二级域名
	// www.dc1.com  顶级域名

	reg := this.getHostRegexp()
	if reg.MatchString(host) {
		matches := reg.FindAllStringSubmatch(host, 1)
		usr := matches[0][1]
		partner = this.GetPartnerByUsr(usr)
	} else {

		partner = new(entity.Partner)
		if this.Context.Db().GetOrm().GetByQuery(partner,
			`SELECT * FROM pt_partner INNER JOIN pt_siteconf
					 ON pt_siteconf.pt_id = pt_partner.id
					 WHERE host='`+host+`'`) != nil {
			partner = nil
		}
	}

	if partner == nil {
		return nil, errors.New("101: no such partner")
	}
	if time.Now().After(partner.Expires) {
		return nil, errors.New("103: partner is expires")
	}

	//清除必要的数据
	//partner.Pwd = ""
	//partner.Secret = ""

	return partner, nil
}

func (this *partnerDao) CheckHost(host string) error {
	//todo:
	return nil
}

/*


def getpage(partnerid,type):
    '获取页面内容'
    row = newdb().fetchone('SELECT `content` FROM pt_page WHERE `ptid`=%(ptid)s and `type`=%(type)s',
                            {
                             'ptid':partnerid,
                             'type':type
                            })

    return '' if row==None else row[0]


def updatepage(partnerid,type,content):
    '更新页面内容'
    data={
          'ptid':partnerid,
          'type':type,
          'content':content,
          'updatetime':utility.timestr()
          }

    row=newdb().query('UPDATE pt_page SET `content`=%(content)s,`updatetime`=%(updatetime)s WHERE `ptid`=%(ptid)s AND `type`=%(type)s',data)
    if row==0:
        newdb().query('INSERT INTO pt_page (`ptid`,`type`,`content`,`updatetime`) VALUES(%(ptid)s,%(type)s,%(content)s,%(updatetime)s)',data)

*/
