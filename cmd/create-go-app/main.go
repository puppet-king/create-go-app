package main

import (
	"flag"
	"os"

	"github.com/puppet-king/create-go-app/internal/config"
	"github.com/puppet-king/create-go-app/internal/scaffold"
	"github.com/puppet-king/create-go-app/internal/ui"
)

func main() {
	opts, err := config.Parse(os.Args[1:], os.Stdin)
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		ui.Error("%v", err)
		os.Exit(1)
	}

	if err := scaffold.Run(opts); err != nil {
		ui.Error("项目创建失败: %v", err)
		os.Exit(1)
	}
}
