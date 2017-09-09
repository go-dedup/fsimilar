#!/bin/sh

set -e

rm -rf /tmp/fsimilar.* 
../fsimilar vec -i test1.lst -S -F > /tmp/fsimilar_shell_test.ref
echo $?
diff -U1 shell_test.ref /tmp/fsimilar_shell_test.ref
echo $?

ls shell_*.tmpl.sh | xargs -t -i sh -c 'diff -U1 {} /tmp/fsimilar.*.{}'
ret=$?
echo $ret
exit $ret
