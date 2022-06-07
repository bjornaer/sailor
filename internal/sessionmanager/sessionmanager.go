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
		portDataStr, err := portData.Serialize()
		fmt.Println(portDataStr)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		err = sm.DBClient.Put(portKey, portDataStr)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.String(http.StatusOK, "Finished updating DB with ports data!")
}

func (sm *SessionManager) GetPort(c *gin.Context) {
	portKey := c.Param("portid")
	portDataSerialized, err := sm.DBClient.Get(portKey)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	portData, err := port.DeserializePort(portDataSerialized)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, portData)
}

func (sm *SessionManager) PutPort(c *gin.Context) {
	var port port.PortData
	c.BindJSON(&port)
	portString, err := port.Serialize()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	err = sm.DBClient.Put(port.Unlocs[0], portString)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.String(http.StatusOK, "Port Created/Updated successfully")
}
