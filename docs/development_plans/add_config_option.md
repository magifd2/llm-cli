# Development Plan: Add Config File Option

## I. Goal
To extend the `llm-cli` application to allow users to specify an arbitrary configuration file path via a command-line option, instead of being limited to the default `~/.config/llm-cli/config.json`.

## II. Feasibility Assessment
The proposed feature is feasible. It primarily involves modifying the `internal/config` package to accept an optional configuration path and integrating a new global flag in the `cmd` package to pass this path.

## III. Development Plan

This plan adheres to the "Safe Refactoring Protocol" outlined in `.gemini/GEMINI.md`.

### Step 1: Establish Baseline
Commit the current stable state of the code. This ensures a clean rollback point if any issues arise.

### Step 2: Implement Changes Incrementally

#### Change 1: Modify `config.Load()` to accept an optional path.
- **Description**: Update the `Load()` function in `internal/config/config.go` to accept an optional `configPath` string argument. If this argument is provided and not empty, it will be used as the configuration file path. Otherwise, the function will fall back to the existing `GetConfigPath()` logic.
- **File**: `internal/config/config.go`
- **Action**:
    - Change `func Load() (*Config, error)` to `func Load(configPath string) (*Config, error)`.
    - Inside `Load()`, modify the line `configPath, err := GetConfigPath()` to:
      ```go
      var actualConfigPath string
      if configPath != "" {
          actualConfigPath = configPath
      } else {
          var err error
          actualConfigPath, err = GetConfigPath()
          if err != nil {
              return nil, err
          }
      }
      // Use actualConfigPath in subsequent operations
      ```
    - Update all existing calls to `config.Load()` throughout the codebase to pass an empty string (`""`) as the `configPath` argument. This ensures no immediate breakage of existing functionality.
- **Verification**: Run `make test` to ensure all existing tests pass.

#### Change 2: Modify `config.Save()` to accept an optional path.
- **Description**: Update the `Save()` method in `internal/config/config.go` to accept an optional `configPath` string argument. Similar to `Load()`, if this argument is provided and not empty, it will be used as the save path. Otherwise, it will use the existing `GetConfigPath()` logic.
- **File**: `internal/config/config.go`
- **Action**:
    - Change `func (c *Config) Save() error` to `func (c *Config) Save(configPath string) error`.
    - Inside `Save()`, modify the line `configPath, err := GetConfigPath()` to:
      ```go
      var actualConfigPath string
      if configPath != "" {
          actualConfigPath = configPath
      } else {
          var err error
          actualConfigPath, err = GetConfigPath()
          if err != nil {
              return err
          }
      }
      // Use actualConfigPath in subsequent operations
      ```
    - Update all existing calls to `config.Save()` throughout the codebase to pass an empty string (`""`) as the `configPath` argument.
- **Verification**: Run `make test` to ensure all existing tests pass.

#### Change 3: Add `--config` flag in `cmd/root.go` and integrate it.
- **Description**: Introduce a new global persistent flag `--config` (shorthand `-c`) in the `rootCmd` to allow users to specify the configuration file path. The value of this flag will then be passed to the `config.Load()` and `config.Save()` functions.
- **File**: `cmd/root.go` and potentially other `cmd/*.go` files where `config.Load()` or `config.Save()` are called.
- **Action**:
    - In `cmd/root.go`, declare a global variable to hold the config file path: `var cfgFile string`.
    - In the `init()` function of `cmd/root.go`, add the persistent flag:
      ```go
      rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/llm-cli/config.json)")
      ```
    - Identify all locations where `config.Load()` is called (e.g., in `cmd/root.go` or other command files). Modify these calls to pass `cfgFile` as the argument: `cfg, err := config.Load(cfgFile)`.
    - Identify all locations where `config.Save()` is called. Modify these calls to pass `cfgFile` as the argument: `err := cfg.Save(cfgFile)`.
- **Verification**: Run `make test` and `make build`. Manually test the new `--config` flag by:
    - Running `llm-cli --config /tmp/my_custom_config.json` and verifying that a new config file is created at `/tmp/my_custom_config.json` if it doesn't exist, or loaded from there if it does.
    - Running `llm-cli` (without `--config`) and verifying it still uses the default path.

### Step 3: Test and Verify
- After each incremental change, run `make test` to ensure no regressions.
- After all changes are implemented, perform a full build using `make build`.
- Manually test the new `--config` flag as described in Change 3's verification step.
- Run `make lint` to ensure code style and quality.

### Step 4: Complete
Once all changes are implemented, thoroughly tested, and verified, commit the final changes with a descriptive commit message following Conventional Commits specification (e.g., `feat: Add --config option for specifying config file path`).
