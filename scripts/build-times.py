# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# license file in the root directory of this source tree.

'''
This script measures the time to build over time for orbit's build command.
'''
import os
import time
import tempfile
import subprocess
from typing import List, Tuple
import matplotlib.pyplot as plt


def main(tags: List[Tuple[str, str]]) -> None:
    init_dir = os.getcwd()
    tag_stats = []
    
    for tag, alias in tags:
        tmpdir = tempfile.mkdtemp()
        print('new tempdir', tmpdir)
        os.chdir(tmpdir)            
        repo_dir = setup_repo(tmpdir, tag)
        os.chdir(repo_dir)
        print(os.getcwd())
        stats = profile_build_cmd(tag)
        stats.set_alias(alias)

        tag_stats.append(stats)
    
    os.chdir(init_dir)    
    plot_from_stats('./build-times.png', tag_stats)

class ProfilerStats:
    tag: str
    runtime_duration: int
    alias: str

    def __init__(self, tag: str, duration: int) -> None:
        self.tag = tag
        self.runtime_duration = duration        

    def set_alias(self, alias: str):
        self.alias = alias
        return self

def plot_from_stats(path: str, stats: List[ProfilerStats]):
    points = []
    lbls = []

    for s in stats:
        points.append(s.runtime_duration)
        if s.alias:
            lbls.append(s.alias)
        else:
            lbls.append(s.tag.replace('commit@', '')[:6] )

    plt.xticks(range(len(points)), lbls)
    plt.ylabel('Duration')
    plt.xlabel('Version No.')
    plt.title('Runtime duration (orbit build cmd)')
    plt.bar(range(len(points)), points) 
    plt.savefig(path)

def profile_build_cmd(tag: str) -> ProfilerStats:
    try:
        subprocess.check_output([f'npm i'], shell=True)
    except:
        raise Exception('npm install failed')
    
    sum_of_times = 0
    for _ in range(5):
        start_time = time.time()    

        try:
            subprocess.check_output([f'./orbit build --pacname=orbitgen --mode=development'], shell=True)
        except:
            raise Exception('build command failed')

        end_time = time.time()
        sum_of_times += end_time - start_time

    return ProfilerStats(tag, sum_of_times / 5)

def setup_repo(dir: str, tag: str) -> str:
    try:        
        if tag == 'latest':
            subprocess.check_output([f'git clone http://github.com/GuyARoss/orbit'], shell=True)
            os.chdir(dir + "/orbit")
        elif "commit@" in tag:
            commit = tag.replace('commit@', '')

            subprocess.check_output(['git clone http://github.com/GuyARoss/orbit'], shell=True)
            os.chdir(dir + "/orbit")
            subprocess.check_output([f'git checkout {commit}'], shell=True)
        else:
            subprocess.check_output([f'git clone --branch {tag} http://github.com/GuyARoss/orbit'], shell=True)
            os.chdir(dir + "/orbit")

    except:
        print('AT CURRENT DIR', dir)
        raise Exception('github clone failed')

    try:        
        subprocess.check_output([f'go build -o ./orbit'], shell=True)
        subprocess.check_output([f'./scripts/link_examples.sh'], shell=True)
    except:
        raise Exception('make example failed')

    return dir + "/orbit/examples/basic-react"

if __name__ == "__main__":
    tags = [
        ('v0.16.0', 'v.16'),
        ('commit@affdd1742be19697ae9f0c693312e118ea33a766',  'Error_Prop'),
        ('commit@a77c1c4a79268acc6e443a7682a7ad156f79fda4',  'lighthouse'),
        ('latest', 'main'),
        ('regression/feb-21', 'WIP')
    ]

    main(
        tags
    )