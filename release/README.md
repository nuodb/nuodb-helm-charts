# The `relman.py` tool

`relman.py` is a tool for managing releases from Git repositories.
It imposes branching and tagging conventions for associating releases with Git commits, can generate changelogs for releases, and can be used to manage releases for multiple deliverables out of the same repository.

## Deliverables

`relman.py` can be used to manage releases for multiple _deliverables_ out of the same Git repository (monorepo).
A deliverable is associated with a subdirectory within the Git repository that contains a `get-version.sh` script to obtain the current version.

## Branching and tagging conventions

Branching and tagging conventions are used to specify which commits the deliverable artifacts for a release should be created from.
Semantic versioning is used with the form `<major>.<minor>.<patch>`.
To summarize semver conventions:

- Patch version releases can introduce bug fixes.
- Minor version releases can introduce backwards-compatible enhancements.
- Major version releases can introduce backwards-incompatible enhancements.

The tag format `v<major>.<minor>.<patch>` is used to denote a release of the root deliverable.
The tag format `<deliverable>/v<major>.<minor>.<patch>` is used to denote a release of a non-root deliverable with path `<deliverable>`.

Major and minor releases are created from the `main` branch, while patch releases are created from release branches.
Branches of the form `v<major>.<minor>-dev` are used to create patch releases for the root deliverable, for a particular major and minor version.
Branches of the form `<deliverable>/v<major>.<minor>-dev` are used to create patch releases for non-root deliverables with path `<deliverable>`, for a particular major and minor version.

## `relman.py` usage

The `relman.py` tool can be used to prepare a release.
Currently it is capable of checking that the current branch follows the tagging conventions using `--check-current` or `--check-tags`, and it can also generate a changelog by scraping commit messages since the last commit using `--create-changelog`.
The changelog files for all releases are stored in `changelog` for the root deliverable and `<deliverable>/changelog` for non-root deliverables.
These files should be checked into version control, perhaps after some manual editing, and can be used to describe the release in GitHub.

## Example: Creating a minor release

To create a minor release of the root deliverable, `relman.py` can be used to create the commit, tag, and release branch as follows:

1. Switch to the `main` branch:
```sh
git checkout main
```
2. Use the `relman.py` tool to check that the current development version is larger than the last release tag, commit the changelog, tag the commit, and create release branch:
```sh
./release/relman.py --check-tags --create-changelog --commit --tag --branch
```
3. Once you've verified that everything is correct, push the commit and tag:
```sh
git push --tags --atomic origin main v<major>.<minor>-dev
```

Creating a major release is identical to creating a minor release.
Creating a patch release is similar, but there is no step to create another branch, and it is only the patch version that is updated.
