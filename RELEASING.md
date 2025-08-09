# Releasing a New Version

This document outlines the step-by-step process for releasing a new version of `llm-cli`.

## 1. Pre-release Checklist

- [ ] Ensure the `main` branch is stable and all tests are passing (`make test`).
- [ ] Ensure there are no known vulnerabilities (`make vulncheck`).

## 2. Release Steps

### Step 1: Update Documentation

1.  **Determine the new version number** (e.g., `v0.1.0`). We follow Semantic Versioning.
2.  **Update `CHANGELOG.md` and `CHANGELOG.ja.md`**: Add a new section for the release, moving changes from `[Unreleased]` to the new version section.
3.  **Update Development Log**: Create a new `docs/development_logs/dev_log-YYYY-MM-DD.md` file detailing the release process and its context.
4.  **Review and Update All Relevant Documentation**: Ensure all documentation reflects the changes in the new release. This includes, but is not limited to:
    *   `README.md` and `README.ja.md`: For user-facing features, configuration, and quick start guides.
    *   `CONTRIBUTING.md`: For any changes in development setup, testing, or contribution guidelines.
    *   `SECURITY.md`: For any updates to security policies or vulnerability reporting.
    *   Any other relevant `docs/` files.
    *   **Double-check for consistency**: Verify that new features, changes, or deprecations are accurately reflected across all relevant documents, including examples and command references.

### Step 2: Commit and Tag

1.  **Commit the changes**:
    ```sh
    git add CHANGELOG.md CHANGELOG.ja.md docs/development_logs/ go.mod go.sum
    git commit -m "feat: Release vX.Y.Z"
    ```
    *(Adjust the commit type e.g., `fix:`, `docs:` as appropriate)*

2.  **Tag the release**: Create a new Git tag on the release commit.
    ```sh
    git tag vX.Y.Z
    ```

### Step 3: Push to GitHub

1.  **Push the commit**:
    ```sh
    git push origin main
    ```

2.  **Push the tag**:
    ```sh
    git push origin vX.Y.Z
    ```
    *(If you need to fix a tag, use `git tag -d vX.Y.Z`, `git push --delete origin vX.Y.Z`, then re-tag and re-push)*

### Step 4: Build and Verify

1.  **Run the full build**:
    ```sh
    make all
    ```
2.  **Verify the version directly from the build output**: Instead of installing, check the version of the newly compiled binary. Because the commit is tagged, the version should be clean (e.g., `vX.Y.Z`). If it shows extra information (e.g., `vX.Y.Z-1-g123abc`), it means the tag was not placed on the latest commit.
    ```sh
    # For macOS
    ./bin/darwin-universal/llm-cli -v

    # For Linux
    # ./bin/linux-amd64/llm-cli -v
    ```
    The output should show the correct, clean version number (e.g., `llm-cli version vX.Y.Z`). This confirms the version was correctly embedded during the build.

### Step 5: Create GitHub Release

1.  Go to the [New Release page](https://github.com/magifd2/llm-cli/releases/new) on GitHub.
2.  Select the tag you just pushed (e.g., `vX.Y.Z`).
3.  Copy the relevant section from `CHANGELOG.md` into the release description.
4.  Upload the release assets from the local `bin/` directory:
    *   `llm-cli-darwin-universal.tar.gz`
    *   `llm-cli-linux-amd64.tar.gz`
    *   `llm-cli-windows-amd64.zip`
5.  Click **"Publish release"**.
