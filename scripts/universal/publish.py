import subprocess

REPO = "cast-onion/cast-onion"
VERSION = "0.1.0"

def run(cmd):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True)

# Homebrew
def update_homebrew():
    print("Update Homebrew formula manually or via CLI")
    print("Formula repo usually: homebrew-core or tap")

# Scoop
def update_scoop():
    print("Scoop manifest update required in bucket repo")
    print("Example: cast-onion.json")

# Winget
def update_winget():
    manifest = f"""
PackageIdentifer: cast.onion
PackageVersion: {VERSION}
Publisher: cast.onion
PackageName: cast-onion
License: MIT
Installers:
  - Architecture: x64
    InstallerType: zip
    InstallerUrl: https://github.com/cast-onion/cast.onion/releases/download/v{VERSION}/cast-onion-windows.zip
    InstallerSha256: REPLACING
"""
    
    with open("winget-manifest.yml", "w") as f:
        f.write(manifest)

    print("Winget manifest generated")

def main():
    update_homebrew()
    update_scoop()
    update_winget()

if __name__ == "__main__":
    main()