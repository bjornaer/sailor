package sessionmanager

import (
	"fmt"
	"net/http"

	"github.com/bjornaer/sailor/internal/db"
	"github.com/bjornaer/sailor/internal/port"
	"github.com/gin-gonic/gin"
)

type SessionManager struct {
	UpdatesFile string
	DBClient    db.DBClient
}

func (sm *SessionManager) ProcessPorts(c *gin.Context) {
	err := port.HandlePorts(sm.UpdatesFile, func(portKey string, portData port.PortData) {
		fmt.Printf("| Port Key: %s | Port Name: %v | \n", portKey, portData.Name)
		portDataStr, err := portData.Serialize()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		sm.DBClient.Put(portKey, portDataStr)
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.String(http.StatusOK, "Finished updating DB with ports data!")
}
