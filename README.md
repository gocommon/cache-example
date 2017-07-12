# cache-example

example for https://github.com/gocommon/cache

```

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

	err = c.Tags(tags).Get(key, info)

	// no err return
	if err == nil {
		log.Println("get from cache")
		return info, nil
	}

	// server err
	if !cache.IsErrNil(err) {
		log.Println("server err")
		return nil, err
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
		err = c.Tags(tags).Get(key, info)

		// no err return
		if err == nil {
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
```