import subprocess
import os
import os.path
import requests

'''
a test of orbit integrations

contracts:
- code that orbit generates is valid and can be ran
- .orbit dist gets created with valid bundles    
- the application runs and the correct set of data is obtained
'''

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
        subprocess.check_output([f'{path}/orbit build --pacname=autogen --auditpage=./page.audit'], shell=True)    

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

def is_application_ran_successfully(path: str) -> bool:
    subprocess.check_output([f'go run {path}/main.go &'], shell=True)

    f = requests.get('http://localhost:3030/')
    assert f.status_code == 200, "status code failure"

    bk_count = len(f.text.split('orbit_bk'))
    print(bk_count)

    return True
    

if __name__ == '__main__':
    current_dir = os.getcwd()

    # TODO: run the orbit build comand before any of these tests get ran (prefer the make cmd for linking)
    path = '../examples/basic-react'
    os.chdir(path)

    tmp_dir = os.getcwd()

    assert is_orbit_gooutput_valid(tmp_dir), "invalid go orbit output"
    assert is_orbit_dist_valid(tmp_dir), "invalid orbit dist"
    assert is_application_ran_successfully(tmp_dir), "application ran"