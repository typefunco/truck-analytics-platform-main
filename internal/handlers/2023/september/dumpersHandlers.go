package september

import (
	"context"
	"log/slog"
	"net/http"
	"truck-analytics-platform/internal/db"

	"github.com/gin-gonic/gin"
)

func NineMonth2023Dumpers6x4(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName string `json:"region_name"`
		FAW        *int   `json:"faw"`
		HOWO       *int   `json:"howo"`
		JAC        *int   `json:"jac"`
		SANY       *int   `json:"sany"`
		SITRAK     *int   `json:"sitrak"`
		TOTAL      int    `json:"total"`
	}

	type TruckAnalyticsResponse struct {
		Data  map[string][]TruckAnalytics `json:"data"`
		Error string                      `json:"error,omitempty"`
	}

	query := `
		WITH base_data AS (
			SELECT 
				truck_analytics_2023_01_12."Federal_district",
				truck_analytics_2023_01_12."Region",
				truck_analytics_2023_01_12."Brand",
				SUM(truck_analytics_2023_01_12."Quantity") as total_sales
			FROM truck_analytics_2023_01_12
			WHERE 
				truck_analytics_2023_01_12."Wheel_formula" = '6x4'
				AND truck_analytics_2023_01_12."Brand" IN ('FAW', 'HOWO', 'JAC', 'SANY', 'SITRAK')
				AND truck_analytics_2023_01_12."Month_of_registration" <= 9
				AND truck_analytics_2023_01_12."Body_type" = 'Самосвал'
				AND truck_analytics_2023_01_12."Mass_in_segment_1" = '32001-40000'
			GROUP BY 
				truck_analytics_2023_01_12."Federal_district", 
				truck_analytics_2023_01_12."Region", 
				truck_analytics_2023_01_12."Brand"
		)
		SELECT 
			"Federal_district",
			COALESCE("Region", "Federal_district") as Region_name,
			MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END) as FAW,
			MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END) as HOWO,
			MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END) as JAC,
			MAX(CASE WHEN "Brand" = 'SANY' THEN total_sales END) as SANY,
			MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END) as SITRAK,
			COALESCE(MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SANY' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END), 0) as TOTAL
		FROM base_data
		GROUP BY 
			"Federal_district",
			"Region"
		ORDER BY 
			"Federal_district",
			CASE 
				WHEN "Region" = "Federal_district" THEN 1 
				ELSE 0 
			END,
			"Region"
	`

	db, err := db.Connect()
	if err != nil {
		slog.Warn("Can't connect to database")
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Failed to execute query: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	dataByDistrict := make(map[string][]TruckAnalytics)
	summaryByDistrict := make(map[string]TruckAnalytics)

	for rows.Next() {
		var ta TruckAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
			&ta.RegionName,
			&ta.FAW,
			&ta.HOWO,
			&ta.JAC,
			&ta.SANY,
			&ta.SITRAK,
			&ta.TOTAL,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
				Error: "Failed to scan row: " + err.Error(),
			})
			return
		}

		// Переводим название региона и федерального округа на английский
		if engName, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = engName
		}
		if engDistrict, ok := districtTranslations[federalDistrict]; ok {
			federalDistrict = engDistrict
		}

		// Обновляем данные по округу
		summary := summaryByDistrict[federalDistrict]

		if ta.FAW != nil {
			if summary.FAW == nil {
				summary.FAW = new(int)
			}
			*summary.FAW += *ta.FAW
		}
		if ta.HOWO != nil {
			if summary.HOWO == nil {
				summary.HOWO = new(int)
			}
			*summary.HOWO += *ta.HOWO
		}
		if ta.JAC != nil {
			if summary.JAC == nil {
				summary.JAC = new(int)
			}
			*summary.JAC += *ta.JAC
		}
		if ta.SANY != nil {
			if summary.SANY == nil {
				summary.SANY = new(int)
			}
			*summary.SANY += *ta.SANY
		}
		if ta.SITRAK != nil {
			if summary.SITRAK == nil {
				summary.SITRAK = new(int)
			}
			*summary.SITRAK += *ta.SITRAK
		}

		summary.TOTAL += ta.TOTAL
		summary.RegionName = federalDistrict
		summaryByDistrict[federalDistrict] = summary

		dataByDistrict[federalDistrict] = append(dataByDistrict[federalDistrict], ta)
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Error iterating over rows: " + err.Error(),
		})
		return
	}

	for district, summary := range summaryByDistrict {
		dataByDistrict[district] = append(dataByDistrict[district], summary)
	}

	ctx.JSON(http.StatusOK, TruckAnalyticsResponse{
		Data: dataByDistrict,
	})
}

func Dumpers6x4WithTotalMarket2023(ctx *gin.Context) {
	// Структура для хранения данных округа
	type DistrictData struct {
		RegionName string `json:"region_name"`
		Faw        *int   `json:"faw"`
		Howo       *int   `json:"howo"`
		Jac        *int   `json:"jac"`
		Sany       *int   `json:"sany"`
		Sitrak     *int   `json:"sitrak"`
		Total      int    `json:"total"`
	}

	// Структура для итогов
	type Summary struct {
		Faw    int `json:"faw"`
		Howo   int `json:"howo"`
		Jac    int `json:"jac"`
		Sany   int `json:"sany"`
		Sitrak int `json:"sitrak"`
		Total  int `json:"total"`
	}

	// Мапа для перевода русских названий федеральных округов на английский
	var regionTranslation = map[string]string{
		"Дальневосточный Федеральный Округ":   "Far Eastern Federal District",
		"Приволжский Федеральный Округ":       "Volga Federal District",
		"Северо-Западный Федеральный Округ":   "North West Federal District",
		"Северо-Кавказский Федеральный Округ": "North Caucasian Federal District",
		"Сибирский Федеральный Округ":         "Siberian Federal District",
		"Уральский Федеральный Округ":         "Ural Federal District",
		"Центральный Федеральный Округ":       "Central Federal District",
		"Южный Федеральный Округ":             "Southern Federal District",
	}

	// SQL-запрос для получения данных по округам
	query := `
		SELECT 
			"Federal_district",
			COALESCE(SUM(CASE WHEN "Brand" = 'FAW' THEN "Quantity" END), 0) AS FAW,
			COALESCE(SUM(CASE WHEN "Brand" = 'HOWO' THEN "Quantity" END), 0) AS HOWO,
			COALESCE(SUM(CASE WHEN "Brand" = 'JAC' THEN "Quantity" END), 0) AS JAC,
			COALESCE(SUM(CASE WHEN "Brand" = 'SANY' THEN "Quantity" END), 0) AS SANY,
			COALESCE(SUM(CASE WHEN "Brand" = 'SITRAK' THEN "Quantity" END), 0) AS SITRAK,
			COALESCE(SUM("Quantity"), 0) AS TOTAL
		FROM truck_analytics_2023_01_12
		WHERE 
			"Wheel_formula" = '6x4'
			AND "Brand" IN ('FAW', 'HOWO', 'JAC', 'SANY', 'SITRAK')
			AND "Month_of_registration" <= 9
			AND "Body_type" = 'Самосвал'
			AND "Mass_in_segment_1" = '32001-40000'
		GROUP BY "Federal_district"
		ORDER BY "Federal_district"
	`

	// Подключение к базе данных
	db, err := db.Connect()
	if err != nil {
		slog.Info("Can't connect to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't connect to database"})
		return
	}
	defer db.Close(context.Background())

	// Запрос к базе данных
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		slog.Info("Failed to execute query:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	var data []DistrictData
	for rows.Next() {
		var item DistrictData
		err := rows.Scan(&item.RegionName, &item.Faw, &item.Howo, &item.Jac, &item.Sany, &item.Sitrak, &item.Total)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		data = append(data, item)
	}
	if rows.Err() != nil {
		slog.Info("Failed to iterate over rows:", rows.Err())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over rows"})
		return
	}

	// Итоговые данные
	var summary Summary
	for _, item := range data {
		if item.Faw != nil {
			summary.Faw += *item.Faw
		}
		if item.Howo != nil {
			summary.Howo += *item.Howo
		}
		if item.Jac != nil {
			summary.Jac += *item.Jac
		}
		if item.Sany != nil {
			summary.Sany += *item.Sany
		}
		if item.Sitrak != nil {
			summary.Sitrak += *item.Sitrak
		}
		summary.Total += item.Total
	}

	// Подготовка JSON-ответа
	response := make(map[string][]DistrictData)

	// Добавляем Summary в начало
	response["Summary"] = []DistrictData{
		{
			RegionName: "Summary",
			Faw:        &summary.Faw,
			Howo:       &summary.Howo,
			Jac:        &summary.Jac,
			Sany:       &summary.Sany,
			Sitrak:     &summary.Sitrak,
			Total:      summary.Total,
		},
	}

	// Добавляем данные по округам после Summary
	for _, item := range data {
		// Перевод имени округа
		translatedRegionName, exists := regionTranslation[item.RegionName]
		if !exists {
			translatedRegionName = item.RegionName // если нет перевода, оставляем как есть
		}

		// Вставляем данные для каждого региона
		if _, exists := response[translatedRegionName]; !exists {
			response[translatedRegionName] = []DistrictData{}
		}
		item.RegionName = translatedRegionName
		response[translatedRegionName] = append(response[translatedRegionName], item)
	}

	// Ответ
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func NineMonth2023Dumpers8x4(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName string `json:"region_name"`
		FAW        *int   `json:"faw"`
		HOWO       *int   `json:"howo"`
		SHACMAN    *int   `json:"shacman"`
		SITRAK     *int   `json:"sitrak"`
		TOTAL      int    `json:"total"`
	}

	// TruckAnalyticsResponse structure for wrapping the response data
	type TruckAnalyticsResponse struct {
		Data  map[string][]TruckAnalytics `json:"data"`
		Error string                      `json:"error,omitempty"`
	}

	query := `
		WITH base_data AS (
			SELECT 
				truck_analytics_2023_01_12."Federal_district",
				truck_analytics_2023_01_12."Region",
				truck_analytics_2023_01_12."Brand",
				SUM(truck_analytics_2023_01_12."Quantity") as total_sales
			FROM truck_analytics_2023_01_12
			WHERE 
				truck_analytics_2023_01_12."Wheel_formula" = '8x4'
				AND truck_analytics_2023_01_12."Brand" IN ('FAW', 'HOWO', 'SHACMAN', 'SITRAK')
				AND truck_analytics_2023_01_12."Month_of_registration" <= 9
				AND truck_analytics_2023_01_12."Body_type" = 'Самосвал'
				AND truck_analytics_2023_01_12."Weight_in_segment_4" = '35001-45000'
			GROUP BY 
				truck_analytics_2023_01_12."Federal_district", 
				truck_analytics_2023_01_12."Region", 
				truck_analytics_2023_01_12."Brand"
		),
		federal_totals AS (
			SELECT 
				"Federal_district",
				"Federal_district" as "Region",
				"Brand",
				SUM(total_sales) as total_sales
			FROM base_data
			GROUP BY "Federal_district", "Brand"
		),
		combined_data AS (
			SELECT * FROM base_data
			UNION ALL
			SELECT * FROM federal_totals
		)
		SELECT 
			"Federal_district",
			COALESCE("Region", "Federal_district") as Region_name,
			MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END) as FAW,
			MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END) as HOWO,
			MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END) as SHACMAN,
			MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END) as SITRAK,
			COALESCE(MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END), 0) as TOTAL
		FROM combined_data
		GROUP BY 
			"Federal_district",
			"Region"
		ORDER BY 
			"Federal_district",
			CASE 
				WHEN "Region" = "Federal_district" THEN 1 
				ELSE 0 
			END,
			"Region"
	`

	db, err := db.Connect()
	if err != nil {
		slog.Warn("Can't connect to database")
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Failed to execute query: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	// Map for grouping data by federal district
	dataByDistrict := make(map[string][]TruckAnalytics)

	// Process query results and group by federal district
	for rows.Next() {
		var ta TruckAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
			&ta.RegionName,
			&ta.FAW,
			&ta.HOWO,
			&ta.SHACMAN,
			&ta.SITRAK,
			&ta.TOTAL,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
				Error: "Failed to scan row: " + err.Error(),
			})
			return
		}

		// Переводим название региона и федерального округа на английский
		if engName, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = engName
		}
		if engDistrict, ok := districtTranslations[federalDistrict]; ok {
			federalDistrict = engDistrict
		}

		// Обработка null значений: если значение nil, ставим 0 в поле TOTAL
		ta.FAW = nullToZero(ta.FAW)
		ta.HOWO = nullToZero(ta.HOWO)
		ta.SHACMAN = nullToZero(ta.SHACMAN)
		ta.SITRAK = nullToZero(ta.SITRAK)

		// Append the region data to the corresponding federal district
		dataByDistrict[federalDistrict] = append(dataByDistrict[federalDistrict], ta)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Error iterating over rows: " + err.Error(),
		})
		return
	}

	// Send the response
	ctx.JSON(http.StatusOK, TruckAnalyticsResponse{
		Data: dataByDistrict,
	})
}

// nullToZero function to convert null values to 0
func nullToZero(val *int) *int {
	if val == nil {
		return nil // Если значение равно nil (NULL в базе данных), возвращаем nil
	}
	return val // Если значение не nil, возвращаем сам указатель
}

func Dumpers8x4WithTotalMarket2023(ctx *gin.Context) {
	// Структура для хранения данных округа
	type DistrictData struct {
		RegionName string `json:"region_name"`
		Faw        *int   `json:"faw"`
		Howo       *int   `json:"howo"`
		Sitrak     *int   `json:"sitrak"`
		Shacman    *int   `json:"shacman"`
		Total      int    `json:"total"`
	}

	// Структура для итогов
	type Summary struct {
		Faw     int `json:"faw"`
		Howo    int `json:"howo"`
		Sitrak  int `json:"sitrak"`
		Shacman int `json:"shacman"`
		Total   int `json:"total"`
	}

	// Мапа для перевода русских названий федеральных округов на английский
	var regionTranslation = map[string]string{
		"Дальневосточный Федеральный Округ":   "Far Eastern Federal District",
		"Приволжский Федеральный Округ":       "Volga Federal District",
		"Северо-Западный Федеральный Округ":   "North West Federal District",
		"Северо-Кавказский Федеральный Округ": "North Caucasian Federal District",
		"Сибирский Федеральный Округ":         "Siberian Federal District",
		"Уральский Федеральный Округ":         "Ural Federal District",
		"Центральный Федеральный Округ":       "Central Federal District",
		"Южный Федеральный Округ":             "Southern Federal District",
	}

	// SQL-запрос для получения данных по округам
	query := `
		SELECT 
			"Federal_district",
			COALESCE(SUM(CASE WHEN "Brand" = 'FAW' THEN "Quantity" END), 0) AS FAW,
			COALESCE(SUM(CASE WHEN "Brand" = 'HOWO' THEN "Quantity" END), 0) AS HOWO,
			COALESCE(SUM(CASE WHEN "Brand" = 'SITRAK' THEN "Quantity" END), 0) AS SITRAK,
			COALESCE(SUM(CASE WHEN "Brand" = 'SHACMAN' THEN "Quantity" END), 0) AS SHACMAN,
			COALESCE(SUM("Quantity"), 0) AS TOTAL
		FROM truck_analytics_2023_01_12
		WHERE 
			"Wheel_formula" = '8x4'
			AND "Brand" IN ('FAW', 'HOWO', 'SITRAK', 'SHACMAN')
			AND "Body_type" = 'Самосвал'
			AND "Month_of_registration" <= 9
			AND "Weight_in_segment_4" = '35001-45000'
		GROUP BY "Federal_district"
		ORDER BY "Federal_district"
	`

	// Подключение к базе данных
	db, err := db.Connect()
	if err != nil {
		slog.Info("Can't connect to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't connect to database"})
		return
	}
	defer db.Close(context.Background())

	// Запрос к базе данных
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		slog.Info("Failed to execute query:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	var data []DistrictData
	for rows.Next() {
		var item DistrictData
		err := rows.Scan(&item.RegionName, &item.Faw, &item.Howo, &item.Sitrak, &item.Shacman, &item.Total)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		data = append(data, item)
	}
	if rows.Err() != nil {
		slog.Info("Failed to iterate over rows:", rows.Err())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over rows"})
		return
	}

	// Итоговые данные
	var summary Summary
	for _, item := range data {
		if item.Faw != nil {
			summary.Faw += *item.Faw
		}
		if item.Howo != nil {
			summary.Howo += *item.Howo
		}
		if item.Sitrak != nil {
			summary.Sitrak += *item.Sitrak
		}
		if item.Shacman != nil {
			summary.Shacman += *item.Shacman
		}

		summary.Total += item.Total
	}

	// Подготовка JSON-ответа
	response := make(map[string][]DistrictData)

	// Добавляем Summary в начало
	response["Summary"] = []DistrictData{
		{
			RegionName: "Summary",
			Faw:        &summary.Faw,
			Howo:       &summary.Howo,
			Sitrak:     &summary.Sitrak,
			Shacman:    &summary.Shacman,
			Total:      summary.Total,
		},
	}

	// Добавляем данные по округам после Summary
	for _, item := range data {
		// Перевод имени округа
		translatedRegionName, exists := regionTranslation[item.RegionName]
		if !exists {
			translatedRegionName = item.RegionName // если нет перевода, оставляем как есть
		}

		// Вставляем данные для каждого региона
		if _, exists := response[translatedRegionName]; !exists {
			response[translatedRegionName] = []DistrictData{}
		}
		item.RegionName = translatedRegionName
		response[translatedRegionName] = append(response[translatedRegionName], item)
	}

	// Ответ
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
