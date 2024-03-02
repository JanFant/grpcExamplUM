say_hello:
	echo "Hello"

proto_build:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usermgmt/usermgmt.proto

run_Server:
	go run "usermgmtServer/usermgmtServer.go"

run_Client:
	go run "usermgmtClient/usermgmtClient.go"