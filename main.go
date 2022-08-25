package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"github.com/elazarl/goproxy"
	"github.com/fatih/color"
)

func startProxy() {
	found := false

	authRe := regexp.MustCompile(`authkey=[^&]+`)
	proxy := goproxy.NewProxyHttpServer()
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", host, port), Handler: proxy}

	stopCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(`(mihoyo|hoyoverse)\.com`))).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if authRe.Match([]byte(req.URL.String())) && !found {
				color.White(req.URL.String())

				clipboard.WriteAll(req.URL.String())
				color.Green("Copied the link to your clipboard, paste it to Paimon.moe!")

				found = true

				setProxy(false)

				go func() {
					time.Sleep(5 * time.Second)
					cancel()
				}()
			}

			return req, nil
		})

	go func() {
		color.Yellow("Please open the Wish History on the game now!")
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Panicln("proxy server", err)
			}
		}
	}()

	<-stopCtx.Done()
	if err := server.Shutdown(stopCtx); err != nil {
		log.Panicln("proxy server:", err)
	}

}

func main() {
	defer func() {
		if err := recover(); err != nil {
			setProxy(false)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		setProxy(false)
		os.Exit(1)
	}()

	setProxy(true)
	startProxy()
	setProxy(false)
}
