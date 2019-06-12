
--require("gate")

function test(agent)
	print(agent:id())
	print(agent:remoteaddr())
	agent:sendmsg("111")
	agent:close()
end



