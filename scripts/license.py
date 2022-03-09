# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# LICENSE file in the root directory of this source tree.
import os

statement = [
" Copyright (c) 2021 Guy A. Ross",
" This source code is licensed under the GNU GPLv3 found in the",
" LICENSE file in the root directory of this source tree.",
]

def original_txt(path):
    original_file = open(path, "r")
    f = original_file.read()

    original_file.close()

    return f   

def write_license_comment(comment_str):
    original = original_txt(filepath)

    f = open(filepath, "w")
    for s in statement:
        f.write(comment_str + s + "\n")

    f.write(original) 
    f.close()

AUDIT_FILE = "./scripts/license_audit.txt"

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

        if any([d in filepath for d in blacklist_dirkeys]) or filepath in audit_lines:
            continue

        extensions = filepath.split(".")
        file_extension = extensions[len(extensions) - 1]

        if file_extension == "go":
            audit.write(filepath + "\n")
            write_license_comment("//")
        elif file_extension == "sh" or file_extension == "py":            
            audit.write(filepath + "\n")
            write_license_comment("#")

audit.write(previous_audit)
audit.close()