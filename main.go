package main

import (
	"fmt"
	"os"

	"templater/internal/executor"
	"templater/internal/fs"
	"templater/internal/template"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "templater",
	Short: "A CLI tool for applying patch-based features to projects",
}

var listCmd = &cobra.Command{
	Use:   "list <template-repo>",
	Short: "Display available features as an ASCII tree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath := args[0]
		fileSystem := fs.OSFileSystem{}

		features, err := template.ListFeatures(fileSystem, repoPath)
		if err != nil {
			return err
		}

		tree := template.RenderTree(features)
		fmt.Print(tree)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status <target-dir>",
	Short: "Show features applied to a target project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := args[0]
		fileSystem := fs.OSFileSystem{}

		applied, err := template.ReadApplied(fileSystem, targetPath)
		if err != nil {
			return err
		}

		if len(applied) == 0 {
			fmt.Println("No features applied.")
			return nil
		}

		fmt.Println("Applied features:")
		for _, feature := range applied {
			fmt.Printf("  - %s\n", feature)
		}
		return nil
	},
}

var (
	dryRun       bool
	featuresFile string
)

var applyCmd = &cobra.Command{
	Use:   "apply <template-repo> <target-dir> [features...]",
	Short: "Apply features and their dependencies to a target project",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		templatePath := args[0]
		targetPath := args[1]
		features := args[2:]

		if featuresFile != "" && len(features) > 0 {
			return fmt.Errorf("cannot use both -f and positional feature arguments")
		}

		fileSystem := fs.OSFileSystem{}

		if featuresFile != "" {
			var err error
			features, err = template.ParseFeaturesFile(fileSystem, featuresFile)
			if err != nil {
				return fmt.Errorf("failed to read features file: %w", err)
			}
		}

		if len(features) == 0 {
			return fmt.Errorf("no features specified")
		}

		if dryRun {
			result, err := template.DryRun(fileSystem, templatePath, targetPath, features)
			if err != nil {
				return err
			}

			fmt.Println("Would apply:")
			for i, feature := range result.WouldApply {
				fmt.Printf("  %d. %s\n", i+1, feature)
			}
			return nil
		}

		exec := executor.NewShellExecutor()
		result, err := template.ApplyFeatures(fileSystem, exec, templatePath, targetPath, features)
		if err != nil {
			return err
		}

		for _, feature := range result.Applied {
			fmt.Printf("Applying %s... done\n", feature)
		}

		appliedCount := len(result.Applied)
		if appliedCount == 1 {
			fmt.Printf("\nApplied 1 feature.")
		} else {
			fmt.Printf("\nApplied %d features.", appliedCount)
		}

		if len(result.AlreadyApplied) > 0 {
			fmt.Printf(" (%d already applied: %s)", len(result.AlreadyApplied), joinFeatures(result.AlreadyApplied))
		}
		fmt.Println()

		allApplied := append(result.AlreadyApplied, result.Applied...)
		if err := template.WriteApplied(fileSystem, targetPath, allApplied); err != nil {
			return fmt.Errorf("failed to update applied.yml: %w", err)
		}

		return nil
	},
}

func joinFeatures(features []string) string {
	if len(features) == 0 {
		return ""
	}
	result := features[0]
	for _, f := range features[1:] {
		result += ", " + f
	}
	return result
}

func init() {
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be applied without applying")
	applyCmd.Flags().StringVarP(&featuresFile, "file", "f", "", "Read features from file (one per line)")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
