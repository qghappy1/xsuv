

--local io = require("io")
--local os = require("os")
local proto = require("proto")

function on_message(role_id, msg_id, smsg) 
	local msg, err = proto.unmarshal(smsg)
	if msg_id == 1 then 
		for _, card in pairs(msg.cards) do 
			print('card.weapon2:'..card.weapon2)
		end 
		local data, err1 = proto.marshal(1, msg)
		if err1 ~= nil then
			print(err1)
		else
			print(data)
		end 
	else
		print(msg.ret)
		local data, err1 = proto.marshal(2, msg)
		if err1 ~= nil then
			print(err1)
		else
			print(data)
		end
	end
end 



