build:
	go build

build-verbose:
	go build -gcflags '-m -l'

test:
	go test

benchmark:
	go test -run NONE -bench . -benchmem 

compare:
	go test -run NONE -bench . -benchmem -tags cbor_comparison
