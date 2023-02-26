# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# license file in the root directory of this source tree.

import signal
import subprocess
import os


def terminate_pid_by_port(port: int) -> str:
    netstat = subprocess.getoutput(f"netstat -nlp | grep {port}")

    if str(port) in netstat:
        p = netstat.split("/main")
        pid = pull_number_from_last(p[len(p) - 2])

        os.kill(int(pid), signal.SIGTERM)


def pull_number_from_last(text: str) -> int:
    t = ""
    for n in range(len(text)):
        if text[len(text) - n - 1].isnumeric():
            t += text[len(text) - n - 1]
        else:
            return t[::-1]

    return t[::-1]
