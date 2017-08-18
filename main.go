package main

import (
	"fmt"
	"log"
	"time"

	"sync"

	"github.com/gocommon/cache"
	"github.com/gocommon/etcd3"
)

func main() {

	id := int64(1)

	flushcache(id)

	conf := map[string]etcd3.Config{
		"default": etcd3.Config{
			Endpoints: []string{"127.0.0.1:2379"},
		},
	}

	if err := etcd3.InitEtcdv3(conf); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			info, err := getTestUserInfoFromCache(id, i)
			if err != nil {
				panic(err)
			}

			log.Println(i, "<< done", info)

			wg.Done()
		}(i)

	}

	wg.Wait()

}

func flushcache(id int64) {
	c := cache.New()

	c.Tags(getTestUserInfoTag(id)).Flush()

}

func getTestUserInfoFromCache(id int64, idx int) (*TestUser, error) {
	var err error

	key := getTestUserInfoKey(id)

	tags := []string{
		getTestUserAllTag(),
		getTestUserInfoTag(id),
	}

	c := cache.New()

	// var info *TestUser

	info := &TestUser{}

	has, err := c.Tags(tags...).Get(key, info)
	if err != nil {
		return nil, err
	}

	if has {
		log.Println(idx, "get from cache")
		return info, nil
	}

	sess, err := etcd3.Session()
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	l := sess.NewLocker(getTestUserInfoTag(id))
	l.Lock()
	defer l.Unlock()

	has, err = c.Tags(tags...).Get(key, info)
	if err != nil {

		return nil, err
	}

	if has {
		log.Println(idx, "get from cache")
		return info, nil
	}

	// if not exists go to get data from db

	log.Println(idx, "get from db")

	info, err = getTestUserInfoFromDB(id)
	if err != nil {
		log.Println(idx, "getTestUserInfoFromDB err")
		return nil, err
	}

	log.Println(idx, "get from done set cache")

	err = c.Tags(tags...).Set(key, info)
	if err != nil {
		log.Println(idx, "Set err")
		return nil, err
	}

	return info, nil

}

func getTestUserInfoFromDB(id int64) (*TestUser, error) {
	time.Sleep(1 * time.Second)
	return &TestUser{1, "weisd"}, nil
}

type TestUser struct {
	ID   int64
	Name string
}

func getTestUserAllTag() string {
	return "getTestUserAllTag"
}

func getTestUserInfoTag(id int64) string {
	return fmt.Sprintf("getTestUserInfoTag:%d", id)
}

func getTestUserInfoKey(id int64) string {
	return fmt.Sprintf("testuserinfo:%d", id)
}
