# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# license file in the root directory of this source tree.

import os

files = os.listdir("./docs")

for f in files:
    if ".html" in f:
        with open("./docs/" + f, "r+") as file:
            text = file.read()

            text = text.replace("$title", "Orbit - " + f.replace(".html", ""))

            file.seek(0)
            file.write(text)
            file.truncate()
