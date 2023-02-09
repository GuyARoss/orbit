import os
from typing import Tuple
import json
import pytest
import subprocess
import requests
from time import sleep
from utils import terminate_pid_by_port
from utils.lighthouse import audits_url


@pytest.fixture(scope="module")
def path():
    temp_dir = os.getcwd()
    print("temp directory", temp_dir)

    example_path = "./examples/basic-react"
    os.chdir(example_path)

    yield os.getcwd()

    os.chdir(temp_dir)


@pytest.fixture(scope="module")
def request_endpoint() -> Tuple[str, int]:
    return ("http://localhost:3030", 3030)


@pytest.fixture(scope="module", autouse=True)
def application(path, request_endpoint):
    terminate_pid_by_port(request_endpoint[1])

    p = subprocess.Popen(["go", "run", f"{path}/main.go"])
    sleep(5)
    yield p
    p.terminate()
    terminate_pid_by_port(request_endpoint[1])


def test_does_csr_application_run_successfully(request_endpoint) -> bool:
    f = requests.get(request_endpoint[0])

    assert f.status_code == 200

    assert 'class="orbit_bk"' in f.text
    assert 'id="orbit_manifest"' in f.text


@pytest.fixture(scope="module")
def lighthouse_audits(request_endpoint, application):
    return audits_url(request_endpoint[0])


def test_lighthouse_important_heuristics(lighthouse_audits):
    important_heuristics = [
        ("first-contentful-paint", 1),
        ("speed-index", 1),
        ("largest-contentful-paint", 0.90),
        ("server-response-time", 0.90),
    ]

    for key, score in important_heuristics:
        assert lighthouse_audits[key]["score"] >= score, key
