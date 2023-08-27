# NuoDB Helm Charts release process

NuoDB Helm Charts releases are delivered as follows:

1. First a changelog is generated that lists the commits in a release.
2. The changelog is committed to version control.
3. The commit is tagged with a release, and for `<major>.<minor>.0` releases (i.e. off of master), a branch is created for patch releases for that `<major>.<minor>` version.
4. A GitHub release is created that includes the changelog and release artifacts, which are the packaged Helm charts.
5. A publishing step is performed to make the new release artifacts downloadable using `helm`.

## Branching and tagging conventions

Branching and tagging conventions are used to specify which commits the deliverable artifacts for a release should be created from.
Semantic versioning is used with the form `<major>.<minor>.<patch>`.
To summarize semver conventions:

- Patch version releases can introduce bug fixes.
- Minor version releases can introduce backwards-compatible enhancements.
- Major version releases can introduce backwards-incompatible enhancements.

The tag format `v<major>.<minor>.<patch>` is used to denote a release.

Major and minor releases are created from the `master` branch, while patch releases are created from release branches.
Branches of the form `v<major>.<minor>-dev` are used to create patch releases for a particular major and minor version.

## `relman.py` usage

The `relman.py` tool can be used to perform steps 1, 2, and 3 of the release process described above.

1. Switch to the `master` branch:
```sh
git checkout master
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
