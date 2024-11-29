### SBUILDTOOLS written in Go. (for your ease and convenience)
- Currently implemented:
    - linter

#### About(linter)
The `sbuild-linter` is a tool that validates your recipes, it checks for all kinds of wrongdoings and makes sure that your file is valid, according to the [SBUILD spec](https://github.com/pkgforge/soarpkgs/blob/main/SBUILD_SPEC.md)
The linter will produce: `{{ yourRecipe }}.validated` and a `{{ yourRecipe }}.pkgver`, these should be used by an implementation of the `sbuild-builder`.


TODO(linter):
- Switch to using "github.com/goccy/go-yaml"
- Convert into library
- Preserve input file's whitespace
- Use AnnotateSource to print errors and warnings that show the user specifically what line is wrong
- Write unit tests (and make a GH action that runs on every commit against those tests) [badge?]
- Flags to disable creation of .validated
- Print time it took to validate each file
- Parallel mode?

TODO(builder):
- Implement using the `sbuild-linter` as library

RFC:
- Add support for templating in the `linter` and `builder` (go tags?)
