package main

import (
	"an-blog/routes"
	"golang.org/x/sync/errgroup"
	"log"
)

var g errgroup.Group

func main() {
	// 初始化全局变量
	routes.InitGlobalVariable()

	//// 前台接口服务
	//g.Go(func() error {
	//	return routes.FrontendServer().ListenAndServe()
	//})

	// 后台接口服务
	g.Go(func() error {
		return routes.BackendServer().ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
