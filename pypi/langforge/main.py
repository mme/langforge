import sys
import subprocess
import platform
import os

def get_path():
    platform_system = platform.system().lower()
    os_name = None
    if platform_system == 'linux':
        os_name = 'linux'
    elif platform_system == 'darwin':
        os_name = 'macos'
    elif platform_system == 'windows':
        os_name = 'windows'
    else:
        print(f'Unsupported platform: {platform_system}', file=sys.stderr)
        sys.exit(1)   
      
    machine = platform.machine()
    arch = None
    if machine in ['x86_64', 'AMD64']:
        arch = 'amd64'
    elif machine in ['arm64', 'aarch64']:
        arch = 'arm64'
    else:
        print(f'Unsupported architecture: {machine}', file=sys.stderr)
        sys.exit(1)    

    package_dir = os.path.dirname(os.path.abspath(__file__))
    return os.path.join(package_dir, 'bin', f'langforge-{os_name}-{arch}')

def main():
    command = get_path()
    process = subprocess.Popen([command] + sys.argv[1:], stdout=sys.stdout, stderr=sys.stderr)
    process.wait()
    sys.exit(process.returncode)

if __name__ == "__main__":
    main()
