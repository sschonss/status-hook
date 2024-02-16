package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
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

	return systemStatus
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obter os status do sistema
	status := getSystemStatus()

	// Converter para JSON
	jsonData, err := json.Marshal(status)
	if err != nil {
		http.Error(w, "Erro ao converter para JSON", http.StatusInternalServerError)
		return
	}

	// Responder com os dados do sistema em JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func main() {
	// Configurar o manipulador para o caminho '/status'
	http.HandleFunc("/status", webhookHandler)

	// Define a porta para ouvir
	port := 8080

	// Inicia o servidor na porta especificada
	fmt.Printf("Servidor rodando em http://localhost:%d/status\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
