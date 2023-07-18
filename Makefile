run:
	go build
	go run .

example:
	go build
	./pool_calc amountOut 155554778123672 0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2 0xdac17f958d2ee523a2206206994597c13d831ec7

lint:
	# go mod verify
	go vet
	golangci-lint run 
	golangci-lint run --presets bugs 
	golangci-lint run --presets style -D tagalign -D tagliatelle -D varnamelen -D depguard -D exhaustruct

