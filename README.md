REGGIE - Regex file copier
---

Reggie is a command line utility for searching files via regex and copying them to
a destination directory. Need to run a script on multiple operating systems that
are running different versions of find? Reggie don't care. Have spaces or weird
characters in your file path? Reggie don't care. Too lazy to look up the syntax
for `find -exec` or `xargs`? Reggie don't care.

Usage
=====
`reggie "(jpg|mov)" ~/mystuff /mnt/backup/mystuff`

Options
=======
dry-run: Don't actually copy the files

verbose: Print the names of all files found

ex:

    reggie --dry-run --verbose "jpg" ~/mystuff /mnt/backup/mystuff

