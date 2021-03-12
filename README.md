# prettycrontab

The goal is to show the user which things are about to be executed via cron, in
a more or less pretty way, and allow the user to deem entries "UNINTERESTING",
in which case they won't be listed.

A sample use-case is given in the `i3blocks-prettycrontab` script, which uses
this script to display the next things that are about to be ran via cron, with
proper label.

Basic usage:

    crontab -l | prettycrontab

See also:

    prettycrontab --help

## How to get it

    go get github.com/mfontani/prettycrontab

## Usage

Given a bog-standard `crontab -l` containing one entry:

    * * * * * /usr/local/bin/foo --bar >/dev/null 2>&1

... it'll display a pretty-printed version of what is going to happen when
(specifically, on the next minute I ran this command on). Note the various
"empty" output redirections are removed from the actual output, as they're
"noise".

    echo '* * * * * /usr/local/bin/foo --bar >/dev/null 2>&1' | TZ=UTC ./prettycrontab
    2021-03-12 09:58:00 +0000 UTC   /usr/local/bin/foo --bar

The `-redir` option keeps the redirection:

    echo '* * * * * /usr/local/bin/foo --bar >/dev/null 2>&1' | TZ=UTC ./prettycrontab -redir
    2021-03-12 10:02:00 +0000 UTC   /usr/local/bin/foo --bar >/dev/null 2>&1

Non-empty output redirections are otherwise kept by default, as one might want
to know/be remembered whether the action logs to a file and whatnot:

    echo '* * * * * /usr/local/bin/foo --bar >>/var/log/bar 2>&1' | TZ=UTC ./prettycrontab
    2021-03-12 10:02:00 +0000 UTC   /usr/local/bin/foo --bar >>/var/log/bar

The `-deltahms` option displays a pretty-printed version of how long until the
command gets executed:

    echo '*/2 * * * * /usr/local/bin/foo --bar >/dev/null 2>&1' | TZ=UTC ./prettycrontab -deltahms
    1m 41s  2021-03-12 10:04:00 +0000 UTC   /usr/local/bin/foo --bar

The `-deltacoarse` displays a more coarse version of the same, i.e. compare:

    echo '* */3 * * * /usr/local/bin/foo --bar >/dev/null 2>&1' | TZ=UTC ./prettycrontab -deltacoarse
    an hour 2021-03-12 12:00:00 +0000 UTC   /usr/local/bin/foo --bar
    echo '* */3 * * * /usr/local/bin/foo --bar >/dev/null 2>&1' | TZ=UTC ./prettycrontab -deltahms
    1h 56m 21s      2021-03-12 12:00:00 +0000 UTC   /usr/local/bin/foo --bar

The timestamp can be removed via `-timestamp=false`:

    echo '* */3 * * * /usr/local/bin/foo --bar >/dev/null 2>&1' | TZ=UTC ./prettycrontab -deltahms -timestamp=false
    1h 55m 16s      /usr/local/bin/foo --bar

You can read `prettycrontab --help` and play with the options.

## Advanced use

The crontab parser handles some comment blocks in a "special" way, and this is
what gives this program the ability to have a much neater/pretty output.

The "parser" works over two things:

### Single cron lines

If a _single cron line_ is preceded by a `## LABEL` entry, i.e.:

    ## LABEL foo bar baz
    * */3 * * * /usr/local/bin/foo --bar >/dev/null 2>&1

... *what comes after the `LABEL`* will be used to "describe" the entry, i.e.:

    printf '## LABEL foo bar\n* */3 * * * /usr/local/bin/foo --bar >/dev/null 2>&1\n' | TZ=UTC ./prettycrontab -deltahms -timestamp=false
    1h 49m 41s      foo bar

Using this method, one can both have a complex command in cron, as well as have
it show up with a neater description when using this program, all at the
expense of a single "double-comment" line.

### Newline-separated blocks

If a newline-separarted block is preceded by a `## UNINTERESTING` entry, it
will not be output by `prettycrontab`.

    ## UNINTERESTING
    * * * * * /usr/local/bin/foo --bar >/dev/null 2>&1 # not shown due to previous
    * * * * * /usr/local/bin/bar --foo >/dev/null 2>&1 # still not shown

The effect of `## UNINTERESTING` cease after a double-newline, i.e:

    printf '## UNINTERESTING\n* * * * * foo --bar\n* * * * * bar --quux\n\n* */3 * * * /usr/local/bin/foo --bar >/dev/null 2>&1\n' | TZ=UTC ./prettycrontab -deltahms -timestamp=false
    1h 45m 19s      /usr/local/bin/foo --bar

## Examples

Here are some examples taken from my own crontab:

    # This runs on my mini pc, but I'm tracking it on the laptop so I can
    # remember that the network will be choppy once this runs
    ## LABEL hourly speedtest
    9     *  *   *   *     true >/dev/null 2>&1

    # Keep track of which packages are installed on this laptop.
    ## UNINTERESTING
    14    *  *   *   *     /home/marco/.local/bin/installed-packages-update.sh >/dev/null 2>&1

    # Fetch email during the working day every few mins
    ## UNINTERESTING
    */2 8-20 *   *   *     /bin/bash -c '( date ; PATH=$PATH:/home/marco/.local/bin /home/marco/.local/bin/gmi-foo pull ) >>/home/marco/tmp/gmi-foo.log 2>&1'
    */2 8-20 *   *   *     /bin/bash -c '( date ; PATH=$PATH:/home/marco/.local/bin /home/marco/.local/bin/gmi-bar pull ) >>/home/marco/tmp/gmi-bar.log 2>&1'
    */2 8-20 *   *   *     /bin/bash -c '( date ; PATH=$PATH:/home/marco/.local/bin /home/marco/.local/bin/gmi-baz pull ) >>/home/marco/tmp/gmi-baz.log 2>&1'

    # Hourly duplicacy backup.
    ## LABEL duplicacy laptop
    46    *  *   *   *     . /home/marco/.duplicacy/env.sh && /home/marco/.local/bin/laptop-backup >/dev/null 2>&1
    # Hourly borg backup. This logs on its own.
    ## LABEL borg laptop
    11    *     *   *   *     /home/marco/.local/bin/borg-backup--borg-laptop cron >/dev/null 2>&1

## Copyright and License

`prettycrontab` is Copyright (c) 2021, Marco Fontani <MFONTANI@cpan.org>

It is released under the MIT license - see the `LICENSE` file in this repository/directory.
