# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# license file in the root directory of this source tree.

import os


def verify_titlekeyword_doesnotexist() -> bool:
    files = os.listdir("./docs")

    for f in files:
        if ".html" in f:
            with open("./docs/" + f, "r+") as file:
                text = file.read()

                if "$title" in text:
                    return False

    return True


if __name__ == "__main__":
    assert (
        verify_titlekeyword_doesnotexist()
    ), "title keyword was found, please run the 'finalize.py' tool found in the website directory"
