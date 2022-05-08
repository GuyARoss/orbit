'''
This script measures the time to build over time for orbit's build command.
'''
import os
import time
import tempfile
import subprocess
from typing import List
import matplotlib.pyplot as plt


def main(tags: List[str]) -> None:
    init_dir = os.getcwd()
    tag_stats = []
    
    for tag in tags:
        tmpdir = tempfile.mkdtemp()
        print('new tempdir', tmpdir)
        os.chdir(tmpdir)            
        repo_dir = setup_repo(tmpdir, tag)
        os.chdir(repo_dir)
        print(os.getcwd())
        stats = profile_build_cmd(tag)
        tag_stats.append(stats)
    
    os.chdir(init_dir)    
    plot_from_stats('./build-times.png', tag_stats)

class ProfilerStats:
    tag: str
    runtime_duration: int

    def __init__(self, tag: str, duration: int) -> None:
        self.tag = tag
        self.runtime_duration = duration        

def plot_from_stats(path: str, stats: List[ProfilerStats]):
    points = []
    lbls = []

    for s in stats:
        points.append(s.runtime_duration)
        lbls.append(s.tag)

    plt.xticks(range(len(points)), [lbl.replace('commit@', '')[:6] for lbl in lbls])
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
            subprocess.check_output([f'./orbit build --pacname=orbitgen'], shell=True)
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
    tags = ['v0.2.0', 'v0.3.0', 'v0.3.2', 'v0.3.6', 'v0.7.0', 'v0.7.1', 'latest']

    # tags = [
    #     'latest',
    #     'commit@6a131ded2e281846a1ca71d87a41ee14c30bcdfa',
    #     'commit@9184fce235309eb26a0451d4facbecd1aa3566bb',
    #     'commit@e5f0835c519cf400c2fd4bbd41f8e5a30fb1b09a',
    #     'commit@4cbfbbb0b4deaf770274b65591bd310f418791d9',
    #     'commit@571d3ced6cec838c622a47201be5420d3ff0ee16',
    #     'commit@1cbf8146637b88c37e8e93338051325ea2077f00',
    #     'commit@4e52e3ce9cd5b0689d00d395caa182afe983debd']
    # tags.reverse()
    main(
        tags
    )