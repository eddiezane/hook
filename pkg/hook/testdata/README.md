# Fire tests

Tests are automatically ran for any file in this folder with a `.hook`
extenstion.

For each test, there should be a corresponding `.http` file for the expected
HTTP request.

Note: When adding new tests, HTTP requests need carriage returns, which vim may
not insert by default. To add a carriage return manually, enter insert mode and
use `ctrl+v` + `enter`.
