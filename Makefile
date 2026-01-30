# xcap Makefile

.PHONY: clean

# clean 清理构建产物和调试信息
clean:
	@echo "Cleaning build artifacts and debug files..."
	@rm -rf xcap_output/
	@rm -rf examples/basic/xcap_output/
	@rm -f *.png
	@rm -f *.log
	@rm -rf bin/
	@rm -rf dist/
	@echo "Done."
