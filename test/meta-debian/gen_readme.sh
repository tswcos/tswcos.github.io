#!/bin/bash

THISDIR=$(dirname $(readlink -f "$0"))
cd $THISDIR

echo "# meta-debian testing result" > README.md
for distro in *; do
	if [ -d $distro ]; then
		echo "## $distro" >> README.md
		for i in $distro/*; do
			if [ -d $i ]; then
				machine=`basename $i`
				echo "- [$machine](./$i)" >> README.md
			fi
		done
		echo "" >> README.md
	fi
done
