package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type currencyBalance struct {
	Currency     string `json:"currency"`
	TotalBalance string `json:"total_balance"`
}

type balanceResponse struct {
	IsAvailable  bool              `json:"is_available"`
	BalanceInfos []currencyBalance `json:"balance_infos"`
}

func FetchBalance(apikey string) (map[string]float64, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.deepseek.com/user/balance", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apikey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回状态码 %d", resp.StatusCode)
	}

	var br balanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&br); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if !br.IsAvailable {
		return nil, fmt.Errorf("API Key 不可用")
	}

	balances := make(map[string]float64)
	for _, info := range br.BalanceInfos {
		var val float64
		if _, err := fmt.Sscanf(info.TotalBalance, "%f", &val); err != nil {
			return nil, fmt.Errorf("解析余额失败: %w", err)
		}
		balances[info.Currency] = val
	}
	return balances, nil
}
