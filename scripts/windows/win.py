import os
import shutil
import subprocess

PACKAGE_NAME = "cast.onion"
BINARY_NAME = "cast-onion-app.exe"
VERSION = "0.1.0"
ARCH = "x64"

BUILD_DIR = f"{PACKAGE_NAME}_{VERSION}_{ARCH}"

ICON_SOURCE = os.path.join("..", "app", "src", "assets", "favicon.png")

BIN_DIR = os.path.join(BUILD_DIR, "bin")
ASSETS_DIR = os.path.join(BUILD_DIR, "assets")

def run(cmd):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True)

def build_rust():
    run("cargo build --release")

def prepare_dirs():
    if os.path.exists(BUILD_DIR):
        shutil.rmtree(BUILD_DIR)

    os.makedirs(BIN_DIR)
    os.makedirs(ASSETS_DIR)

def copy_binary():
    src = os.path.join("target", "release", BINARY_NAME)
    dst = os.path.join(BIN_DIR, BINARY_NAME)
    shutil.copy(src, dst)

def copy_icon():
    dst = os.path.join(ASSETS_DIR, "cast-onion.png")
    shutil.copy(ICON_SOURCE, dst)

def write_metadata():
    content = f"""Name: {PACKAGE_NAME}
Version: {VERSION}
Author: cast.onion organization
Description: The official cast.onion client to connect to the cast.onion radio network
"""
    
    with open(os.path.join(BUILD_DIR, "metadata.txt"), "w") as f:
        f.write(content)

def create_batch_launcher():
    """
    Optional: simple double-click launcher
    """
    launcher_path = os.path.join(BUILD_DIR, "run.bat")

    content = f"""@echo off
cd /d %~dp0
bin\\{BINARY_NAME}
pause
"""
    
    with open(launcher_path, "w") as f:
        f.write(content)

def package_zip():
    shutil.make_archive(BUILD_DIR, "zip", BUILD_DIR)
    print(f"\nBuilt {BUILD_DIR.zip}")

def main():
    build_rust()

    try:
        run(f"strip target/release/{BINARY_NAME}")
    except Exception:
        print("strip not available, skipping")
    
    prepare_dirs()
    copy_binary()
    copy_icon()
    write_metadata()
    create_batch_launcher()
    package_zip()

if __name__ == "__main__":
    main()