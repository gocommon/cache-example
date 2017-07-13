package main

import (
	"fmt"
	"log"
	"time"

	"sync"

	"github.com/gocommon/cache"
	"github.com/gocommon/cache/locker"
)

func main() {

	id := int64(1)

	flushcache(id)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			info, err := getTestUserInfoFromCache(id)
			if err != nil {
				panic(err)
			}

			log.Println(info)

			wg.Done()
		}()

	}

	wg.Wait()

}

func flushcache(id int64) {
	c := cache.NewCache()
	tags := []string{
		getTestUserInfoTag(id),
	}
	c.Flush(tags)

}

func getTestUserInfoFromCache(id int64) (*TestUser, error) {
	var err error

	key := getTestUserInfoKey(id)

	tags := []string{
		getTestUserAllTag(),
		getTestUserInfoTag(id),
	}

	c := cache.NewCache(cache.UseLocker(true))

	// var info *TestUser

	info := &TestUser{}

	has, err := c.Tags(tags).Get(key, info)
	if err != nil {
		return nil, err
	}

	if has {
		log.Println("get from cache")
		return info, nil
	}

	// if not exists go to get data from db

	l := c.NewLocker(key)

	// lock
GETLOCK:
	err = l.Lock()

	if locker.IsErrLockFailed(err) {
		// wait
		time.Sleep(500 * time.Millisecond)
		// get again
		log.Println("get data again")

		has, err := c.Tags(tags).Get(key, info)
		if err != nil {
			return nil, err
		}

		if has {
			log.Println("get from cache")
			return info, nil
		}

		// if empty goto lock
		goto GETLOCK
	} else if err != nil {
		return nil, err
	}

	defer l.Unlock()

	// get lock

	log.Println("get lock")
	log.Println("get from db")

	info, err = getTestUserInfoFromDB(id)
	if err != nil {
		log.Println("getTestUserInfoFromDB err")
		return nil, err
	}

	log.Println("get from done set cache")

	err = c.Tags(tags).Set(key, info)
	if err != nil {
		log.Println("Set err")
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
