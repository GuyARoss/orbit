import os
import time
import tempfile
import subprocess
from datetime import datetime
from typing import List
from requests.auth import HTTPBasicAuth
'''
This script measures the time to build over time for orbit's build command.
'''

def main() -> None:
    tags = ['v0.7.1', 'v0.7.2', 'v0.3.6', 'v0.3.2', 'v0.3.0', 'v0.2.0']

    tag_stats = []
    for tag in tags[0:1]:
        tmpdir = tempfile.mkdtemp()
        print('new tempdir', tmpdir)
        os.chdir(tmpdir)            
        repo_dir = setup_repo(tmpdir, tag)
        os.chdir(repo_dir)
        print(os.getcwd())
        stats = profile_build_cmd(tag)
        tag_stats.append(stats)
        
    print(tag_stats)
    # TODO: sort tag stats
    # 5 - display in matplotlib    

class ProfilerStats:
    tag: str
    runtime_duration: int

    def __init__(self, tag: str, duration: int) -> None:
        self.version_tag = tag
        self.runtime_duration = duration        

def profile_build_cmd(tag: str) -> ProfilerStats:
    try:
        subprocess.check_output([f'npm i'], shell=True)
    except:
        raise Exception('npm install failed')
    
    start_time = time.time()    

    try:
        subprocess.check_output([f'./orbit build --pacname=orbitgen'], shell=True)
    except:
        raise Exception('build command failed')

    end_time = time.time()

    return ProfilerStats(tag, end_time - start_time)

def setup_repo(dir: str, tag: str) -> str:
    try:
        subprocess.check_output([f'git clone --branch {tag} http://github.com/GuyARoss/orbit'], shell=True)
    except:
        raise Exception('github clone failed')

    try:
        subprocess.check_output([f'go build -o ./orbit'], shell=True)
    except:
        raise Exception('make example failed')

    return dir + "/orbit/examples/basic-react"

def create_working_dir():    
    pass

if __name__ == "__main__":
    main()