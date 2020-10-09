BUILD := build
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
GREP := grep --exclude-dir=".git" -hIioPrs
LDFLAGS := -s -w
OUT := $(BUILD)/$(GOOS)/$(GOARCH)
PKG := $(shell $(GREP) "module\s+\K.+" go.mod)
SRC := $(shell find . -name "*.go" -exec dirname {} + | sort -u)
SRCDEPS := $(shell go list -deps $(SRC) | grep -iPs "^git")
TEST := $(shell find . -name "*_test.go" -exec dirname {} + | sort -u)
VERS := $(shell $(GREP) "const\s+Version\s+\=\s+\"\K[^\"]+" .)

all: build

%: %-default
	@true

ifeq ($(shell ls -d cmd 2>/dev/null), cmd)
    include gomk/cmd.mk
else
    include gomk/so.mk
endif

clean-default: fmt
	@rm -rf "$(BUILD)" go.sum
	@[[ ! -f go.mod ]] || go mod tidy

clena-default: clean

cyclo-default: havego
	@which gocyclo >/dev/null 2>&1 || \
	    go get -u github.com/fzipp/gocyclo
	@gocyclo -over 15 . || echo -n

dir-default:
	@mkdir -p "$(OUT)"

fmt-default: havego
	@go fmt $(SRC) >/dev/null

gen-default: havego
	@go generate $(SRC)

havego-default:
	@which go >/dev/null 2>&1

ineffassign-default: havego
	@which ineffassign >/dev/null 2>&1 || \
	    go get -u github.com/gordonklaus/ineffassign
	@ineffassign . || echo -n

lint-default: havego
	@which golint >/dev/null 2>&1 || \
	    go get -u golang.org/x/lint/golint
	@golint $(SRC)

push-default:
	@git tag "v$(VERS)"
	@git push
	@git push --tags

reportcard-default: fmt cyclo ineffassign lint simplify vet

simplify-default: havego
	@gofmt -s -w $(SRC)

sloc-default: havego
	@which sloc >/dev/null 2>&1 || \
	    go get -u github.com/bytbox/sloc/sloc
	@sloc .

test-default: havego
	@go clean --testcache
	@for i in $(TEST); do \
	    go test -v $(PKG)/$${i##./}; \
	done

updatedeps-default: havego
	@for dep in $(SRCDEPS); do \
	    go get -u -v $$dep; \
	done
	@rm -f go.sum
	@[[ ! -f go.mod ]] || go mod tidy

updatereportcard-default: havego
	@go get -u github.com/fzipp/gocyclo
	@go get -u github.com/gordonklaus/ineffassign
	@go get -u golang.org/x/lint/golint
	@rm -f go.sum
	@[[ ! -f go.mod ]] || go mod tidy

vet-default: havego
	@go vet $(SRC) || echo -n

yank-default:
	@git tag -d "v$(VERS)"
	@git push -d origin "v$(VERS)"
