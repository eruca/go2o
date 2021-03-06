package goclient

import (
	"com/domain/interface/member"
	"com/ording/entity"
	"errors"
	"fmt"
	"ops/cf/net/jsv"
)

type memberClient struct {
	conn *jsv.TCPConn
}

func (this *memberClient) Login(usr, pwd string) (b bool, token string, msg string) {
	var result jsv.Result
	err := this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"usr":"%s","pwd":"%s"}>>Member.Login`,
		usr, pwd)), &result, 128)
	if err != nil {
		return false, "", err.Error()
	}

	jsv.Println(result.Data)
	if result.Result {
		return result.Result, result.Data.(string), result.Message
	} else {
		return result.Result, "", result.Message
		return result.Result, "", result.Message
	}
}

func (this *memberClient) Verify(memberId int, token string) (r jsv.Result) {
	var result jsv.Result
	err := this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s"}>>Member.Verify`,
		memberId, token)), &result, -1)
	if err != nil {
		r.Message = err.Error()
		r.Result = false
	}
	return result
}

func (this *memberClient) GetMember(memberId int, token string) (a *member.ValueMember, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s"}>>Member.GetMember`,
		memberId, token)), &result, 512)
	if err != nil {
		return nil, err
	}
	a = &member.ValueMember{}
	err = jsv.UnmarshalMap(result.Data, a)
	return a, err
}

func (this *memberClient) GetMemberAccount(memberId int, token string) (
	a *member.Account, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s"}>>Member.GetMemberAccount`,
		memberId, token)), &result, 256)
	if err != nil {
		return nil, err
	}
	a = &member.Account{}
	err = jsv.UnmarshalMap(result.Data, a)
	return a, err
}

func (this *memberClient) GetBankInfo(memberId int, token string) (
	a *member.BankInfo, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s"}>>Member.GetBankInfo`,
		memberId, token)), &result, 256)
	if err != nil {
		return nil, err
	}
	a = &member.BankInfo{}
	err = jsv.UnmarshalMap(result.Data, &a)
	return a, err
}

func (this *memberClient) GetBindPartner(memberId int, token string) (
	a *entity.Partner, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s"}>>Member.GetBindPartner`,
		memberId, token)), &result, 512)
	if err != nil {
		return nil, err
	}
	a = &entity.Partner{}
	err = jsv.UnmarshalMap(result.Data, &a)
	return a, err
}

func (this *memberClient) SaveMember(m *member.ValueMember, token string) (b bool, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s","json":%s}>>Member.SaveMember`,
		m.Id, token, jsv.MarshalString(m))), &result, -1)
	if err != nil {
		return false, err
	}
	return result.Result, err
}

func (this *memberClient) GetDeliverAddrs(memberId int, token string) (addrs []entity.DeliverAddress, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s"}>>Member.GetDeliverAddrs`,
		memberId, token)), &result, 256)
	if err != nil {
		return nil, err
	}
	addrs = []entity.DeliverAddress{}
	err = jsv.UnmarshalMap(result.Data, &addrs)
	return addrs, err
}

func (this *memberClient) GetDeliverAddrById(memberId int, token string, addr_id int) (
	e *entity.DeliverAddress, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s","addr_id":"%d"}>>Member.GetDeliverAddrById`,
		memberId, token, addr_id)), &result, 256)
	if err != nil {
		return nil, err
	}
	e = &entity.DeliverAddress{}
	err = jsv.UnmarshalMap(result.Data, &e)
	return e, err
}

func (this *memberClient) SaveDeliverAddr(memberId int, token string, e *entity.DeliverAddress) (b bool, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s","json":%s}>>Member.SaveDeliverAddr`,
		memberId, token, jsv.MarshalString(e))), &result, -1)
	if err != nil {
		return false, err
	}
	return result.Result, err
}

func (this *memberClient) DeleteDeliverAddr(memberId int, token string, addr_id int) (b bool, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"member_id":"%d","token":"%s","addr_id":"%d"}>>Member.DeleteDeliverAddr`,
		memberId, token, addr_id)), &result, -1)
	if err != nil {
		return false, err
	}
	if !result.Result {
		return false, errors.New(result.Data.(string))
	}
	return true, nil
}
