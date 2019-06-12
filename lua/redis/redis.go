
package redis

import (
	redis2 "xsuv/redis"
	"github.com/yuin/gopher-lua"
)

var api = map[string]lua.LGFunction{
	"open": openRedis,
}

var luaRedisMethods = map[string]lua.LGFunction{
	"is_exist": exists,
	"del": del,
	"type": type_,
	"keys": keys,
	"random_key": randomKey,
	"rename": rename,
	"db_size": dbSize,
	"expire": expire,
	"ttl": ttl,
	"move": move,
	"flush": flush,
	"set": set,
	"get": get,
	"getset": getset,
	"mget": mget,
	"setnx": setnx,
	"setex": setex,
	"mset": mset,
	"msetnx": msetnx,
	"incr": incr,
	"incrby": incrby,
	"decr": decr,
	"decrby": decrby,
	"append": append_,
	"hset": hset,
	"hget": hget,
	"hmget": hmget,
	"hincrby": hincrby,
	"hexists": hexists,
	"hdel": hdel,
	"hlen": hlen,
	"hkeys": hkeys,
	"hgetall": hgetall,
}

const (
	luaRedisTypeName = "redis"
)

func Preload(L *lua.LState) {
	registerRedisDBType(L)
	L.PreloadModule("redis", load)
}

func load(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

func RegisterRedisValue(L *lua.LState, name string, db *redis2.RedisDB)  {
	ud := L.NewUserData()
	ud.Value = db
	ud.Metatable = L.NewTypeMetatable(name)
	L.SetField(ud.Metatable, "__index", L.SetFuncs(L.NewTable(), luaRedisMethods))
	L.SetGlobal(name, ud)
}

func registerRedisDBType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaRedisTypeName)
	L.SetGlobal(luaRedisTypeName, mt)
	// static attributes
	//L.SetField(mt, "new", L.NewFunction(openRedis))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), luaRedisMethods))
}

func openRedis(L *lua.LState) int{
	connectStr := L.CheckString(1)
	db, err := redis2.NewRedisDB(connectStr);
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = db
	L.SetMetatable(ud, L.GetTypeMetatable(luaRedisTypeName))
	L.Push(ud)
	return 1
}

func checkRedis(L *lua.LState) *redis2.RedisDB {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*redis2.RedisDB); ok {
		return v
	}
	L.ArgError(1, "not redis type")
	return nil
}

func exists(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	ok, err := p.Exists(key)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(ok))
	return 1
}

func del(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	ok, err := p.Del(key)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(ok))
	return 1
}

func type_(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	s, err := p.Type(key)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(s))
	return 1
}

func randomKey(L *lua.LState) int{
	p := checkRedis(L)
	s, err := p.Randomkey()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(s))
	return 1
}

func rename(L *lua.LState) int{
	p := checkRedis(L)
	src := L.CheckString(2)
	dst := L.CheckString(3)
	err := p.Rename(src, dst)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func dbSize(L *lua.LState) int{
	p := checkRedis(L)
	size, err := p.Dbsize()
	if err != nil {
		L.Push(lua.LNumber(size))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(size))
	return 1
}

func expire(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	time := L.CheckInt64(3)
	ok, err := p.Expire(key, time)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(ok))
	return 1
}

func ttl(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	time, err := p.Ttl(key)
	if err != nil {
		L.Push(lua.LNumber(-1))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(time))
	return 1
}

func move(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	dbnum := L.CheckInt(3)
	ok, err := p.Move(key, dbnum)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(ok))
	return 1
}

func flush(L *lua.LState) int{
	p := checkRedis(L)
	all := L.CheckBool(2)
	err := p.Flush(all)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func set(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	value := L.CheckString(3)
	err := p.Set(key, []byte(value))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func get(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	s, err := p.Get(key)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(s))
	return 1
}

func getset(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	val := L.CheckString(3)
	s, err := p.Getset(key, []byte(val))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(s))
	return 1
}

func mget(L *lua.LState) int{
	p := checkRedis(L)
	keys := L.CheckTable(2)
	arg := make([]string, 0)
	max := keys.MaxN()
	for i := 1; i<=max; i++ {
		v := keys.RawGetInt(i)
		if s, ok := v.(lua.LString); ok {
			arg = append(arg, string(s))
		}
	}
	res, err := p.Mget(arg...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	t := new(lua.LTable)
	for _, r := range res {
		t.Append(lua.LString(string(r)))
	}
	L.Push(t)
	return 1
}

func setnx(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	val := L.CheckString(3)
	ok, err := p.Setnx(key, []byte(val))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(ok))
	return 1
}

func setex(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	time := L.CheckInt64(3)
	val := L.CheckString(4)
	err := p.Setex(key, time, []byte(val))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func mset(L *lua.LState) int{
	p := checkRedis(L)
	arg := L.CheckTable(2)
	mapping := make(map[string][]byte)
	arg.ForEach(func(key, val lua.LValue){
		k, ok := key.(lua.LString)
		v, ok1 := val.(lua.LString)
		if ok && ok1 {
			mapping[string(k)] = []byte(v)
		}
	})
	err := p.Mset(mapping)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func msetnx(L *lua.LState) int{
	p := checkRedis(L)
	arg := L.CheckTable(2)
	mapping := make(map[string][]byte)
	arg.ForEach(func(key, val lua.LValue){
		k, ok := key.(lua.LString)
		v, ok1 := val.(lua.LString)
		if ok && ok1 {
			mapping[string(k)] = []byte(v)
		}
	})
	ok, err := p.Msetnx(mapping)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(ok))
	return 1
}

func incr(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	val, err := p.Incr(key)
	if err != nil {
		L.Push(lua.LNumber(val))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(val))
	return 1
}

func incrby(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	val := L.CheckInt64(3)
	ret, err := p.Incrby(key, val)
	if err != nil {
		L.Push(lua.LNumber(ret))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(ret))
	return 1
}

func decr(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	val, err := p.Decr(key)
	if err != nil {
		L.Push(lua.LNumber(val))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(val))
	return 1
}

func decrby(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	val := L.CheckInt64(3)
	ret, err := p.Decrby(key, val)
	if err != nil {
		L.Push(lua.LNumber(ret))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(ret))
	return 1
}

func append_(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	val := L.CheckString(3)
	err := p.Append(key, []byte(val))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func hset(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	field := L.CheckString(3)
	val := L.CheckString(4)
	ok, err := p.Hset(key, field, []byte(val))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(ok))
	return 1
}

func hget(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	field := L.CheckString(3)
	val, err := p.Hget(key, field)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(val)))
	return 1
}

func hmget(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	t := L.CheckTable(3)
	arg := make([]string, 0)
	max := t.MaxN()
	for i := 1; i<=max; i++ {
		v := t.RawGetInt(i)
		if s, ok := v.(lua.LString); ok {
			arg = append(arg, string(s))
		}
	}
	res, err := p.Hmget(key, arg...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	t1 := new(lua.LTable)
	for _, r := range res {
		t1.Append(lua.LString(string(r)))
	}
	L.Push(t1)
	return 1
}

func hincrby(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	field := L.CheckString(3)
	val := L.CheckInt64(4)
	res, err := p.Hincrby(key, field, val)
	if err != nil {
		L.Push(lua.LNumber(res))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(res))
	return 1
}

func hexists(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	field := L.CheckString(3)
	ok, err := p.Hexists(key, field)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(ok))
	return 1
}

func hdel(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	field := L.CheckString(3)
	ok, err := p.Hdel(key, field)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(ok))
	return 1
}

func hlen(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	num, err := p.Hlen(key)
	if err != nil {
		L.Push(lua.LNumber(num))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(num))
	return 1
}

func hkeys(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	keys, err := p.Hkeys(key)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	t := new(lua.LTable)
	for _, r := range keys {
		t.Append(lua.LString(string(r)))
	}
	L.Push(t)
	return 1
}

func hgetall(L *lua.LState) int{
	p := checkRedis(L)
	key := L.CheckString(2)
	mapping := make(map[string]string)
	err := p.Hgetall(key, mapping)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	t1 := new(lua.LTable)
	for k, v := range mapping {
		t1.RawSet(lua.LString(k), lua.LString(v))
	}
	L.Push(t1)
	return 1
}

func keys(L *lua.LState) int{
	p := checkRedis(L)
	pattern := L.CheckString(2)
	res, err := p.Keys(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	t := new(lua.LTable)
	for _, r := range res {
		t.Append(lua.LString(r))
	}
	L.Push(t)
	return 1
}