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
2017/08/16 17:33:40 0 start>>>>>
2017/08/16 17:33:40 1 start>>>>>
2017/08/16 17:33:40 0 get lock 2017-08-16 17:33:40.316834982 +0800 CST
2017/08/16 17:33:40 0 get from db 2017-08-16 17:33:40.316875691 +0800 CST
2017/08/16 17:33:40 2 start>>>>>
2017/08/16 17:33:40 3 start>>>>>
2017/08/16 17:33:40 4 start>>>>>
2017/08/16 17:33:40 5 start>>>>>
2017/08/16 17:33:40 6 start>>>>>
2017/08/16 17:33:40 7 start>>>>>
2017/08/16 17:33:40 8 start>>>>>
2017/08/16 17:33:40 9 start>>>>>
2017/08/16 17:33:40 1 get data again 2017-08-16 17:33:40.370094687 +0800 CST
2017/08/16 17:33:40 1 miss goto GETLOCK 2017-08-16 17:33:40.37307862 +0800 CST
2017/08/16 17:33:40 2 get data again 2017-08-16 17:33:40.373432433 +0800 CST
2017/08/16 17:33:40 2 miss goto GETLOCK 2017-08-16 17:33:40.375455519 +0800 CST
2017/08/16 17:33:40 0 get from db done set cache 2017-08-16 17:33:40.376908568 +0800 CST
2017/08/16 17:33:40 3 get data again 2017-08-16 17:33:40.379526761 +0800 CST
2017/08/16 17:33:40 0 ====done==== &{1 weisd} 2017-08-16 17:33:40.380615076 +0800 CST
2017/08/16 17:33:40 3 get from cache 2017-08-16 17:33:40.383589909 +0800 CST
2017/08/16 17:33:40 3 ====done==== &{1 weisd} 2017-08-16 17:33:40.383603495 +0800 CST
2017/08/16 17:33:40 4 get data again 2017-08-16 17:33:40.385465321 +0800 CST
2017/08/16 17:33:40 4 get from cache 2017-08-16 17:33:40.387349697 +0800 CST
2017/08/16 17:33:40 4 ====done==== &{1 weisd} 2017-08-16 17:33:40.387364105 +0800 CST
2017/08/16 17:33:40 5 get data again 2017-08-16 17:33:40.392656672 +0800 CST
2017/08/16 17:33:40 5 get from cache 2017-08-16 17:33:40.394686621 +0800 CST
2017/08/16 17:33:40 5 ====done==== &{1 weisd} 2017-08-16 17:33:40.394701635 +0800 CST
2017/08/16 17:33:40 6 get data again 2017-08-16 17:33:40.396078692 +0800 CST
2017/08/16 17:33:40 6 get from cache 2017-08-16 17:33:40.399355424 +0800 CST
2017/08/16 17:33:40 6 ====done==== &{1 weisd} 2017-08-16 17:33:40.399373761 +0800 CST
2017/08/16 17:33:40 7 get data again 2017-08-16 17:33:40.402670814 +0800 CST
2017/08/16 17:33:40 7 get from cache 2017-08-16 17:33:40.405372319 +0800 CST
2017/08/16 17:33:40 7 ====done==== &{1 weisd} 2017-08-16 17:33:40.405389775 +0800 CST
2017/08/16 17:33:40 8 get data again 2017-08-16 17:33:40.407656199 +0800 CST
2017/08/16 17:33:40 8 get from cache 2017-08-16 17:33:40.410130508 +0800 CST
2017/08/16 17:33:40 8 ====done==== &{1 weisd} 2017-08-16 17:33:40.410144082 +0800 CST
2017/08/16 17:33:40 9 get data again 2017-08-16 17:33:40.412209662 +0800 CST
2017/08/16 17:33:40 9 get from cache 2017-08-16 17:33:40.414075883 +0800 CST
2017/08/16 17:33:40 9 ====done==== &{1 weisd} 2017-08-16 17:33:40.414102745 +0800 CST
2017/08/16 17:33:40 1 get data again 2017-08-16 17:33:40.426284965 +0800 CST
2017/08/16 17:33:40 2 get data again 2017-08-16 17:33:40.426966117 +0800 CST
2017/08/16 17:33:40 1 get from cache 2017-08-16 17:33:40.448540359 +0800 CST
2017/08/16 17:33:40 1 ====done==== &{1 weisd} 2017-08-16 17:33:40.448563874 +0800 CST
2017/08/16 17:33:40 2 get from cache 2017-08-16 17:33:40.448539244 +0800 CST
2017/08/16 17:33:40 2 ====done==== &{1 weisd} 2017-08-16 17:33:40.448580754 +0800 CST
```