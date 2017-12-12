package controllers

import (
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"strconv"
	"cart/models/redis"
)

type loginInfo struct {
	UserId string
}

type good struct {
	SkuId string
	Amount int
}

type cart struct {
	Goods []good
}

func AddCart(w http.ResponseWriter, r *http.Request) {
	var item good
	var c cart
	var logined bool
	var login loginInfo
	var exist bool
	var err error


	r.ParseForm()

	item.SkuId = r.Form.Get("skuid")
	if item.SkuId == "" {
		log.Fatalln("skuid is nil")
	}

	item.Amount, err = strconv.Atoi(r.Form.Get("amount"))
	if err != nil {
		log.Fatalln(err)
	}

	// 判断是否登录
	token := r.Form.Get("token")
	// for test
	token = "t"
	if token != "" {
		jToken := fmt.Sprintf("{\"token\":\"%s\"}", token)
		body := strings.NewReader(jToken)
		resp, err := http.Post("http://localhost:8019/getUserIdByToken", "text/plain", body)
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal(b, &login)
		if err != nil {
			log.Fatalln(err)
		}
		if login.UserId != "" {
			logined = true
		} else {
			logined = false
		}
	} else {
		logined = false
	}

	if logined {
		userCart := "cart_" + login.UserId
		conn := redis.Pool.Get()
		// 判断商品是否已存在
		_, err := redis.Int(conn.Do("HGET", userCart, item.SkuId))
		if err != nil {
			_, err = conn.Do("HSET", userCart, item.SkuId, item.Amount)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			// 商品存在则增加数量
			_, err = conn.Do("HINCRBY", userCart, item.SkuId, item.Amount)
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		cookie, _ := r.Cookie("cart")
		if cookie != nil {
			//获取购物车
			value := strings.Replace(cookie.Value, "{", "{\"", -1)
			value = strings.Replace(value, ":", "\":", -1)
			value = strings.Replace(value, ",A", "\",\"A", -1)
			value = strings.Replace(value, "d\":", "d\":\"", -1)
			err = json.Unmarshal([]byte(value), &c)
			if err != nil {
				log.Fatalln(err)
			}
			// 判断商品是否已存在，存在则增加数量
			for k, g := range c.Goods {
				if g.SkuId == item.SkuId {
					exist = true
					c.Goods[k].Amount = c.Goods[k].Amount + item.Amount
					break
				}
			}
		} else {
			cookie = new(http.Cookie)
			cookie.Name = "cart"
			cookie.Path = "/"
			cookie.MaxAge = 24*60*60
		}
		if !exist {
			c.Goods = append(c.Goods, item)
		}
		cB, err := json.Marshal(c)
		if err != nil {
			log.Fatalln(err)
		}
		cookie.Value = string(cB)
		http.SetCookie(w, cookie)
		fmt.Fprint(w, cookie.Value)
	}
}
