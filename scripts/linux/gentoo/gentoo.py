import os
import shutil
import subprocess

PACKAGE_NAME = "cast-onion"
BINARY_NAME = "cast-onion-app"
VERSION = "0.1.0"

CATEGORY = "media-sound"
OVERLAY = os.path.expanduser("~/gentoo-overlay")

ICON_SOURCE = os.path.expanduser("..", "app", "src", "assets", "favicon.png")

def run(cmd):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True)

def build_rust():
    run("cargo build --release")

def strip_binary():
    run(f"strip target/release/{BINARY_NAME}")

def prepare_overlay():
    pkg_dir = os.path.join(OVERLAY, CATEGORY, PACKAGE_NAME)

    if os.path.exists(pkg_dir):
        shutil.rmtree(pkg_dir)

    os.makedirs(pkg_dir, exist_ok=True)
    return pkg_dir

def write_ebuild(pkg_dir):
    ebuild_path = os.path.join(pkg_dir, f"{PACKAGE_NAME}- {VERSION}.ebuild")

    ebuild = f"""# Copyright 2026
# Distributed under the terms of the GNU General Public License v2

EAPI=8

DESCRIPTION="cast.onion Radio Client"
HOMEPAGE=""
SRC_URI=""

LICENSE="custom"
SLOT="0"
KEYWORDS="~amd64"
IUSE=""

DEPEND=""
RDEPEND="media-libs/alsa-lib x11-libs/libX11 x11-libs/libxcb mebia-libs/mesa"

S="{{
    WORKDIR
}}"

src_unpack() {{
    mkdir =p "$S"
}}

src_install() {{
    dobin "{BINARY_NAME}"

    insito /usr/share/icons/hicolor/256x256/apps
    doins cast-onion.png

    cat <<EOF > cast.onion.desktop
[Desktop Entry]
Name=cast.onion
Comment=cast.onion Radio Client
Exec=/usr/bin/{BINARY_NAME}
Icon=cast-onion
Terminal=false
Type=Application
Categories=AudioVideo;Network;
EOF

    insito /usr/share/applications
    doins cast.onion.desktop
}}
"""
    
    with open(ebuild_path, "w") as f:
        f.write(ebuild)

    return ebuild_path

def manifest(pkg_dir):
    os.chdir(pkg_dir)
    run(f"ebuild {PACKAGE_NAME}-{VERSION}.ebuild manifest")

def copy_files(pkg_dir):
    shutil.copy(f"target/release/{BINARY_NAME}", os.path.join(pkg_dir, BINARY_NAME))
    shutil.copy(ICON_SOURCE, os.path.join(pkg_dir, "cast-onion.png"))



def install_package():
    run(f"emerge {CATEGORY}/{PACKAGE_NAME}")

def main():
    build_rust()
    strip_binary()

    pkg_dir = prepare_overlay()
    copy_files(pkg_dir)

    write_ebuild(pkg_dir)
    manifest(pkg_dir)
    
    install_package()

    print("\nGentoo package installed via emerge")

if __name__ == "__main__":
    main()