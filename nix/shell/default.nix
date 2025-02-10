_: {
  perSystem = {
    config,
    pkgs,
    ...
  }: let
    shell = pkgs.callPackage ./shell.nix {};
    inherit (pkgs.llvmPackages_latest) stdenv bintools clang;

  in {
    devShells.default = (shell.override { inherit stdenv; } ).overrideAttrs (_oa: {
      nativeBuildInputs =
        [
          clang
          bintools
        ] ++ _oa.nativeBuildInputs
          ++ config.pre-commit.settings.enabledPackages;

        _oa.nativeBuildInputs
      shellHook = ''
        ${config.pre-commit.installationScript}
        ${_oa.shellHook}
      '';
    });
  };
}
