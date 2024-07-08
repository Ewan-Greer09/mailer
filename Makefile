e2e:
	curl --header "Content-Type: application/json" \
	--request POST \
	--data '{"communication_channel":"email","subject":"password-reset","to":"user1@gmail.com", "from":"provider1@gmail.com", "reply_to":"no-reply@gmail.com", "metadata": {"source":"curl/makefile"}}' \
	http://localhost:3000/api/send/password-reset