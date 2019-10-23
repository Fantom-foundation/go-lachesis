
# build
.PHONY : build
build :
	go build -o build/lachesis ./cmd/lachesis

#test
.PHONY : test
test :
	go test ./...

#clean
.PHONY : clean
clean :
	rm ./build/lachesis
