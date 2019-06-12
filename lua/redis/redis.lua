

--local io = require("io")
--local os = require("os")
local redis = require("redis")

function test()
	local db, err = redis.open("111.230.46.154:6379|1|qq20160101")
	if err ~= nil then 
		print(err)
		return 
	end 
	local role, ok 
	role, err = db:hgetall('1000000')
	if err ~= nil then 
		print(err)
		return 
	end 	
	for k, v in pairs(role) do 
		print(k..':'..v)
	end 
	print('------------------')
	err = db:hset('1000000', 'Power', '20')
	if err ~= nil then 
		print(err)
		return 
	end
end

function testExist()
	local db, err = redis.open("111.230.46.154:6379|1|qq20160101")
	if err ~= nil then
		print(err)
		return
	end
	local ok
	ok, err = db:is_exist('1000002')
	if err ~= nil then
		print(err)
		return
	end
	print(ok)
end

testExist()



