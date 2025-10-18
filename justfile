alias r := run
alias b := build
alias f := fmt
alias i := install
alias l := lint
alias t := test
alias up := upgrade
alias ti := tidy
alias v := vuln
alias us := update-snapshots

default:
    just --choose

run check:
    go run . {{ check }}

build:
    go build -ldflags "-w -s" .

fmt:
    gofumpt -l -w .

install:
    go install -ldflags "-w -s" .

lint:
    golangci-lint run

test:
    go test -v -count 1 ./...

upgrade:
    go get -u ./...

tidy:
    go mod tidy

vuln:
    govulncheck ./...

update-snapshots $UPDATE_SNAPS='true':
  just test
