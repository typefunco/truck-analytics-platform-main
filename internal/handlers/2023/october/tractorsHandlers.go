package october

import (
	"context"
	"log/slog"
	"net/http"
	"truck-analytics-platform/internal/db"

	"github.com/gin-gonic/gin"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func TenMonth2023Tractors4x2(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName  string `json:"region_name"`
		DONGFENG    *int   `json:"dongfeng"`
		FAW         *int   `json:"faw"`
		FOTON       *int   `json:"foton"`
		JAC         *int   `json:"jac"`
		SHACMAN     *int   `json:"shacman"`
		SITRAK      *int   `json:"sitrak"`
		TOTAL       int    `json:"total"`
		TotalMarket int    `json:"total_market"`
	}

	type TruckAnalyticsResponse struct {
		Data  *orderedmap.OrderedMap[string, []TruckAnalytics] `json:"data"`
		Error string                                           `json:"error,omitempty"`
	}

	// SQL запрос для получения данных
	query := `
        WITH base_data AS (
            SELECT 
                truck_analytics_2023_01_12."Federal_district",
                truck_analytics_2023_01_12."Region",
                truck_analytics_2023_01_12."Brand",
                SUM(truck_analytics_2023_01_12."Quantity") as total_sales
            FROM truck_analytics_2023_01_12
            WHERE 
                truck_analytics_2023_01_12."Wheel_formula" = '4x2'
                AND truck_analytics_2023_01_12."Brand" IN ('DONGFENG', 'FAW', 'FOTON', 'JAC', 'SHACMAN', 'SITRAK')
                AND truck_analytics_2023_01_12."Body_type" = 'Седельный тягач'
                AND truck_analytics_2023_01_12."Exact_mass" = 18000
				AND truck_analytics_2023_01_12."Month_of_registration" <= 10
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
            MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END) as DONGFENG,
            MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END) as FAW,
            MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END) as FOTON,
            MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END) as JAC,
            MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END) as SHACMAN,
            MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END) as SITRAK,
            COALESCE(MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END), 0) +
            COALESCE(MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END), 0) +
            COALESCE(MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END), 0) +
            COALESCE(MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END), 0) +
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

	// Соединение с базой данных
	db, err := db.Connect()
	if err != nil {
		slog.Warn("Can't connect to database")
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Can't connect to database"})
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to execute query: " + err.Error()})
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

	// Обработка данных из результата SQL запроса
	for rows.Next() {
		var ta TruckAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
			&ta.RegionName,
			&ta.DONGFENG,
			&ta.FAW,
			&ta.FOTON,
			&ta.JAC,
			&ta.SHACMAN,
			&ta.SITRAK,
			&ta.TOTAL,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to scan row: " + err.Error()})
			return
		}

		// Перевод федерального округа, если есть в маппинге
		if translatedDistrict, ok := districtTranslations[federalDistrict]; ok {
			federalDistrict = translatedDistrict
		}

		// Перевод региона, если есть в маппинге
		if translatedRegion, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = translatedRegion
		}

		// Рассчитываем общий рынок
		ta.TotalMarket = nulltoZero(ta.DONGFENG) + nulltoZero(ta.FAW) + nulltoZero(ta.FOTON) + nulltoZero(ta.JAC) + nulltoZero(ta.SHACMAN) + nulltoZero(ta.SITRAK)

		// Добавляем данные в соответствующий округ
		if existing, ok := dataByDistrict.Get(federalDistrict); ok {
			dataByDistrict.Set(federalDistrict, append(existing, ta))
		}
	}

	// Проверка на ошибки при итерации
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Error iterating over rows: " + err.Error()})
		return
	}

	// Отправка ответа
	ctx.JSON(http.StatusOK, TruckAnalyticsResponse{
		Data: dataByDistrict,
	})
}

func TenTractors4x2WithTotalMarket2023(ctx *gin.Context) {
	// Структура для хранения данных округа
	type DistrictData struct {
		RegionName string `json:"region_name"`
		Dongfeng   *int   `json:"dongfeng"`
		Faw        *int   `json:"faw"`
		Foton      *int   `json:"foton"`
		Jac        *int   `json:"jac"`
		Shacman    *int   `json:"shacman"`
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
		"Северо-Западный Федеральный Округ":   "North West",
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
			COALESCE(SUM(CASE WHEN "Brand" = 'DONGFENG' THEN "Quantity" END), 0) AS DONGFENG,
			COALESCE(SUM(CASE WHEN "Brand" = 'FAW' THEN "Quantity" END), 0) AS FAW,
			COALESCE(SUM(CASE WHEN "Brand" = 'FOTON' THEN "Quantity" END), 0) AS FOTON,
			COALESCE(SUM(CASE WHEN "Brand" = 'JAC' THEN "Quantity" END), 0) AS JAC,
			COALESCE(SUM(CASE WHEN "Brand" = 'SHACMAN' THEN "Quantity" END), 0) AS SHACMAN,
			COALESCE(SUM(CASE WHEN "Brand" = 'SITRAK' THEN "Quantity" END), 0) AS SITRAK,
			COALESCE(SUM("Quantity"), 0) AS TOTAL
		FROM truck_analytics_2023_01_12
		WHERE 
			"Wheel_formula" = '4x2'
			AND "Brand" IN ('DONGFENG', 'FAW', 'FOTON', 'JAC', 'SHACMAN', 'SITRAK')
			AND "Month_of_registration" <= 10
			AND "Body_type" = 'Седельный тягач'
			AND "Exact_mass" = 18000
		GROUP BY "Federal_district"
		ORDER BY "Federal_district"
	`

	// Подключение к базе данных
	db, err := db.Connect()
	if err != nil {
		slog.Info("Can't connect to database:", err)
		ctx.JSON(http.StatusInternalServerError, Response{Error: "Can't connect to database"})
		return
	}
	defer db.Close(context.Background())

	// Запрос к базе данных
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		slog.Info("Failed to execute query:", err)
		ctx.JSON(http.StatusInternalServerError, Response{Error: "Failed to execute query"})
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
	summary.Dongfeng = new(int)
	summary.Faw = new(int)
	summary.Foton = new(int)
	summary.Jac = new(int)
	summary.Shacman = new(int)
	summary.Sitrak = new(int)

	// Обработка данных из результата SQL запроса
	for rows.Next() {
		var item DistrictData
		err := rows.Scan(&item.RegionName, &item.Dongfeng, &item.Faw, &item.Foton, &item.Jac, &item.Shacman, &item.Sitrak, &item.Total)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, Response{Error: "Failed to scan row"})
			return
		}

		// Обновляем суммарные значения
		if item.Dongfeng != nil {
			*summary.Dongfeng += *item.Dongfeng
		}
		if item.Faw != nil {
			*summary.Faw += *item.Faw
		}
		if item.Foton != nil {
			*summary.Foton += *item.Foton
		}
		if item.Jac != nil {
			*summary.Jac += *item.Jac
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
		ctx.JSON(http.StatusInternalServerError, Response{Error: "Failed to iterate over rows"})
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

func TenMonth2023Tractors6x4(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName string `json:"region_name"`
		DONGFENG   *int   `json:"dongfeng"`
		FAW        *int   `json:"faw"`
		FOTON      *int   `json:"foton"`
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
				truck_analytics_2023_01_12."Wheel_formula" = '6x4'
				AND truck_analytics_2023_01_12."Brand" IN ('DONGFENG', 'FAW', 'FOTON', 'HOWO', 'SHACMAN', 'SITRAK')
				AND truck_analytics_2023_01_12."Month_of_registration" <= 10
				AND truck_analytics_2023_01_12."Body_type" = 'Седельный тягач'
				AND truck_analytics_2023_01_12."Exact_mass" = 25000
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
			MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END) as DONGFENG,
			MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END) as FAW,
			MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END) as FOTON,
			MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END) as HOWO,
			MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END) as SHACMAN,
			MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END) as SITRAK,
			COALESCE(MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END), 0) +
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

	// Соединение с базой данных
	db, err := db.Connect()
	if err != nil {
		slog.Warn("Can't connect to database")
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Can't connect to database"})
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to execute query: " + err.Error()})
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

	for rows.Next() {
		var ta TruckAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
			&ta.RegionName,
			&ta.DONGFENG,
			&ta.FAW,
			&ta.FOTON,
			&ta.HOWO,
			&ta.SHACMAN,
			&ta.SITRAK,
			&ta.TOTAL,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to scan row: " + err.Error()})
			return
		}

		// Перевод федерального округа, если есть в маппинге
		if translatedDistrict, ok := districtTranslations[federalDistrict]; ok {
			federalDistrict = translatedDistrict
		}

		// Перевод региона, если есть в маппинге
		if translatedRegion, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = translatedRegion
		}

		// Рассчитываем общий рынок
		ta.TOTAL = nulltoZero(ta.DONGFENG) + nulltoZero(ta.FAW) + nulltoZero(ta.FOTON) + nulltoZero(ta.HOWO) + nulltoZero(ta.SHACMAN) + nulltoZero(ta.SITRAK)

		// Добавляем данные в соответствующий округ
		if existing, ok := dataByDistrict.Get(federalDistrict); ok {
			dataByDistrict.Set(federalDistrict, append(existing, ta))
		}
	}

	// Проверка на ошибки при итерации
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Error iterating over rows: " + err.Error()})
		return
	}

	// Отправка ответа
	ctx.JSON(http.StatusOK, TruckAnalyticsResponse{
		Data: dataByDistrict,
	})
}

func TenTractors6x4WithTotalMarket2023(ctx *gin.Context) {
	// Структура для хранения данных округа
	type DistrictData struct {
		DistrictName string `json:"region_name"`
		DONGFENG     *int   `json:"dongfeng"`
		FAW          *int   `json:"faw"`
		FOTON        *int   `json:"foton"`
		HOWO         *int   `json:"howo"`
		SHACMAN      *int   `json:"shacman"`
		SITRAK       *int   `json:"sitrak"`
		TOTAL        int    `json:"total"`
	}

	type Response struct {
		Data  *orderedmap.OrderedMap[string, []DistrictData] `json:"data"`
		Error string                                         `json:"error,omitempty"`
	}
	// Мапа для перевода русских названий федеральных округов на английский
	var regionTranslation = map[string]string{
		"Дальневосточный Федеральный Округ":   "Far East",
		"Приволжский Федеральный Округ":       "Volga",
		"Северо-Западный Федеральный Округ":   "North West",
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
			COALESCE(SUM(CASE WHEN "Brand" = 'DONGFENG' THEN "Quantity" END), 0) AS DONGFENG,
			COALESCE(SUM(CASE WHEN "Brand" = 'FAW' THEN "Quantity" END), 0) AS FAW,
			COALESCE(SUM(CASE WHEN "Brand" = 'FOTON' THEN "Quantity" END), 0) AS FOTON,
			COALESCE(SUM(CASE WHEN "Brand" = 'HOWO' THEN "Quantity" END), 0) AS HOWO,
			COALESCE(SUM(CASE WHEN "Brand" = 'SHACMAN' THEN "Quantity" END), 0) AS SHACMAN,
			COALESCE(SUM(CASE WHEN "Brand" = 'SITRAK' THEN "Quantity" END), 0) AS SITRAK,
			COALESCE(SUM("Quantity"), 0) AS TOTAL
		FROM truck_analytics_2023_01_12
		WHERE
			"Wheel_formula" = '6x4'
			AND "Brand" IN ('DONGFENG', 'FAW', 'FOTON', 'HOWO', 'SHACMAN', 'SITRAK')
			AND "Month_of_registration" <= 10
			AND "Body_type" = 'Седельный тягач'
			AND "Exact_mass" = 25000
		GROUP BY "Federal_district"
		ORDER BY "Federal_district"
	`

	db, err := db.Connect()
	if err != nil {
		slog.Info("Can't connect to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't connect to database"})
		return
	}
	defer db.Close(context.Background())

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		slog.Info("Failed to execute query:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	dataByDistrict := orderedmap.New[string, []DistrictData]()
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

	for _, district := range customOrder {
		dataByDistrict.Set(district, []DistrictData{})
	}

	summary := DistrictData{
		DistrictName: "Summary",
		DONGFENG:     new(int),
		FAW:          new(int),
		FOTON:        new(int),
		HOWO:         new(int),
		SHACMAN:      new(int),
		SITRAK:       new(int),
		TOTAL:        0,
	}

	for rows.Next() {
		var ba DistrictData

		err := rows.Scan(
			&ba.DistrictName,
			&ba.DONGFENG,
			&ba.FAW,
			&ba.FOTON,
			&ba.HOWO,
			&ba.SHACMAN,
			&ba.SITRAK,
			&ba.TOTAL,
		)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}

		translatedRegionName, exists := regionTranslation[ba.DistrictName]
		if !exists {
			translatedRegionName = ba.DistrictName
		}

		if existing, exists := dataByDistrict.Get(translatedRegionName); exists {
			ba.DistrictName = translatedRegionName
			dataByDistrict.Set(translatedRegionName, append(existing, ba))
		}

		// Агрегируем данные для Summary
		*summary.DONGFENG += *ba.DONGFENG
		*summary.FAW += *ba.FAW
		*summary.FOTON += *ba.FOTON
		*summary.HOWO += *ba.HOWO
		*summary.SHACMAN += *ba.SHACMAN
		*summary.SITRAK += *ba.SITRAK
		summary.TOTAL += ba.TOTAL
	}

	if rows.Err() != nil {
		slog.Info("Failed to iterate over rows:", rows.Err())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over rows"})
		return
	}

	// Добавляем итоговые данные в Summary
	dataByDistrict.Set("Summary", []DistrictData{summary})

	ctx.JSON(http.StatusOK, Response{
		Data: dataByDistrict,
	})
}
