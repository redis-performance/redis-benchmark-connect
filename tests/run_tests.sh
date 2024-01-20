#!/bin/bash

[[ $VERBOSE == 1 ]] && set -x

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
ROOT=$(cd $HERE/.. && pwd)

#----------------------------------------------------------------------------------------------

help() {
	cat <<-END
		Run flow tests.

		[ARGVARS...] run_tests.sh [--help|help]

		Argument variables:

		TLS=0|1             Run tests with TLS enabled 

	END
}

#----------------------------------------------------------------------------------------------

run_tests() {
	local title="$1"
	if [[ -n $title ]]; then
		printf "Running $title:\n\n"
	fi

	cd $ROOT/tests

	local E=0
	{
		$OP $ROOT/redis_benchmark_connect $TEST_ARGS
		((E |= $?))
	} || true

	return $E
}

#----------------------------------------------------------------------------------------------

[[ $1 == --help || $1 == help ]] && {
	help
	exit 0
}

#----------------------------------------------------------------------------------------------

TLS_KEY=$ROOT/tests/tls/redis.key
TLS_CERT=$ROOT/tests/tls/redis.crt
TLS_CACERT=$ROOT/tests/tls/ca.crt
OSS_STANDALONE=1

TEST_ARGS=" "
[[ $TLS == 1 ]] && TEST_ARGS+=" --port 6380 --certFile $TLS_CERT --certKey $TLS_KEY --caCertFile $TLS_CACERT --tlsSkipVerify --tls"

cd $ROOT/tests

E=0
[[ $OSS_STANDALONE == 1 ]] && {
	(TLS_KEY=$TLS_KEY TLS_CERT=$TLS_CERT TLS_CACERT=$TLS_CACERT TEST_ARGS="${TEST_ARGS}" run_tests "tests on OSS standalone")
	((E |= $?))
} || true

exit $E
