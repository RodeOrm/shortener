package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rodeorm/shortener/internal/core"
)

/*
APIStatsHandler - обработчик для GET /api/internal/stats, возвращающий в ответ объект:

	{
	  "urls": <int>, // количество сокращённых URL в сервисе
	  "users": <int> // количество пользователей в сервисе
	}
*/
func (h *httpServer) APIStatsHandler(w http.ResponseWriter, r *http.Request) {

	/*
	   При запросе эндпоинта /api/internal/stats нужно проверять, что переданный в заголовке запроса X-Real-IP IP-адрес клиента входит в доверенную подсеть,
	   в противном случае возвращать статус ответа 403 Forbidden.
	   При пустом значении переменной trusted_subnet доступ к эндпоинту должен быть запрещён для любого входящего запроса.
	*/

	trusted, err := core.CheckNet(r, h.Config.ServerConfig.TrustedSubnet)
	if err != nil || !trusted {
		h.ForbiddenHandler(w, r)
	}

	s, err := h.StatStorage.SelectStatistic()
	if err != nil {
		handleError(w, err, "APIStatsHandler 1")
		h.badRequestHandler(w, r)
	}

	bodyBytes, err := json.Marshal(s)
	if err != nil {
		handleError(w, err, "APIStatsHandler 2")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bodyBytes))
}
