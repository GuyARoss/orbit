'''
a test of orbit integrations

contracts:
- code that orbit generates is valid and can be ran
- .orbit dist gets created with valid bundles    
- the application runs and the correct set of data is obtained
'''
import subprocess
import os
import os.path
from time import sleep
import requests
import signal


def is_orbit_gooutput_valid(path: str) -> bool:
    '''
        is_orbit_gooutput_valid returns true if the project is successfully ran
    '''        
    try:
        subprocess.check_output([f'go build {path}/main.go'], shell=True)
        return True
    except:
        return False
    
def is_orbit_dist_valid(path: str) -> bool:
    '''
        is_orbit_dist_valid returns true if the dist directory gets computed correctly 
    '''    
    try:
        subprocess.check_output([f'{path}/orbit build --pacname=orbitgen --auditpage=./page.audit'], shell=True)    

        # read the page audit
        f = open("./page.audit","r")
        lines = f.readlines()

        assert lines[0] == "audit: components\n", "audit file should be component audit"

        for i in lines[1:]:
            bundle = i.split(' ')[1].strip()

            assert os.path.isfile(f'{path}/.orbit/dist/{bundle}.js'), f"bundle file was not created {path}/.orbit/dist/{bundle}.js"

        return True
    except:
        return False

def pull_number_from_last(text: str) -> int:
    t = ""
    for n in range(len(text)):        
        if text[len(text) - n -1].isnumeric():
            t += text[len(text) - n -1]
        else:
            return t[::-1]

    return t[::-1]
    
def terminate_port_pid(port: int) -> str:
    netstat = subprocess.getoutput(f"netstat -nlp | grep {port}")

    if "3030" in netstat:
        p = netstat.split('/main')
        pid = pull_number_from_last(p[len(p) - 2])
        
        os.kill(int(pid), signal.SIGTERM)
        
def is_application_ran_successfully(path: str) -> bool:
    try:
        terminate_port_pid(3030)

        p = subprocess.Popen(['go', 'run', f'{path}/main.go'])
        sleep(5)

        f = requests.get('http://localhost:3030/')
        assert f.status_code == 200, "status code failure"

        bk_count = len(f.text.split('orbit_bk'))
        
        p.terminate()
        terminate_port_pid(3030)

        return True
    except:
        return False


if __name__ == '__main__':    
    current_dir = os.getcwd()

    # TODO: run the orbit build comand before any of these tests get ran (prefer the make cmd for linking)
    path = './examples/basic-react'
    os.chdir(path)

    tmp_dir = os.getcwd()

    assert is_orbit_gooutput_valid(tmp_dir), "invalid go orbit output"
    assert is_orbit_dist_valid(tmp_dir), "invalid orbit dist"
    assert is_application_ran_successfully(tmp_dir), "application ran"

    print('integration contracts completed successfully')