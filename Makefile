e2e:
	curl -H "Content-Type: application/json" \
	-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJERVZFTE9QTUVOVCBLRVkiLCJuYW1lIjoiRXdhbiBHcmVlciIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoyNTQxMTYxNTExfQ.2xpRU9rmIoNpk9syA_6Bgc3Alrbu3tp4Xm0-s-MymTs" \
	--request POST \
	--data '{"communication_channel":"email","to":"user1@gmail.com", "from":"provider1@gmail.com", "reply_to":"no-reply@gmail.com", "metadata": {"source":"curl/makefile"}, "message_datafields":{"subject":"password-reset", "first_name":"John", "last_name":"Doe"}}' \
	http://localhost:3000/api/send/password-reset
