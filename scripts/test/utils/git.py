# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# license file in the root directory of this source tree.

import subprocess


def clone_repo(clone_path: str, tag: str) -> str:
    if tag == "latest":
        subprocess.getoutputs([f"git clone {clone_path}"])
    elif "commit@" in tag:
        commit = tag.replace("commit@", "")

        subprocess.getoutput([f"git clone {clone_path}"])

        subprocess.getoutput([f"git checkout {commit}"])
    else:
        subprocess.getoutput([f"git clone --branch {tag} {clone_path}"])

    path_split = clone_path.split("/")
    return path_split[len(path_split) - 1]


def version_abbrev(abbrev: str):
    return subprocess.getoutput(f"git describe --abbrev={abbrev}")
