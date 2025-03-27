{pkgs ? import <nixpkgs> {}}:
with pkgs; let
  inherit (xorg) libX11 libXrandr libXinerama libXcursor libXi libXext;
  harfbuzzWithCoreText = harfbuzz.override {withCoreText = stdenv.isDarwin;};
  inherit (llvmPackages_latest) stdenv clang bintools;
in
  with python3Packages;
    mkShell.override {inherit stdenv;} rec {
      buildInputs =
        [
          harfbuzzWithCoreText
          ncurses
          lcms2
          xxHash
          simde
          go_1_24
        ]
        ++ lib.optionals stdenv.isDarwin [
          apple-sdk_15
          libpng
          python3
          zlib
        ]
        ++ lib.optionals stdenv.isLinux [
          fontconfig
          libunistring
          libcanberra
          libX11
          libXrandr
          libXinerama
          libXcursor
          libxkbcommon
          libXi
          libXext
          wayland-protocols
          wayland
          openssl
          dbus
        ]
        ++ checkInputs;

      nativeBuildInputs =
        [
          ncurses
          pkg-config
          sphinx
          furo
          sphinx-copybutton
          sphinxext-opengraph
          sphinx-inline-tabs
          sphinx-autobuild
          matplotlib
          clang
          bintools
        ]
        ++ lib.optionals stdenv.isDarwin [
          imagemagick
          libicns # For the png2icns tool.
        ]
        ++ lib.optionals stdenv.isLinux [
          wayland-scanner
        ];

      propagatedBuildInputs = lib.optional stdenv.isLinux libGL;

      checkInputs = [
        pillow
      ];

      # Configure hardening options using kitty setup.py CLI
      hardeningDisable = ["all"];

      shellHook = lib.optionalString stdenv.isLinux ''
        export CC=clang
        export CXX=clang++
        export KITTY_EGL_LIBRARY='${lib.getLib libGL}/lib/libEGL.so.1'
        export KITTY_STARTUP_NOTIFICATION_LIBRARY='${libstartup_notification}/lib/libstartup-notification-1.so'
        export KITTY_CANBERRA_LIBRARY='${libcanberra}/lib/libcanberra.so'
      '';
    }
