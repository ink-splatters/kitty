{pkgs ? import <nixpkgs> {}}:
with pkgs; let
  inherit (lib) optional optionals optionalString;
  inherit (xorg) libX11 libXrandr libXinerama libXcursor libXi libXext;
  harfbuzzWithCoreText = harfbuzz.override {withCoreText = stdenv.isDarwin;};
  inherit (llvmPackages_latest) clang bintools stdenv;
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
          matplotlib
        ]
        ++ optionals stdenv.isDarwin [
          apple-sdk_15
          libpng
          zlib
        ]
        ++ optionals stdenv.isLinux [
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
          cairo #
        ]
        ++ lib.optionals stdenv.hostPlatform.isLinux [
          wayland-scanner
        ]
        ++ checkInputs;

      nativeBuildInputs =
        [
          clang
          bintools
          ncurses
          pkg-config
          sphinx
          furo
          sphinx-copybutton
          sphinxext-opengraph
          sphinx-inline-tabs
          sphinx-autobuild
        ]
        ++ optionals stdenv.isDarwin [
          imagemagick
          libicns # For the png2icns tool.
        ];

      propagatedBuildInputs = optional stdenv.isLinux libGL;

      checkInputs = [
        pillow
      ];

      # Hardening control moved to setup.py
      hardeningDisable = ["all"];

      shellHook = ''
        # Add fonts by hand
        if [ ! -e ./fonts/SymbolsNerdFontMono-Regular.ttf ]; then
          mkdir fonts
          cp "${nerd-fonts.symbols-only}/share/fonts/truetype/NerdFonts/Symbols/SymbolsNerdFontMono-Regular.ttf" ./fonts/
        fi
        '' + optionalString stdenv.isLinux ''
          export KITTY_EGL_LIBRARY='${lib.getLib libGL}/lib/libEGL.so.1'
          export KITTY_STARTUP_NOTIFICATION_LIBRARY='${libstartup_notification}/lib/libstartup-notification-1.so'
          export KITTY_CANBERRA_LIBRARY='${libcanberra}/lib/libcanberra.so'
          export KITTY_FONTCONFIG_LIBRARY='${fontconfig.lib}/lib/libfontconfig.so'
        '';
    }
