
// Package runtime provides isolation and resource control for RunInk node execution
package runtime

import (
        "bytes"
        "fmt"
        "os"
        "os/exec"
        "path/filepath"
        "syscall"
)

// NewExecutor creates a new executor with default configuration
func NewExecutor() *Executor {
        return &Executor{
                Config: ExecutorConfig{
                        Namespaces: DefaultNamespaces,
                        CgroupName: fmt.Sprintf("runink-%d", os.Getpid()),
                },
        }
}

// Executor manages the execution of commands in isolated environments
type Executor struct {
        Config ExecutorConfig
}

// SetCommand sets the command to execute
func (e *Executor) SetCommand(cmd []string) *Executor {
        e.Config.Command = cmd
        return e
}

// SetLimits sets the resource limits
func (e *Executor) SetLimits(limits Limits) *Executor {
        e.Config.Limits = limits
        return e
}

// SetChrootDir sets the chroot directory
func (e *Executor) SetChrootDir(dir string) *Executor {
        e.Config.ChrootDir = dir
        return e
}

// SetWorkDir sets the working directory inside the chroot
func (e *Executor) SetWorkDir(dir string) *Executor {
        e.Config.WorkDir = dir
        return e
}

// SetEnv sets the environment variables
func (e *Executor) SetEnv(env []string) *Executor {
        e.Config.Env = env
        return e
}

// SetNamespaces sets the namespaces to unshare
func (e *Executor) SetNamespaces(ns int) *Executor {
        e.Config.Namespaces = ns
        return e
}

// SetCgroupName sets the cgroup name
func (e *Executor) SetCgroupName(name string) *Executor {
        e.Config.CgroupName = name
        return e
}

// ExecInSandbox executes a command in an isolated environment with resource limits
// This is a simplified version that demonstrates the concept
func ExecInSandbox(cmd []string, limits Limits, chrootDir string) error {
        // Create a new executor
        executor := NewExecutor().
                SetCommand(cmd).
                SetLimits(limits).
                SetChrootDir(chrootDir)

        // Execute the command
        result, err := executor.Execute()
        if err != nil {
                return err
        }

        // Check the result
        if result.ExitCode != 0 {
                return fmt.Errorf("command failed with exit code %d: %s", 
                        result.ExitCode, string(result.Stderr))
        }

        return nil
}

// Execute runs the command in an isolated environment
func (e *Executor) Execute() (ExecutorResult, error) {
        result := ExecutorResult{}

        // Validate configuration
        if len(e.Config.Command) == 0 {
                return result, fmt.Errorf("no command specified")
        }

        // Prepare the command
        command := e.Config.Command[0]
        args := []string{}
        if len(e.Config.Command) > 1 {
                args = e.Config.Command[1:]
        }

        // Create a temporary directory for the chroot if not specified
        tempDir := ""
        if e.Config.ChrootDir == "" {
                var err error
                tempDir, err = os.MkdirTemp("", "runink-chroot-")
                if err != nil {
                        return result, fmt.Errorf("failed to create temporary directory: %v", err)
                }
                defer os.RemoveAll(tempDir)
                e.Config.ChrootDir = tempDir

                // Prepare the chroot environment
                if err := PrepareChroot(e.Config.ChrootDir); err != nil {
                        return result, fmt.Errorf("failed to prepare chroot: %v", err)
                }
        }

        // Create command
        cmd := exec.Command(command, args...)

        // Set up stdout and stderr capture
        var stdout, stderr bytes.Buffer
        cmd.Stdout = &stdout
        cmd.Stderr = &stderr

        // Set environment variables
        if len(e.Config.Env) > 0 {
                cmd.Env = e.Config.Env
        }

        // Set up process attributes for isolation
        cmd.SysProcAttr = &syscall.SysProcAttr{
                Chroot:     e.Config.ChrootDir,
                Cloneflags: uintptr(e.Config.Namespaces),
        }

        // Set working directory if specified
        if e.Config.WorkDir != "" {
                cmd.Dir = e.Config.WorkDir
        } else {
                cmd.Dir = "/"
        }

        // Start the command
        if err := cmd.Start(); err != nil {
                return result, fmt.Errorf("failed to start command: %v", err)
        }

        // Apply cgroup limits
        if e.Config.CgroupName != "" {
                if err := ApplyCgroup(e.Config.CgroupName, cmd.Process.Pid, e.Config.Limits); err != nil {
                        // Don't fail the command if cgroup setup fails, just log the error
                        fmt.Fprintf(os.Stderr, "Warning: failed to apply cgroup limits: %v\n", err)
                }
                // Clean up cgroup when done
                defer CleanupCgroup(e.Config.CgroupName)
        }

        // Wait for the command to complete
        err := cmd.Wait()
        result.Stdout = stdout.Bytes()
        result.Stderr = stderr.Bytes()

        // Get exit code
        if err != nil {
                if exitErr, ok := err.(*exec.ExitError); ok {
                        result.ExitCode = exitErr.ExitCode()
                }
                result.Error = err
        }

        return result, nil
}

// RunCommand is a simplified function to run a command with isolation
// This is a convenience function for simple use cases
func RunCommand(command []string, limits Limits) (string, error) {
        // Create a temporary directory for the chroot
        tempDir, err := os.MkdirTemp("", "runink-chroot-")
        if err != nil {
                return "", fmt.Errorf("failed to create temporary directory: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Prepare the chroot environment
        if err := PrepareChroot(tempDir); err != nil {
                return "", fmt.Errorf("failed to prepare chroot: %v", err)
        }

        // Copy the command binary to the chroot
        binPath := filepath.Join(tempDir, "bin", filepath.Base(command[0]))
        if err := copyFile(command[0], binPath); err != nil {
                return "", fmt.Errorf("failed to copy binary: %v", err)
        }

        // Create a new executor
        executor := NewExecutor().
                SetCommand(append([]string{"/bin/" + filepath.Base(command[0])}, command[1:]...)).
                SetLimits(limits).
                SetChrootDir(tempDir)

        // Execute the command
        result, err := executor.Execute()
        if err != nil {
                return "", err
        }

        // Return the output
        if result.ExitCode != 0 {
                return string(result.Stdout), fmt.Errorf("command failed with exit code %d: %s", 
                        result.ExitCode, string(result.Stderr))
        }

        return string(result.Stdout), nil
}

// Helper function to copy a file
func copyFile(src, dst string) error {
        // Ensure the destination directory exists
        if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
                return err
        }

        // Read the source file
        data, err := os.ReadFile(src)
        if err != nil {
                return err
        }

        // Write to the destination file
        if err := os.WriteFile(dst, data, 0755); err != nil {
                return err
        }

        return nil
}
