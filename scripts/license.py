# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# license file in the root directory of this source tree.

import os
import sys

statement = [
    " Copyright (c) 2021 Guy A. Ross",
    " This source code is licensed under the GNU GPLv3 found in the",
    " license file in the root directory of this source tree.",
]


def original_txt(path):
    original_file = open(path, "r")
    f = original_file.read()

    original_file.close()

    return f


def write_license_comment(comment_str, filepath, offset=0):
    original = original_txt(filepath).split("\n")

    f = open(filepath, "w")
    for s in statement:
        f.write(comment_str + s + "\n")

    f.write("\n")

    f.write("\n".join(original[offset:]))
    f.close()


AUDIT_FILE = "./scripts/license_audit.txt"


def update_all(prev_linecount):
    previous_audit = original_txt(AUDIT_FILE)
    audit_lines = previous_audit.split()

    for l in audit_lines:
        extensions = l.split(".")
        file_extension = extensions[len(extensions) - 1]

        if file_extension == "go":
            write_license_comment("//", l, prev_linecount)
        elif file_extension == "sh" or file_extension == "py":
            write_license_comment("#", l, prev_linecount)


def write_all():
    previous_audit = original_txt(AUDIT_FILE)
    audit_lines = previous_audit.split()

    audit = open(AUDIT_FILE, "w")

    # blacklist_dirkeys, provides a blacklist for the following:
    # - the output of the embded directory (we don't mind what license the users the tools output uses)
    # - examples, as the examples can be used how ever
    blacklist_dirkeys = ["/examples/", "/embed/"]

    for path, subdirs, files in os.walk("./"):
        files.append(path)

        for name in files:
            if name != path:
                filepath = os.path.join(path, name)
            else:
                filepath = name

            if (
                any([d in filepath for d in blacklist_dirkeys])
                or filepath in audit_lines
            ):
                continue

            extensions = filepath.split(".")
            file_extension = extensions[len(extensions) - 1]

            if file_extension == "go":
                audit.write(filepath + "\n")
                write_license_comment("//", filepath)
            elif file_extension == "sh" or file_extension == "py":
                audit.write(filepath + "\n")
                write_license_comment("#", filepath)

    audit.write(previous_audit)
    audit.close()


if __name__ == "__main__":
    if sys.argv[1] == "update":
        # you may want to change this value depending how many lines the license takes up.
        update_all(len(statement))

    if sys.argv[1] == "write":
        write_all()
