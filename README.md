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
2017/08/16 17:28:21 9 get lock 2017-08-16 17:28:21.027368439 +0800 CST
2017/08/16 17:28:21 9 get from db 2017-08-16 17:28:21.027554005 +0800 CST
2017/08/16 17:28:21 1 get data again 2017-08-16 17:28:21.090332791 +0800 CST
2017/08/16 17:28:21 8 get data again 2017-08-16 17:28:21.090371038 +0800 CST
2017/08/16 17:28:21 6 get data again 2017-08-16 17:28:21.090389003 +0800 CST
2017/08/16 17:28:21 5 get data again 2017-08-16 17:28:21.09033649 +0800 CST
2017/08/16 17:28:21 3 get data again 2017-08-16 17:28:21.090429781 +0800 CST
2017/08/16 17:28:21 9 get from db done set cache 2017-08-16 17:28:21.090437407 +0800 CST
2017/08/16 17:28:21 0 get data again 2017-08-16 17:28:21.090355625 +0800 CST
2017/08/16 17:28:21 2 get data again 2017-08-16 17:28:21.090365172 +0800 CST
2017/08/16 17:28:21 4 get data again 2017-08-16 17:28:21.090374836 +0800 CST
2017/08/16 17:28:21 7 get data again 2017-08-16 17:28:21.090429778 +0800 CST
2017/08/16 17:28:21 7 get from cache 2017-08-16 17:28:21.101164768 +0800 CST
2017/08/16 17:28:21 7 ====done==== &{1 weisd} 2017-08-16 17:28:21.101187775 +0800 CST
2017/08/16 17:28:21 2 get from cache 2017-08-16 17:28:21.101354708 +0800 CST
2017/08/16 17:28:21 2 ====done==== &{1 weisd} 2017-08-16 17:28:21.101362138 +0800 CST
2017/08/16 17:28:21 0 get from cache 2017-08-16 17:28:21.101442976 +0800 CST
2017/08/16 17:28:21 0 ====done==== &{1 weisd} 2017-08-16 17:28:21.101449058 +0800 CST
2017/08/16 17:28:21 1 get from cache 2017-08-16 17:28:21.101531593 +0800 CST
2017/08/16 17:28:21 1 ====done==== &{1 weisd} 2017-08-16 17:28:21.101538815 +0800 CST
2017/08/16 17:28:21 4 miss goto GETLOCK 2017-08-16 17:28:21.101794079 +0800 CST
2017/08/16 17:28:21 8 get from cache 2017-08-16 17:28:21.101950373 +0800 CST
2017/08/16 17:28:21 8 ====done==== &{1 weisd} 2017-08-16 17:28:21.101960922 +0800 CST
2017/08/16 17:28:21 6 get from cache 2017-08-16 17:28:21.10198521 +0800 CST
2017/08/16 17:28:21 6 ====done==== &{1 weisd} 2017-08-16 17:28:21.101993858 +0800 CST
2017/08/16 17:28:21 3 get from cache 2017-08-16 17:28:21.102049059 +0800 CST
2017/08/16 17:28:21 3 ====done==== &{1 weisd} 2017-08-16 17:28:21.102056277 +0800 CST
2017/08/16 17:28:21 5 get from cache 2017-08-16 17:28:21.102069568 +0800 CST
2017/08/16 17:28:21 5 ====done==== &{1 weisd} 2017-08-16 17:28:21.102076384 +0800 CST
2017/08/16 17:28:21 9 ====done==== &{1 weisd} 2017-08-16 17:28:21.103265864 +0800 CST
2017/08/16 17:28:21 4 get lock 2017-08-16 17:28:21.103668621 +0800 CST
2017/08/16 17:28:21 4 get from db 2017-08-16 17:28:21.103682676 +0800 CST
2017/08/16 17:28:21 4 get from db done set cache 2017-08-16 17:28:21.16779999 +0800 CST
2017/08/16 17:28:21 4 ====done==== &{1 weisd} 2017-08-16 17:28:21.172120355 +0800 CST
```