{
  description = "a fast and minimal static site generator";

  inputs.nixpkgs.url = "github:nixos/nixpkgs";

  outputs =
    { self
    , nixpkgs
    ,
    }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      overlay = final: prev: {
        vite = self.packages.${prev.system}.vite;
      };
      nixosModule = import ./module.nix;
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          vite = pkgs.buildGoModule {
            name = "vite";
            rev = "master";
            src = ./.;

            vendorHash = "sha256-jZO2ZX5Ik3TxBWMkq4TkA3TZvzGTQsuKRNKZFQt3gac=";
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.vite);
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [
              go
            ];
          };
        });
    };
}
