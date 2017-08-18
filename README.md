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
```

test console log
```
2017/08/18 13:36:40 4 get from db
2017/08/18 13:36:41 4 get from done set cache
2017/08/18 13:36:41 4 << done &{1 weisd}
2017/08/18 13:36:41 2 get from cache
2017/08/18 13:36:41 2 << done &{1 weisd}
2017/08/18 13:36:41 3 get from cache
2017/08/18 13:36:41 3 << done &{1 weisd}
2017/08/18 13:36:41 9 get from cache
2017/08/18 13:36:41 9 << done &{1 weisd}
2017/08/18 13:36:41 1 get from cache
2017/08/18 13:36:41 1 << done &{1 weisd}
2017/08/18 13:36:41 6 get from cache
2017/08/18 13:36:41 6 << done &{1 weisd}
2017/08/18 13:36:41 8 get from cache
2017/08/18 13:36:41 7 get from cache
2017/08/18 13:36:41 8 << done &{1 weisd}
2017/08/18 13:36:41 7 << done &{1 weisd}
2017/08/18 13:36:41 0 get from cache
2017/08/18 13:36:41 0 << done &{1 weisd}
2017/08/18 13:36:41 5 get from cache
2017/08/18 13:36:41 5 << done &{1 weisd}
```