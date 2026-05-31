import os
import subprocess

APP_NAME = "cast-onion-app"
VERSION = "0.1.0"

def run(cmd):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True)

def build_rust():
    run("cargo build --release")

def build_snap():
    run("snapcraft")

def main():
    build_rust()
    run(f"strip target/release/{APP_NAME}-app")

    build_snap()

    print(f"\nBuild snap package for {APP_NAME}")

if __name__ == "__main__":
    main()