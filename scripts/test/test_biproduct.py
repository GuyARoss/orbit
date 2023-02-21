import pytest
import os
import subprocess


@pytest.fixture(scope="module")
def path():
    temp_dir = os.getcwd()

    example_path = "./examples/basic-react"
    os.chdir(example_path)

    yield os.getcwd()

    os.chdir(temp_dir)


@pytest.fixture(autouse=True)
def orbit_run_on_example(path):
    subprocess.check_output(
        [f"{path}/orbit build --pacname=orbitgen --auditpage=./page.audit"], shell=True
    )


def test_can_compile_autogen(path):
    subprocess.check_output([f"go build {path}/main.go"], shell=True)


def test_is_orbit_dist_valid(path):
    f = open("./page.audit", "r")
    lines = f.readlines()

    assert lines[0] == "audit: components\n", "audit file should be component audit"

    for i in lines[1:]:
        bundle = i.split(" ")[1].strip()
        assert os.path.isfile(
            f"{path}/.orbit/dist/{bundle}.js"
        ), f"bundle file was not created {path}/.orbit/dist/{bundle}.js"
