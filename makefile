default:
	go build -o vite

install:
	install -Dm755 vite $(DESTDIR)$(PREFIX)/bin/vite

uninstall:
	@rm -f $(DESTDIR)$(PREFIX)/bin/vite
