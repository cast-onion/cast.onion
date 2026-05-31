import os
import shutil
import subprocess

PACKAGE_NAME = "cast-onion"
BINARY_NAME = "cast-onion-app"
VERSION = "0.1.0"
RELEASE = "1"
ARCH = "x86_64"

BUILD_ROOT = os.path.abspath("rpmbuild")
BUILD_DIR = os.path.join(BUILD_ROOT, "BUILD")
RPMS_DIR = os.path.join(BUILD_ROOT, "RPMS")
SOURCES_DIR = os.path.join(BUILD_ROOT, "SOURCES")
SPECS_DIR = os.path.join(BUILD_ROOT, "SPECS")

ICON_SOURCE = os.path.join("..", "app", "src", "assets", "favicon.png")


def run(cmd):
    print(f"> {cmd}")
    subprocess.check_call(cmd, shell=True)


def build_rust():
    run("cargo build --release")


def prepare_rpm_dirs():
    for d in [BUILD_DIR, RPMS_DIR, SOURCES_DIR, SPECS_DIR]:
        os.makedirs(d, exist_ok=True)


def copy_binary():
    src = f"target/release/{BINARY_NAME}"
    dst = os.path.join(SOURCES_DIR, BINARY_NAME)
    shutil.copy(src, dst)
    os.chmod(dst, 0o755)


def copy_icon():
    dst = os.path.join(SOURCES_DIR, "cast-onion.png")
    shutil.copy(ICON_SOURCE, dst)


def write_spec():
    spec_content = f"""Name:           {PACKAGE_NAME}
Version:        {VERSION}
Release:        {RELEASE}%{{?dist}}
Summary:        cast.onion Radio Client

License:        Proprietary
BuildArch:      {ARCH}

Requires:       alsa-lib, libX11, libxcb, mesa-libGL

%description
The official cast.onion client to connect to the cast.onion radio network.

%prep

%build

%install
mkdir -p %{{buildroot}}/usr/bin
mkdir -p %{{buildroot}}/usr/share/icons/hicolor/256x256/apps
mkdir -p %{{buildroot}}/usr/share/applications

install -m 755 %{{_sourcedir}}/{BINARY_NAME} %{{buildroot}}/usr/bin/{BINARY_NAME}
install -m 644 %{{_sourcedir}}/cast-onion.png %{{buildroot}}/usr/share/icons/hicolor/256x256/apps/cast-onion.png

cat > %{{buildroot}}/usr/share/applications/cast.onion.desktop << EOF
[Desktop Entry]
Name=cast.onion
Comment=cast.onion Radio Client
Exec=/usr/bin/{BINARY_NAME}
Icon=cast-onion
Terminal=false
Type=Application
Categories=AudioVideo;Network;
EOF

%post
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache -f /usr/share/icons/hicolor || true
fi

if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database || true
fi

%files
/usr/bin/{BINARY_NAME}
/usr/share/icons/hicolor/256x256/apps/cast-onion.png
/usr/share/applications/cast.onion.desktop

%changelog
* Thu Jan 01 2026 cast.onion <dev> - {VERSION}-1
- Fedora RPM build
"""

    spec_path = os.path.join(SPECS_DIR, f"{PACKAGE_NAME}.spec")
    with open(spec_path, "w") as f:
        f.write(spec_content)


def build_rpm():
    run(
        f"rpmbuild --define '_topdir {BUILD_ROOT}' "
        f"-ba {SPECS_DIR}/{PACKAGE_NAME}.spec"
    )


def main():
    build_rust()
    run(f"strip target/release/{BINARY_NAME}")

    prepare_rpm_dirs()
    copy_binary()
    copy_icon()
    write_spec()

    build_rpm()

    print(f"\nRPM built in: {RPMS_DIR}")


if __name__ == "__main__":
    main()