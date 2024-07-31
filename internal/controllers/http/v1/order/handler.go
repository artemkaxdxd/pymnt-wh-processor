package order

import (
	"backend/config"
	requestOrder "backend/internal/controllers/http/request/order"
	"backend/internal/controllers/http/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h handler) getOrders(c *gin.Context) {
	limit, offset, err := h.parseLimitOffset(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err.Error()))
		return
	}

	filters, err := h.parseFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err.Error()))
		return
	}

	orders, svcCode, err := h.svc.GetOrders(c.Request.Context(), limit, offset, filters)
	if err != nil {
		c.JSON(config.ServiceCodeToHttpStatus(svcCode), response.NewErr(svcCode, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.New(svcCode).AddKey("orders", orders))
}

func (handler) parseLimitOffset(c *gin.Context) (limit, offset int, err error) {
	limitQuery := c.DefaultQuery("limit", "10")
	if limit, err = strconv.Atoi(limitQuery); err != nil {
		return
	}

	offsetQuery := c.DefaultQuery("offset", "0")
	offset, err = strconv.Atoi(offsetQuery)
	return
}

func (handler) parseFilters(c *gin.Context) (map[string]any, error) {
	filters := make(map[string]any)

	if query := c.Query("user_id"); query != "" {
		userID, err := uuid.Parse(query)
		if err != nil {
			return nil, err
		}

		filters["user_id"] = userID
	}

	if query := c.DefaultQuery("sort_by", "created_at"); query != "" {
		if query != "created_at" && query != "updated_at" {
			return nil, config.ErrInvalidFilters
		}

		filters["sort_by"] = query
	}
	if query := c.DefaultQuery("sort_order", "desc"); query != "" {
		if query != "asc" && query != "desc" {
			return nil, config.ErrInvalidFilters
		}

		filters["sort_order"] = query
	}

	statusQuery := c.Query("status")
	if statusQuery != "" {
		filters["status"] = strings.Split(statusQuery, ",")
	}

	isFinalQuery := c.Query("is_final")
	if isFinalQuery != "" {
		isFinal, err := strconv.ParseBool(isFinalQuery)
		if err != nil {
			return nil, err
		}

		filters["is_final"] = isFinal
	}

	if (statusQuery == "" && isFinalQuery == "") ||
		(statusQuery != "" && isFinalQuery != "") {
		return nil, config.ErrInvalidFilters
	}

	return filters, nil
}

func (h handler) processOrderPayment(c *gin.Context) {
	var body requestOrder.OrderEvent
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err.Error()))
		return
	}

	svcCode, err := h.svc.ProcessOrderEvent(c.Request.Context(), body)
	if err != nil {
		c.JSON(config.ServiceCodeToHttpStatus(svcCode), response.NewErr(svcCode, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.New(config.CodeOK).SetDescription(config.MsgCreateOK))
}

func (h handler) subscribeToOrders(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err.Error()))
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	eventChan := h.svc.Subscribe(orderID)

	for event := range eventChan {
		h.l.Info("Event status: ", event.OrderStatus)
		c.SSEvent("order_event", event)
		c.Writer.Flush()
	}

	c.Writer.WriteString("event: close\n\n")
	c.Writer.Flush()
}
