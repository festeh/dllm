install_server:
		@echo "Building dllm server"
		@cd backend/cmd/server && echo $(shell cat backend/version.txt)
		@cd backend/cmd/server && go build -o dllm_server -ldflags "-X dllm.VERSION=$(shell cat backend/version.txt)"
		@echo "Moving binary to build folder..."
		@mkdir -p build
		@mv backend/cmd/server/dllm_server build/
		@echo "Done."


install_dllm:
		@echo "Building dllm terminal client"
		@cd backend/cmd/terminal && go build -o dllm .
		@echo "Moving binary to build folder..."
		@mkdir -p ~/.local/bin
		@mv backend/cmd/terminal/dllm ~/.local/bin/
		@echo "Done."
