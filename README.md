# EcoMap

TODO

# Development

The following section focuses on the development part of the project, including prerequisites, how to build and run the code, and how to contribute.

### Table of Contents

- [Prerequisites](#prerequisites)
  - [Web App](#web-app)
  - [Android App](#android-app)
- [Web App](#web-app-1)
- [Android App](#android-app-1)
- [Contributing](#contributing)

## Prerequisites

### Web App

The web application runs on [Node.js](https://nodejs.org/) version 20.

Install the dependencies inside the `web` directory with:

```shell
npm install
```

### Android App

The Android application targets Android SDK 34.

Install [Android Studio](https://developer.android.com/studio) before you start contributing to the Android project.

## Web App

The web application can be found in the `web` directory. It uses [Svelte](https://svelte.dev/).

The web application contains several scripts to lint, build and run the project. To check the available scripts, run the following command inside the `web` directory:

```shell
npm run
```

## Android App

The Android application can be found in the `android` directory.

Use [gradlew](android/gradlew) to lint, build and run the project from the command line.

## Contributing

### Branches

- The [main branch](https://github.com/Goncalo-Marques/ecomap/tree/main) contains the production code
- To develop a new feature or fix a bug, a new branch should be created based on the main branch

### Issues

- The features and bugs should exist as a [GitHub issue](https://github.com/Goncalo-Marques/ecomap/issues) with an appropriate description
- The status of the issues can be seen in the associated [GitHub project](https://github.com/users/Goncalo-Marques/projects/2)

### Commits

- Git commits should follow [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/)

### Pull Requests

- To merge the code into production, a pull request should be opened following the existing template to describe it, as well as the appropriate labels
- To merge the pull request, the code must pass all GitHub action checks, be approved by at least one of the code owners, and be squashed and merged

### Deployments

- After the code is merged into the main branch, there is a GitHub action that automatically builds and deploys the code to production

### Releases

- To release a new version of the project, [semantic versioning](https://semver.org/) should be followed
