package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newThemesCmd(comp *components) *cobra.Command {
	return &cobra.Command{
		Use:   "themes",
		Short: "列出所有可用的视觉主题",
		RunE: func(cmd *cobra.Command, args []string) error {
			compLoader, err := comp.themeLoader.Get()
			if err != nil {
				return fmt.Errorf("获取主题加载器失败: %w", err)
			}
			themes := compLoader.GetThemes()
			if len(themes) == 0 {
				fmt.Fprintln(os.Stderr, "没有找到任何主题")
				return nil
			}
			// 输出表格（人类可读）
			fmt.Fprintf(os.Stdout, "%-20s %-20s %s\n", "ID", "名称", "描述")
			for _, t := range themes {
				fmt.Fprintf(os.Stdout, "%-20s %-20s %s\n", t.ID, t.Name, t.Description)
			}
			return nil
		},
	}
}
