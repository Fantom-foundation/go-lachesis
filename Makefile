BUILD_TAGS?=lachesis

SUBDIRS := src/.
TARGETS := build proto clean
SUBDIR_TARGETS := $(foreach t,$(TARGETS),$(addsuffix $t,$(SUBDIRS)))

# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	glide install

# install compiles and places the binary in GOPATH/bin
install:
	go install --ldflags '-extldflags "-static"' \
		--ldflags "-X github.com/Fantom-foundation/go-lachesis/src/version.GitCommit=`git rev-parse HEAD`" \
		./cmd/lachesis
	go install --ldflags '-extldflags "-static"' \
		--ldflags "-X github.com/Fantom-foundation/go-lachesis/src/version.GitCommit=`git rev-parse HEAD`" \
		./cmd/network

# build compiles and places the binary in /build
build:
	CGO_ENABLED=0 go build \
		--ldflags "-X github.com/Fantom-foundation/go-lachesis/src/version.GitCommit=`git rev-parse HEAD`" \
		-o build/lachesis ./cmd/lachesis/main.go
	CGO_ENABLED=0 go build \
		--ldflags "-X github.com/Fantom-foundation/go-lachesis/src/version.GitCommit=`git rev-parse HEAD`" \
		-o build/network ./cmd/network/

# dist builds binaries for all platforms and packages them for distribution
dist:
	@BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/dist.sh'"

test:
	glide novendor | grep -v -e "^\.$$" | GOCACHE=off xargs go test -race -timeout 45s

# clean up and generate protobuf files
proto: clean

clean:

.PHONY: $(TARGETS) $(SUBDIR_TARGETS) vendor install dist test

# static pattern rule, expands into:
# all clean : % : foo/.% bar/.%
$(TARGETS) : % : $(addsuffix %,$(SUBDIRS))

# here, for foo/.all:
#   $(@D) is foo
#   $(@F) is .all, with leading period
#   $(@F:.%=%) is just all
$(SUBDIR_TARGETS) :
	@$(MAKE) -C $(@D) $(@F:.%=%)
