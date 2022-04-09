# Changelog
All notable changes to this project will be documented in this file.

## 0.7.1 (April 9, 2022)

Bugfixes:
- [#82](https://github.com/GuyARoss/orbit/issues/82): Fix for windows support

## 0.7.0 (April 8, 2022)

Enhancements:
- [#49](https://github.com/GuyARoss/orbit/issues/49): Added new tool for visualization of the dependency graphs output by the `depout` flag.
- [#10](https://github.com/GuyARoss/orbit/issues/10): Added support for recomputing dependency map during repack.
- [#59](https://github.com/GuyARoss/orbit/issues/59): Added experimental support for react server-side-rendering.
- [#11](https://github.com/GuyARoss/orbit/issues/11): Added support for vanilla js components with micro-frontends.

Bugfixes:
- [#51](https://github.com/GuyARoss/orbit/issues/51): Fix bug that prevented from detecting correct page extension.
- [#55](https://github.com/GuyARoss/orbit/issues/55): Fix issue that applied duplicate pages to the generated go output.

## 0.3.6 (March 20, 2022)

Bugfixes:
- [#45](https://github.com/GuyARoss/orbit/issues/45): Fix JS lexer from having duplicate extensions in the file output.
- [#41](https://github.com/GuyARoss/orbit/issues/41): Fix JS lexer bug that caused non jsx files to fail parsing on jsx rules.
- [#32](https://github.com/GuyARoss/orbit/issues/32): Fix dev command to detect new pages.
- [#38](https://github.com/GuyARoss/orbit/issues/38): Fix for various `libout` slug parsing bugs.


## 0.3.2 (March 13, 2022)

Bugfixes:
- [#36](https://github.com/GuyARoss/orbit/issues/36): Fix `libout` issue that caused cached urls to be applied twice to the web bundle.
- [#29](https://github.com/GuyARoss/orbit/issues/29): Fix issue that caused a browser error to be thrown when a component was processed to quickly.

## 0.3.0 (March 11, 2022)

Enhancements:
- [#30](https://github.com/GuyARoss/orbit/issues/30): Added support for page hot reloading for dev command.

## 0.2.0 (March 8, 2022)
initial release
