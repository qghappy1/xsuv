

--local io = require("io")
--local os = require("os")
local proto = require("proto")

function on_message(role_id, msg_id, msg) 
	if msg_id == 1 then 
		for _, card in msg.Cards() do 
			--print(card)
			print('card.weapon2:'..card:GetWeapon2())
			for m, s in card.Skills() do 
				print('skill.level', s:GetLevel())
				s.Id = s.Id + 10
				for z, id in s.Id() do 
					print('skill.id:'..id)
				end 
			end
		end

		local s = CardBagSC()
		local c = CardInfo()
		s.Cards = CardInfos()
		s.Cards = s.Cards + c 
		local data, err = proto.marshal(1, s)
		if err ~= nil then 
			print(err)
		else 
			print('marshal:'..data)
		end 
	else
		print(msg.ret)
	end
end 



