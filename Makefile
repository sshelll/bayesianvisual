install:
	@go install ./cmd/bayesianvisual
	@if [ ! -d ~/.oh-my-zsh/custom/plugins/bayesianvisual ]; then \
		mkdir -p ~/.oh-my-zsh/custom/plugins/bayesianvisual; \
	fi
	@bayesianvisual completion zsh > ~/.oh-my-zsh/custom/plugins/bayesianvisual/bayesianvisual.plugin.zsh
