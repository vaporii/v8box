inputs: { config, lib, pkgs, ... }: {
  options.v8box = {
    enable = lib.mkEnableOption "v8box";
  };

  config = lib.mkIf config.v8box.enable {
    environment.systemPackages = [
      (import ./. { inherit pkgs; })
    ];
  };
}
