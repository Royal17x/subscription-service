package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Royal17x/subscription-service/internal/model"
	"github.com/go-chi/chi/v5"
)

// @Summary	Создать подписку
// @Tags		subscriptions
// @Accept		json
// @Produce	json
// @Param		request	body		model.SubscriptionRequest	true	"Данные подписки"
// @Success	201		{object}	model.SubscriptionResponse
// @Failure	400		{object}	map[string]string
// @Failure	500		{object}	map[string]string
// @Router		/subscriptions [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.Create(r.Context(), &req)
	if err != nil {
		h.log.Error("create subscription", "error", err)
		writeError(w, errorToStatus(err), err.Error())
		return
	}
	writeJson(w, http.StatusCreated, resp)
}

// @Summary	Получить подписку по ID
// @Tags		subscriptions
// @Produce	json
// @Param		id	path		int	true	"ID подписки"
// @Success	200	{object}	model.SubscriptionResponse
// @Failure	400	{object}	map[string]string
// @Failure	404	{object}	map[string]string
// @Failure	500	{object}	map[string]string
// @Router		/subscriptions/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("get subscription by id", "id", id, "error", err)
		writeError(w, errorToStatus(err), err.Error())
		return
	}

	writeJson(w, http.StatusOK, resp)
}

// @Summary	Обновить подписку
// @Tags		subscriptions
// @Accept		json
// @Produce	json
// @Param		id		path		int							true	"ID подписки"
// @Param		request	body		model.SubscriptionRequest	true	"Данные подписки"
// @Success	200		{object}	model.SubscriptionResponse
// @Failure	400		{object}	map[string]string
// @Failure	404		{object}	map[string]string
// @Failure	500		{object}	map[string]string
// @Router		/subscriptions/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req model.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		h.log.Error("update subscription", "id", id, "error", err)
		writeError(w, errorToStatus(err), err.Error())
		return
	}
	writeJson(w, http.StatusOK, resp)
}

// @Summary	Список подписок
// @Tags		subscriptions
// @Produce	json
// @Param		user_id			query		string	false	"Фильтр по user ID"
// @Param		service_name	query		string	false	"Фильтр по названию сервиса"
// @Success	200				{array}		model.SubscriptionResponse
// @Failure	500				{object}	map[string]string
// @Router		/subscriptions [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filter := model.SubscriptionFilter{}

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		filter.UserID = &userID
	}
	if serviceName := r.URL.Query().Get("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		h.log.Error("list subscriptions", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusOK, resp)
}

// @Summary	Суммарная стоимость подписок за период
// @Tags		subscriptions
// @Produce	json
// @Param		date_from		query		string	false	"Начало периода MM-YYYY"
// @Param		date_to			query		string	false	"Конец периода MM-YYYY"
// @Param		user_id			query		string	false	"Фильтр по user ID"
// @Param		service_name	query		string	false	"Фильтр по названию сервиса"
// @Success	200				{object}	map[string]int64
// @Failure	400				{object}	map[string]string
// @Failure	500				{object}	map[string]string
// @Router		/subscriptions/total-cost [get]
func (h *Handler) TotalCost(w http.ResponseWriter, r *http.Request) {
	filter := model.TotalCostFilter{}

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		filter.UserID = &userID
	}
	if serviceName := r.URL.Query().Get("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")
	if dateFrom == "" || dateTo == "" {
		writeError(w, http.StatusBadRequest, "date_from and date_to are required")
		return
	}
	from, err := parseQueryDate(dateFrom)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date_to: expected MM-YYYY")
		return
	}
	to, err := parseQueryDate(dateTo)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date_to: expected MM-YYYY")
		return
	}
	if !to.After(from) {
		writeError(w, http.StatusBadRequest, "date to must be after date_from")
		return
	}

	filter.DateFrom = from
	filter.DateTo = to

	total, err := h.service.TotalCost(r.Context(), filter)
	if err != nil {
		h.log.Error("total cost", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusOK, map[string]int64{"total_cost": total})
}

// @Summary	Удалить подписку
// @Tags		subscriptions
// @Param		id	path	int	true	"ID подписки"
// @Success	204
// @Failure	400	{object}	map[string]string
// @Failure	404	{object}	map[string]string
// @Failure	500	{object}	map[string]string
// @Router		/subscriptions/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.log.Error("delete subscription", "id", id, "error", err)
		writeError(w, errorToStatus(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func parseQueryDate(s string) (t time.Time, err error) {
	return time.Parse("01-2006", s)
}
