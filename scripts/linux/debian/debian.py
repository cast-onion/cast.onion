import os
import shutil
import subprocess

PACKAGE_NAME = "cast.onion"
BINARY_NAME = "cast-onion-app"
VERSION = "0.1.0"
ARCH = "amd64"

BUILD_DIR = f"{PACKAGE_NAME}_{VERSION}_{ARCH}"

ICON_SOURCE = os.path.join("..", "app", "src", "assets", "favicon.png")

ICON_DEST = os.path.join(
    BUILD_DIR,
    "usr/share/icons/hicolor/256x256/apps"
)

DESKTOP_DIR = os.path.join(
    BUILD_DIR,
    "usr/share/applications"
)

DEBIAN_DIR = os.path.join(BUILD_DIR, "DEBIAN")
BIN_DIR = os.path.join(BUILD_DIR, "usr", "bin")


def run(cmd):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True)


def build_rust():
    run("cargo build --release")


def prepare_dirs():
    if os.path.exists(BUILD_DIR):
        shutil.rmtree(BUILD_DIR)

    os.makedirs(DEBIAN_DIR)
    os.makedirs(BIN_DIR)
    os.makedirs(ICON_DEST, exist_ok=True)
    os.makedirs(DESKTOP_DIR, exist_ok=True)


def copy_binary():
    src = f"target/release/{BINARY_NAME}"
    dst = os.path.join(BIN_DIR, BINARY_NAME)
    shutil.copy(src, dst)
    os.chmod(dst, 0o755)


def copy_icon():
    dst = os.path.join(ICON_DEST, "cast-onion.png")
    shutil.copy(ICON_SOURCE, dst)


def write_control():
    control_content = f"""Package: {PACKAGE_NAME}
Version: {VERSION}
Section: utils
Priority: optional
Architecture: {ARCH}
Maintainer: cast.onion organization <1500forlifejahh@gmail.com>
Depends: libasound2, libx11-6, libxcb1, libgl1
Description: The official cast.onion client to connect to the cast.onion radio network
"""

    with open(os.path.join(DEBIAN_DIR, "control"), "w") as f:
        f.write(control_content)


def write_desktop_file():
    desktop_content = f"""[Desktop Entry]
Name=cast.onion
Comment=cast.onion Radio Client
Exec=/usr/bin/{BINARY_NAME}
Icon=cast-onion
Terminal=false
Type=Application
Categories=AudioVideo;Network;
"""

    path = os.path.join(DESKTOP_DIR, "cast.onion.desktop")
    with open(path, "w") as f:
        f.write(desktop_content)


def post_install():
    postinst_path = os.path.join(DEBIAN_DIR, "postinst")

    content = """#!/bin/sh
set -e

if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache -f /usr/share/icons/hicolor || true
fi

if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database || true
fi
"""

    with open(postinst_path, "w") as f:
        f.write(content)

    os.chmod(postinst_path, 0o755)


def build_deb():
    run(f"dpkg-deb --build {BUILD_DIR}")


def main():
    build_rust()
    run(f"strip target/release/{BINARY_NAME}")
    prepare_dirs()
    copy_binary()
    write_control()
    copy_icon()
    write_desktop_file()
    post_install()
    build_deb()

    print(f"\nBuilt {BUILD_DIR}.deb")


if __name__ == "__main__":
    main()