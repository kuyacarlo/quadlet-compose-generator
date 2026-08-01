NAME = hitch-generator

build:
	go build -o $(NAME) .

install: build
	install -Dpm 0755 $(NAME) $(HOME)/.config/systemd/user-generators/$(NAME)

test:
	go test -v ./...

clean:
	rm -f $(NAME)
