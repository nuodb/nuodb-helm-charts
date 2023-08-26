#!/usr/bin/env python

import logging
import os
import shutil
import subprocess
import sys
import unittest

from relman import run, LOGGER

LOGGER.setLevel(logging.DEBUG)

def coverage_enabled():
    try:
        import coverage # type: ignore # pylint: disable=missing-imports
        return True
    except ImportError:
        return False

class RelmanTest(unittest.TestCase):

    relman_dir = os.path.dirname(os.path.abspath(__file__))
    tmp_dir = os.path.join(relman_dir, 'tmp')
    test_repo = os.path.join(tmp_dir, 'test_repo')
    coverage_enabled = coverage_enabled()

    def path(self, *args):
        return os.path.join(self.test_repo, *args)

    @classmethod
    def runcmd(cls, *args, **kwargs):
        if 'cwd' not in kwargs:
            kwargs['cwd'] = cls.test_repo
        if 'stderr' not in kwargs:
            kwargs['stderr'] = subprocess.STDOUT
        try:
            return run(*args, **kwargs)
        except subprocess.CalledProcessError as e:
            LOGGER.error('---\n%s\n---', e.output.decode('utf-8'))
            raise

    def relman(self, *args, **kwargs):
        relman_script = self.path('release', 'relman.py')
        if self.coverage_enabled:
            args = (sys.executable, '-m', 'coverage', 'run', '-a', relman_script) + args
            return self.runcmd(*args, **kwargs)
        else:
            return self.runcmd(sys.executable, relman_script, *args, **kwargs)

    def create_file(self, content, *path, **kwargs):
        abspath = self.path(*path)
        with open(abspath, 'w') as f:
            f.write(content)
        # set mode for file
        mode = kwargs.get('mode')
        if mode:
            os.chmod(abspath, mode)
        # stage file
        if kwargs.get('stage'):
            self.runcmd('git', 'add', abspath)

    def commit(self, msg):
        self.runcmd('git', 'commit', '-m', msg)

    def set_version(self, version, *path):
        path += ('get-version.sh',)
        self.create_file('#!/bin/sh\n\necho "{}"\n'.format(version), *path, mode=0o755, stage=True)

    @classmethod
    def clean_tmp(cls, path=tmp_dir):
        # invoke 'git clean' on real repo to clean up test files
        run('git', 'clean', '-Xff', '--', path)

    @classmethod
    def setUpClass(cls):
        # create directory to store code coverage information if coverage module is available
        if cls.coverage_enabled:
            # clean and re-create directory
            cls.coverage_dir = os.path.join(cls.relman_dir, 'coverage_info')
            cls.clean_tmp(cls.coverage_dir)
            os.makedirs(cls.coverage_dir)

    def setUp(self):
        self.clean_tmp()

        # create release directory and copy relman script
        os.makedirs(self.path('release'))
        self.runcmd('git', 'init')
        shutil.copy(os.path.join(self.relman_dir, 'relman.py'), self.path('release'))

        # create initial commit with relman script
        self.runcmd('git', 'add', self.path('release'))
        self.create_file('.coveragerc', '.gitignore', stage=True)
        self.commit('Initial commit to test repo')

        # create get-version.sh for root project
        self.set_version('1.0.0')
        os.makedirs(self.path('changelogs'))
        self.create_file('', 'changelogs', '.keep', stage=True)
        self.commit('Creating root project')

        # create sub-project and get-version.sh
        os.makedirs(self.path('subproject'))
        self.set_version('1.0.0', 'subproject')
        os.makedirs(self.path('subproject', 'changelogs'))
        self.create_file('', 'subproject', 'changelogs', '.keep', stage=True)
        self.commit('Creating sub-project (#1)')

        # configure code coverage if module is available
        if self.coverage_enabled:
            self.create_file("""
[run]
data_file = {}

[html]
directory = {}
""".format(os.path.join(self.coverage_dir, '.coverage'), os.path.join(self.coverage_dir, 'html')), '.coveragerc')

    MIN_COVERAGE_PCT = 90

    @classmethod
    def tearDownClass(cls):
        # aggregate coverage data and generate html coverage report if coverage is enabled
        if cls.coverage_enabled:
            cls.runcmd(sys.executable, '-m', 'coverage', 'html')
            cls.runcmd(sys.executable, '-m', 'coverage', 'report', '--fail-under=' + str(cls.MIN_COVERAGE_PCT))

    def check_changelogs(self, root_changelog, subproject_changelog, is_show=False):
        # changelog for root project should include all commits
        if is_show:
            self.assertIn('Generating changelog for version 1.0.0:', root_changelog)
        self.assertIn('# Changelog [v1.0.0](https://github.com/<org>/<repo>/tree/v1.0.0)', root_changelog)
        self.assertIn('## [Full Changelog](https://github.com/<org>/<repo>/compare/...v1.0.0)', root_changelog)
        self.assertIn('Initial commit to test repo', root_changelog)
        self.assertIn('Creating root project', root_changelog)
        self.assertIn(r'Creating sub-project [\#1](https://github.com/<org>/<repo>/pull/1)', root_changelog)

        # changelog for sub-project should only include commits that change files in that directory
        if is_show:
            self.assertIn('Generating changelog for version 1.0.0 of subproject:', subproject_changelog)
        self.assertIn('# Changelog [subproject/v1.0.0](https://github.com/<org>/<repo>/tree/subproject/v1.0.0)', subproject_changelog)
        self.assertIn('## [Full Changelog](https://github.com/<org>/<repo>/compare/...subproject/v1.0.0)', subproject_changelog)
        self.assertNotIn('Initial commit to test repo', subproject_changelog)
        self.assertNotIn('Creating root project', subproject_changelog)
        self.assertIn(r'Creating sub-project [\#1](https://github.com/<org>/<repo>/pull/1)', root_changelog)

    def check_negative(self, fn, output):
        try:
            fn()
            self.fail('Expected command to fail')
        except subprocess.CalledProcessError as e:
            self.assertIn(output, e.output.decode('utf-8'))

    def testShowChangelog(self):
        # show changelog for current version of root project
        root_changelog = self.relman('--show-changelog')

        # show changelog for current version of sub-project
        subproject_changelog = self.relman('--path', 'subproject', '--show-changelog')

        # check that changelogs contain the expected entries
        self.check_changelogs(root_changelog, subproject_changelog, is_show=True)

    def testCreateChangelog(self):
        # create changelog for current version of root project and read it from file
        self.relman('--create-changelog')
        with open(self.path('changelogs', 'v1.0.0.md')) as f:
            root_changelog = f.read()

        # create changelog for current version of sub-project and read it from file
        self.relman('--path', 'subproject', '--create-changelog')
        with open(self.path('subproject', 'changelogs', 'v1.0.0.md')) as f:
            subproject_changelog = f.read()

        # check that changelogs contain the expected entries
        self.check_changelogs(root_changelog, subproject_changelog)

        # negative test: try to create changelog with unstaged changes
        self.check_negative(
                lambda: self.relman('--create-changelog'),
                'Cannot write "{}" because it has uncommitted changes'.format(self.path('changelogs', 'v1.0.0.md')))

        # negative test: try to create changelog for sub-project with unstaged changes
        self.check_negative(
                lambda: self.relman('--path', 'subproject', '--create-changelog'),
                'Cannot write "{}" because it has uncommitted changes'.format(self.path('subproject', 'changelogs', 'v1.0.0.md')))

        # remove unstaged changelog files
        self.runcmd('git', 'clean', '-f')

        # create changelog files again, this time committing them and tagging as release
        self.relman('--create-changelog', '--commit', '--tag')
        self.relman('--path', 'subproject', '--create-changelog', '--commit', '--tag')

        # check that changelog files are identical to previously created ones
        with open(self.path('changelogs', 'v1.0.0.md')) as f:
            self.assertEqual(root_changelog, f.read())

        # create changelog for current version of sub-project and read it from file
        with open(self.path('subproject', 'changelogs', 'v1.0.0.md')) as f:
            self.assertEqual(subproject_changelog, f.read())

    def testCreateAndCheckReleases(self):
        # create changelog commit and tag release
        self.relman('--create-changelog', '--commit', '--tag')
        self.assertEqual('v1.0.0', self.runcmd('git', 'tag', '--points-at').strip())

        # create changelog commit and tag release for sub-project
        self.relman('--path', 'subproject', '--create-changelog', '--commit')
        self.assertEqual('', self.runcmd('git', 'tag', '--points-at').strip())

        # commit can be tagged separately
        self.relman('--path', 'subproject', '--tag')
        self.assertEqual('subproject/v1.0.0', self.runcmd('git', 'tag', '--points-at').strip())

        # negative test: try to create changelog with the same previous and current tag
        self.check_negative(
                lambda: self.relman('--show-changelog'),
                'Cannot create changelog because previous and current tag are both v1.0.0')

        # negative test: try to create changelog for sub-project with the same previous and current tag
        self.check_negative(
                lambda: self.relman('--path', 'subproject', '--show-changelog'),
                'Cannot create changelog because previous and current tag are both subproject/v1.0.0')

        # bump version of root project
        self.set_version('1.1.0')
        self.commit('Bumping minor version of root project')

        # bump version of sub-project
        self.set_version('2.0.0', 'subproject')
        self.commit('Bumping major version of sub-project')

        # show changelog and check that it contains only the commits since last release tag
        root_changelog = self.relman('--show-changelog')
        self.assertIn('Generating changelog for version 1.1.0:', root_changelog)
        self.assertIn('# Changelog [v1.1.0](https://github.com/<org>/<repo>/tree/v1.1.0)', root_changelog)
        self.assertIn('## [Full Changelog](https://github.com/<org>/<repo>/compare/v1.0.0...v1.1.0)', root_changelog)
        self.assertIn('Bumping minor version of root project', root_changelog)
        self.assertIn('Bumping major version of sub-project', root_changelog)
        self.assertNotIn(r'Creating sub-project [\#1](https://github.com/<org>/<repo>/pull/1)', root_changelog)

        # show changelog and check that it contains only the commits since last release tag for sub-project
        subproject_changelog = self.relman('--path', 'subproject', '--show-changelog')
        self.assertIn('Generating changelog for version 2.0.0 of subproject:', subproject_changelog)
        self.assertIn('# Changelog [subproject/v2.0.0](https://github.com/<org>/<repo>/tree/subproject/v2.0.0)', subproject_changelog)
        self.assertIn('## [Full Changelog](https://github.com/<org>/<repo>/compare/subproject/v1.0.0...subproject/v2.0.0)', subproject_changelog)
        self.assertIn('Bumping major version of sub-project', subproject_changelog)
        self.assertNotIn('Bumping minor version of root project', subproject_changelog)
        self.assertNotIn(r'Creating sub-project [\#1](https://github.com/<org>/<repo>/pull/1)', subproject_changelog)

        # use --check-current and --check-tags to check current version and release tags
        self.relman('--check-current')
        self.relman('--check-tags')
        self.relman('--path', 'subproject', '--check-current')
        self.relman('--path', 'subproject', '--check-tags')

        # negative test: set current to a different version from release tag on HEAD
        self.relman('--create-changelog', '--commit', '--tag')
        self.set_version('1.2.0')
        self.check_negative(
                lambda: self.relman('--check-current'),
                'Release tag v1.1.0 does not match current version 1.2.0')

        # remove release tag
        self.runcmd('git', 'tag', '--delete', 'v1.1.0')
        self.relman('--check-current')

        # negative test: try to tag a release with uncommitted changes
        self.check_negative(
                lambda: self.relman('--tag'),
                'Cannot create release tag because there are uncommitted changes')

        # negative test: set current to lower version than last release tag
        self.set_version('0.9.0')
        self.check_negative(
                lambda: self.relman('--check-current'),
                'Release tag v1.0.0 has later version than current version 0.9.0')

        # revert version change
        self.runcmd('git', 'reset', 'HEAD', '--hard')

        # negative test: create release tag with lower version than previous release tag
        self.runcmd('git', 'tag', 'v0.9.0', 'HEAD~')
        self.check_negative(
                lambda: self.relman('--check-tags'),
                'Release tag v1.0.0 has later version than subsequent release tag v0.9.0')

    def testBranchConventions(self):
        # negative test: release tag to non-0 patch version on main branch
        self.runcmd('git', 'tag', 'v0.9.9')
        self.set_version('0.9.9')
        self.check_negative(
                lambda: self.relman('--check-current'),
                'Current version 0.9.9 on branch master has non-0 patch version')

        # revert version change and remove tag
        self.runcmd('git', 'reset', 'HEAD', '--hard')
        self.runcmd('git', 'tag', '--delete', 'v0.9.9')

        # create change, commit, release tag, and release branch
        self.relman('--create-changelog', '--commit', '--tag', '--branch')

        # negative test: set non-0 patch version on main branch
        self.set_version('1.0.1')
        self.check_negative(
                lambda: self.relman('--check-current'),
                'Current version 1.0.1 on branch master has non-0 patch version')

        # revert version change and switch to release branch
        self.runcmd('git', 'reset', 'HEAD', '--hard')
        self.runcmd('git', 'checkout', 'v1.0-dev')

        self.relman('--check-tags')

        # negative test: set mismatching version on release branch
        self.set_version('1.1.0')
        self.check_negative(
                lambda: self.relman('--check-current'),
                'Current version 1.1.0 on branch v1.0-dev has mismatching major.minor version')

        # revert version change
        self.runcmd('git', 'reset', 'HEAD', '--hard')

        # negative test: use --branch on release branch
        self.check_negative(
                lambda: self.relman('--branch'),
                'Cannot create a release branch off of branch v1.0-dev')

        # bump patch version of root project
        self.set_version('1.0.1')
        self.commit('Bumping patch version of root project')
        self.relman('--create-changelog', '--commit', '--tag')

    def testShowVersion(self):
        # create change, commit, release tag, and release branch
        self.relman('--create-changelog', '--commit', '--tag', '--branch')

        # check version and semver output
        self.assertEqual('1.0.0', self.relman('--show-head-version').strip())
        self.assertEqual(['1.0.0', '1.0', '1'], self.relman('--show-head-semver').split())

        # switch to release branch
        self.runcmd('git', 'checkout', 'v1.0-dev')

        # bump patch version
        self.set_version('1.0.1')
        self.commit('Bumping patch version')
        self.relman('--create-changelog', '--commit', '--tag')

        # check version and semver output
        self.assertEqual('1.0.1', self.relman('--show-head-version').strip())
        self.assertEqual(['1.0.1', '1.0'], self.relman('--show-head-semver').split())


if __name__ == '__main__':
    unittest.main()
