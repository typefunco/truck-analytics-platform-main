package september

import (
	"context"
	"log/slog"
	"net/http"
	"truck-analytics-platform/internal/db"

	"github.com/gin-gonic/gin"
	orderedmap "github.com/wk8/go-ordered-map/v2"
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
		Data  *orderedmap.OrderedMap[string, []TruckAnalytics] `json:"data"`
		Error string                                           `json:"error,omitempty"`
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

	// Используем orderedmap для хранения данных
	dataByDistrict := orderedmap.New[string, []TruckAnalytics]()

	// Определяем порядок округов
	customOrder := []string{
		"Central",
		"North West",
		"Volga",
		"South",
		"Ural",
		"Siberia",
		"Far East",
	}

	// Предзаполняем карту пустыми значениями для каждого округа
	for _, district := range customOrder {
		dataByDistrict.Set(district, []TruckAnalytics{})
	}

	// Суммируем данные по округам
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

		// Добавляем данные региона в соответствующий федеральный округ
		if existing, ok := dataByDistrict.Get(federalDistrict); ok {
			dataByDistrict.Set(federalDistrict, append(existing, ta))
		}
	}

	// Проверка на ошибки итерации
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Error iterating over rows: " + err.Error(),
		})
		return
	}

	// Добавляем суммарные данные в карту
	for district, summary := range summaryByDistrict {
		if existing, ok := dataByDistrict.Get(district); ok {
			dataByDistrict.Set(district, append(existing, summary))
		}
	}

	// Отправляем ответ
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

	// Структура для ответа
	type Response struct {
		Data  *orderedmap.OrderedMap[string, []DistrictData] `json:"data"`
		Error string                                         `json:"error,omitempty"`
	}

	// Мапа для перевода русских названий федеральных округов на английский
	var regionTranslation = map[string]string{
		"Дальневосточный Федеральный Округ":   "Far East",
		"Приволжский Федеральный Округ":       "Volga",
		"Северо-Западный Федеральный Округ":   "North West ",
		"Северо-Кавказский Федеральный Округ": "North Caucasian",
		"Сибирский Федеральный Округ":         "Siberia",
		"Уральский Федеральный Округ":         "Ural",
		"Центральный Федеральный Округ":       "Central",
		"Южный Федеральный Округ":             "South",
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

	// Создаем ordered map для хранения данных с нужным порядком
	dataByDistrict := orderedmap.New[string, []DistrictData]()

	// Определяем порядок федеральных округов
	customOrder := []string{
		"Summary",
		"Central",
		"North West",
		"Volga",
		"South",
		"North Caucasian",
		"Ural",
		"Siberia",
		"Far East",
	}

	// Инициализируем пустые слайсы для каждого округа
	for _, district := range customOrder {
		dataByDistrict.Set(district, []DistrictData{})
	}

	// Добавляем Summary в начало
	dataByDistrict.Set("Summary", []DistrictData{})

	var summary DistrictData
	summary.RegionName = "Summary"
	summary.Faw = new(int)
	summary.Howo = new(int)
	summary.Jac = new(int)
	summary.Sany = new(int)
	summary.Sitrak = new(int)

	for rows.Next() {
		var item DistrictData
		err := rows.Scan(&item.RegionName, &item.Faw, &item.Howo, &item.Jac, &item.Sany, &item.Sitrak, &item.Total)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}

		// Обновляем суммарные значения
		if item.Faw != nil {
			*summary.Faw += *item.Faw
		}
		if item.Howo != nil {
			*summary.Howo += *item.Howo
		}
		if item.Jac != nil {
			*summary.Jac += *item.Jac
		}
		if item.Sany != nil {
			*summary.Sany += *item.Sany
		}
		if item.Sitrak != nil {
			*summary.Sitrak += *item.Sitrak
		}
		summary.Total += item.Total

		// Перевод имени округа
		translatedRegionName, exists := regionTranslation[item.RegionName]
		if !exists {
			translatedRegionName = item.RegionName
		}

		// Добавляем данные в соответствующий округ
		if existing, exists := dataByDistrict.Get(translatedRegionName); exists {
			item.RegionName = translatedRegionName
			dataByDistrict.Set(translatedRegionName, append(existing, item))
		}
	}
	if rows.Err() != nil {
		slog.Info("Failed to iterate over rows:", rows.Err())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over rows"})
		return
	}

	// Добавляем суммарные данные
	if summaryData, exists := dataByDistrict.Get("Summary"); exists {
		dataByDistrict.Set("Summary", append(summaryData, summary))
	}

	// Отправка ответа
	ctx.JSON(http.StatusOK, Response{
		Data: dataByDistrict,
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

	type TruckAnalyticsResponse struct {
		Data  *orderedmap.OrderedMap[string, []TruckAnalytics] `json:"data"`
		Error string                                           `json:"error,omitempty"`
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

	dataByDistrict := orderedmap.New[string, []TruckAnalytics]()

	customOrder := []string{
		"Central",
		"North West",
		"Volga",
		"South",
		"Ural",
		"Siberia",
		"Far East",
	}

	for _, district := range customOrder {
		dataByDistrict.Set(district, []TruckAnalytics{})
	}

	// Храним только одну сводку по округу
	summaryByDistrict := make(map[string]TruckAnalytics)

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

		if engName, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = engName
		}
		if engDistrict, ok := districtTranslations[federalDistrict]; ok {
			federalDistrict = engDistrict
		}

		// Если это итог по федеральному округу, сохраняем его отдельно
		if ta.RegionName == federalDistrict {
			summaryByDistrict[federalDistrict] = ta
			continue // Пропускаем добавление в основной список
		}

		// Добавляем только данные по регионам
		if existing, ok := dataByDistrict.Get(federalDistrict); ok {
			dataByDistrict.Set(federalDistrict, append(existing, ta))
		}
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Error iterating over rows: " + err.Error(),
		})
		return
	}

	// Добавляем итоговые данные в конец списка каждого округа
	for _, district := range customOrder {
		if summary, exists := summaryByDistrict[district]; exists {
			if existing, ok := dataByDistrict.Get(district); ok {
				dataByDistrict.Set(district, append(existing, summary))
			}
		}
	}

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

	// Структура для ответа
	type Response struct {
		Data  *orderedmap.OrderedMap[string, []DistrictData] `json:"data"`
		Error string                                         `json:"error,omitempty"`
	}

	// Мапа для перевода русских названий федеральных округов на английский
	var regionTranslation = map[string]string{
		"Дальневосточный Федеральный Округ":   "Far East",
		"Приволжский Федеральный Округ":       "Volga",
		"Северо-Западный Федеральный Округ":   "North West ",
		"Северо-Кавказский Федеральный Округ": "North Caucasian",
		"Сибирский Федеральный Округ":         "Siberia",
		"Уральский Федеральный Округ":         "Ural",
		"Центральный Федеральный Округ":       "Central",
		"Южный Федеральный Округ":             "South",
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

	// Создаем ordered map для хранения данных с нужным порядком
	dataByDistrict := orderedmap.New[string, []DistrictData]()

	// Определяем порядок федеральных округов
	customOrder := []string{
		"Summary",
		"Central",
		"North West",
		"Volga",
		"South",
		"North Caucasian",
		"Ural",
		"Siberia",
		"Far East",
	}

	// Инициализируем пустые слайсы для каждого округа
	for _, district := range customOrder {
		dataByDistrict.Set(district, []DistrictData{})
	}

	// Добавляем Summary в начало
	dataByDistrict.Set("Summary", []DistrictData{})

	var summary DistrictData
	summary.RegionName = "Summary"
	summary.Faw = new(int)
	summary.Howo = new(int)
	summary.Shacman = new(int)
	summary.Sitrak = new(int)

	for rows.Next() {
		var item DistrictData
		err := rows.Scan(&item.RegionName, &item.Faw, &item.Howo, &item.Sitrak, &item.Shacman, &item.Total)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		// Обновляем суммарные значения
		if item.Faw != nil {
			*summary.Faw += *item.Faw
		}
		if item.Howo != nil {
			*summary.Howo += *item.Howo
		}
		if item.Shacman != nil {
			*summary.Shacman += *item.Shacman
		}
		if item.Sitrak != nil {
			*summary.Sitrak += *item.Sitrak
		}
		summary.Total += item.Total

		// Перевод имени округа
		translatedRegionName, exists := regionTranslation[item.RegionName]
		if !exists {
			translatedRegionName = item.RegionName
		}

		// Добавляем данные в соответствующий округ
		if existing, exists := dataByDistrict.Get(translatedRegionName); exists {
			item.RegionName = translatedRegionName
			dataByDistrict.Set(translatedRegionName, append(existing, item))
		}
	}
	if rows.Err() != nil {
		slog.Info("Failed to iterate over rows:", rows.Err())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over rows"})
		return
	}

	// Добавляем суммарные данные
	if summaryData, exists := dataByDistrict.Get("Summary"); exists {
		dataByDistrict.Set("Summary", append(summaryData, summary))
	}

	// Отправка ответа
	ctx.JSON(http.StatusOK, Response{
		Data: dataByDistrict,
	})
}
