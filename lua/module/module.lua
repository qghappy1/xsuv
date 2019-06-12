
local os = require("os")

function test()
	local t = os.time()
	print(TestGame1:name())
	TestGame1:sendmsg(LoginType, 0, tostring(os.time()))
	local msg = TestLogin1:rpcmsg(GameType, 0, tostring(os.time()))
	print('rpc:'..msg)
	print('ms:'..(os.time()-t))
end 

test()
test()



