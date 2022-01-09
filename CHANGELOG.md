# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.9.0] - 2022-01-09
### Added
- Methods `Abs(x)`, `MinValue()`, `MaxValue()`.
### Changed
- It is 25% faster.

## [1.8.0] - 2022-01-05
### Added
- Support for 2-bit exponent.
### Removed
- All preconfigured methods.

## [1.7.0] - 2022-01-05
### Changed
- Tests now work through the new API.
- Single-module package.

## [1.6.0] - 2022-01-04
### Added
- Fully customizable object-oriented API.
### Deprecated
- All preconfigured methods.

## [1.5.0] - 2022-01-01
### Removed
- `Default`, `14` and `m11x3` types.

### Changed
- Type `defaultD` is now `12`.
- Type `unsigned` renamed to `12u`.
- Type `14d` renamed to `14`.
- Type `m11x3d` renamed to `15x3`.

## [1.4.0] - 2022-01-01
### Changed
- `EncodeDelta[TypeName]` renamed to `GetIntegerDelta[TypeName]`.
- `DecodeDelta[TypeName]` renamed to `UseIntegerDelta[TypeName]`.

## [1.3.1] - 2022-01-01
### Changed
- Even though everything works fine, I've added an explicit limit
  on the maximum value of the mantissa
  to ensure that it never gets rounded incorrectly.

## [1.3.0] - 2022-01-01
### Deprecated
- `Default`, `14` and `m11x3` types.

## [1.2.0] - 2022-01-01
### Added
- New types `defaultD`, `14d`, `m11x3d` with different bits ordering 
suitable for delta encoding.
- Functions `EncodeDelta[TypeName]`, `DecodeDelta[TypeName]`.

### Changed
- Accuracy has increased. The mantissa
  &ndash; `binarySignificand` variable &ndash;
  is now rounded to the nearest integer
  ([git:18cfb1a9](https://github.com/georgy7/toyfloat/commit/18cfb1a9a1ef6ed719e7a208cc8add4975643049#diff-143de6cb31239060551e2b97d13f56c5567d10caf0a112671f77f8f40a82caa9L107)).
- Function `getExponent` has been changed to match the rounding change
  ([git:28458380](https://github.com/georgy7/toyfloat/commit/284583809808e712e83bf2521a4c913547eac45d#diff-143de6cb31239060551e2b97d13f56c5567d10caf0a112671f77f8f40a82caa9R123)).

## [1.1.0] - 2021-12-25
### Added
- 15-bit type with 3-bit exponent (`m11x3`).
