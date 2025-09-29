.PHONY: dev build run clean

dev:
	npm run dev

build:
	npm run build

run:
	npm run start

clean:
	rm -rf bin apps/web/dist