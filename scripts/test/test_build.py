import os
import tempfile
import utils.git as git_util
import utils.orbit_bindings as orbit


def test_compare_build_cmd_with_latest_version():
    start_dir = os.getcwd()

    example_path = "./examples/basic-react"
    os.chdir(example_path)

    current_build_stats = orbit.e2e_measure_build_cmd()
    latest_version = git_util.version_abbrev("0")

    tmpdir = tempfile.mkdtemp()
    os.chdir(tmpdir)
    clone_path = git_util.clone_repo("http://github.com/GuyARoss/orbit", latest_version)

    os.chdir(f"./{clone_path}")
    orbit.link_examples()
    os.chdir(f"./examples/basic-react")

    latest_version_stats = orbit.e2e_measure_build_cmd()

    os.chdir(start_dir)

    assert round(current_build_stats["avg_build_time"] + 0.2, 1) >= round(
        latest_version_stats["avg_build_time"], 1
    )
    assert current_build_stats["failure_rate"] <= latest_version_stats["failure_rate"]
