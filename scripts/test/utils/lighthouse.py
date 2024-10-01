# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# license file in the root directory of this source tree.

import subprocess
import json

def audits_url(endpoint):    
    output = subprocess.getoutput(
        f'lighthouse {endpoint} --quiet --chrome-flags="--headless" --preset=desktop --output=json'
    )
    json_output = json.loads(output)

    return json_output["audits"]
