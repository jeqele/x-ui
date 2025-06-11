package controller

import (
	"x-ui/web/service"

	"github.com/gin-gonic/gin"
)

type APIController struct {
	BaseController
	inboundController *InboundController
	Tgbot             service.Tgbot
	serverService     service.ServerService
}

func NewAPIController(g *gin.RouterGroup) *APIController {
	a := &APIController{}
	a.initRouter(g)
	return a
}

func (a *APIController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/xui/API/inbounds")
	g.Use(a.checkLogin)

	a.inboundController = NewInboundController(g)

	inboundRoutes := []struct {
		Method  string
		Path    string
		Handler gin.HandlerFunc
	}{
		{"GET", "/createbackup", a.createBackup},
		{"POST", "/importbackup", a.importBackup},
		{"GET", "/", a.inboundController.getInbounds},
		{"GET", "/get/:id", a.inboundController.getInbound},
		{"GET", "/getClientTraffics/:email", a.inboundController.getClientTraffics},
		{"GET", "/getClientTrafficsById/:id", a.inboundController.getClientTrafficsById},
		{"POST", "/add", a.inboundController.addInbound},
		{"POST", "/del/:id", a.inboundController.delInbound},
		{"POST", "/update/:id", a.inboundController.updateInbound},
		{"POST", "/addClient", a.inboundController.addInboundClient},
		{"POST", "/:id/delClient/:clientId", a.inboundController.delInboundClient},
		{"POST", "/updateClient/:clientId", a.inboundController.updateInboundClient},
		{"POST", "/:id/resetClientTraffic/:email", a.inboundController.resetClientTraffic},
		{"POST", "/resetAllTraffics", a.inboundController.resetAllTraffics},
		{"POST", "/resetAllClientTraffics/:id", a.inboundController.resetAllClientTraffics},
		{"POST", "/delDepletedClients/:id", a.inboundController.delDepletedClients},
		{"POST", "/onlines", a.inboundController.onlines},
	}

	for _, route := range inboundRoutes {
		g.Handle(route.Method, route.Path, route.Handler)
	}
}

func (a *APIController) createBackup(c *gin.Context) {
	a.Tgbot.SendBackupToAdmins()
}

func (a *APIController) importBackup(c *gin.Context) {
	// Get the file from the request body
	file, _, err := c.Request.FormFile("db")
	if err != nil {
		jsonMsg(c, "Error reading db file", err)
		return
	}
	defer file.Close()
	
	// Import it
	err = a.serverService.ImportDB(file)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	
	// Restart Xray service
	err = a.serverService.RestartXrayService()
	if err != nil {
		jsonMsg(c, "Import successful, but failed to restart Xray", err)
		return
	}
	
	jsonObj(c, "Import DB successful", nil)
}
