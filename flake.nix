{
  description = "go flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }@inputs:
  let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in
  {
    packages.${system}.default = (import ./nix { inherit pkgs; });

    nixosModules.default = self.nixosModules.v8box;
    nixosModules.v8box = import ./nix/nixos.nix inputs;

    homeManagerModules.default = self.homeManagerModules.v8box;
    homeManagerModules.v8box = import ./nix/home-manager.nix inputs;

    devShells.x86_64-linux.default = pkgs.mkShell {
      packages = with pkgs; [
        go
      ];
    };
  };
}
