#!/usr/bin/env python

import argparse
import datetime
import logging
import os
import re
import subprocess
import sys

LOGGER = logging.Logger(__name__)
LOGGER.addHandler(logging.StreamHandler(sys.stdout))
LOGGER.setLevel(logging.INFO)

def get_abspath(*from_root_dir):
    """
    Get absolute path of file that is specified relative to the root of the Git repository.
    """

    abs_file = os.path.abspath(__file__)
    # skip the directory that contains relman.py, which cannot be root of Git repository
    abs_dir = os.path.dirname(abs_file)
    git_dir = os.path.dirname(abs_dir)
    # traverse up directory tree until .git directory is found
    while not os.path.isdir(os.path.join(git_dir, '.git')):
        parent_dir = os.path.dirname(git_dir)
        # make sure we do not go up to root directory
        if git_dir == parent_dir:
            raise RuntimeError(abs_file + ' was not invoked within a Git repository')
        git_dir = parent_dir
    return os.path.join(git_dir, *from_root_dir)

def run(*args, **kwargs):
    """
    Execute a command and returns standard output.
    """

    LOGGER.debug('> %s', subprocess.list2cmdline(args))
    if 'cwd' not in kwargs:
        kwargs['cwd'] = get_abspath('.')
    out = subprocess.check_output(list(args), **kwargs).decode('utf-8')
    LOGGER.debug('---\n%s\n---', out)
    return out


class GitRepo(object):

    def __init__(self, path):
        self.releasable = Releasable(path)

    @classmethod
    def get_tags(cls, ref):
        return run('git', 'tag', '--sort=committerdate', '--merged', ref).split()

    @classmethod
    def get_tags_on(cls, ref='HEAD'):
        return run('git', 'tag', '--points-at', ref).split()

    @classmethod
    def get_current_branch(cls):
        return run('git', 'rev-parse', '--abbrev-ref', 'HEAD').strip()

    @classmethod
    def get_commit_range(cls, ref1, ref2, *files):
        refs = ref2 if not ref1 else '{}..{}'.format(ref1, ref2)
        return run('git', 'log', '--pretty=format:%h %s', refs, '--', *files).splitlines()

    @classmethod
    def get_uncommitted_changes(cls, *files):
        return run('git', 'status', '--porcelain', '--', *files).strip()

    @classmethod
    def has_uncommitted_changes(cls, *files):
        return cls.get_uncommitted_changes(*files) != ''

    @classmethod
    def stage(cls, *files):
        return run('git', 'add', '--', *files)

    @classmethod
    def commit(cls, msg):
        return run('git', 'commit', '-m', msg)

    @classmethod
    def tag(cls, tag):
        return run('git', 'tag', tag)

    @classmethod
    def get_branches(cls):
        for branch in run('git', 'branch').split():
            if branch != '*':
                yield branch

    @classmethod
    def get_remote_urls(cls, remote='origin'):
        try:
            return run('git', 'remote', 'get-url', remote, '--all', stderr=subprocess.STDOUT).split()
        except subprocess.CalledProcessError as e:
            LOGGER.debug('---\n%s\n---', e.output.decode('utf-8'))
            return []

    @classmethod
    def get_default_branch(cls):
        try:
            ret = run('git', 'config', '--get', 'init.defaultBranch').strip()
            if ret != '':
                return ret
        except subprocess.CalledProcessError:
            pass
        # if main exists, assume it is the default branch
        if 'main' in cls.get_branches():
            return 'main'
        # otherwise, assume master
        return 'master'

    @classmethod
    def create_branch(cls, branch):
        return run('git', 'branch', branch)

    def check_current(self, check_all=False):
        """
        Check that the current version is consistent with the most recent
        release tag. If check_all is True, check that all release tags are
        consistent. All release tags must be in ascending order and must follow
        conventions described in README.md.
        """

        release_branch = self.releasable.get_release_branch()
        previous = None
        LOGGER.debug('Checking release tags for branch %s', release_branch.branch)
        # check that current version matches branch conventions; main cannot have non-0 patch releases, and release branches cannot have non-matching branches
        release_branch.check_current_version(self.releasable)
        for release, is_head in self.releasable.get_releases():
            # check that the current version is larger than tag, or identical to it if HEAD is tagged
            release.check_version(is_head)
            if check_all:
                # check that tag matches branch conventions
                release_branch.check_release(release)
                # check that release tags are in ascending order
                if previous is not None:
                    previous.check_before(release)
                previous = release

    def show_head_versions(self, semver=False):
        for tag in self.get_tags_on():
            release = self.releasable.get_release(tag)
            if release is not None:
                sys.stdout.write(release.version + '\n')
                # this is also the most recent release of <major>.<minor>
                if semver:
                    sys.stdout.write('{}.{}\n'.format(release.semver[0], release.semver[1]))
                    # if we are on main, then this is also the most recent
                    # release of <major>
                    current_branch = GitRepo.get_current_branch()
                    if current_branch == GitRepo.get_default_branch():
                        sys.stdout.write('{}\n'.format(release.semver[0]))


    def show_changelogs(self):
        release_desc = self.releasable.version
        if self.releasable.tag_prefix != 'v':
            release_desc += ' of {}'.format(self.releasable.path)
        LOGGER.info('Generating changelog for version %s:\n%s', release_desc, self.releasable.get_changelog())

    def create_changelogs(self, commit=False):
        self.releasable.create_changelog(commit)

    def tag_releases(self):
        # check current version
        self.check_current()

        # make sure there are no uncommitted changes
        uncommitted_changes = self.get_uncommitted_changes()
        if uncommitted_changes.strip() != '':
            raise RuntimeError('Cannot create release tag because there are uncommitted changes:\n' + uncommitted_changes)

        # make sure we are on main or a release branch
        release_branch = self.releasable.get_release_branch()
        if release_branch.branch != self.get_default_branch() and release_branch.semver is None:
            raise RuntimeError('Cannot create release tag on branch ' + release_branch.branch)

        # tag HEAD with the release tag
        GitRepo.tag(self.releasable.tag_prefix + self.releasable.version)

    def create_branches(self):
        # check current version
        self.check_current()

        # make sure we are on main branch
        current_branch = GitRepo.get_current_branch()
        if current_branch != GitRepo.get_default_branch():
            raise RuntimeError('Cannot create a release branch off of branch ' + current_branch)

        # get latest release
        release = None
        on_head = False
        for release, on_head in self.releasable.get_releases():
            pass
        # make sure release exists and is on HEAD
        if release is None:
            raise RuntimeError('No release for "{}"'.format(self.releasable.path))
        if not on_head:
            raise RuntimeError('Latest release tag {} is not on HEAD'.format(release.tag))
        # make sure patch version is 0
        if release.semver[-1] != 0:
            raise RuntimeError('Cannot create release branch from non-0 patch release tag ' + release.tag)
        # create release branch
        GitRepo.create_branch('{}{}.{}-dev'.format(self.releasable.tag_prefix, *release.semver[:-1]))


class Changelog(object):
    """
    Helper class to create a changelog from the last release on the current
    branch. It is assumed that the release will be tagged with the current
    development version for the releasable.
    """

    def __init__(self, releasable):
        self.releasable = releasable

    RE_REPO_URL = r'^(git@|https://)(www\.|)github.com[:/](.*)\.git$'

    @classmethod
    def extract_repo_url(cls):
        # if environment variable was specified, use it
        if 'REPO_URL' in os.environ:
            return os.environ['REPO_URL']
        # otherwise, extract URL from remote
        for url in GitRepo.get_remote_urls():
            m = re.match(cls.RE_REPO_URL, url)
            if m:
                return 'https://{}github.com/{}'.format(m.group(2), m.group(3))
        # if there are no remotes, just return a placeholder
        return 'https://github.com/<org>/<repo>'

    @classmethod
    def get_repo_url(cls):
        if not hasattr(cls, 'REPO_URL'):
            cls.REPO_URL = cls.extract_repo_url()
        return cls.REPO_URL

    RE_COMMIT_MSG = r'(.*) \(#([0-9]+)\)$'
    CHANGELOG_ENTRY_FMT_BASIC = '- [`{1}`]({0}/commit/{1}) {2}'
    CHANGELOG_ENTRY_FMT = CHANGELOG_ENTRY_FMT_BASIC + ' [\\#{3}]({0}/pull/{3})'
    CHANGELOG_FMT = """
# Changelog [{current}]({repo_url}/tree/{current}) ({date})

## [Full Changelog]({repo_url}/compare/{previous}...{current})

{changelog_entries}
"""

    def create(self):
        """
        Return the formatted changelog content for the current version that is
        being prepared for release.
        """

        # inject values into format string for changelog
        repo_url = self.get_repo_url()
        current = self.releasable.tag_prefix + self.releasable.version
        previous = ''
        for release, _ in self.releasable.get_releases():
            previous = release.tag

        # do not create changelog if previous and current tag are the same
        if previous == current:
            raise RuntimeError('Cannot create changelog because previous and current tag are both ' + previous)

        changelog_entries = '\n'.join(self.get_changelog_entries(previous, 'HEAD', self.releasable.path))
        date = datetime.date.today().strftime('%Y-%m-%d')
        return self.CHANGELOG_FMT.format(repo_url=repo_url, current=current, previous=previous, changelog_entries=changelog_entries, date=date)

    @classmethod
    def get_changelog_entries(cls, ref1, ref2, *files):
        changelog_entries = []
        for commit in GitRepo.get_commit_range(ref1, ref2, *files):
            # format is <sha> <msg>
            sha, msg = commit.split(' ', 1)
            # PRs that are merged using the GitHub UI append (#<PR number>) to the first line of the commit message
            m = re.match(cls.RE_COMMIT_MSG, msg)
            if m:
                msg = m.group(1)
                pr = m.group(2)
                changelog_entries.append(cls.CHANGELOG_ENTRY_FMT.format(cls.get_repo_url(), sha, msg, pr))
            else:
                changelog_entries.append(cls.CHANGELOG_ENTRY_FMT_BASIC.format(cls.get_repo_url(), sha, msg))
        return changelog_entries


class Releasable(object):
    """
    Class that represents a product that is releasable from code within this
    Git repository. The convention for releasables is that the commit that a
    release will be created from is tagged with
    "[directory/]v<major>.<minor>.<patch>[-comment]". The directory should
    contain a `get-version.sh` that returns the current version in development.
    The comment is disregarded.

    Releases with patch version 0 are created from main, and a release branch
    with the same major and minor version and the format
    "[directory/]v<major>.<minor>-dev" should also be branched off of main at
    this commit. All releases on release branches must have the same major and
    minor version.
    """

    @classmethod
    def get_tag_prefix(cls, path):
        if not path or path == '.' or path == '/':
            return 'v'
        abspath = get_abspath(path)
        if not os.path.exists(abspath):
            raise RuntimeError('Directory "{}" does not exist'.format(abspath))
        return path.rstrip('/') + '/v'

    def __init__(self, path):
        self.tag_prefix = self.get_tag_prefix(path)
        self.path = os.path.join(*path.split('/'))
        self.version = run(get_abspath(self.path, 'get-version.sh')).strip()
        self.semver = self.get_semver(self.version)

    SEMVER_PATTERN = r'([0-9]+)\.([0-9]+)\.([0-9]+)'

    @classmethod
    def get_semver(cls, version):
        m = re.match(cls.SEMVER_PATTERN, version)
        if m:
            return int(m.group(1)), int(m.group(2)), int(m.group(3))

    def get_release(self, tag):
        """
        Return a release described by the supplied tag, or None if the tag does
        not describe a release.
        """

        if tag.startswith(self.tag_prefix):
            version = tag.lstrip(self.tag_prefix)
            semver = self.get_semver(version)
            if semver is not None:
                LOGGER.debug('Parsed version tuple %s from tag %s', semver, tag)
                return ReleaseTag(self, tag, version, semver)

    def get_releases(self, branch=None):
        """
        Get all releases on the supplied branch, or the current branch if it is
        not specified.
        """

        if branch is None:
            branch = GitRepo.get_current_branch()
        head_tags = GitRepo.get_tags_on()
        for tag in GitRepo.get_tags(branch):
            release = self.get_release(tag)
            if release is not None:
                yield release, tag in head_tags

    REL_BRANCH_PATTERN = r'([0-9]+)\.([0-9]+)-dev$'

    def get_release_branch(self, branch=None):
        if branch is None:
            branch = GitRepo.get_current_branch()
        if branch.startswith(self.tag_prefix):
            version = branch.lstrip(self.tag_prefix)
            m = re.match(self.REL_BRANCH_PATTERN, version)
            if m:
                semver = (int(m.group(1)), int(m.group(2)))
                LOGGER.debug('Parsed version tuple %s from branch %s', semver, branch)
                return ReleaseBranch(self, branch, semver)
        return ReleaseBranch(self, branch, None)

    def get_changelog(self):
        return Changelog(self).create()

    def create_changelog(self, commit=False):
        # make sure there are no uncommitted changes to changelog file
        changelog_file = os.path.join(self.path, 'changelogs', 'v{}.md'.format(self.version))
        changelog_file = os.path.abspath(changelog_file)
        if GitRepo.has_uncommitted_changes(changelog_file):
            raise RuntimeError('Cannot write "{}" because it has uncommitted changes'.format(changelog_file))

        # write changelog file
        changelog = self.get_changelog()
        release_desc = self.version
        if self.tag_prefix != 'v':
            release_desc += ' of {}'.format(self.path)
        LOGGER.info('Writing changelog for version %s to file "%s"', release_desc, changelog_file)
        with open(changelog_file, 'w') as f:
            f.write(changelog)

        # if commit=True, commit the changelog
        if commit:
            # check that there are no other uncommitted changes
            uncommitted_changes = GitRepo.get_uncommitted_changes()
            if GitRepo.get_uncommitted_changes(changelog_file) != uncommitted_changes:
                raise RuntimeError('Unexpected uncommitted changes:\n' + uncommitted_changes)

            # commit the changelog
            GitRepo.stage(changelog_file)
            GitRepo.commit('Create changelog for ' + release_desc)


class ReleaseTag(object):

    def __init__(self, releasable, tag, version, semver):
        self.releasable = releasable
        self.tag = tag
        self.version = version
        self.semver = semver

    def check_version(self, is_head):
        # if tag is on HEAD, make sure that current version is identical to it (not including prefix); otherwise, make sure that current version is at least as large
        if is_head and self.version != self.releasable.version:
            raise RuntimeError('Release tag {} does not match current version {}'.format(self.tag, self.releasable.version))
        elif self.semver > self.releasable.semver:
            raise RuntimeError('Release tag {} has later version than current version {}'.format(self.tag, self.releasable.version))

    def check_before(self, release_after):
        # check that self is less than release tag on later commit
        if self.semver > release_after.semver:
            raise RuntimeError('Release tag {} has later version than subsequent release tag {}'.format(self.tag, release_after.tag))


class ReleaseBranch(object):
    """
    Class that describes a release branch. If the branch name does not contain
    the major and minor components of a semantic version in the format
    described above, then it is assumed to be a branch off of main, and the
    rules for the main branch are applied.
    """

    def __init__(self, releasable, branch, semver):
        self.releasable = releasable
        self.branch = branch
        self.semver = semver
        self.tags_on_main = set()
        if self.semver is not None:
            self.tags_on_main = GitRepo.get_tags(GitRepo.get_default_branch())

    def check_current_version(self, releasable):
        # assume that any branch that does not have format "<prefix><major>.<minor>-dev" is a branch off of main
        if self.semver is None and releasable.semver[-1] != 0:
            raise RuntimeError('Current version {} on branch {} has non-0 patch version'.format(releasable.version, self.branch))
        # do not allow current version to have different major.minor
        if self.semver is not None and releasable.semver[:-1] != self.semver:
            raise RuntimeError('Current version {} on branch {} has mismatching major.minor version'.format(releasable.version, self.branch))

    def check_release(self, release):
        # assume that any branch that does not have format "<prefix><major>.<minor>-dev" is a branch off of main
        if self.semver is None and release.semver[-1] != 0:
            raise RuntimeError('Release tag {} on branch {} has non-0 patch version'.format(release.tag, self.branch))

        if self.semver is not None:
            # do not allow release branch to have different major.minor version tag on it unless it is before the branch was created
            if release.semver[:-1] != self.semver and release.tag not in self.tags_on_main:
                raise RuntimeError('Release tag {} on branch {} has mismatching major.minor version'.format(release.tag, self.branch))
            # require that release branch was created with patch 0
            if release.semver == self.semver + (0,) and release.tag not in self.tags_on_main:
                raise RuntimeError('Release branch {} was not branched with tag {}'.format(self.branch, release.tag))


def main():
    parser = argparse.ArgumentParser(description='Release Manager')
    parser.add_argument('--debug', action='store_true', help='log debug messages')
    parser.add_argument('--path', default='.', help='the path of the releasable')
    parser.add_argument('--check-current', action='store_true', help='check that current version is at least as large as all release tags on the current branch and identical to version in release tag on HEAD, if there is one')
    parser.add_argument('--check-tags', action='store_true', help='check that release tags on the current branch are in ascending order and follow tagging conventions')
    mtx_group = parser.add_mutually_exclusive_group()
    mtx_group.add_argument('--show-head-version', action='store_true', help='if there is a release tag on HEAD, output the version')
    mtx_group.add_argument('--show-head-semver', action='store_true', help='if there is a release tag on HEAD, output the version and any prefixes for which it is the latest')
    mtx_group.add_argument('--show-changelog', action='store_true', help='output changelog from last release tag on current branch for all releasables')
    mtx_group.add_argument('--create-changelog', action='store_true', help='create changelog from last release tag on current branch for all releasables')
    parser.add_argument('--commit', action='store_true', help='commit changelog file created by --create-changelog; has no effect if --create-changelog was not specified')
    parser.add_argument('--tag', action='store_true', help='tag HEAD with next release tag')
    parser.add_argument('--branch', action='store_true', help='create release branch off of HEAD')
    args = parser.parse_args()

    if args.debug:
        LOGGER.setLevel(logging.DEBUG)

    git = GitRepo(args.path)
    if args.check_current or args.check_tags:
        git.check_current(args.check_tags)

    if args.show_head_version or args.show_head_semver:
        git.show_head_versions(args.show_head_semver)
    elif args.create_changelog:
        git.create_changelogs(args.commit)
    elif args.show_changelog:
        git.show_changelogs()

    if args.tag:
        git.tag_releases()

    if args.branch:
        git.create_branches()

if __name__ == '__main__':
    main()
