# Releasing

How to publish a new version of `gphotos-uploader-cli`. Only repository admins can do this: `v*` tags are protected, and the release workflow needs a maintainer's approval.

## What gets published

Pushing a `v*` tag runs the [Create Release](.github/workflows/create-release.yml) workflow. It runs the tests, then [GoReleaser](.goreleaser.yml):

- builds the Linux and macOS binaries and attaches them to a GitHub release, with a `checksums.txt` file.
- writes the release notes from the PR titles since the previous tag. `docs:`, `ci:` and `test:` PRs are left out.
- updates the Homebrew formula in [gphotosuploader/homebrew-tap](https://github.com/gphotosuploader/homebrew-tap).

[CHANGELOG.md](CHANGELOG.md) is written by hand and isn't generated. Each PR with user-facing changes adds an entry under `## Unreleased`.

## Steps

1. **Pick the version.** Follow [Semantic Versioning](https://semver.org/), based on what's under `## Unreleased` in `CHANGELOG.md`:
   - major: something removed or changed in a way that breaks existing configs or commands.
   - minor: new features or deprecations.
   - patch: only fixes and dependency updates.

2. **Prepare the release in a PR** titled `docs: prepare release vX.Y.Z`:
   - In `CHANGELOG.md`, rename `## Unreleased` to `## X.Y.Z` and add a new, empty `## Unreleased` above it.
   - Check that every user-facing change since the last release has an entry. `git log --oneline vPREVIOUS..main` lists them.
   - For a major release, add the migration steps to [UPGRADING.md](UPGRADING.md).

3. **Merge the PR**, and wait for the checks on `main` to pass.

4. **Tag `main` and push the tag:**

   ```sh
   git switch main && git pull
   git tag -a vX.Y.Z -m "Version X.Y.Z"
   git push origin vX.Y.Z
   ```

5. **Approve the deployment.** The workflow waits for approval on the `release` environment. Approve it from the workflow run in the Actions tab.

6. **Check the result:**
   - the [release page](https://github.com/gphotosuploader/gphotos-uploader-cli/releases) has the binaries, `checksums.txt` and the release notes.
   - the formula in the Homebrew tap shows the new version.

## Release candidates

To test a release before publishing it, tag a release candidate such as `vX.Y.Z-rc1` (message `Version X.Y.Z RC1`). GoReleaser marks it as a pre-release, so it isn't shown as the latest release, and doesn't update the Homebrew formula.

## If something goes wrong

- **Tests or GoReleaser fail before publishing:** fix the problem in a PR. Then publish the next patch version. Don't move or re-push the existing tag.
- **A broken release was published:** publish a new patch version with the fix. Don't delete the published release, because users and the Homebrew formula may already point to it.
