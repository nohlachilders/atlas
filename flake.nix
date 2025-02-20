{
    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
        systems.url = "github:nix-systems/default";
        flake-utils = {
            url = "github:numtide/flake-utils";
            inputs.systems.follows = "systems";
        };

        devenv.url = "github:cachix/devenv";
        gomod2nix = {
            url = "github:nix-community/gomod2nix";
            inputs.nixpkgs.follows = "nixpkgs";
        };
    };

    outputs = { self, nixpkgs, devenv, flake-utils, gomod2nix, ... } @ inputs:
        flake-utils.lib.eachDefaultSystem (system: let
            pkgs = import nixpkgs {
                inherit system;
                overlays = [
                    gomod2nix.overlays.default
                ];
            };
        in {
            packages = {
                devenv-up = self.devShells.${system}.default.config.procfileScript;
                devenv-test = self.devShells.${system}.default.config.test;
                gomod2nix = inputs.gomod2nix.default;
                default = self.bin;
            };

            devShells.default = devenv.lib.mkShell {
                inherit inputs pkgs;
                modules = [
                    ({pkgs, config, ... }: {
                        # stuff goes here
                        languages.go = {
                            enable = true;
                            enableHardeningWorkaround = true;
                        };

                        services.postgres = {
                            listen_addresses = "127.0.0.1";
                            enable = true;
                            createDatabase = false;
                        };

                        packages = with pkgs; [
                            goose
                            sqlc
                            air
                            gopls
                            delve
                            gomod2nix.packages.${system}.default
                        ];
                    })
                ];
            };
        }
        );
}
