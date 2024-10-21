build:
	go build -o ./bin/novus main.go

# This should only run in the CI pipeline to set latest version before release
git_tag=$(shell git describe --tags --abbrev=0)
update-assets-version:
	sed -i '' "s/%RELEASE_VERSION%/$(git_tag)/g" ./assets/nginx/html/404.html
	sed -i '' "s/%RELEASE_VERSION%/$(git_tag)/g" ./assets/nginx/html/502.html
	sed -i '' "s/%RELEASE_VERSION%/$(git_tag)/g" ./assets/nginx/html/index.html