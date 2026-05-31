import os
import shutil
import subprocess

APP_NAME = "cast.onion"
BINARY_NAME = "cast-onion-app"
VERSION = "0.1.0"

APPDIR = f"{APP_NAME}.AppDir"
APPIMAGE_NAME = f"{APP_NAME}-{VERSION}.AppImage"

ICON_PATH = os.path.join("..", "..", "..", "app", "src", "assets", "favicon.png")

def run(cmd):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True)

def build_rust():
    run("cargo build --release")

def prepare_appdir():
    if os.path.exists(APPDIR):
        shutil.rmtree(APPDIR)

    os.makedirs(f"{APPDIR}/usr/bin")
    os.makedirs(f"{APPDIR}/usr/share/icons/hicolor/256x256/apps")

def copy_binary():
    shutil.copy(
        f"target/release/{BINARY_NAME}"
    )

def write_desktop():
    desktop = f"""[Desktop Entry]
Name=cast.onion
Exec={BINARY_NAME}
Icon=cast-onion
Type=Application
Categories=AudioVideo;Network;
"""
    
    with open(f"{APPDIR}/{APP_NAME}.desktop", "w") as f:
        f.write(desktop)

def copy_icon():
    shutil.copy(
        ICON_PATH,
        f"{APPDIR}/usr/share/icons/hicolor/256x256/apps/{APP_NAME}.png"
    )

def build_appimage():
    run(f"appimagetool {APPDIR} {APPIMAGE_NAME}")

def main():
    build_rust()
    run(f"strip target/release/{BINARY_NAME}")

    prepare_appdir()
    copy_binary()
    write_desktop()
    copy_icon()

    build_appimage()

    print(f"\nBuilt {APPIMAGE_NAME}")

if __name__ == "__main__":
    main()