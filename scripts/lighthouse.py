import json
import subprocess

important_heuristics = [
    'first-contentful-paint',
    'speed-index',
    'largest-contentful-paint',
    'server-response-time'
]

def format_heuristic(d):
    return f"{d['numericValue']} {d['numericUnit']}"

def main():
    # TODO: make sure that the service is running. 
    print('started, please note that the server needs to be started before running this script.')

    output = subprocess.getoutput('lighthouse http://localhost:3030 --quiet --chrome-flags="--headless" --preset=desktop --output=json')
    json_output = json.loads(output)
    
    for h in json_output['audits'].keys():
        s = json_output['audits'][h]['score']
        if s and s < .90:
            print(f"[failed] '{h}': {s}")
            
            assert h not in important_heuristics, f"important heuristic failed {h}: {s}"
        else:
            print(f"[passed] '{h}': {s}")


if __name__ == "__main__":
    main()