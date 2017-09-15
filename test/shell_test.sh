#!/bin/sh

set -e

# test --ext first
rm -rf /tmp/fsimilar.*
../fsimilar vec -i test1.lst -S -F -e .mkv > /tmp/fsimilar_shell_test_mkv.ref
echo $?
diff -U1 shell_test_mkv.ref /tmp/fsimilar_shell_test_mkv.ref
echo $?
diff -U1 shell_ln.tmpl_mkv.sh /tmp/fsimilar.*.shell_ln.tmpl.sh
echo $?

# test regular next
rm -rf /tmp/fsimilar.*
../fsimilar vec -i test1.lst -S -F > /tmp/fsimilar_shell_test.ref
echo $?
diff -U1 shell_test.ref /tmp/fsimilar_shell_test.ref
echo $?

ls shell_*.tmpl.sh | xargs -t -i sh -c 'diff -U1 {} /tmp/fsimilar.*.{}'
ret=$?
echo $ret
exit $ret
