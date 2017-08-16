# cache-example

example for https://github.com/gocommon/cache

```

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
```

test console log
```
2017/07/12 17:27:51 get lock
2017/07/12 17:27:51 get from db
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get from done set cache
2017/07/12 17:27:52 &{1 weisd}
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get data again
2017/07/12 17:27:52 get from cache
2017/07/12 17:27:52 &{1 weisd}
2017/07/12 17:27:52 get from cache
2017/07/12 17:27:52 &{1 weisd}
2017/07/12 17:27:52 get from cache
2017/07/12 17:27:52 &{1 weisd}
2017/07/12 17:27:52 get from cache
2017/07/12 17:27:52 &{1 weisd}
2017/07/12 17:27:52 get from cache
2017/07/12 17:27:52 &{1 weisd}
2017/07/12 17:27:52 get from cache
2017/07/12 17:27:52 &{1 weisd}
2017/07/12 17:27:52 get from cache
2017/07/12 17:27:52 &{1 weisd}
2017/07/12 17:27:52 get from cache
2017/07/12 17:27:52 get from cache
2017/07/12 17:27:52 &{1 weisd}
2017/07/12 17:27:52 &{1 weisd}
```