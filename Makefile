.PHONY: all run epub_test_v2 epub_test_v3 paste_test

GO_CMD = go run .
TEST_DIR = ./test_file

all: run

run:
	go run .

epub_v2:
	$(GO_CMD) $(TEST_DIR)/test_v2.epub

epub_v3:
	$(GO_CMD) $(TEST_DIR)/test_v3.epub

paste:
	wl-copy < $(TEST_DIR)/test_for_paste.txt
	$(GO_CMD) --paste
