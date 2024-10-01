# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# license file in the root directory of this source tree.

from typing import Dict
import subprocess
import time


def e2e_measure_build_cmd(samples=3) -> Dict[str, int]:
    # TODO if .orbit exists, delete it
    try:
        subprocess.check_output([f"npm i"], shell=True)
    except:
        raise Exception("npm install failed, does node exist on the host machine?")

    sum_of_times = 0
    failure_rate = 0
    for _ in range(samples):
        start_time = time.time()
        try:
            res = subprocess.check_output([f"./orbit build --package_name=orbitgen"], shell=True)
            if b"unknown flag" in res:
                res = subprocess.check_output([f"./orbit build --pacname=orbitgen"], shell=True)

        except:            
            failure_rate += 1

        end_time = time.time()
        sum_of_times += end_time - start_time

    return {
        "failure_rate": failure_rate,
        "avg_build_time": sum_of_times / samples,
    }


def link_examples():
    try:
        subprocess.check_output([f"go build -o ./orbit"], shell=True)
        subprocess.check_output([f"./scripts/link_examples.sh"], shell=True)
    except:
        raise Exception("make example failed")
