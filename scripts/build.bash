#!/bin/bash

for target in runner; do
	echo "Build> matryoshka/${target}"
	docker build -t matryoshka/$target \
		-f deployments/docker/$target.Dockerfile .
	ret=$?
	if [ $ret -ne 0 ]; then
		echo "failed to build matryoshka/${target}" >&2
		exit $ret
	fi
done

echo "\nBuilding language images..."
for target in c golang; do
	echo "Build> matryoshka/${target}"
	docker build -t matryoshka/$target \
		-f deployments/docker/languages/$target.Dockerfile .
	ret=$?
	if [ $ret -ne 0 ]; then
		echo "failed to build matryoshka/${target}" >&2
		exit $ret
	fi
done

