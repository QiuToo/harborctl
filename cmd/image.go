package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xiaolonglong/harborctl/internal/client"
)

var imageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"repo", "repository"},
	Short:   "Manage Harbor repositories and images",
	Long:    `Commands for managing Harbor repositories and images.`,
}

var repoListCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "List repositories in a project",
	Long:  `List all repositories in a Harbor project.

Examples:
  harborctl image list my-project
  harborctl repo list library`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		}

		repos, err := getRepositories(projectName)
		if err != nil {
			return err
		}

		if len(repos) == 0 {
			if projectName != "" {
				fmt.Printf("No repositories found in project '%s'.\n", projectName)
			} else {
				fmt.Println("No repositories found.")
			}
			return nil
		}

		printRepos(repos, projectName)
		return nil
	},
}

var repoInspectCmd = &cobra.Command{
	Use:   "inspect <project>/<repository>",
	Short: "Show repository details",
	Long:  `Display detailed information about a repository.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName, repoName, err := parseProjectRepo(args[0])
		if err != nil {
			return err
		}

		artifacts, err := harborClient.ListArtifacts(projectName, repoName)
		if err != nil {
			return err
		}

		fmt.Printf("Repository: %s/%s\n", projectName, repoName)
		fmt.Printf("Artifacts: %d\n", len(artifacts))

		if len(artifacts) > 0 {
			fmt.Println("\nTags:")
			for _, a := range artifacts {
				tags := "<none>"
				if len(a.Tags) > 0 {
					tags = strings.Join(a.Tags, ", ")
				}
				digest := a.Digest
				if len(digest) > 20 {
					digest = digest[:20] + "..."
				}
				fmt.Printf("  %s - %s (size: %d)\n", tags, digest, a.Size)
			}
		}

		return nil
	},
}

var imageListTagsCmd = &cobra.Command{
	Use:   "tags <project>/<repository>",
	Short: "List image tags",
	Long:  `List all tags for a specific image repository.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName, repoName, err := parseProjectRepo(args[0])
		if err != nil {
			return err
		}

		artifacts, err := harborClient.ListArtifacts(projectName, repoName)
		if err != nil {
			return err
		}

		if len(artifacts) == 0 {
			fmt.Printf("No artifacts found for %s/%s\n", projectName, repoName)
			return nil
		}

		fmt.Printf("Tags for %s/%s:\n", projectName, repoName)
		tagsMap := make(map[string][]string)
		for _, a := range artifacts {
			for _, tag := range a.Tags {
				tagsMap[tag] = append(tagsMap[tag], a.Digest)
			}
		}

		for tag, digests := range tagsMap {
			fmt.Printf("  %s (%d manifests)\n", tag, len(digests))
		}

		return nil
	},
}

var imageDeleteCmd = &cobra.Command{
	Use:   "delete <project>/<repository>:<tag>",
	Short: "Delete a specific image tag",
	Long:  `Delete a specific image by repository and tag.

Examples:
  harborctl image delete my-project/my-image:latest`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag := args[0]

		lastColon := strings.LastIndex(tag, ":")
		if lastColon == -1 {
			return fmt.Errorf("invalid image reference. Use format: project/repository:tag")
		}

		ref := tag[:lastColon]
		tagValue := tag[lastColon+1:]

		parts := strings.Split(ref, "/")
		if len(parts) < 2 {
			return fmt.Errorf("invalid image reference. Use format: project/repository:tag")
		}

		projectName := parts[0]
		repoName := strings.TrimPrefix(ref, projectName+"/")

		artifacts, err := harborClient.ListArtifacts(projectName, repoName)
		if err != nil {
			return err
		}

		var digest string
		for _, a := range artifacts {
			for _, t := range a.Tags {
				if t == tagValue {
					digest = a.Digest
					break
				}
			}
			if digest != "" {
				break
			}
		}

		if digest == "" {
			return fmt.Errorf("tag '%s' not found", tagValue)
		}

		if err := harborClient.DeleteArtifact(projectName, repoName, digest); err != nil {
			return err
		}

		fmt.Printf("Image %s deleted successfully!\n", tag)
		return nil
	},
}

var repoCleanCmd = &cobra.Command{
	Use:   "clean <project> [--keep <n>]",
	Short: "Clean up unused image tags",
	Long:  `Clean up unused image tags in a repository.

Keep the latest N tags and delete older ones.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keep, _ := cmd.Flags().GetInt("keep")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		projectName := args[0]

		repos, err := harborClient.ListRepositories(projectName)
		if err != nil {
			return err
		}

		totalDeleted := 0
		for _, repo := range repos {
			repoName := strings.TrimPrefix(repo.Name, projectName+"/")

			artifacts, err := harborClient.ListArtifacts(projectName, repoName)
			if err != nil {
				continue
			}

			sorted := make([]client.Artifact, len(artifacts))
			copy(sorted, artifacts)

			toDelete := sorted
			if keep > 0 && len(sorted) > keep {
				toDelete = sorted[keep:]
			}

			deleted := 0
			for _, a := range toDelete {
				err := harborClient.DeleteArtifact(projectName, repoName, a.Digest)
				if err == nil {
					deleted++
				}
			}
			totalDeleted += deleted

			if !dryRun && deleted > 0 {
				fmt.Printf("Repository %s: deleted %d old tags\n", repoName, deleted)
			}
		}

		if dryRun {
			fmt.Printf("[DRY RUN] Would delete %d old tags\n", totalDeleted)
		} else {
			fmt.Printf("Total deleted: %d tags\n", totalDeleted)
		}

		return nil
	},
}

func getRepositories(projectName string) ([]client.Repository, error) {
	if projectName != "" {
		return harborClient.ListRepositories(projectName)
	}

	projects, err := harborClient.ListProjects("", 1, 100)
	if err != nil {
		return nil, err
	}

	var allRepos []client.Repository
	for _, p := range projects {
		repos, err := harborClient.ListRepositories(p.Name)
		if err != nil {
			continue
		}
		allRepos = append(allRepos, repos...)
	}

	return allRepos, nil
}

func parseProjectRepo(ref string) (string, string, error) {
	ref = strings.TrimSuffix(ref, ":")

	lastColon := strings.LastIndex(ref, ":")
	if lastColon != -1 {
		ref = ref[:lastColon]
	}

	parts := strings.Split(ref, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid reference. Use format: project/repository")
	}

	projectName := parts[0]
	repoName := strings.TrimPrefix(ref, projectName+"/")

	return projectName, repoName, nil
}

func printRepos(repos []client.Repository, projectName string) {
	fmt.Printf("%-40s %-10s %-10s %s\n", "REPOSITORY", "PULLS", "STARS", "SIZE")
	fmt.Printf("%-40s %-10s %-10s %s\n", strings.Repeat("-", 40), strings.Repeat("-", 10), strings.Repeat("-", 10), strings.Repeat("-", 8))

	for _, r := range repos {
		size := formatBytes(r.Size)
		owner := ""
		if projectName != "" {
			owner = projectName + "/"
		} else {
			parts := strings.SplitN(r.Name, "/", 2)
			if len(parts) == 2 {
				owner = parts[0] + "/"
			}
		}
		repoName := strings.TrimPrefix(r.Name, owner)
		fmt.Printf("%-40s %-10d %-10d %s\n", owner+repoName, r.PullCount, r.StarCount, size)
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.AddCommand(repoListCmd)
	imageCmd.AddCommand(repoInspectCmd)
	imageCmd.AddCommand(imageListTagsCmd)
	imageCmd.AddCommand(imageDeleteCmd)
	imageCmd.AddCommand(repoCleanCmd)

	repoCleanCmd.Flags().IntP("keep", "k", 3, "Number of recent tags to keep")
	repoCleanCmd.Flags().BoolP("dry-run", "d", false, "Show what would be deleted")
}