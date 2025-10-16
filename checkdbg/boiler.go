package checkdbg

import (
	// "korrectkm/spaserver"
	// "korrectkm/spaserver/templates"
	"github.com/mechiko/utility"
)

// func (c *Checks) GuideGtins() error {
// 	if ss, err := c.app.Repo().DbZnak().GtinAll(); err != nil {
// 		return err
// 	} else {
// 		c.app.Logger().Debugf("%v", len(ss))
// 	}
// 	return nil
// }

// func (c *Checks) AttachLite() error {
// 	dbFile := c.app.Repo().Dbs().Lite().File()
// 	if id, err := c.app.Repo().DbZnak().AttachLite(dbFile, "introduced", "0104810014020552215+L2JPj", "empty"); err != nil {
// 		return err
// 	} else {
// 		c.app.Logger().Debugf("заказ ид %d", id)
// 	}
// 	return nil
// }

// func (c *Checks) NewTemplates() (err error) {
// 	t := templates.New(c.app, true)
// 	err = t.LoadTemplates()
// 	return err
// }

// func (c *Checks) CheckSPAServer() (err error) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			err = fmt.Errorf("panic CheckSPAServer %v", r)
// 		}
// 	}()
// 	port := c.app.Configuration().SpaPort
// 	if port == "" || port == "auto" {
// 		if portFree, err := utility.GetFreePort(); err == nil {
// 			port = fmt.Sprintf("%d", portFree)
// 			// порт не записываем в файл конфигурации остается только в модели приложения
// 			c.app.Config().Set("spaport", port, false)
// 		}
// 	}
// 	c.app.Logger().Infof("spa http port %s", port)
// 	host := c.app.Configuration().Hostname
// 	if host == "" {
// 		host = "127.0.0.1"
// 	}
// 	spaServer := spaserver.New(c.app, port, true)
// 	urlTest := url.URL{
// 		Scheme: "http",
// 		Host:   fmt.Sprintf("%s:%s", host, port),
// 		Path:   "dbinfo",
// 	}
// 	spaServer.Start()
// 	c.app.Logger().Infof("spa url %s", urlTest.String())
// 	c.app.Open(urlTest.String())
// 	time.Sleep(120 * time.Second)
// 	spaServer.Shutdown()
// 	return fmt.Errorf("dumb")
// }

func (c *Checks) ParseZnak(znak string) string {
	z, _ := utility.ParseCisInfo(znak)
	return z.Serial
}
