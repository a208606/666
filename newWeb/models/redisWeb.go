package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

func init() {

	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		beego.Info("链接redis失败")
		return
	}
	defer conn.Close()

	err = conn.Send("set", "c1", "hello word")
	err = conn.Flush()
	reply, err := conn.Receive()

	str1, err := redis.String(reply, err)
	fmt.Println("str", str1)

	reply1, err := conn.Do("get", "c1")
	str, err := redis.String(reply1, err)
	fmt.Println("str", str)

	reply2, err := conn.Do("mset", "t1", "56", "t2", "60", "t3", "qian")

	beego.Info("reply2", reply2)

	reply3, err := conn.Do("mget",  "t3","t2","t1" )

	result, err := redis.Values(reply3, err)

	var t1, t2 int
	var t3 string
	str2, err := redis.Scan(result, &t3, &t2, &t1)

	fmt.Println(str2)

	fmt.Println("result", t3, t1, t2)

}
