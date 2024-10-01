# Copyright (c) 2021 Guy A. Ross
# This source code is licensed under the GNU GPLv3 found in the
# license file in the root directory of this source tree.

import pytest
import os
import subprocess

def _read_all_lines(path):
    f = open(path, "r")
    lines = f.readlines()
    f.close()

    return lines

class Test_ReactOrbitServer:
    audit_path: str = "./page.audit"

    @pytest.fixture(scope="class")
    def path(self):
        temp_dir = os.getcwd()

        example_path = "./examples/basic-react"
        os.chdir(example_path)

        yield os.getcwd()

        os.chdir(temp_dir)

    @pytest.fixture(autouse=True)
    def orbit_run_on_example(self, path):
        try:
            subprocess.check_output(
                [f"{path}/orbit build --package_name=orbitgen --audit_path={self.audit_path}"], shell=True
            )
        except:
            subprocess.check_output(
                [f"{path}/orbit build --pacname=orbitgen --auditpage={self.audit_path}"], shell=True
            )

    def test_can_compile_autogen(self, path):
        subprocess.check_output([f"go build {path}/main.go"], shell=True)

    def test_is_orbit_dist_valid(self, path):
        lines = _read_all_lines(self.audit_path)

        assert lines[0] == "audit: components\n", "audit file should be component audit"

        for i in lines[1:]:
            bundle = i.split(" ")[1].strip()
            assert os.path.isfile(
                f"{path}/.orbit/dist/{bundle}.js"
            ), f"bundle file was not created {path}/.orbit/dist/{bundle}.js"

class Test_ReactSPA:
    audit_path: str = './page.audit'
    spa_out_dir: str = './dist'

    @pytest.fixture(scope="class")
    def path(self):
        temp_dir = os.getcwd()

        example_path = "./examples/spa"
        os.chdir(example_path)

        yield os.getcwd()

        os.chdir(temp_dir)

    @pytest.fixture(autouse=True)    
    def orbit_run_on_example(self, path):
        subprocess.check_output([f"{path}/orbit build --spa_entry_path=./pages/app.jsx --audit_path={self.audit_path} --spa_out_dir={self.spa_out_dir}"], shell=True)

    def test_is_dist_found(self, path):
        lines = _read_all_lines(self.audit_path)
        assert lines[0] == "audit: components\n", "audit file should be component audit"

        bundle_id = lines[1].split(" ")[1].strip()
        res = os.listdir(self.spa_out_dir)

        assert f"{bundle_id}.js" in res, "bundle not found in the dist"
        assert "index.html" in res, "index not found in the dist"
            
