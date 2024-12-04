package september

import (
	"context"
	"log/slog"
	"net/http"
	"truck-analytics-platform/internal/db"

	"github.com/gin-gonic/gin"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func NineMonth2024Dumpers6x4(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName  string `json:"region_name"`
		FAW         *int   `json:"faw"`
		HOWO        *int   `json:"howo"`
		JAC         *int   `json:"jac"`
		SANY        *int   `json:"sany"`
		SITRAK      *int   `json:"sitrak"`
		SHACMAN     *int   `json:"shacman"`
		DONGFENG    *int   `json:"dongfeng"`
		TOTAL       int    `json:"total"`
		TotalMarket int    `json:"total_market"` // Добавлено поле для total market
	}

	// Структура для обертки ответа
	type TruckAnalyticsResponse struct {
		Data  *orderedmap.OrderedMap[string, []TruckAnalytics] `json:"data"`
		Error string                                           `json:"error,omitempty"`
	}

	// SQL запрос
	query := `
		WITH base_data AS (
			SELECT 
				truck_analytics_2024_01_09."Federal_district",
				truck_analytics_2024_01_09."Region",
				truck_analytics_2024_01_09."Brand",
				SUM(truck_analytics_2024_01_09."Quantity") as total_sales
			FROM truck_analytics_2024_01_09
			WHERE 
				truck_analytics_2024_01_09."Wheel_formula" = '6x4'
				AND truck_analytics_2024_01_09."Brand" IN ('FAW', 'HOWO', 'JAC', 'SANY', 'SITRAK', 'SHACMAN', 'DONGFENG')
				AND truck_analytics_2024_01_09."Body_type" = 'Самосвал'
				AND truck_analytics_2024_01_09."Mass_in_segment_1" = '32001-40000'
			GROUP BY 
				truck_analytics_2024_01_09."Federal_district", 
				truck_analytics_2024_01_09."Region", 
				truck_analytics_2024_01_09."Brand"
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
			MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END) as JAC,
			MAX(CASE WHEN "Brand" = 'SANY' THEN total_sales END) as SANY,
			MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END) as SITRAK,
			MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END) as SHACMAN,
			MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END) as DONGFENG,
			COALESCE(MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SANY' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END), 0) as TOTAL
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

	// Соединение с базой данных
	db, err := db.Connect()
	if err != nil {
		slog.Warn("Can't connect to database")
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		response := TruckAnalyticsResponse{
			Error: "Failed to execute query: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	defer rows.Close()

	// Создаем упорядоченную карту для данных
	dataByDistrict := orderedmap.New[string, []TruckAnalytics]()

	// Определяем пользовательский порядок округов
	customOrder := []string{
		"Central",
		"North West",
		"Volga",
		"South",
		"North Caucasian",
		"Ural",
		"Siberia",
		"Far East",
	}

	// Предзаполняем карту пустыми значениями для каждого округа
	for _, district := range customOrder {
		dataByDistrict.Set(district, []TruckAnalytics{})
	}

	// Обработка результатов и группировка по федеральному округу
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
			&ta.SHACMAN,
			&ta.DONGFENG,
			&ta.TOTAL,
		)
		if err != nil {
			response := TruckAnalyticsResponse{
				Error: "Failed to scan row: " + err.Error(),
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		// Переводим федеральный округ и регион на английский
		if translated, ok := regionTranslations[federalDistrict]; ok {
			federalDistrict = translated
		}

		if translated, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = translated
		}

		// Рассчитываем общий рынок как сумму всех брендов для данного региона
		ta.TotalMarket = nullToZero(ta.FAW) + nullToZero(ta.HOWO) + nullToZero(ta.JAC) + nullToZero(ta.SANY) + nullToZero(ta.SITRAK) + nullToZero(ta.SHACMAN) + nullToZero(ta.DONGFENG)

		// Добавляем данные в соответствующий округ
		if existing, ok := dataByDistrict.Get(federalDistrict); ok {
			dataByDistrict.Set(federalDistrict, append(existing, ta))
		}
	}

	// Проверка на ошибки при итерации
	if err := rows.Err(); err != nil {
		response := TruckAnalyticsResponse{
			Error: "Error iterating over rows: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// Отправка ответа
	response := TruckAnalyticsResponse{
		Data: dataByDistrict,
	}
	ctx.JSON(http.StatusOK, response)
}

func Dumpers6x4WithTotalMarket2024(ctx *gin.Context) {
	// Структура для хранения данных округа
	type DistrictData struct {
		RegionName string `json:"region_name"`
		Faw        *int   `json:"faw"`
		Howo       *int   `json:"howo"`
		Jac        *int   `json:"jac"`
		Sany       *int   `json:"sany"`
		Sitrak     *int   `json:"sitrak"`
		Shacman    *int   `json:"shacman"`
		Dongfeng   *int   `json:"dongfeng"`
		Total      int    `json:"total"`
	}

	// Структура для ответа
	type Response struct {
		Data  *orderedmap.OrderedMap[string, []DistrictData] `json:"data"`
		Error string                                         `json:"error,omitempty"`
	}

	// Мапа для перевода русских названий федеральных округов на английский
	var regionTranslation = map[string]string{
		"Центральный Федеральный Округ":       "Central",
		"Северо-Западный Федеральный Округ":   "North West",
		"Южный Федеральный Округ":             "South",
		"Северо-Кавказский Федеральный Округ": "North Caucasian",
		"Приволжский Федеральный Округ":       "Volga",
		"Уральский Федеральный Округ":         "Ural",
		"Сибирский Федеральный Округ":         "Siberia",
		"Дальневосточный Федеральный Округ":   "Far East",
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
			COALESCE(SUM(CASE WHEN "Brand" = 'SHACMAN' THEN "Quantity" END), 0) AS SHACMAN,
			COALESCE(SUM(CASE WHEN "Brand" = 'DONGFENG' THEN "Quantity" END), 0) AS DONGFENG,
			COALESCE(SUM("Quantity"), 0) AS TOTAL
		FROM truck_analytics_2024_01_09
		WHERE 
			"Wheel_formula" = '6x4'
			AND "Brand" IN ('FAW', 'HOWO', 'JAC', 'SANY', 'SITRAK', 'SHACMAN', 'DONGFENG')
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
	summary.Shacman = new(int)
	summary.Dongfeng = new(int)

	for rows.Next() {
		var item DistrictData
		err := rows.Scan(&item.RegionName, &item.Faw, &item.Howo, &item.Jac, &item.Sany, &item.Sitrak, &item.Shacman, &item.Dongfeng, &item.Total)
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

		if item.Shacman != nil {
			*summary.Shacman += *item.Shacman
		}

		if item.Dongfeng != nil {
			*summary.Dongfeng += *item.Dongfeng
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

func NineMonth2024Dumpers8x4(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName  string `json:"region_name"`
		FAW         *int   `json:"faw"`
		HOWO        *int   `json:"howo"`
		SHACMAN     *int   `json:"shacman"`
		SITRAK      *int   `json:"sitrak"`
		TOTAL       int    `json:"total"`
		TotalMarket int    `json:"total_market"`
	}

	type TruckAnalyticsResponse struct {
		Data  *orderedmap.OrderedMap[string, []TruckAnalytics] `json:"data"`
		Error string                                           `json:"error,omitempty"`
	}

	query := `
		WITH base_data AS (
			SELECT 
				truck_analytics_2024_01_09."Federal_district",
				truck_analytics_2024_01_09."Region",
				truck_analytics_2024_01_09."Brand",
				SUM(truck_analytics_2024_01_09."Quantity") as total_sales
			FROM truck_analytics_2024_01_09
			WHERE 
				truck_analytics_2024_01_09."Wheel_formula" = '8x4'
				AND truck_analytics_2024_01_09."Brand" IN ('FAW', 'HOWO', 'SHACMAN', 'SITRAK')
				AND truck_analytics_2024_01_09."Body_type" = 'Самосвал'
				AND truck_analytics_2024_01_09."Weight_in_segment_4" = '35001-45000'
			GROUP BY 
				truck_analytics_2024_01_09."Federal_district", 
				truck_analytics_2024_01_09."Region", 
				truck_analytics_2024_01_09."Brand"
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

func Dumpers8x4WithTotalMarket2024(ctx *gin.Context) {
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
		"Центральный Федеральный Округ":       "Central",
		"Северо-Западный Федеральный Округ":   "North West",
		"Южный Федеральный Округ":             "South",
		"Северо-Кавказский Федеральный Округ": "North Caucasian",
		"Приволжский Федеральный Округ":       "Volga",
		"Уральский Федеральный Округ":         "Ural",
		"Сибирский Федеральный Округ":         "Siberia",
		"Дальневосточный Федеральный Округ":   "Far East",
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
		FROM truck_analytics_2024_01_09
		WHERE 
			"Wheel_formula" = '8x4'
			AND "Brand" IN ('FAW', 'HOWO', 'SITRAK', 'SHACMAN')
			AND "Body_type" = 'Самосвал'
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
