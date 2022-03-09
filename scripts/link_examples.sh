# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# LICENSE file in the root directory of this source tree.
echo "\n"

for file in examples/*; do
    rm -rf "$file"/orbit

    cp orbit "$file"/orbit

    echo "linked $file"
done
