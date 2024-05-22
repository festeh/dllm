install_server:
		@echo "Building dllm server"
		@cd backend/cmd/server && go build -o server .
		@echo "Moving binary to build folder..."
		@mkdir -p build
		@mv backend/cmd/server/server build/
		@echo "Done."
