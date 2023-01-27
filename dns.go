package main

import ( //导入dns库
	"fmt"

	"github.com/miekg/dns"
)

type LastCnamer interface {
	LastCname(domain string) (lastCname string, err error)
}

type LastCnameResolver string

func (l LastCnameResolver) LastCname(domain string) (lastCname string, err error) {
	ns := string(l)
	//设置DNS类型
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(domain, dns.TypeA)
	r, _, err := c.Exchange(&m, ns) //发送DNS包
	if err != nil {
		return lastCname, err
	}
	// Last CNAME
	for _, ans := range r.Answer { //返回正确的信息
		cname, ok := ans.(*dns.CNAME)
		if ok {
			lastCname = cname.Target
		}
	}
	return lastCname, nil
}

type LastCnameStub map[string]string

func (m LastCnameStub) LastCname(domain string) (string, error) {
	return m[domain], nil
}

func main() {
	m := make(map[string]string)
	m["www.shiep.edu.cn."] = "www.shiep.edu.cn." //设置域名
	var l LastCnamer = LastCnameStub(m)
	cname, _ := l.LastCname("www.shiep.edu.cn.")
	fmt.Println("Stub", cname)

	l = LastCnameResolver("8.8.8.8:53") //模拟域名服务器地址
	cname, err := l.LastCname("www.shiep.edu.cn.")
	if err != nil {
		panic(err)
	}
	fmt.Println("Resolved", cname)
}
