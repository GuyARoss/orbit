# !/bin/python3.8
import os
from typing import NoReturn

SUPPORTED_TARGETS = {
    'win32': {
        'os': 'windows',
        'arch': 'amd64',
    },
    'linux': {
        'os': 'linux',
        'arch': 'amd64',
    },
    'darwin': {
        'os': 'darwin',
        'arch': 'amd64'
    }
}


def golang_build_target(name: str, os_target: str, arch_target: os) -> NoReturn:
    os.environ['GOOS'] = os_target
    os.environ['GOARCH'] = arch_target

    os.system(f'go build -o ./bin/exec/{name} ./main.go')


# build_targets_from_map
# @returns bool denoting failure or success of the build operation
def build_targets_from_map() -> bool:
    try:
        for k in SUPPORTED_TARGETS:
            (os_t, arch) = SUPPORTED_TARGETS[k].values()
            golang_build_target(k, os_t, arch)

        return True
    except:
        return False


def main() -> NoReturn:
    build_targets_from_map()


if __name__ == '__main__':
    main()
