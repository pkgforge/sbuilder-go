### SBUILDTOOLS written in Go. (for your ease and convenience)
- Currently implemented:
    - linter

#### About(linter)
The `sbuild-linter` is a tool that validates your recipes, it checks for all kinds of wrongdoings and makes sure that your file is valid, according to the [SBUILD spec](https://github.com/pkgforge/soarpkgs/blob/main/SBUILD_SPEC.md)
The linter will produce: `{{ yourRecipe }}.validated` and a `{{ yourRecipe }}.pkgver`, these should be used an implementation of the `sbuild-builder`.
