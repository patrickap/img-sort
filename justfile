PROJECT_NAME := "patrickap/img-sort"
PROJECT_VERSION := "VERSION"

[private]
get-version:
  @cat {{PROJECT_VERSION}}

[private]
set-version version:
  @echo {{version}} > {{PROJECT_VERSION}}

[private]
backup-version:
  @cp {{PROJECT_VERSION}} {{PROJECT_VERSION}}.bak

[private]
restore-version:
  @cp {{PROJECT_VERSION}}.bak {{PROJECT_VERSION}}

[private]
go-build:
  @go mod download
  @GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'github.com/patrickap/img-sort/m/v2/cmd.version=v$(just get-version)'" -o ./build/darwin-amd64/img-sort
  @GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'github.com/patrickap/img-sort/m/v2/cmd.version=v$(just get-version)'" -o ./build/darwin-arm64/img-sort
  @GOOS=linux GOARCH=amd64 go build -ldflags "-X 'github.com/patrickap/img-sort/m/v2/cmd.version=v$(just get-version)'" -o ./build/linux-amd64/img-sort
  @GOOS=linux GOARCH=arm64 go build -ldflags "-X 'github.com/patrickap/img-sort/m/v2/cmd.version=v$(just get-version)'" -o ./build/linux-arm64/img-sort

[private]
git-publish:
  @git add . -- ':!{{PROJECT_VERSION}}.bak'
  @git commit -m "chore(release): v$(just get-version)"
  @git push
  @git tag -a "v$(just get-version)" -m "Release v$(just get-version)"
  @git push --tags origin

[private]
clean-up:
  @rm {{PROJECT_VERSION}}.bak

[private]
release-patch:
  @just backup-version
  @just set-version $(just get-version | awk -F. -v OFS=. '{$3++; print}')
  @just go-build && just git-publish || just restore-version
  @just clean-up

[private]
release-minor:
  @just backup-version
  @just set-version $(just get-version | awk -F. -v OFS=. '{$2++; $3=0; print}')
  @just go-build && just git-publish || just restore-version
  @just clean-up

[private]
release-major:
  @just backup-version
  @just set-version $(just get-version | awk -F. -v OFS=. '{$1++; $2=0; $3=0; print}')
  @just go-build && just git-publish || just restore-version
  @just clean-up

release type:
  @just release-{{type}}
