package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// SystemStatus representa os status do sistema
type SystemStatus struct {
	CPU      float64 `json:"cpu"`
	Memory   float64 `json:"memory"`
	Disk     float64 `json:"disk"`
	Temperature float64 `json:"temperature"`
}

func getSystemStatus() SystemStatus {
	// Obter informações da CPU
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println("Erro ao obter informações da CPU:", err)
	}

	// Obter informações da memória
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Erro ao obter informações de memória:", err)
	}

	// Obter informações do disco
	diskStat, err := disk.Usage("/")
	if err != nil {
		log.Println("Erro ao obter informações do disco:", err)
	}

	// Obter informações da temperatura (nota: pode não funcionar em todos os sistemas)
	temperature, err := host.SensorsTemperatures()
	if err != nil {
		log.Println("Erro ao obter informações de temperatura:", err)
	}

	systemStatus := SystemStatus{
		CPU:        cpuPercent[0],
		Memory:     vmStat.UsedPercent,
		Disk:       diskStat.UsedPercent,
		Temperature: temperature[0].Temperature,
	}

	if len(temperature) > 0 {
    	systemStatus.Temperature = temperature[0].Temperature
	} else {
    	log.Println("Nenhuma informação de temperatura disponível.")
	}

	return systemStatus
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	status := getSystemStatus()

	jsonData, err := json.Marshal(status)
	if err != nil {
		http.Error(w, "Erro ao converter para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/status", webhookHandler)

	port := 8080

	fmt.Printf("Servidor rodando em http://localhost:%d/status\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
