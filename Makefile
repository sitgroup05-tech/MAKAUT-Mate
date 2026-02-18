





.PHONY: compile migrate 


compile : 

build:
	uv syc
	go mod download

build_docs:
	go install github.com/swaggo/swag/cmd/swag@latest
	cd backend ; swag init --parseDependency --parseInternal -g cmd/main.go -o ../docs 

server: rpc-server
	echo "hello world"
	
rpc-server : 
	cd backend/python-rpc;uv run main.py &
	


frontend : 


migrate : 
	uv run alembic up