package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var profile, region, bucket, prefix string
	var bucketDelete, listObjects bool
	var rootCmd = &cobra.Command{
		Use:   "s3-cleaner",
		Short: "A cli to clean and destroy s3 bucket",
		Long:  `A cli to clean and destroy s3 bucket.`,
		Run: func(cmd *cobra.Command, args []string) {
			bucketDelete, _ = cmd.Flags().GetBool("bucket-delete")
			listObjects, _ = cmd.Flags().GetBool("list-objects")
			s3Cleanup(&profile, &region, &bucket, &prefix, bucketDelete, listObjects)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "AWS Profile (Required)")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS Region (Required)")
	rootCmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", "", "AWS S3 bucket (Required)")
	rootCmd.PersistentFlags().StringVarP(&prefix, "prefix", "f", "", "AWS S3 bucket Prefix (Optional)")
	rootCmd.PersistentFlags().Bool("bucket-delete", false, "Delete the bucket after cleaning (Optional)")
	rootCmd.PersistentFlags().Bool("list-objects", false, "List bucket objects (Optional)")
	rootCmd.MarkFlagsRequiredTogether("profile", "region", "bucket")
	rootCmd.MarkFlagsMutuallyExclusive("bucket-delete", "list-objects")

	rootCmd.Execute()
}
