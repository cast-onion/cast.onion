import os
import shutil
import subprocess

PACKAGE_NAME = "cast-onion"
BINARY_NAME = "cast-onion-app"
VERSION = "0.1.0"
ARCH = "x86_64"

ICON_SOURCE = os.path.join("..", "app", "src", "assets", "favicon.png")


def run(cmd):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True)


def build_rust():
    run("cargo build --release")


def strip_binary():
    run(f"strip target/release/{BINARY_NAME}")


def prepare_pkg_dir():
    if os.path.exists("pkgbuild"):
        shutil.rmtree("pkgbuild")
    os.makedirs("pkgbuild")


def copy_files():
    shutil.copy(f"target/release/{BINARY_NAME}", "pkgbuild/")
    shutil.copy(ICON_SOURCE, "pkgbuild/cast-onion.png")


def write_pkgbuild():
    pkgbuild = f"""# Maintainer: cast.onion <1500forlifejahh@gmail.com>

pkgname={PACKAGE_NAME}
pkgver={VERSION}
pkgrel=1
pkgdesc="cast.onion Radio Client"
arch=('{ARCH}')
url=""
license=('custom')
depends=('alsa-lib' 'libx11' 'libxcb' 'mesa')
source=('{BINARY_NAME}' 'cast-onion.png')
noextract=()

package() {{
    install -Dm755 "{BINARY_NAME}" "$pkgdir/usr/bin/{BINARY_NAME}"

    install -Dm644 "cast-onion.png" \
        "$pkgdir/usr/share/icons/hicolor/256x256/apps/cast-onion.png"

    install -Dm644 /dev/stdin \
        "$pkgdir/usr/share/applications/cast.onion.desktop" << EOF
[Desktop Entry]
Name=cast.onion
Comment=cast.onion Radio Client
Exec=/usr/bin/{BINARY_NAME}
Icon=cast-onion
Terminal=false
Type=Application
Categories=AudioVideo;Network;
EOF
}}
"""
    with open("pkgbuild/PKGBUILD", "w") as f:
        f.write(pkgbuild)


def build_package():
    os.chdir("pkgbuild")
    run("makepkg -f")


def main():
    build_rust()
    strip_binary()

    prepare_pkg_dir()
    copy_files()
    write_pkgbuild()

    build_package()

    print("\nArch package built via makepkg")


if __name__ == "__main__":
    main()