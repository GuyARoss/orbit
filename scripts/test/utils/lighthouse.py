import subprocess
import json


def audits_url(endpoint):
    output = subprocess.getoutput(
        f'lighthouse {endpoint} --quiet --chrome-flags="--headless" --preset=desktop --output=json'
    )
    json_output = json.loads(output)

    return json_output["audits"]
