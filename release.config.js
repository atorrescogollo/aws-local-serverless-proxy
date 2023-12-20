module.exports = {
  "branches": [
    "main",
  ],
  "ci": false,
  "plugins": [
    ["@semantic-release/exec", {
      "prepareCmd": "echo '${nextRelease.version}' > .version && sed -i 's|^\\(LABEL org.opencontainers.image.version\\) .*|\\1 \"v${nextRelease.version}\"|g' Dockerfile",
      "shell": true,
    }],
    ["@semantic-release/commit-analyzer", {
      "preset": "conventionalcommits"
    }],
    ["@semantic-release/release-notes-generator", {
      "preset": "conventionalcommits"
    }],
    ["@semantic-release/changelog", {
      "changelogFile": "CHANGELOG.md",
      "changelogTitle": "# Changelog\n\nAll notable changes to this project will be documented in this file.",
    }],
    ["@semantic-release/git", {
      "assets": [
        "CHANGELOG.md",
        "Dockerfile",
      ],
      "message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}",
    }],
    ["@semantic-release/github", {
      "assets": [
        { "path": "dist/*" },
      ],
      "successComment": "This ${issue.pull_request ? 'PR is included' : 'issue has been resolved'} in version ${nextRelease.version} :tada:",
    }],
  ]
}
