help:
	@echo "test             run test"
	@echo "lint             run lint"
	@echo "gen_pinyin_dict  gen pinyin dict"

.PHONY: test
test:
	@echo "run test"
	@go test -v -cover

.PHONY: gen_pinyin_dict
gen_pinyin_dict:
	@go run tools/gen_pinyin_dict.go tools/pinyin-data/pinyin.txt pinyin_dict.go

.PHONY: lint
lint:
	gofmt -s -w . pinyin tools
	golint .
	golint pinyin
	golint tools
