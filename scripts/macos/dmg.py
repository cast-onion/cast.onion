import os
import shutil
import subprocess

APP_NAME = "cast.onion"
BINARY_NAME = "cast-onion-app"
VERSION = "0.1.0"

APP_DIR = f"{APP_NAME.app}"
DMG_NAME = f"{APP_NAME}-{VERSION}.dmg"

CONTENTS_DIR = os.path.join(APP_DIR, "Contents")
MACOS_DIR = os.path.join(CONTENTS_DIR, "MacOS")
RESOURCES_DIR = os.path.join(CONTENTS_DIR, "Resources")

ICON_SOURCE = os.path.join("..", "app", "src", "assets", "favicon.png")
ICON_ICNS = os.path.join(RESOURCES_DIR, "icon.icns")

def run(cmd):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True)

def build_rust():
    run("cargo build --release")

def prepare_app():
    if os.path.exists(APP_DIR):
        shutil.rmtree(APP_DIR)

    os.makedirs(MACOS_DIR)
    os.makedirs(RESOURCES_DIR)

def copy_binary():
    src = f"target/release/{BINARY_NAME}"
    dst = os.path.join(MACOS_DIR, BINARY_NAME)
    shutil.copy(src, dst)
    os.chmod(dst, 0o755)

def convert_icon():
    iconset_dir = "icon.iconset"

    if os.path.exists(iconset_dir):
        shutil.rmtree(iconset_dir)

    os.makedirs(iconset_dir)

    sizes = [16, 32, 64, 128, 256, 512]
    for size in sizes:
        run(f"sips -z {size} {size} {ICON_SOURCE} --out {iconset_dir}/icon_{size}x{size}.png")
        run(f"sips -z {size*2} {size*2} {ICON_SOURCE} --out {iconset_dir}/icon_{size}x{size}@2x.png")

    run(f"iconutils -c icns {iconset_dir} -o {ICON_ICNS}")
    shutil.rmtree(iconset_dir)

def write_info_plust():
    plist = f"""<?xml version="1.0 encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
"http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>{APP_NAME}</string>
    <key>CFBundleExecutable</key>
    <string>{BINARY_NAME}</string>
    <key>CFBundleIdentifier</key>
    <string>onion.cast.app</string>
    <key>CFBundleVersion</key>
    <string>{VERSION}</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
</dict>
</plist>
"""
    with open(os.path.join(CONTENTS_DIR, "Info.plist"), "w") as f:
        f.write(plist)

def build_dmg():
    temp_dir = "dmg_temp"

    if os.path.exits(temp_dir):
        shutil.rmtree(temp_dir)

    os.makedirs(temp_dir)

    # Copy app into DMG staging folder
    shutil.copytree(APP_DIR, os.path.join(temp_dir, APP_DIR))

    # Add Applications shortcut
    os.symlink("/Applications", os.path.join(temp_dir, "Applications"))

    # Create DMG
    run(f"hdiutil create -volume {APP_NAME} -srcfolder {temp_dir} -ov -format UDZO {DMG_NAME}")

    shutil.rmtree(temp_dir)

def main():
    build_rust()
    run(f"strip target/release/{BINARY_NAME}")

    prepare_app()
    copy_binary()
    convert_icon()
    write_info_plust()

    print(f"\nBuild {APP_DIR}")

    build_dmg()
    print(f"Build {DMG_NAME}")

if __name__ == "__main__":
    main()