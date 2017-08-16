package main

import (
	"fmt"
	"log"
	"time"

	"sync"

	"github.com/gocommon/cache"
	"github.com/gocommon/cache/locker"
)

var c cache.Cacher

func main() {

	c = cache.NewCache(cache.UseLocker(true))

	id := int64(1)

	flushcache(id)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			info, err := getTestUserInfoFromCache(id, i)
			if err != nil {
				panic(err)
			}

			log.Println(i, "====done====", info, time.Now())

			wg.Done()
		}(i)

	}

	wg.Wait()

}

func flushcache(id int64) {

	c.Tags(getTestUserInfoTag(id)).Flush()

}

func getTestUserInfoFromCache(id int64, idx int) (*TestUser, error) {
	var err error

	key := getTestUserInfoKey(id)

	tags := []string{
		getTestUserAllTag(),
		getTestUserInfoTag(id),
	}

	// var info *TestUser

	var info *TestUser

	has, err := c.Tags(tags...).Get(key, &info)
	if err != nil {
		return nil, err
	}

	if has {
		log.Println(idx, "get from cache", time.Now())
		return info, nil
	}

	// if not exists go to get data from db

	l := c.Locker(key)

	// lock
GETLOCK:
	err = l.Lock()

	if locker.IsErrLockFailed(err) {
		// wait
		time.Sleep(50 * time.Millisecond)
		// get again
		log.Println(idx, "get data again", time.Now())

		has, err := c.Tags(tags...).Get(key, &info)
		if err != nil {
			return nil, err
		}

		if has {
			log.Println(idx, "get from cache", time.Now())
			return info, nil
		}

		log.Println(idx, "miss goto GETLOCK", time.Now())

		// if empty goto lock
		goto GETLOCK
	} else if err != nil {
		return nil, err
	}

	defer l.Unlock()

	// get lock

	log.Println(idx, "get lock", time.Now())
	log.Println(idx, "get from db", time.Now())

	info, err = getTestUserInfoFromDB(id)
	if err != nil {
		log.Println(idx, "getTestUserInfoFromDB err", time.Now())
		return nil, err
	}

	log.Println(idx, "get from db done set cache", time.Now())

	err = c.Tags(tags...).Set(key, info)
	if err != nil {
		log.Println(idx, "Set err", time.Now())
		return nil, err
	}

	return info, nil

}

func getTestUserInfoFromDB(id int64) (*TestUser, error) {
	time.Sleep(60 * time.Millisecond)
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
