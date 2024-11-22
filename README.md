### SBUILDTOOLS written in Go. (for your ease and convenience)
- Currently implemented:
    - linter

#### About(linter)
The `sbuild-linter` is a tool that validates your recipes, it checks for all kinds of wrongdoings and makes sure that your file is valid, according to the [SBUILD spec](https://github.com/pkgforge/soarpkgs/blob/main/SBUILD_SPEC.md)
The linter will produce: `{{ yourRecipe }}.validated` and a `{{ yourRecipe }}.pkgver`, these should be used an implementation of the `sbuild-builder`.


TODO(linter):
- Convert into library
- Preserve input file's whitespace
- Write unit tests (and make a GH action that runs on every commit against those tests) [badge?]
- Flags to disable: {shellcheck, creation of .pkgver, creation of .validated}

TODO(builder):
- Implement using the `sbuild-linter` as library
