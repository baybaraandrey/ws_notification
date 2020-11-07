package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/mem"
)

// NewMonitorHandler register handlers for certain routes
func NewMonitorHandler(
	r *mux.Router,
) {
	handler := monitorHandler{}
	r.HandleFunc("/_api/", handler.APIok)
}

// monitorHandler represent API service status
type monitorHandler struct{}

// @Summary Application monitor
// @Description Application monitor
// @Produce  plain
// @Success 200
// @Router /_api/ [get]
// @Tags app-monitor
func (h *monitorHandler) APIok(w http.ResponseWriter, r *http.Request) {

	metrics := Gather()
	metrics.Time = time.Now().UTC().Format(time.RFC3339)

	b, err := json.MarshalIndent(metrics, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(b))
}

func Gather() *Metrics {
	v, _ := mem.VirtualMemory()
	return &Metrics{
		Hostname: "web",
		VirtualMemoryStat: VirtualMemory{
			Total:       v.Total,
			Free:        v.Free,
			UsedPercent: v.UsedPercent,
		},
	}
}
