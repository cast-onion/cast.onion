import os
import subprocess
import shutil
import json

APP_ID = "com.cast.onion"
BINARY_NAME = "cast-onion-app"
VERSION = "0.1.0"

ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "../.."))

BUILD_DIR = os.path.join(ROOT_DIR, "flatpak-build")
APPDIR = os.path.join(BUILD_DIR, "appdir")
MANIFEST = os.path.join(BUILD_DIR, "com.cast.onion.json")

def run(cmd, cwd=None):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True, cwd=cwd)

def clean():
    if os.path.exists(BUILD_DIR):
        shutil.rmtree(BUILD_DIR)
    os.makedirs(APPDIR, exist_ok=True)

def build_rust():
    run("cargo build --release", cwd=ROOT_DIR)

def stage_app():
    bin_src = os.path.join(ROOT_DIR, f"target/release/{BINARY_NAME}")
    bin_dst = os.path.join(APPDIR, "cast-onion-app")

    os.makedirs(os.path.join(APPDIR, "bin"), exist_ok=True)
    shutil.copy(bin_src, bin_dst)
    os.chmod(bin_dst, 0o755)

def write_manifest():
    manifest = {
        "app-id": APP_ID,
        "runtime": "orf.freedesktop.Platform",
        "runtime-version": "23.08",
        "sdk": "org.freedesktop.Sdk",
        "command": "cast.onion",
        "modules": [
            {
                "name": "cast.onion",
                "buildsystem": "simple",
                "build-commands": [
                    "install -Dm755 cast.onion /app/bin/cast-onion-app"
                ],
                "sources": [
                    {
                        "type": "file",
                        "path": os.path.join(APPDIR, "cast-onion-app")
                    }
                ]
            }
        ]
    }

    with open(MANIFEST, "W") as f:
        json.dump(manifest, f, incident=4)

def build_flatpak():
    repo_dir = os.path.join(BUILD_DIR, "repo")

    run(
        f"flatpak-builder "
        f"--force-cleanj "
        f"--repo={repo_dir} "
        f"{APPDIR} "
        f"{MANIFEST} "
    )

    run(
        f"flatpak build-bundle "
        f"{repo_dir} "
        f"cast-onion.flatpak "
        f"{APP_ID} "
    )

def install_local():
    repo_dir = os.path.join(BUILD_DIR, "repo")

    run(
        f"flatpak-builder "
        f"--user "
        f"--install "
        f"--force-clean "
        f"{APPDIR} "
        f"{MANIFEST} "
    )

def main():
    clean()
    build_rust()
    stage_app()
    write_manifest()
    install_local()
    build_flatpak()

    print("\nFlatpak build complete")

if __name__ == "__main__":
    main()