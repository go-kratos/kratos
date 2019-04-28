# Contribution Guidelines

## Introduction

This document explains how to contribute changes to the Ego project. It assumes you have followed the README.md and [API Document](https://github.com/go-ego/gpy/blob/master/docs/doc.md). <!--Sensitive security-related issues should be reported to [security@Ego.io](mailto:security@Ego.io.)-->

## Bug reports

Please search the issues on the issue tracker with a variety of keywords to ensure your bug is not already reported.

If unique, [open an issue](https://github.com/go-ego/gpy/issues/new) and answer the questions so we can understand and reproduce the problematic behavior.

The burden is on you to convince us that it is actually a bug in Ego. This is easiest to do when you write clear, concise instructions so we can reproduce the behavior (even if it seems obvious). The more detailed and specific you are, the faster we will be able to help you. Check out [How to Report Bugs Effectively](http://www.chiark.greenend.org.uk/~sgtatham/bugs.html).

Please be kind, remember that Ego comes at no cost to you, and you're getting free help.

## Discuss your design

The project welcomes submissions but please let everyone know what you're working on if you want to change or add something to the Ego repositories.

Before starting to write something new for the Ego project, please [file an issue](https://github.com/go-ego/gpy/issues/new). Significant changes must go through the [change proposal process](https://github.com/go-ego/proposals) before they can be accepted.

This process gives everyone a chance to validate the design, helps prevent duplication of effort, and ensures that the idea fits inside the goals for the project and tools. It also checks that the design is sound before code is written; the code review tool is not the place for high-level discussions.

## Testing redux

Before sending code out for review, run all the tests for the whole tree to make sure the changes don't break other usage and keep the compatibility on upgrade. You must be test on Mac, Windows, Linux and other. You should install the CLI for Circle CI, as we are using the server for continous testing.

## Code review

In addition to the owner, Changes to Ego must be reviewed before they are accepted, no matter who makes the change even if it is a maintainer. We use GitHub's pull request workflow to do that and we also use [LGTM](http://lgtm.co) to ensure every PR is reviewed by vz or least 2 maintainers.


## Sign your work

The sign-off is a simple line at the end of the explanation for the patch. Your signature certifies that you wrote the patch or otherwise have the right to pass it on as an open-source patch. 

## Maintainers

To make sure every PR is checked, we got team maintainers. A maintainer should be a contributor of Ego and contributed at least 4 accepted PRs. 

## Owners

Since Ego is a pure community organization without any company support, Copyright 2016 The go-ego Project Developers.


## Versions

Ego has the `master` branch as a tip branch and has version branches such as `v0.30.0`. `v0.40.0` is a release branch and we will tag `v0.40.0` for binary download. If `v0.40.0` has bugs, we will accept pull requests on the `v0.40.0` branch and publish a `v0.40.1` tag, after bringing the bug fix also to the master branch.

Since the `master` branch is a tip version, if you wish to use Ego in production, please download the latest release tag version. All the branches will be protected via GitHub, all the PRs to every branch must be reviewed by two maintainers and must pass the automatic tests.

## Copyright

Code that you contribute should use the standard copyright header:

```
// Copyright 2016 The go-ego Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-ego/gpy/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.
```

Files in the repository contain copyright from the year they are added to the year they are last changed. If the copyright author is changed, just paste the header below the old one.
