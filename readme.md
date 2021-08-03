## Orbit SSR
Golang SSR framework featuring native react support with zero setup.

### Contributing

#### Pushing
- creating release ` /c/Users/guyal/Downloads/goreleaser_Windows_x86_64/goreleaser.exe release --snapshot`
- publish `npm publish --access=public`
- update tag `git tag -f -a v1.1.1`

### todo:
- dev server with hot reload
- allow interchangeable bundlers rather than explicitly using webpack.
- make cli pretty

~ build step
1. build 3 different versions of the application (windows, macOS, linux)
2. store those in a bin along with the assets + node_modules


~ install step
1. keep the target version that the system requires, delete the rest? 

