package cmd

import (
	"net/http"

	"code.vegaprotocol.io/go-wallet/cmd/icon"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"

	"go.uber.org/zap"
)

type systrayApp struct {
	log    *zap.Logger
	srv    *wallet.Service
	cproxy *consoleProxy
	// this is use to notify of any error from
	// both the wallet service and proxy
	errChan <-chan error
}

func newSystrayApp(rootPaths string) (*systrayApp, error) {
	cfg, err := wallet.LoadConfig(rootArgs.rootPath)
	if err != nil {
		return nil, err
	}

	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	srv, err := wallet.NewService(log, cfg, rootArgs.rootPath)
	if err != nil {
		return nil, err
	}

	cproxy := newConsoleProxy(log, cfg.Console.LocalPort, cfg.Console.URL, cfg.Nodes.Hosts[0], Version)

	errChan := make(chan error, 1)

	go func() {
		err := srv.Start()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	go func() {
		err := cproxy.Start()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	return &systrayApp{
		log:    log,
		srv:    srv,
		cproxy: cproxy,
	}, nil
}

func (a *systrayApp) onReady() {
	systray.SetIcon(icon.Data)

	openConsole := systray.AddMenuItem("Open Console", "Open the Vega Console")
	editConfig := systray.AddMenuItem("Preferences...", "Edit the Vega wallet configuration")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit Vega", "Quit the vega wallet")

	for {
		select {
		case <-openConsole.ClickedCh:
			open.Run(a.cproxy.GetBrowserURL())
		case <-editConfig.ClickedCh:
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func (a *systrayApp) onExit() {
	err := a.srv.Stop()
	if err != nil {
		a.log.Error("error stopping wallet http server", zap.Error(err))
	}
	err = a.cproxy.Stop()
	if err != nil {
		a.log.Error("error stopping console proxy server", zap.Error(err))
	}
}

func systrayStart(rootPath string) error {
	app, err := newSystrayApp(rootPath)
	if err != nil {
		return err
	}
	systray.Run(app.onReady, app.onExit)

	return nil
}

func checkConfig(rootPath string) error {
	return nil
}
