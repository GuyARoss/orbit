import subprocess
'''
a test of orbit integrations

contracts:
- code that orbit generates is valid and can be ran
- .orbit dist gets created with valid bundles    
- the application runs and the correct set of data is obtained
'''

def is_orbit_gooutput_valid() -> bool:
    '''
        is_orbit_gooutput_valid determines if we can successfully run the go build cmd        
    '''
    path = '../examples/basic-react/main.go'
    
    try:
        subprocess.check_output([f'go build {path}'], shell=True)
        return True
    except:
        return False
    
def is_orbit_dist_valid() -> bool:
    '''
        is_orbit_dist_valid determines if the dist directory gets computed correctly after 
    '''
    path = '../examples/basic-react/.orbit/dist'

    

    pass

def is_application_ran_successfully() -> bool:
    pass

if __name__ == '__main__':
    # TODO: run the orbit build comand before any of these tests get ran (prefer the make cmd for linking)
    assert is_orbit_gooutput_valid(), "invalid go orbit output"